package cli

import (
	"encoding/json"
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/database"
	log "github.com/corgi-kx/logcustom"
)

func (cli *Cli) importWalletByMnemonicword(mnemonicword string) {
	mnemonicwords := []string{}
	err := json.Unmarshal([]byte(mnemonicword), &mnemonicwords)
	if err != nil {
		log.Error("json err:", err)
	}

	bd := database.New()
	wallets := block.NewWallets(bd)
	address, privkey, mnemonicWord := wallets.GenerateWallet(bd, block.CreateBitcoinKeysByMnemonicWord, mnemonicwords)
	fmt.Println("助记词：", mnemonicWord)
	fmt.Println("私钥：", privkey)
	fmt.Println("地址：", address)
}
