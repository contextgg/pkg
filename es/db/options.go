package db

import (
	"fmt"
)

type Options struct {
	Host     string
	Port     int
	DbName   string
	User     string
	Password string
	Debug    bool
	Recreate bool
	Migrate  bool
}

func (o *Options) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", o.User, o.Password, o.Host, o.Port, o.DbName)
}

func NewOptions() *Options {
	return &Options{
		Host: "localhost",
		Port: 5432,
	}
}

type OptionFunc func(*Options)

func WithDbHost(host string) OptionFunc {
	return func(o *Options) {
		if host != "" {
			o.Host = host
		}
	}
}

func WithDbPort(port int) OptionFunc {
	return func(o *Options) {
		o.Port = port
	}
}

func WithDbName(dbName string) OptionFunc {
	return func(o *Options) {
		o.DbName = dbName
	}
}

func WithDbUser(user string) OptionFunc {
	return func(o *Options) {
		o.User = user
	}
}

func WithDbPassword(password string) OptionFunc {
	return func(o *Options) {
		o.Password = password
	}
}

func WithDebug(debug bool) OptionFunc {
	return func(o *Options) {
		o.Debug = debug
	}
}

func Recreate(recreate bool) OptionFunc {
	return func(o *Options) {
		o.Recreate = recreate
	}
}
