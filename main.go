package main

import (
	"fmt"
	"os"

	"github.com/mheers/k3dnifi/cmd"
	"github.com/mheers/k3dnifi/nifi"
	"github.com/mheers/k3droot/helpers"
	"github.com/sirupsen/logrus"
)

// build flags
var (
	VERSION    string
	BuildTime  string
	CommitHash string
	GoVersion  string
	GitTag     string
	GitBranch  string
)

func main() {
	_, err := nifi.Init()
	if err != nil {
		panic(err)
	}

	k3d, err := helpers.IsK3d()
	if err != nil {
		panic(err)
	}
	if !k3d {
		fmt.Println("Not a k3d cluster")
		os.Exit(1)
	}

	cmd.VERSION = VERSION
	cmd.BuildTime = BuildTime
	cmd.CommitHash = CommitHash
	cmd.GoVersion = GoVersion
	cmd.GitTag = GitTag
	cmd.GitBranch = GitBranch

	// execeute the command
	err = cmd.Execute()
	if err != nil {
		logrus.Fatalf("Execute failed: %+v", err)
	}
}
