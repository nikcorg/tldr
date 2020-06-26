package main

type runnable interface {
	ParseArgs(subcommand string, args ...string) error
	Execute(subcommand string, args ...string) error
	Help(subcommand string, args ...string)
}
