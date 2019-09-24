package cli

import (
	"fmt"
	block "myCode/public_blockchain/part7-network/blc"
)

func (cli *Cli) resetUTXODB() {
	bc := block.NewBlockchain()
	defer bc.BD.Close()
	utxos := block.UTXOHandle{bc}
	utxos.ResetUTXODataBase()
	fmt.Println("已重置UTXO数据库")
}
