package cli

import (
	"bufio"
	"fmt"
	log "github.com/corgi-kx/logcustom"
	"os"
	"strconv"
	"strings"
)

type Cli struct {
}

//打印帮助提示
func printUsage() {
	fmt.Println("----------------------------------------------------------------------------- ")
	fmt.Println("Usage:")
	fmt.Println("\thelp                                              打印命令行说明")
	fmt.Println("\tgenesis  -a DATA  -v DATA                         生成创世区块")
	fmt.Println("\tsetRewardAddr -a DATA                             设置挖矿奖励地址")
	fmt.Println("\tgenerateWallet                                    创建新钱包")
	fmt.Println("\timportMnword -m DATA                              根据助记词导入钱包")
	fmt.Println("\tprintAllWallets                                   查看本地存在的钱包信息")
	fmt.Println("\tprintAllAddr                                      查看本地存在的地址信息")
	fmt.Println("\tgetBalance  -a DATA                               查看用户余额")
	fmt.Println("\ttransfer -from DATA -to DATA -amount DATA         进行转账操作")
	fmt.Println("\tprintAllBlock                                     查看所有区块信息")
	fmt.Println("\tresetUTXODB                                       遍历区块数据，重置UTXO数据库")
	fmt.Println("------------------------------------------------------------------------------")
}

func New() *Cli {
	return &Cli{}
}

func (cli *Cli) Run() {
	printUsage()
	go cli.startNode()
	cli.ReceiveCMD()
}

//获取用户输入
func (cli Cli) ReceiveCMD() {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}
		cli.userCmdHandle(sendData)
	}
}

//用户输入命令的解析
func (cli Cli) userCmdHandle(data string) {
	//去除命令前后空格
	data = strings.TrimSpace(data)
	var cmd string
	var context string
	if strings.Contains(data, " ") {
		cmd = data[:strings.Index(data, " ")]
		context = data[strings.Index(data, " ")+1:]
	} else {
		cmd = data
	}
	switch cmd {
	case "help":
		printUsage()
	case "genesis":
		address := getSpecifiedContent(data, "-a", "-v")
		value := getSpecifiedContent(data, "-v", "")
		v, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}
		cli.genesis(address, v)
	case "generateWallet":
		cli.generateWallet()
	case "setRewardAddr":
		addrss := getSpecifiedContent(data, "-a", "")
		cli.setRewardAddress(addrss)
	case "importMnword":
		mnemonicword := getSpecifiedContent(data, "-m", "")
		cli.importWalletByMnemonicword(mnemonicword)
	case "printAllAddr":
		cli.printAllAddress()
	case "printAllWallets":
		cli.printAllWallets()
	case "printAllBlock":
		cli.printAllBlock()
	case "getBalance":
		address := getSpecifiedContent(data, "-a", "")
		cli.getBalance(address)
	case "resetUTXODB":
		cli.resetUTXODB()
	case "transfer":
		fromString := (context[strings.Index(context, "-from")+len("-from") : strings.Index(context, "-to")])
		toString := strings.TrimSpace(context[strings.Index(context, "-to")+len("-to") : strings.Index(context, "-amount")])
		amountString := strings.TrimSpace(context[strings.Index(context, "-amount")+len("-amount"):])
		cli.transfer(fromString, toString, amountString)
	default:
		fmt.Println("无此命令!")
		printUsage()
	}
}

//返回data字符串中,标签为tag的内容
func getSpecifiedContent(data, tag, end string) string {
	if end != "" {
		return strings.TrimSpace(data[strings.Index(data, tag)+len(tag) : strings.Index(data, end)])
	}
	return strings.TrimSpace(data[strings.Index(data, tag)+len(tag):])
}
