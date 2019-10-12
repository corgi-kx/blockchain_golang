package cli

import (
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/database"
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
