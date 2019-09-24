package cli

import block "myCode/public_blockchain/part7-network/blc"

func (cli *Cli) printAllBlock() {
	bc := block.NewBlockchain()
	defer bc.BD.Close()
	bc.PrintAllBlockInfo()
}
