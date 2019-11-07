package network

import "github.com/libp2p/go-libp2p-core/host"

//p2p
var (
	RendezvousString = "meetme"
	ProtocolID       = "/chain/1.1.0"
	ListenHost       = "0.0.0.0"
	ListenPort       = "3001"
	localHost        host.Host
	localAddr        string
)

//交易池默认大小
var TradePoolLength = 2

//版本信息 默认0
const versionInfo = byte(0x00)

//发送数据的头部多少位为命令
const prefixCMDLength = 12

type command string

const (
	cVersion     command = "version"
	cGetHash     command = "getHash"
	cHashMap     command = "hashMap"
	cGetBlock    command = "getBlock"
	cBlock       command = "block"
	cTransaction command = "transaction"
	cMyError     command = "myError"
)
