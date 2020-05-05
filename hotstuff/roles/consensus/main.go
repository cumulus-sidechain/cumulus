package main

import (
	"os"

	"github.com/cumulus-sidechain/cumulus/hotstuff/roles/consensus/node"
)

func main() {
	defer os.Exit(0)

	cmd := node.CommandLine{}
	cmd.Run()
}
