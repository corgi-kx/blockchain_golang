package network

import "github.com/libp2p/go-libp2p-core/host"

//p2p相关,程序启动时,会被配置文件所替换
var (
	RendezvousString = "meetme"
	ProtocolID       = "/chain/1.1.0"
	ListenHost       = "0.0.0.0"
	ListenPort       = "3001"
	localHost        host.Host
	localAddr        string
)

//交易池
var tradePool = Transactions{}

//交易池默认大小
var TradePoolLength = 2

//版本信息 默认0
const versionInfo = byte(0x00)

//发送数据的头部多少位为命令
const prefixCMDLength = 12

type command string

//网络通讯互相发送的命令
const (
	cVersion     command = "version"
	cGetHash     command = "getHash"
	cHashMap     command = "hashMap"
	cGetBlock    command = "getBlock"
	cBlock       command = "block"
	cTransaction command = "transaction"
	cMyError     command = "myError"
)
