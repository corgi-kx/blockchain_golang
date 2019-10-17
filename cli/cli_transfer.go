package cli

import (
	"encoding/json"
	block "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/database"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
	"github.com/corgi-kx/blockchain_golang/send"
	"unsafe"
)

func (cli *Cli) transfer(from, to string, amount string) {
	bc := block.NewBlockchain()
	defer bc.BD.Close()

	//判断一下是否已生成创世区块
	if bc.BD.View([]byte(block.LastBlockHashMapping), database.BlockBucket) == nil {
		log.Fatal("还没有生成创世区块，不可进行转账操作 !")
	}
	if bc.BD.View([]byte(block.RewardAddrMapping), database.AddrBucket) == nil {
		log.Fatal("没有设置挖矿奖励地址，请前往设置!")
	}
	fromSlice := []string{}
	toSlice := []string{}
	amountSlice := []int{}
	err := json.Unmarshal([]byte(from), &fromSlice)
	if err != nil {
		log.Fatal("json err:", err)
	}
	err = json.Unmarshal([]byte(to), &toSlice)
	if err != nil {
		log.Fatal("json err:", err)
	}
	err = json.Unmarshal([]byte(amount), &amountSlice)
	if err != nil {
		log.Fatal("json err:", err)
	}
	if len(fromSlice) != len(toSlice) || len(fromSlice) != len(amountSlice) {
		log.Fatal("转账数组长度不一致")
	}

	for i, v := range fromSlice {
		if !block.IsVaildBitcoinAddress(v) {
			log.Errorf(" %s,地址格式不正确！已将此笔交易剔除\n", v)
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
				amountSlice = append(amountSlice[:i], amountSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
				amountSlice = append(amountSlice[:i])
			}
		}
	}

	for i, v := range toSlice {
		if !block.IsVaildBitcoinAddress(v) {
			log.Errorf(" %s,地址格式不正确！已将此笔交易剔除\n", v)
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
				amountSlice = append(amountSlice[:i], amountSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
				amountSlice = append(amountSlice[:i])
			}
		}
	}
	for i, v := range amountSlice {
		if v<0 {
			log.Error("转账金额不可小于0，已将此笔交易剔除")
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
				amountSlice = append(amountSlice[:i], amountSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
				amountSlice = append(amountSlice[:i])
			}
		}
	}
	ts:=bc.CreateTransaction(fromSlice,toSlice,amountSlice)
	if len(ts) == 0 {
		log.Fatal("没有通过验证的交易，无法发送交易至中心节点。")
	}

	nTs:=[]send.Transaction{}
	for _,v:=range ts {
		nv:=(*send.Transaction)(unsafe.Pointer(&v))
		nTs = append(nTs,*nv)
	}
	nTss:=send.Transactions{nTs}
	//nTss:=(*send.Transactions)(unsafe.Pointer(&tss))   //强制类型转换，将block下的Transactions结构体转换为send下的Transactions结构体
	send.SendTransToCenterNode(nTss)
}
