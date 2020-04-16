package main

type runnable interface {
	Execute(subcommand string, args ...string) error
	Help(subcommand string, args ...string)
}
