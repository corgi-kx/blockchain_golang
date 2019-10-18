package cli

import (
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
)

func (cli *Cli) setRewardAddress(address string) {
	bc := block.NewBlockchain()
	bc.SetRewardAddress(address)
	fmt.Printf("已设置地址%s为挖矿奖励地址！\n", address)
}
