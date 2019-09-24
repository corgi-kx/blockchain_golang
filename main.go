package main

import (
	block "myCode/public_blockchain/part7-network/cli"
)

func main() {
	cli := block.New()
	cli.Run()
}
