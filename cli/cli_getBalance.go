package cli

import (
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
)

func (cli *Cli) getBalance(address string) {
	bc := block.NewBlockchain()
	balance := bc.GetBalance(address)
	fmt.Printf("地址:%s的余额为：%d\n", address, balance)
}
