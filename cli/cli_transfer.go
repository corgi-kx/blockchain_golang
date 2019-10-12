package cli

import block "github.com/corgi-kx/blockchain_golang/blc"

func (cli *Cli) transfer(from, to string, amount string) {
	bc := block.NewBlockchain()
	defer bc.BD.Close()
	bc.Transfer(from, to, amount)
}
