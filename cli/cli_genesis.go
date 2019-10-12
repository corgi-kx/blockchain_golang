package cli

import block "github.com/corgi-kx/blockchain_golang/blc"

func (cli *Cli) genesis(address string, value int) {
	bc := block.NewBlockchain()
	defer bc.BD.Close()
	bc.CreataGenesisTransaction(address, value)
}
