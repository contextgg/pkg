package es

import (
	"sync"
)

type Config interface {
	Aggregate(fn EntityFunc, opts ...EntityOption) *AggregateConfig
	Saga(s Saga) *SagaConfig
}

type config struct {
	sync.RWMutex

	Aggregates []*AggregateConfig
}

func (cfg *config) Aggregate(fn EntityFunc, opts ...EntityOption) *AggregateConfig {
	cfg.Lock()
	defer cfg.Unlock()

	all := append(opts, EntityFactory(fn))
	return &AggregateConfig{
		opts: NewEntityOptions(all),
	}
}

func (cfg *config) Saga(s Saga) *SagaConfig {
	cfg.Lock()
	defer cfg.Unlock()

	return &SagaConfig{
		s: s,
	}
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
