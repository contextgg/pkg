package es

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/contextgg/pkg/types"
)

// CommandRegistry stores the handlers for commands
type CommandRegistry interface {
	SetHandler(CommandHandler, ...Command)
	GetHandler(Command) (CommandHandler, error)
	NewCommand(string) (Command, error)
}

// NewCommandRegistry creates a new CommandRegistry
func NewCommandRegistry() CommandRegistry {
	return &commandRegistry{
		handlers:      make(map[string]CommandHandler),
		typesRegistry: types.NewRegistry(),
	}
}

type commandRegistry struct {
	sync.RWMutex
	handlers      map[string]CommandHandler
	typesRegistry types.Registry
}

func (r *commandRegistry) SetHandler(handler CommandHandler, cmds ...Command) {
	r.Lock()
	defer r.Unlock()

	for _, cmd := range cmds {
		entry := types.EntryFromType(cmd, true)
		if err := r.typesRegistry.Add(entry); err != nil {
			panic(err)
		}
		for _, name := range entry.Names {
			r.handlers[name] = handler
		}
	}
}

func (r *commandRegistry) GetHandler(cmd Command) (CommandHandler, error) {
	if cmd == nil {
		return nil, errors.New("You need to supply a command")
	}

	name := types.GetTypeName(cmd)
	handler, ok := r.handlers[name]
	if !ok {
		return nil, fmt.Errorf("Cannot find %s in registry", name)
	}
	return handler, nil
}

func (r *commandRegistry) NewCommand(name string) (Command, error) {
	names := []string{name}
	if strings.HasSuffix(strings.ToLower(name), "command") {
		names = append(names, name[:len(name)-7])
	}

	entry, ok := types.GetFirstByNames(r.typesRegistry, names)
	if !ok {
		return nil, fmt.Errorf("Cannot find %s in registry", name)
	}

	obj := entry.Factory()
	return obj.(Command), nil
}
