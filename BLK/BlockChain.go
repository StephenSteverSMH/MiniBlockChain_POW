package BLK

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"time"
)
const DBNAME = "blockchain.db"
const DBPATH = "D:/go语言项目/src/blockchain/"
const BLOCK_TABLE_NAME = "blocks"

type Blockchain struct {
	Tip []byte //最新的区块
	DB *bolt.DB
}
//判断db文件是否存在（即区块链是否已经存在）
func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
//创建带有创始区块的区块链；如果已经存在，则获得已经存在的区块链
//to代表 coinbase交易目标是哪个用户
func CreateBlockChainWithGenesisBlock(to string) *Blockchain{
	//如果数据文件已存在
	if Exists(DBPATH + DBNAME){
		db, err := bolt.Open(DBPATH + DBNAME, 0600, nil)
		if err != nil{
			log.Fatal(err)
		}
		var tipBlock Block
		err = db.View(func(tx *bolt.Tx) error{
			b:= tx.Bucket([]byte("BlockBucket"))
			if b!=nil{
				//获取最新区块
				tipBlock.Deserialize(string(b.Get([]byte("l"))))
			}
			//返回nil，以便数据库处理相应操作
			return nil
		})
		if err!=nil{
			log.Panic("CreateBlockChainWithGenesisBlock error")
		}
		GLOBAL_UTXOS.ReadFromDB_UTXOS(db)
		return &Blockchain{tipBlock.Hash, db}
	}

	//如果数据文件不存在，重新创建区块链
	var genesisBlock *Block
	db, err := bolt.Open(DBPATH + DBNAME, 0600, nil)
	if err != nil{
		log.Fatal(err)
	}
	//创建初始区块，并将其放入数据库中
	err=db.Update(func(tx *bolt.Tx) error{
		b, err:= tx.CreateBucket([]byte("BlockBucket"))
		if err!= nil{
			log.Panic(err)
		}
		if b!=nil{
			genesisBlock = CreateGenesisBlock(to)
			b.Put(genesisBlock.Hash, []byte(genesisBlock.Serialize()))
			b.Put([]byte("l"), []byte(genesisBlock.Serialize()))
		}
		//返回nil，以便数据库处理相应操作
		return nil
	})

	//更新全局UTXO
	GLOBAL_UTXOS.UpdateUTXO(GLOBAL_TXPOOL)
	//将最新UTXO写入数据库
	GLOBAL_UTXOS.WriteToDB_UTXOS(db)
	//清空全局交易池
	GLOBAL_TXPOOL = new(TXPOOL)
	return &Blockchain{genesisBlock.Hash, db}
}

//增加区块至区块链
func (blc *Blockchain)AddBlockToBlockChain(to string){
	//最新区块
	var tipBlock Block
	err:=blc.DB.View(func(tx *bolt.Tx) error{
		b:= tx.Bucket([]byte("BlockBucket"))
		if b!=nil{
			//获取最新区块
			tipBlock.Deserialize(string(b.Get(blc.Tip)))
		}
		//返回nil，以便数据库处理相应操作
		return nil
	})
	//构建coinbase交易
	NewCoinbaseTX(to)
	//新建区块(挖矿)
	newBlock := NewBlock(tipBlock.Height+1, tipBlock.Hash)
	//将新区块放入数据库中
	err=blc.DB.Update(func(tx *bolt.Tx) error{
		b:= tx.Bucket([]byte("BlockBucket"))
		if b!=nil{
			b.Put(newBlock.Hash, []byte(newBlock.Serialize()))
			b.Put([]byte("l"), []byte(newBlock.Serialize()))
		}
		//返回nil，以便数据库处理相应操作
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	//更新区块链
	blc.Tip = newBlock.Hash
	//更新全局UTXO
	GLOBAL_UTXOS.UpdateUTXO(GLOBAL_TXPOOL)
	//将最新UTXO写入数据库
	GLOBAL_UTXOS.WriteToDB_UTXOS(blc.DB)
	//清空全局交易池
	GLOBAL_TXPOOL = new(TXPOOL)
}

//输出所有区块的信息
func (blc *Blockchain)PrintChain(){
	var block Block
	blc.DB.View(func(tx *bolt.Tx) error{
		b:= tx.Bucket([]byte("BlockBucket"))
		if b!=nil{
			hashTemp := blc.Tip
			if hashTemp == nil{
				fmt.Println("区块链中没有区块")
			}else {
				for{
					//获取最新区块
					(&block).Deserialize(string(b.Get(hashTemp)))
					//输出该区块
					fmt.Println("----------------------------------------------------")
					fmt.Println("区块的高度:",block.Height)
					fmt.Println("区块HASH:",block.Hash)
					fmt.Println("区块的上一个区块的HASH",block.PrevBlockHash)
					fmt.Println("区块的Nonce值:",block.Nonce)
					fmt.Println("区块的出矿时间:",time.Unix(block.Timestamp,0).Format("Mon Jan 2 15:04:05 2006"))
					fmt.Println("区块中的数据:",string(block.Data))
					fmt.Println("----------------------------------------------------")
					//比较父区块链HASH是否为0
					hashInt := big.NewInt(0)
					hashInt.SetBytes(block.PrevBlockHash)
					if hashInt.Cmp(big.NewInt(0)) == 0 {
						break
					}
					//循环
					hashTemp = block.PrevBlockHash
				}
			}
		}
			//返回nil，以便数据库处理相应操作
			return nil
		})
}
