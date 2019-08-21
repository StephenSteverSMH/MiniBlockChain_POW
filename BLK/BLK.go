package BLK

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Height int64
	PrevBlockHash []byte
	Data []byte
	Timestamp int64
	Hash []byte
	Nonce int64
}
func (block *Block) setHash(){
	//Height int64 -> []byte
	heightBytes := IntToHex(block.Height)
	//TimeStamp int64 -> base:2, string
	//TimeStamp string -> []byte
	timeString := strconv.FormatInt(block.Timestamp, 2)
	timeBytes := []byte(timeString)
	//拼接所有属性
	blockBytes := bytes.Join([][]byte{heightBytes, block.PrevBlockHash, block.Data, timeBytes, block.Hash}, []byte{})
	//生成hash
	hash := sha256.Sum256(blockBytes)
	//给该区块赋值
	block.Hash = hash[:]
}
//创建区块链
func NewBlock(height int64, prevBlockHash []byte)  *Block{
	block := &Block{
		height,
		prevBlockHash,
		[]byte{},
		time.Now().Unix(),
		nil,
		0}
	//block.setHash()
	//生成工作量证明对象
	pow := NewProofOfWork(block)
	//执行工作量证明（挖矿）
	hash, nonce := pow.Run()
	//成功挖到矿，赋值给新区块
	block.Hash = hash
	block.Nonce = nonce
	//将全局交易池中交易打包入该区块
	block.Data = []byte(GLOBAL_TXPOOL.Serialize())
	return block
}
//生成创始区块
func CreateGenesisBlock(to string) *Block{
	initPrevBlockHash := [64]byte{}
	NewCoinbaseTX(to)
	return NewBlock(1, initPrevBlockHash[:])
}

//JSON序列化区块
func (block* Block) Serialize() string{
	data, err := json.Marshal(block)
	if err!=nil{
		fmt.Printf("json.Marshal,err:",err);
	}
	return string(data)
}

//JSON反序列化区块
func (block* Block) Deserialize(str string){
	err := json.Unmarshal([]byte(str), block)
	if err!=nil{
		fmt.Printf("json.Marshal,err:",err);
	}
}