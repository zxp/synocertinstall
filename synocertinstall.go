package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	version, tagName, branch, commitID, buildTime string
)

func main() {
	version = fmt.Sprintf("Version: %s, Branch: %s, Build: %s, Build time: %s", tagName, branch, commitID, buildTime)

	fmt.Println("Synology NAS certification install tool")

	flag.Usage = func() {
		fmt.Println(version)
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	var listFlag bool
	flag.BoolVar(&listFlag, "list", false, "list applications")

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Parse()
}
