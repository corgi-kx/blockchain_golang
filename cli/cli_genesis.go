package cli

import block "myCode/public_blockchain/part7-network/blc"

func (cli *Cli) genesis(address string, value int) {
	bc := block.NewBlockchain()
	defer bc.BD.Close()
	bc.CreataGenesisTransaction(address, value)
}
