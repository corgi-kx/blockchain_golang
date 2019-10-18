package cli

import (
	"encoding/json"
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/database"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
)

func  (cli *Cli)importWalletByMnemonicword(j string){
	mnemonicwords:=[]string{}
	err := json.Unmarshal([]byte(j), &mnemonicwords)
	if err != nil {
		log.Fatal("json err:", err)
	}

	bd := database.New(nodeID)
	wallets := block.NewWallets(bd)
	address,privkey,mnemonicWord := wallets.GenerateWallet(bd,block.CreateBitcoinKeysByMnemonicWord,mnemonicwords)
	fmt.Println("助记词：",mnemonicWord)
	fmt.Println("私钥：", privkey)
	fmt.Println("地址：", address)
}
