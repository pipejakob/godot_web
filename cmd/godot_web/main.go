package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pipejakob/godot_web"
)

var port = flag.Int("port", 8000, "port to listen on")
var show_version = flag.Bool("version", false, "print the current version and exit")

func main() {
	flag.Parse()

	if flag.NArg() != 0 {
		fmt.Printf("Unexpected argument(s): %q\n", flag.Args())
		usage()
		os.Exit(1)
	}

	if *show_version {
		fmt.Printf("godot_web %s\n", version())
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("error getting working directory: %v", err))
	}

	server := godot_web.New(dir, *port)
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("error running server: %v", err))
	}
}

func usage() {
	fmt.Println("Usage: godot_web [OPTIONS]")
	flag.PrintDefaults()
}
