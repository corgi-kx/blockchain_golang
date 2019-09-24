package cli

import (
	"fmt"
	block "myCode/public_blockchain/part7-network/blc"
	"myCode/public_blockchain/part7-network/database"
)

func (cli *Cli) printAllAddress() {
	bd := database.New(nodeID)
	defer bd.Close()
	addressList := block.GetAllAddress(bd)
	fmt.Println("===========================================")
	fmt.Println("已生成地址：")
	for _, v := range *addressList {
		fmt.Println(string(v))
	}
	fmt.Println("===========================================")

}
