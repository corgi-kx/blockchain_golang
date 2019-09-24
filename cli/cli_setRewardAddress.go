package cli

import (
	"fmt"
	block "myCode/public_blockchain/part7-network/blc"
)

func (cli *Cli) setRewardAddress(address string) {
	bc := block.NewBlockchain()
	defer bc.BD.Close()
	bc.SetRewardAddress(address)
	fmt.Printf("已设置地址%s为挖矿奖励地址！\n", address)
}
