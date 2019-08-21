package BLK

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
	"math/big"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func NewWallet() *Wallet{
	privateKey, pubKey := NewKeyPair()
	return &Wallet{privateKey, pubKey}
}

func NewKeyPair() (ecdsa.PrivateKey, []byte){
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err!=nil{
		log.Panic("ecc GenerateKey error")
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}
//根据用户公钥推算用户地址
func (wallet *Wallet)GetAddress() string{
	//先对公钥进行ripemd160
	hash:=ripemd160.New()
	hash.Write(wallet.PublicKey)
	hash160_pk := hash.Sum(nil)
	//对hash160_pk进行两次sha256，取前4字节，形成校验码
	temp := sha256.Sum256(hash160_pk)
	check := sha256.Sum256(temp[:])
	//用户地址前缀1字节为0x00
	var version [1]byte = [1]byte{0x00}
	//拼接前缀+hash160_pk+校验码
	address := append(version[:], hash160_pk...)
	address = append(address, check[:4]...)
	//对校验码进行base58处理，并返回
	return Encode(address, BitcoinAlphabet)
}
//判断用户地址是否合法
func IsValidAddress(address string) bool{
	b_addresss, err := Decode(address, BitcoinAlphabet)
	if err!=nil{
		log.Panic("base58 decode error")
	}
	if len(b_addresss)!=25{
		return false
	}
	version := b_addresss[0:1]
	hash160_pk := b_addresss[1:21]
	check := b_addresss[21:]

	temp := sha256.Sum256(hash160_pk)
	new_check := sha256.Sum256(temp[:])

	if version[0] != 0x00{
		return false
	}
	b1 := big.NewInt(0)
	b2 := big.NewInt(0)
	b1.SetBytes(check)
	b2.SetBytes(new_check[:4])
	if b1.Cmp(b2)!=0{
		return false
	}
	return true
}