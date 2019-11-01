package block

import (
	"math"
	"os"
	"strconv"
)

//当前网络中，区块最新高度
var NewestBlockHeight int

var nodeID,_ = strconv.Atoi(os.Getenv("NODE_ID")) //节点id

//挖矿奖励代币数量
const tokenRewardNum = 25

//奖励地址在数据库中的键
const RewardAddrMapping = "rewardAddress"

//最新区块Hash在数据库中的键
const LastBlockHashMapping = "lastHash"

//钱包地址在数据库中的键
const addrListMapping = "addressList"

//公链版本信息默认为0
const version = byte(0x00)

//两次sha256(公钥hash)后截取的字节数量
const checkSum = 4

//挖矿难度值
const targetBits = 21

//随机数不能超过的最大值
const maxInt = math.MaxInt64
