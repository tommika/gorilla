<!--
Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
-->

xflags
======

This package provides a powerful declarative model for Go's built-in command
line processing package (flags.)

Rather than programmatically building up a FlagSet, do it with grace and ease using
structured tags:

```go
type MyOptions struct {
	Verbose   bool     // use defaults
	Debug     bool     `flag:"d|debug,Enable debug output"`
	Port      string   `flag:"p|port,Port to listen on,8080"`
	Timeout   int      `flag:",Time-out in seconds,10"`
	Hosts     []string `flag:"*"` // remaining args collected here
	internal  int      `flag:"-"` // ignore this one
}
...
opts := MyOptions{}
if err := xflags.ParseArgs(os.Args[0],&opts,os.Args[1:]); err != nil {
	os.Exit(1)
}
if opts.Debug {
	fmt.Printf("Debug is enabled")
}
for _, host := range opts.Hosts {
	...
}
...
