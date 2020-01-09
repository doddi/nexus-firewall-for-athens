package main

import (
	"flag"
	"fmt"
	"nexus-firewall-for-athens/cmd"
	"nexus-firewall-for-athens/ossindex"
	"nexus-firewall-for-athens/validate"
	"os"
)

func main() {
	var runAs, using string
	var port int
	var failBuild bool

	extractParameters(&runAs, &using, &port, &failBuild)

	validator := validateUsing(using)

	server := startServerType(runAs, port, failBuild, validator)

	fmt.Println("Starting as", runAs)
	server.Handle()
	fmt.Println("Finished")
}

func extractParameters(runAs *string, using *string, port *int, failBuild *bool) {
	flag.StringVar(runAs, "run", "lambda", "run as either \"lambda\" (default), \"server\"")
	flag.StringVar(using, "using", "ossindex", "Use \"ossindex\" or \"nexusiq\", (default) \"ossindex\"")
	// Server configuration
	flag.IntVar(port, "port", 8080, "port to use when running in server mode (default: 8080)")
	flag.BoolVar(failBuild, "failBuild", true, "set to true to return 403 on security vulnerability discovery")
	flag.Parse()
}

func startServerType(runAs string, port int, failBuild bool, validator validate.Validator) cmd.Environment {
	var server cmd.Environment

	switch runAs {
	case "lambda":
		server = cmd.LambdaServer{Validator: validator}
		break
	case "server":
		server = cmd.LocalServer{Port: port, FailBuild: failBuild, Validator: validator}
		break
	default:
		fmt.Println("Unknown runAs defined {}, exiting", runAs)
		os.Exit(1)
	}
	return server
}

func validateUsing(using string) validate.Validator {
	switch using {
	case "ossindex":
		return ossindex.OssIndex{}
	case "nexusiq":
		fmt.Println("Unimplemented")
		os.Exit(1)
	default:
		fmt.Println("Unknown validator")
		os.Exit(1)
	}
	return nil
}
