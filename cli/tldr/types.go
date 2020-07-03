package main

type runnable interface {
	Init()
	ParseArgs(subcommand string, args ...string) error
	Execute(subcommand string, args ...string) error
	Help(subcommand string, args ...string)
}
