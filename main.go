package main

import (
	"blockchain/BLK"
)

func main() {
	//var cli *BLK.CLI = &BLK.CLI{nil}
	//cli.Run()
	//wallet := BLK.NewWallet()
	//address := wallet.GetAddress()
	//fmt.Println(BLK.IsValidAddress(address))
	test()

}

func test(){
	blk := BLK.CreateBlockChainWithGenesisBlock("smh")
	wallets := BLK.NewWallets()
	//wallet1 := BLK.NewWallet()
	//wallets.AddWallet(wallet1)
	wallets.ReadFromDB(blk.DB)
	wallets.AddressList()
	//val, _ := wallets.GetWallet("1PfGgLvFeRkBMzvHKkrKBEJPTGGz3C92nD")
	//fmt.Println(BLK.IsValidAddress("1PfGgLvFeRkBMzvHKkrKBEJPTGGz3C92nD"))
	//fmt.Println((*val).PrivateKey.D.String())
	defer blk.DB.Close()

}