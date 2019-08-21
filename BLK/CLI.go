package BLK

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	BLC *Blockchain
}
//输出帮助
func (cli* CLI)printHelp(){
	fmt.Println("Usage:")
	fmt.Println("\taddBlock -data <data> //挖矿")
	fmt.Println("\tprintChain //输出所有区块详细信息")
	fmt.Println("\tcreateBlockChain -data <data>//创建创始区块和区块链")
}
//创建区块链并挖掘创始区块
func (cli *CLI)createBlockChain(to string)  (existed bool){
	if Exists(DBPATH+DBNAME){
		return false
	}
	cli.BLC = CreateBlockChainWithGenesisBlock(to)
	cli.BLC.DB.Close()
	return true
}
//添加区块
func (cli *CLI)addBlock(to string){
	if Exists(DBPATH+DBNAME){
		//获得已经存在的区块链
		cli.BLC = CreateBlockChainWithGenesisBlock("no used")
		cli.BLC.AddBlockToBlockChain(to)
		cli.BLC.DB.Close()
	} else {
		fmt.Println("请先创建区块链")
	}
}
//显示所有区块
func (cli *CLI)printChain(){
	if Exists(DBPATH+DBNAME){
		//获得已经存在的区块链
		cli.BLC = CreateBlockChainWithGenesisBlock("no used")
		cli.BLC.PrintChain()
		cli.BLC.DB.Close()
	} else {
		fmt.Println("请先创建区块链")
	}

}
//发送交易
func (cli *CLI)sendTransaction(from []string, to []string, amount []int64){
	if Exists(DBPATH+DBNAME){
		NewSimpleTransaction(from[0], to, amount)
	}else {
		fmt.Println("请先创建区块链")
	}
}
//运行CLI
func (cli *CLI)Run(){
	if len(os.Args)==1{
		fmt.Println("not right usage")
		cli.printHelp()
		return
	}
	mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)
	toFlag := mineCmd.String("to", "discard", "挖矿")

	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	flagCreateBlockChain := createBlockChainCmd.String("data", "genesis data", "给创始区块增加数据")

	sendTransactionCmd := flag.NewFlagSet("sendTX", flag.ExitOnError)
	fromFlag := sendTransactionCmd.String("from", "[\"smh\"]","交易发送方")
	recvFlag := sendTransactionCmd.String("to", "[\"smh\"]","交易接受方")
	amountFlag := sendTransactionCmd.String("amount", "[0]","交易金额")

	switch os.Args[1] {
		case "help":
			cli.printHelp()
		case "createBlockChain":
			err := createBlockChainCmd.Parse(os.Args[2:])
			if err!=nil{
				log.Panic("addBlockCmd parse error")
			}
		case "mine":
			err := mineCmd.Parse(os.Args[2:])
			if err!=nil{
				log.Panic("mine parse error")
			}
		case "printChain":
			err :=printChainCmd.Parse(os.Args[2:])
			if err!=nil{
				log.Panic("printChainCmd parse error")
			}
		case "sendTX":
			err := sendTransactionCmd.Parse(os.Args[2:])
			if err!=nil{
				log.Panic("sendTransaction parse error")
			}
		default:
			fmt.Println("not right usage")
			cli.printHelp()
			os.Exit(1)
	}
	if mineCmd.Parsed(){
		cli.addBlock(*toFlag)
	}
	if printChainCmd.Parsed(){
		cli.printChain()
	}
	if sendTransactionCmd.Parsed(){
		var from []string
		var to []string
		var amount []int64
		json.Unmarshal([]byte(*fromFlag), &from)
		json.Unmarshal([]byte(*recvFlag), &to)
		json.Unmarshal([]byte(*amountFlag), &amount)
		fmt.Println(from, to, amount)
	}
	if createBlockChainCmd.Parsed(){
		ok := cli.createBlockChain(*flagCreateBlockChain)
		if ok{
			fmt.Println("创建区块链成功")
		}else {
			fmt.Println("创建区块链失败, 可能数据文件已经存在")
		}
	}
}