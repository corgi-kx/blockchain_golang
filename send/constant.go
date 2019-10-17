package send

//版本信息 默认0
const versionInfo = byte(0x00)
//发送数据的头部多少位为命令
const prefixCMDLength = 12

//交易池默认大小
const tradePoolLength = 3

type command string

const (
	cVersion command = "version"
	cGetHash command = "getHash"
	cHashMap command = "hashMap"
	cGetBlock command = "getBlock"
	cBlock  command = "block"
	cTransaction  command = "transaction"
	cMining command = "mining"
)

