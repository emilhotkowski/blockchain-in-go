package main

import (
	"github.com/emilhotkowski/blockchain-in-go/pkg/blockchain"
	"github.com/emilhotkowski/blockchain-in-go/pkg/cli"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.Db.Close()

	cli := cli.CLI{bc}
	cli.Run()
}
