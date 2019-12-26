package block

import (
	"crypto/sha256"
	"fmt"
	"github.com/corgi-kx/blockchain_golang/util"
	"os"
	"testing"
	"time"
)

func TestGetBitcoinKeys(t *testing.T) {
	t.Log("测试获取比特币公私钥，并生成地址")
	{
		keys := NewBitcoinKeys([]string{})
		address := keys.getAddress()
		t.Log("\t地址为：", string(address))
		t.Log("\t地址格式是否正确：", IsVaildBitcoinAddress(string(address)))
	}
}

func TestSign(t *testing.T) {
	t.Log("测试数字签名是否可用")
	{
		bk := NewBitcoinKeys([]string{})
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

func TestMnemonicWord(t *testing.T) {
	t.Log("测试中文助记词")
	{
		k := NewBitcoinKeys([]string{})
		t.Log(k.MnemonicWord)
		t.Log(k.PrivateKey)
		t.Log(k.PublicKey)
	}
}

func TestReadTxt(t *testing.T) {
	file, err := os.Open("D:/programming/golang/GOPATH/src/github.com/corgi-kx/blockchain_golang/blc/chinese_mnemonic_world.txt")
	if err != nil {
		t.Error(err)
	}
	b := make([]byte, 6)
	file.ReadAt(b, 3)
	println(string(b))
	//s:=[]string{}
	//设置助记词文件为utf-8 一个中文两个字节
	//因为种子最高40个字节，所以取十对
	//for i := 0;i<10;i++ {
	//	n,err:=rand.Int(rand.Reader,big.NewInt(5948))  //词库一共5949对词语，顾此设置随机数最高5948
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	b:=make([]byte,4)
	//	_,err=file.ReadAt(b,n.Int64()*5)
	//	if err != nil && err !=io.EOF {
	//		t.Error(err)
	//	}
	//	s=append(s,string(b))
	//}
	file.Close()
}
