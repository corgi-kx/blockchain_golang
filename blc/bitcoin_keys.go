package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"golang.org/x/crypto/ripemd160"
	"math/big"
	log "myCode/public_blockchain/part7-network/logcustom"
	"myCode/public_blockchain/part7-network/util"
)

type bitcoinKeys struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func newBitcoinKeys() *bitcoinKeys {
	b := &bitcoinKeys{nil, nil}
	b.newKeyPair()
	return b
}

func (b *bitcoinKeys) newKeyPair() {
	curve := elliptic.P256()
	var err error
	b.PrivateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	b.PublicKey = append(b.PrivateKey.PublicKey.X.Bytes(), b.PrivateKey.PublicKey.Y.Bytes()...)
}

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

func generatePublicKeyHash(publicKey []byte) []byte {
	sha256PubKey := sha256.Sum256(publicKey)
	r := ripemd160.New()
	r.Write(sha256PubKey[:])
	ripPubKey := r.Sum(nil)
	return ripPubKey
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

func isVaildBitcoinAddress(address string) bool {
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

func ellipticCurveSign(privKey *ecdsa.PrivateKey, hash []byte) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		log.Panic("EllipticCurveSign:", err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

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
