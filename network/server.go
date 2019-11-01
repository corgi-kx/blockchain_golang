package network

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	block "github.com/corgi-kx/blockchain_golang/blc"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	"os"
	"strings"
	"time"
)

//在P2P网络中已发现的节点池
//key:节点ID  value:节点详细信息
var peerPool = make(map[string]peer.AddrInfo)
var ctx = context.Background()
var localHost host.Host
var localAddr string
var listenPort int
var tradePool = Transactions{}
var send = Send{}

func StartNode(id int) {
	//先获取本地区块最新高度
	bc := block.NewBlockchain()
	block.NewestBlockHeight = bc.GetLastBlockHeight()
	listenPort = id
	log.Infof("[*] 监听IP地址: %s 端口号: %d", LISTEN_HOST, listenPort)
	r := rand.Reader
	// 为本地节点创建RSA密钥对
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		log.Panic(err)
	}
	// 创建本地节点地址信息
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", LISTEN_HOST, listenPort))
	//传入地址信息，RSA密钥对信息，生成libp2p本地host信息
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		log.Panic(err)
	}
	//写入全局变量本地主机信息
	localHost = host
	//写入全局变量本地P2P节点地址详细信息
	localAddr = fmt.Sprintf( "/ip4/%s/tcp/%v/p2p/%s", LISTEN_HOST, listenPort, host.ID().Pretty())
	log.Infof("[*] 你的P2P地址信息: %s", localAddr)
	//启动监听本地端口，并且传入一个处理流的函数，当本地节点接收到流的时候回调处理流的函数
	host.SetStreamHandler(protocol.ID(PROTOCOL_ID), handleStream)
	//寻找p2p网络并加入到节点池里
	go findP2PPeer()
	//检测节点池,如果发现新节点则打印到屏幕
	go func() {
		currentPeerPoolNum := 0
		for {
			peerPoolNum := len(peerPool)
			if peerPoolNum != 0 && peerPoolNum > currentPeerPoolNum{
				log.Info("----------------------已检测到新P2P节点,当前节点池存在的节点-----------------------------")
				for _,v := range peerPool {
					log.Info("|   ",v,"   |")
				}
				log.Info("-----------------------------------------------------------------------------------------")
				currentPeerPoolNum = peerPoolNum
			}
			time.Sleep(time.Second)
		}
	}()

	//启一个go程去向其他p2p节点发送高度信息，来进行更新区块数据
	go func() {
		//如果节点池中还未存在节点的话,一直循环 直到发现已连接节点
		for {
			if len(peerPool) == 0 {
				time.Sleep(time.Second)
				continue
			}else {
				break
			}
		}
		send.SendVersionToPeers(block.NewestBlockHeight)
	}()
	go receiveUserCmd()
	select {} //wait here
}

//启动mdns寻找p2p网络 并等节点连接
func findP2PPeer() {
	peerChan := initMDNS(ctx, localHost, RENDEZVOUS_STRING)
	for {
		peer := <-peerChan // will block untill we discover a peer
		//将发现的节点加入节点池
		peerPool[fmt.Sprint(peer.ID)] = peer
	}
}

//通过固定格式的地址信息,构建出P2P节点信息对象
func buildPeerInfoByAddr(addrs string) peer.AddrInfo {
	///ip4/0.0.0.0/tcp/9000/p2p/QmUyYpeMSqZp4oNMhANdG6sGeckWiGpBnzfCNvP7Pjgbvg
	p2p:=strings.TrimSpace(addrs[strings.Index(addrs,"/p2p")+len("/p2p/"):])
	ipTcp := addrs[:strings.Index(addrs,"/p2p/")]
	//通过ip与端口获得multiAddr
	multiAddr,err:=multiaddr.NewMultiaddr(ipTcp)
	if err != nil {
		log.Debug(err)
	}
	//拼接成multiAddr数组
	m:=[]multiaddr.Multiaddr{multiAddr}
	//获得host.ID
	id,err := peer.IDB58Decode(p2p)
	if err != nil {
		log.Error(err)
	}
	//传入host.ID , multiAddr数组 拼接成P2P节点信息对象
	return peer.AddrInfo{peer.ID(id),m}
}

func receiveUserCmd() {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}
		fmt.Print(sendData)
		userCmdHandle(sendData)
	}
}

func userCmdHandle(data string) {
	//去除命令前后空格
	//transfer -from [\"17L6dPWj4P4zUqwuYV3wAJcbfogwxVnefo\"] -to [\"1JZRuVD91Jgk3DCVqe216ZTSDueyizz5ZX\"] -amount [50]
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
	case "transfer":
		fromString := strings.TrimSpace(context[strings.Index(context,"-from") + len("-from"):strings.Index(context,"-to")])
		toString := strings.TrimSpace(context[strings.Index(context,"-to") + len("-to") :strings.Index(context,"-amount")])
		amountString := strings.TrimSpace(context[strings.Index(context,"-amount")+len("-amount"):])
		blc:=block.NewBlockchain()
		blc.CreateTransaction(fromString,toString,amountString,send)
	default:
		fmt.Println("无此命令!")
	}
}