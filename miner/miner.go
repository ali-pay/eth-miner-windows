package miner

import (
	"fmt"
	"ieth/conf"
	"ieth/util"
	"os/exec"
)

var Mining = "stop" //当前挖矿的状态，分为:cpu gpu console stop

//初始化链
func newData() {
	if util.Exists("./data/geth") {
		util.Info.Println("本地存在数据，开启Geth")
		return
	}
	util.Info.Println("本地不存在数据, 初始化链: newData()")
	args := fmt.Sprintf(`./bin/geth.exe --datadir data init genesis.json`)
	exec.Command("cmd.exe", "/c", "start "+args).Run()
}

//打开控制台
func OpenConsole(networkid, port, etherbase, bootnodes string) {
	Mining = "console"
	//判断是否初始化
	newData()
	//打开控制台，不能加双引号，不然那段字符串会无效
	args := fmt.Sprintf(`./bin/geth.exe --datadir data --syncmode %s --networkid %s --port %s --rpc --rpcapi db,eth,net,web3,personal,miner --allow-insecure-unlock --etherbase %s --bootnodes %s console`, conf.Syncmode, networkid, port, etherbase, bootnodes)
	exec.Command("cmd.exe", "/c start "+args).Run()
}

//cpu挖矿
func CpuMiner() {
	Mining = "cpu"
	//判断是否初始化
	newData()
	//开始挖矿
	exec.Command("cmd.exe", "/c", "start .\\bin\\geth.vbs").Run()
}

//显卡挖矿
func GpuMiner() {
	Mining = "gpu"
	//开始挖矿
	exec.Command("cmd.exe", "/c", "start ./bin/miner.vbs").Run()
}

//杀死挖矿进程
func KillMiner() {
	util.Info.Println("杀死挖矿进程: killMiner()")
	exec.Command("taskkill.exe", "/f", "/im", "geth.exe").Run()
	exec.Command("taskkill.exe", "/f", "/im", "ethminer.exe").Run()
}
