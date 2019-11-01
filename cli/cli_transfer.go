package cli

import (
	block "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/network"
)

func (cli *Cli) transfer(from, to string, amount string) {
	blc:=block.NewBlockchain()
	s:=network.Send{}
	blc.CreateTransaction(from,to,amount,s)
}
