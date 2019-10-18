package cli

import block "github.com/corgi-kx/blockchain_golang/blc"

func (cli *Cli) printAllBlock() {
	bc := block.NewBlockchain()
	bc.PrintAllBlockInfo()
}
