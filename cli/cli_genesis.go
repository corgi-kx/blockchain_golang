package cli

import (
	block "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/network"
)

func (cli *Cli) genesis(address string, value int) {
	bc := block.NewBlockchain()
	bc.CreataGenesisTransaction(address, value, network.Send{})
}
