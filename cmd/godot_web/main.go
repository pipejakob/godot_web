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
var dir string // the last argument (optional)

func usage() {
	fmt.Fprintln(flag.CommandLine.Output(), "Usage: godot_web [OPTIONS] [DIR]\n\n"+
		"If DIR is omitted, the current working directory is used.\n\n"+
		"Options:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	validateFlags()

	if *showVersion {
		fmt.Printf("godot_web %s\n", version())
		return
	}

	server := godot_web.New(dir, *port, *allowExternal, *tlsCert, *tlsKey)
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("error running server: %v", err))
	}
}

func validateFlags() {
	out := flag.CommandLine.Output()

	if flag.NArg() > 1 {
		fmt.Fprintf(out, "Too many trailing arguments: %q\n", flag.Args())
		usage()
		os.Exit(1)
	}

	if flag.NArg() == 1 {
		dir = flag.Arg(0)

		if err := verifyDir(dir); err != nil {
			fmt.Fprintf(out, "Problem with directory %q: %v\n", dir, err)
			usage()
			os.Exit(1)
		}
	} else {
		var err error
		if dir, err = os.Getwd(); err != nil {
			panic(fmt.Errorf("error getting working directory: %v", err))
		}
	}

	if *tlsCert != "" || *tlsKey != "" {
		if *tlsCert == "" || *tlsKey == "" {
			fmt.Fprintln(out, "If you specify one of --tls-cert or --tls-key, you must specify them both.")
			usage()
			os.Exit(1)
		}

		if !*allowExternal {
			fmt.Fprintln(out, "The flags --tls-cert and --tls-key can only be used when --external is also given.")
			usage()
			os.Exit(1)
		}
	}
}

func verifyDir(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("not a directory")
	}

	return nil
}
