package es

import (
	"sync"
)

type Config interface {
	Aggregate(fn EntityFunc, opts ...EntityOption) *AggregateConfig
	Saga(s Saga) *SagaConfig

	GetAggregates() []*AggregateConfig
	GetSagas() []*SagaConfig
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
		EntityOptions: NewEntityOptions(all),
	}

	cfg.Aggregates = append(cfg.Aggregates, out)
	return out
}

func (cfg *config) Saga(s Saga) *SagaConfig {
	cfg.Lock()
	defer cfg.Unlock()

	out := &SagaConfig{
		Saga: s,
	}

	cfg.Sagas = append(cfg.Sagas, out)
	return out
}

func (cfg *config) GetAggregates() []*AggregateConfig {
	cfg.RLock()
	defer cfg.RUnlock()

	return cfg.Aggregates
}
func (cfg *config) GetSagas() []*SagaConfig {
	cfg.RLock()
	defer cfg.RUnlock()

	return cfg.Sagas
}

type SagaConfig struct {
	sync.RWMutex
	Saga

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
	EntityOptions

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
