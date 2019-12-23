package main

import (
	"flag"
	"fmt"
	"nexus-firewall-for-athens/cmd"
	"os"
)

func main() {
	var runAs string
	var port int

	flag.StringVar(&runAs, "run", "lambda", "run as either \"lambda\" (default), \"server\"")
	// Server configuration
	flag.IntVar(&port, "port", 8080, "port to use when running in server mode (default: 10000)")
	flag.Parse()

	fmt.Println("Starting as", runAs)

	switch runAs {
	case "lambda":
		cmd.HandleLambda()
		break
	case "server":
		cmd.HandleRequests(port)
		break
	default:
		fmt.Println("Unknown runAs defined {}, exiting", runAs)
		os.Exit(1)
	}
	fmt.Println("Finished")
}
