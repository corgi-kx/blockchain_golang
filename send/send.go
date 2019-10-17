/*
	为什么要拆出个send包呢，老老实实的放到network包里不好吗？
	因为blc包生成区块后要调用SendMessage发送区块高度信息到中心节点，
	如果不拆的话，blc与network包会涉及到相互引用问题
*/
package send

import (
	log "github.com/corgi-kx/blockchain_golang/logcustom"
	"net"
	"os"
)

var centerNode = "localhost:3000"
/*
	挖矿后像其他节点发送高度信息
	中心节点向已知节点发送高度信息
	其他节点向中心节点发送高度信息
*/
var nodeID = os.Getenv("NODE_ID")
var localAddr = "localhost:" + nodeID
func SendVersionToCenterNode(lastHeight int) {
		newV:=version{versionInfo,lastHeight,localAddr}
		data:=jointMessage(cVersion,newV.serialize())
		SendMessage(data,centerNode)
}


//中心节点更新区块链后调用此方法向已知节点发送高度信息
func SendVersionToKnows(lastHeight int,knowNodesMap map[string]string) {
	if knowNodesMap != nil {
		for k,_ := range knowNodesMap {
			newV:=version{versionInfo,lastHeight,centerNode}
			data:=jointMessage(cVersion,newV.serialize())
			SendMessage(data,k)
			log.Debugf("中心节点已向%s节点发送高度信息\n",k)
		}
	}
}

//像中心节点发送交易信息
func SendTransToCenterNode(tss Transactions) {
	for i,_:=range tss.Ts {
		tss.Ts[i].AddrFrom = localAddr
	}
	log.Tracef("已发送%d笔交易到中心节点",len(tss.Ts))
	data:=jointMessage(cTransaction,tss.Serialize())
	SendMessage(data,centerNode)
}
func SendMessage(data []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Warn(err)
		return
	}
	defer conn.Close()
	_, err2 := conn.Write(data)
	if err2 != nil {
		log.Panic(err)
	}
}

//默认前十二位为命令名称
func jointMessage(cmd command, content []byte) []byte {
	b := make([]byte, prefixCMDLength)
	for i, v := range []byte(cmd) {
		b[i] = v
	}
	joint := make([]byte, 0)
	joint = append(b, content...)
	return joint
}

