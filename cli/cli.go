package cli

import (
	"flag"
	"fmt"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
	"os"
)

var (
	nodeID = ""
)

type Cli struct {
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tgenesis  -a DATA  -v DATA                         生成创世区块")
	fmt.Println("\tgenerateWallet                                    生成钱包")
	fmt.Println("\tsetRewardAddr -a DATA                             设置挖矿奖励地址")
	fmt.Println("\tibMnemonicword -m                                   根据助记词导入钱包")
	fmt.Println("\tprintAllAddr                                      查看本地存在的地址信息")
	fmt.Println("\tprintAllWallets                                   查看本地存在的钱包信息")
	fmt.Println("\tgetBalance  -a DATA                               查看用户余额")
	fmt.Println("\ttransfer -from DATA -to DATA -amount DATA         进行转账操作")
	fmt.Println("\tresetUTXOdb                                       遍历区块数据，重置UTXO数据库")
	fmt.Println("\tprintAllBlock                                     打印所有区块信息")
	os.Exit(0)
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
	}
}

func getNodeID() {
	id := os.Getenv("NODE_ID")
	if id == "" {
		fmt.Println("请设置节点端口号！")
		os.Exit(1)
	}
	nodeID = id
	log.Debug("NODE_ID:", nodeID)

}

func New() *Cli {
	return &Cli{}
}

func (cli *Cli) Run() {
	isValidArgs()
	getNodeID()
	flagGenesisBlockchain := flag.NewFlagSet("genesis", flag.ExitOnError)
	cmdGenesisBlockchainAddress := flagGenesisBlockchain.String("a", "", "添加创世区块用户地址")
	cmdGenesisBlockchainValue := flagGenesisBlockchain.Int("v", 0, "添加创世区块金额数量")
	flagGenerateWallet := flag.NewFlagSet("generateWallet", flag.ExitOnError)
	flagSetRewardAddress := flag.NewFlagSet("setRewardAddr", flag.ExitOnError)
	flagStartNode := flag.NewFlagSet("startNode", flag.ExitOnError)
	cmdSetRewardAddress := flagSetRewardAddress.String("a", "", "设置挖矿奖励地址")
	flagImportWalletByMnemonicword := flag.NewFlagSet("ibMnemonicword", flag.ExitOnError)
	cmdMnemonicword := flagImportWalletByMnemonicword.String("m", "", "助记词信息")
	flagPrintAllAddress := flag.NewFlagSet("printAllAddr", flag.ExitOnError)
	flagPrintAllWallets := flag.NewFlagSet("printAllWallets", flag.ExitOnError)
	flagGetBalance := flag.NewFlagSet("getBalance", flag.ExitOnError)
	cmdGetBalance := flagGetBalance.String("a", "", "查看用户余额")
	flagTransfer := flag.NewFlagSet("transfer", flag.ExitOnError)
	cmdTransferFrom := flagTransfer.String("from", "", "发送转账账户")
	cmdTransferTo := flagTransfer.String("to", "", "接收转账账户")
	cmdTransferAmount := flagTransfer.String("amount", "", "转账金额")
	flagResetUTXODatabase := flag.NewFlagSet("resetUTXOdb", flag.ExitOnError)
	flagPrintBlock := flag.NewFlagSet("printAllBlock", flag.ExitOnError)

	switch os.Args[1] {
	case "genesis":
		err := flagGenesisBlockchain.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "generateWallet":
		err := flagGenerateWallet.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "setRewardAddr":
		err := flagSetRewardAddress.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startNode":
		err := flagStartNode.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "ibMnemonicword":
		err := flagImportWalletByMnemonicword.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printAllAddr":
		err := flagPrintAllAddress.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printAllWallets":
		err := flagPrintAllWallets.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBalance":
		err := flagGetBalance.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "transfer":
		err := flagTransfer.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "resetUTXOdb":
		err := flagResetUTXODatabase.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printAllBlock":
		err := flagPrintBlock.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
	}

	if flagGenesisBlockchain.Parsed() {
		if *cmdGenesisBlockchainAddress == "" || *cmdGenesisBlockchainValue == 0 {
			printUsage()
		}
		cli.genesis(*cmdGenesisBlockchainAddress, *cmdGenesisBlockchainValue)
	} else if flagGenerateWallet.Parsed() {
		cli.generateWallet()
	} else if flagSetRewardAddress.Parsed() {
		if *cmdSetRewardAddress == "" {
			printUsage()
		}
		cli.setRewardAddress(*cmdSetRewardAddress)
	} else if flagStartNode.Parsed() {
		cli.startNode()
	} else if flagImportWalletByMnemonicword.Parsed() {
		if *cmdMnemonicword == "" {
			printUsage()
		}
		cli.importWalletByMnemonicword(*cmdMnemonicword)
	}else if flagPrintAllAddress.Parsed() {
		cli.printAllAddress()
	}  else if flagPrintAllWallets.Parsed() {
		cli.printAllWallets()
	}else if flagGetBalance.Parsed() {
		if *cmdGetBalance == "" {
			printUsage()
		}
		cli.getBalance(*cmdGetBalance)
	} else if flagTransfer.Parsed() {
		if *cmdTransferFrom == "" || *cmdTransferTo == "" || *cmdTransferAmount == "" {
			printUsage()
		}
		cli.transfer(*cmdTransferFrom, *cmdTransferTo, *cmdTransferAmount)
	} else if flagResetUTXODatabase.Parsed() {
		cli.resetUTXODB()
	} else if flagPrintBlock.Parsed() {
		cli.printAllBlock()
	}
}
