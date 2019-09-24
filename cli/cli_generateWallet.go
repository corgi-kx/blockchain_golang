package cli

import (
	"fmt"
	block "myCode/public_blockchain/part7-network/blc"
	"myCode/public_blockchain/part7-network/database"
)

func (cli *Cli) generateWallet() {
	bd := database.New(nodeID)
	defer bd.Close()
	wallets := block.NewWallets(bd)
	address := wallets.GenerateWallet(bd)
	fmt.Println("新生成的地址为：", address)
}
