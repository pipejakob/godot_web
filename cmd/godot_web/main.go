package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pipejakob/godot_web"
)

var port = flag.Int("port", 8000, "port to listen on")
var allowExternal = flag.Bool("external", false, "allow external traffic")
var tlsCert = flag.String("tls-cert", "", "path to certificate file (in PEM format)")
var tlsKey = flag.String("tls-key", "", "path to private key file (in PEM format)")
var showVersion = flag.Bool("version", false, "print the current version and exit")

func main() {
	flag.Parse()

	validateFlags()

	if *showVersion {
		fmt.Printf("godot_web %s\n", version())
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("error getting working directory: %v", err))
	}

	server := godot_web.New(dir, *port, *allowExternal, *tlsCert, *tlsKey)
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("error running server: %v", err))
	}
}

func validateFlags() {
	if flag.NArg() != 0 {
		fmt.Printf("Unexpected argument(s): %q\n", flag.Args())
		usage()
		os.Exit(1)
	}

	if *tlsCert != "" || *tlsKey != "" {
		if *tlsCert == "" || *tlsKey == "" {
			fmt.Println("If you specify one of --tls-cert or --tls-key, you must specify them both")
			usage()
			os.Exit(1)
		}

		if !*allowExternal {
			fmt.Println("The flags --tls-cert and --tls-key can only be used when --external is also given")
			usage()
			os.Exit(1)
		}
	}
}

func usage() {
	fmt.Println("Usage: godot_web [OPTIONS]")
	flag.PrintDefaults()
}
