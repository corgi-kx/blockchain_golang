package cli

import block "myCode/public_blockchain/part7-network/blc"

func (cli *Cli) transfer(from, to string, amount string) {
	bc := block.NewBlockchain()
	defer bc.BD.Close()
	bc.Transfer(from, to, amount)
}
