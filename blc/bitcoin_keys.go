package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"github.com/corgi-kx/blockchain_golang/util"
	log "github.com/corgi-kx/logcustom"
	"math/big"
	"os"
)

type bitcoinKeys struct {
	PrivateKey   *ecdsa.PrivateKey
	PublicKey    []byte
	MnemonicWord []string
}

//创建公私钥实例
func NewBitcoinKeys(nothing []string) *bitcoinKeys {
	b := &bitcoinKeys{nil, nil, nil}
	b.MnemonicWord = getChineseMnemonicWord()
	b.newKeyPair()
	return b
}

//根据助记词创建公私钥
func CreateBitcoinKeysByMnemonicWord(mnemonicWord []string) *bitcoinKeys {
	if len(mnemonicWord) != 7 {
		log.Error("助记词格式不正确，应为七对中文双字词语")
		return nil
	}
	for _, v := range mnemonicWord {
		if len(v) != 6 {
			log.Error("助记词格式不正确，应为七对中文双字词语")
			return nil
		}
	}

	b := &bitcoinKeys{nil, nil, nil}
	b.MnemonicWord = mnemonicWord
	b.newKeyPair()
	return b
}

//根据中文助记词生成公私钥对
func (b *bitcoinKeys) newKeyPair() {
	curve := elliptic.P256()
	var err error
	buf := bytes.NewReader(b.jointSpeed())
	b.PrivateKey, err = ecdsa.GenerateKey(curve, buf)
	if err != nil {
		log.Panic(err)
	}
	b.PublicKey = append(b.PrivateKey.PublicKey.X.Bytes(), b.PrivateKey.PublicKey.Y.Bytes()...)
}

//将助记词拼接成字节数组，并截取前40位
func (b bitcoinKeys) jointSpeed() []byte {
	bs := make([]byte, 0)
	for _, v := range b.MnemonicWord {
		bs = append(bs, []byte(v)...)
	}
	return bs[:40]
}

//获取中文种子
func getChineseMnemonicWord() []string {
	file, err := os.Open(ChineseMnwordPath)
	//file,err:=os.Open("D:/programming/golang/GOPATH/src/github.com/corgi-kx/blockchain_golang/blc/chinese_mnemonic_world.txt")
	if err != nil {
		log.Panic(err)
	}
	s := []string{}
	//因为种子最高40个字节，所以就取7对词语，7*2*3 = 42字节，返回后在截取前40位
	for i := 0; i < 7; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(5948)) //词库一共5949对词语，顾此设置随机数最高5948
		if err != nil {
			log.Panic(err)
		}
		b := make([]byte, 6)
		_, err = file.ReadAt(b, n.Int64()*7+3) //从文件的具体位置读取 防止乱码
		if err != nil {
			log.Panic(err)
		}
		s = append(s, string(b))
	}
	file.Close()
	return s
}

//私钥长度为32字节
const privKeyBytesLen = 32

//获取私钥
func (keys *bitcoinKeys) GetPrivateKey() string {
	d := keys.PrivateKey.D.Bytes()
	b := make([]byte, 0, privKeyBytesLen)
	priKey := paddedAppend(privKeyBytesLen, b, d)
	//base58加密
	return string(util.Base58Encode(priKey))
}

func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

//通过公钥获得地址
func (b *bitcoinKeys) getAddress() []byte {
	//1.ripemd160(sha256(publickey))
	ripPubKey := generatePublicKeyHash(b.PublicKey)
	//2.最前面添加一个字节的版本信息获得 versionPublickeyHash
	versionPublickeyHash := append([]byte{version}, ripPubKey[:]...)
	//3.sha256(sha256(versionPublickeyHash))  取最后四个字节的值
	tailHash := checkSumHash(versionPublickeyHash)
	//4.拼接最终hash versionPublickeyHash + checksumHash
	finalHash := append(versionPublickeyHash, tailHash...)
	//进行base58加密
	address := util.Base58Encode(finalHash)
	return address
}

//序列化
func (b *bitcoinKeys) serliazle() []byte {
	var result bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

//反序列化
func (v *bitcoinKeys) Deserialize(d []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	gob.Register(elliptic.P256())
	err := decoder.Decode(v)
	if err != nil {
		log.Panic(err)
	}
}

func generatePublicKeyHash(publicKey []byte) []byte {
	sha256PubKey := sha256.Sum256(publicKey)
	r := util.NewRipemd160()
	r.Reset()
	r.Write(sha256PubKey[:])
	ripPubKey := r.Sum(nil)
	return ripPubKey
	return nil
}

func getPublicKeyHashFromAddress(address string) []byte {
	addressBytes := []byte(address)
	fullHash := util.Base58Decode(addressBytes)
	publicKeyHash := fullHash[1 : len(fullHash)-checkSum]
	return publicKeyHash
}

func checkSumHash(versionPublickeyHash []byte) []byte {
	versionPublickeyHashSha1 := sha256.Sum256(versionPublickeyHash)
	versionPublickeyHashSha2 := sha256.Sum256(versionPublickeyHashSha1[:])
	tailHash := versionPublickeyHashSha2[:checkSum]
	return tailHash
}

//判断是否是有效的比特币地址
func IsVaildBitcoinAddress(address string) bool {
	adddressByte := []byte(address)
	fullHash := util.Base58Decode(adddressByte)
	if len(fullHash) != 25 {
		return false
	}
	prefixHash := fullHash[:len(fullHash)-checkSum]
	tailHash := fullHash[len(fullHash)-checkSum:]
	tailHash2 := checkSumHash(prefixHash)
	if bytes.Compare(tailHash, tailHash2[:]) == 0 {
		return true
	} else {
		return false
	}
}

//通过公钥信息获得地址
func GetAddressFromPublicKey(publickey []byte) string {
	if publickey == nil {
		return ""
	}
	b := bitcoinKeys{PublicKey: publickey}
	return string(b.getAddress())
}

//通过公钥信息获得地址
func GetAddressFromPublicKeyHash(publickeyHash []byte) string {
	//2.最前面添加一个字节的版本信息获得 versionPublickeyHash
	versionPublickeyHash := append([]byte{version}, publickeyHash[:]...)
	//3.sha256(sha256(versionPublickeyHash))  取最后四个字节的值
	tailHash := checkSumHash(versionPublickeyHash)
	//4.拼接最终hash versionPublickeyHash + checksumHash
	finalHash := append(versionPublickeyHash, tailHash...)
	//进行base58加密
	address := util.Base58Encode(finalHash)
	return string(address)
}

//使用私钥进行数字签名
func ellipticCurveSign(privKey *ecdsa.PrivateKey, hash []byte) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		log.Panic("EllipticCurveSign:", err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

//使用公钥进行签名验证
func ellipticCurveVerify(pubKey []byte, signature []byte, hash []byte) bool {
	//拆分签名 得到 r,s
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])
	//拆分公钥字节数组，得到公钥对象
	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubKey)
	x.SetBytes(pubKey[:(keyLen / 2)])
	y.SetBytes(pubKey[(keyLen / 2):])
	curve := elliptic.P256()
	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	//传入公钥，要验证的信息，以及签名
	if ecdsa.Verify(&rawPubKey, hash, &r, &s) == false {
		return false
	}
	return true
}
