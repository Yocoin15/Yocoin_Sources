// Authored and revised by YOC team, 2016-2018
// License placeholder #1

// Command bzzhash computes a swarm tree hash.
package main

import (
	"fmt"
	"os"

	"github.com/Yocoin15/Yocoin_Sources/cmd/utils"
	"github.com/Yocoin15/Yocoin_Sources/swarm/storage"
	"gopkg.in/urfave/cli.v1"
)

func hash(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) < 1 {
		utils.Fatalf("Usage: swarm hash <file name>")
	}
	f, err := os.Open(args[0])
	if err != nil {
		utils.Fatalf("Error opening file " + args[1])
	}
	defer f.Close()

	stat, _ := f.Stat()
	chunker := storage.NewTreeChunker(storage.NewChunkerParams())
	key, err := chunker.Split(f, stat.Size(), nil, nil, nil)
	if err != nil {
		utils.Fatalf("%v\n", err)
	} else {
		fmt.Printf("%v\n", key)
	}
}
