package block

import (
	"crypto/sha256"
	"fmt"
	"myCode/public_blockchain/part7-network/util"
	"testing"
	"time"
)

func TestGetBitcoinKeys(t *testing.T) {
	t.Log("测试获取比特币公私钥，并生成地址")
	{
		keys := newBitcoinKeys()
		address := keys.getAddress()
		t.Log("\t地址为：",string(address))
		t.Log("\t地址格式是否正确：",isVaildBitcoinAddress(string(address)))
	}
}

func TestSign(t *testing.T) {
	t.Log("测试数字签名是否可用")
	{
		bk := newBitcoinKeys()
		hash := sha256.Sum256(util.Int64ToBytes(time.Now().UnixNano()))
		fmt.Printf("\t签名hash:%x\n签名hash长度:%d\n", hash, len(hash))
		signature := ellipticCurveSign(bk.PrivateKey, hash[:])
		verifyhash := append(hash[:], []byte("\t知道为什么这么长的验证信息也会通过吗？因为这个椭圆曲线只验证信息的前256位也就是前32字节！！！根据当时传入的elliptic.P256()有关！！！！")...)
		fmt.Printf("\t验证hash:%x\n验证hash长度:%d\n:", verifyhash, len(verifyhash))
		if ellipticCurveVerify(bk.PublicKey, signature, verifyhash) {
			t.Log("\t签名信息验证通过")
		} else {
			t.Fatal("\t签名信息验证失败！！！")
		}
	}
}
