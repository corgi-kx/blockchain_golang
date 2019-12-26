package cli

import (
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/network"
)

func (cli Cli) transfer(from, to, amount string) {
	blc := block.NewBlockchain()
	blc.CreateTransaction(from, to, amount, network.Send{})
	fmt.Println("已执行转帐命令")
}
