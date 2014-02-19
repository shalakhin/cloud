package main

import (
	"flag"
	"fmt"
)

const (
	// VERSION of package
	VERSION = "VERSION\n\t 0.1\n\n"
	// DESCRIPTION what is package for
	DESCRIPTION = "DESCRIPTION\n\tSync your data with cloud storages (like Amazon S3, Rackspace CloudFiles etc.)\n\n"
	// COMMANDS lists commands available
	COMMANDS = "COMMANDS\n\tinit\tinitialize .cloudrc and .cloud files\n\tsync\tsynchronize folder with the cloud\n\thelp\tshow this message\n"
)

func usage() {
	fmt.Printf(DESCRIPTION + VERSION + COMMANDS)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	switch {
	case args[0] == "init":
		initConfigs()
	case args[0] == "sync":
		fmt.Println("Nothing here yet")
	case args[0] == "help":
		usage()
	default:
		usage()
	}
}
