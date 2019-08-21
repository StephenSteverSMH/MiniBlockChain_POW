package BLK

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	Block *Block //当前验证的区块
	target *big.Int //大数存储
}
var diff uint= 10
//进行工作量证明
func (proofOfWork *ProofOfWork) Run() ([]byte, int64){
	var nonce int64= 0
	hashInt := big.NewInt(0)
	var hash [32]byte
	for{
		//准备数据
		dataBytes := proofOfWork.perpareData(nonce)
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		//fmt.Println(hashInt.String())
		//挖到矿了
		if proofOfWork.target.Cmp(hashInt) == 1{
			fmt.Println(hashInt)
			break;
		}
		nonce = nonce+1
	}
	return hash[:], nonce
}
//根据nonce计算hash
func (pow* ProofOfWork) perpareData(nonce int64) []byte{
	data := bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		pow.Block.Data,
		IntToHex(pow.Block.Timestamp),
		IntToHex(nonce),
		IntToHex(int64(pow.Block.Height))}, []byte{})
	return data
}

//创建新的工作量证明对象
func NewProofOfWork(block* Block) *ProofOfWork{
	target := big.NewInt(1)
	target.Lsh(target, 256 - diff)
	return &ProofOfWork{block, target}
}

//判断工作量证明是否有效
func (proofOfWork *ProofOfWork) IsValid()  bool{
	hashInt := big.NewInt(0)
	hashInt.SetBytes(proofOfWork.Block.Hash)
	if proofOfWork.target.Cmp(hashInt) == 1 {
		return true;
	}
	return false;
}