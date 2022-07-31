package es

import (
	"sync"
)

type Config interface {
	Aggregate(fn EntityFunc, opts ...EntityOption) *AggregateConfig
	Saga(s Saga) *SagaConfig

	GetEntities() []Entity
}

type config struct {
	sync.RWMutex

	Aggregates []*AggregateConfig
	Sagas      []*SagaConfig
}

func (cfg *config) Aggregate(fn EntityFunc, opts ...EntityOption) *AggregateConfig {
	cfg.Lock()
	defer cfg.Unlock()

	all := append(opts, EntityFactory(fn))
	out := &AggregateConfig{
		opts: NewEntityOptions(all),
	}

	cfg.Aggregates = append(cfg.Aggregates, out)
	return out
}

func (cfg *config) Saga(s Saga) *SagaConfig {
	cfg.Lock()
	defer cfg.Unlock()

	out := &SagaConfig{
		s: s,
	}

	cfg.Sagas = append(cfg.Sagas, out)
	return out
}

func (cfg *config) GetEntities() []Entity {
	cfg.RLock()
	defer cfg.RUnlock()

	var entities []Entity
	for _, a := range cfg.Aggregates {
		obj := a.opts.Factory("")
		entities = append(entities, obj)
	}
	return entities
}

type SagaConfig struct {
	sync.RWMutex

	s      Saga
	events []interface{}
}

func (cfg *SagaConfig) Events(evts ...interface{}) *SagaConfig {
	cfg.Lock()
	defer cfg.Unlock()

	cfg.events = append(cfg.events, evts...)
	return cfg
}

type AggregateConfig struct {
	sync.RWMutex

	opts     EntityOptions
	commands []Command
}

func (cfg *AggregateConfig) Commands(cmds ...Command) *AggregateConfig {
	cfg.Lock()
	defer cfg.Unlock()

	cfg.commands = append(cfg.commands, cmds...)
	return cfg
}

func NewConfig() Config {
	return &config{}
}
