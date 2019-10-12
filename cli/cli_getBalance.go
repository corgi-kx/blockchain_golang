package cli

import (
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
)

func (cli *Cli) getBalance(address string) {
	bc := block.NewBlockchain()
	defer bc.BD.Close()
	balance := bc.GetBalance(address)
	fmt.Printf("用户:%s的余额为：%d\n", address, balance)
}
