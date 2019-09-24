package block

import (
	"crypto/sha256"
	"fmt"
	"myCode/public_blockchain/part7-network/util"
	"testing"
	"time"
)

func TestGetBitcoinKeys(t *testing.T) {
	keys := newBitcoinKeys()
	address := keys.getAddress()
	t.Log(string(address))
	t.Log(isVaildBitcoinAddress(string(address)))
}

func TestSign(t *testing.T) {
	bk := newBitcoinKeys()
	hash := sha256.Sum256(util.Int64ToBytes(time.Now().UnixNano()))
	fmt.Printf("签名hash:%x\n签名hash长度:%d\n", hash, len(hash))
	signature := ellipticCurveSign(bk.PrivateKey, hash[:])
	verifyhash := append(hash[:], []byte("知道为什么这么长的验证信息也会通过吗？因为这个椭圆曲线只验证信息的前256位也就是前32字节！！！根据当时传入的elliptic.P256()有关！！！！")...)
	fmt.Printf("验证hash:%x\n验证hash长度:%d\n:", verifyhash, len(verifyhash))
	if ellipticCurveVerify(bk.PublicKey, signature, verifyhash) {
		t.Log("签名信息验证通过")
	} else {
		t.Log("签名信息验证失败！！！")
	}
}
