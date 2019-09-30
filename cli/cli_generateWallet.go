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
	address,privkey,mnemonicWord := wallets.GenerateWallet(bd)
	fmt.Println("助记词：",mnemonicWord)
	fmt.Println("私钥：", privkey)
	fmt.Println("地址：", address)


}
