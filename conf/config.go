package conf

import (
	"bufio"
	"fmt"
	"ieth/util"
	"io"
	"os"
	"strings"
)

var (
	confPath = "挖矿配置.ini"

	gethVbs  = "./bin/geth.vbs"
	minerVbs = "./bin/miner.vbs"

	gethBat  = "./bin/geth.bat"
	minerBat = "./bin/miner.bat"

	GethLog  = "./log/geth.log"
	MinerLog = "./log/miner.log"

	Coinbase     = "0xeb22459524804361ab700f2552b066a1392b80ab"
	Port         = "61910"
	Networkid    = "100"
	Bootnodes    = "enode://5e00f29f43637107067e5c0cf94004fd8f9475d4d4f86be7501acb20f54a3ad3d2d70ec5c660dbb86656560a5e19892bf1cf38c65ae8202e54a2eb7fe144bbb1@192.168.105.215:61910"
	Minerthreads = "1"
	Syncmode     = "full"
)

//读取配置文件
func ReadConfig() error {
	//第一次启动：配置文件不存在则初始化它
	if !util.Exists(confPath) {
		util.Info.Println("ini文件不存在")

		//写入ini文件
		if err := SaveConfig(); err != nil {
			return err
		}
		return nil //Config初始化完成
	}

	util.Info.Println("ini文件存在")
	//第N+1次启动：读取配置文件
	f, err := os.Open(confPath)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		//去除两端的空格
		s := strings.TrimSpace(string(b))
		//判断=号的位置
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		//取等号左边的key值
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		//取等号右边的value值
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}

		//util.Info.Println(key,value)
		//读取配置数据
		switch key {
		case "钱包地址":
			Coinbase = value
		case "网络ID":
			Networkid = value
		case "端口":
			Port = value
		case "挖矿线程":
			Minerthreads = value
		case "区块同步":
			if Syncmode != value { //full != full
				Syncmode = "fast"
			}
		case "启动节点":
			Bootnodes = value
		}
	}
	util.Info.Printf("ini文件读取：\r\n钱包地址=%s\r\n网络ID=%s\r\n端口=%s\r\n挖矿线程=%s\r\n区块同步=%s\r\n启动节点=%s\r\n", Coinbase, Networkid, Port, Minerthreads, Syncmode, Bootnodes)
	return nil
}

//创建新的配置文件
func SaveConfig() error {
	f, err := os.OpenFile(confPath, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	defer f.Close()

	//坑爹的记事本，会文件开头添加0xefbbbf（十六进制）的字符，此处添加换行解决问题 https://blog.csdn.net/u013281361/article/details/65631820
	str := fmt.Sprintf("----\r\n钱包地址=%s\r\n网络ID=%s\r\n端口=%s\r\n挖矿线程=%s\r\n区块同步=%s\r\n启动节点=%s\r\n", Coinbase, Networkid, Port, Minerthreads, Syncmode, Bootnodes)
	if _, err = f.Write([]byte(str)); err != nil {
		return err
	}
	util.Info.Println("ini文件保存成功")
	return nil
}

//创建新的挖矿脚本
func GethBat() error {
	f, err := os.OpenFile(gethBat, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	defer f.Close()
	bat := fmt.Sprintf(`
@echo off
start /b ./bin/geth.exe --datadir ./data --syncmode %s --networkid %s --port %s --rpc --rpcapi "db,eth,net,web3,personal,miner" --allow-insecure-unlock --etherbase %s --bootnodes %s console >%s 2>&1
`, Syncmode, Networkid, Port, Coinbase, Bootnodes, GethLog)
	if _, err = f.Write([]byte(bat)); err != nil {
		return err
	}
	util.Info.Println("geth.bat文件创建成功")
	return nil
}

//创建新的挖矿脚本
func MinerBat() error {
	f, err := os.OpenFile(minerBat, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	defer f.Close()
	bat := fmt.Sprintf(`
@echo off
start /b ./bin/ethminer.exe -G -P http://127.0.0.1:8545 >%s 2>&1
`, MinerLog)
	if _, err = f.Write([]byte(bat)); err != nil {
		return err
	}
	util.Info.Println("miner.bat文件创建成功")
	return nil
}

//创建新的配置以及挖矿脚本
func NewFile() error {
	//更新钱包地址
	util.Info.Println("更新钱包地址")

	//写入ini文件
	if err := SaveConfig(); err != nil {
		return err
	}
	//写入geth.vbs文件
	if err := GethVbs(); err != nil {
		return err
	}
	//写入miner.vbs文件
	if err := MinerVbs(); err != nil {
		return err
	}
	//写入geth.bat文件
	if err := GethBat(); err != nil {
		return err
	}
	//写入miner.bat文件
	if err := MinerBat(); err != nil {
		return err
	}

	return nil
}

//创建新的空log
func NewLog() error {
	//geth.log
	f, err := os.OpenFile(GethLog, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write([]byte{}); err != nil {
		return err
	}
	util.Info.Println("geth.log文件创建成功")

	//miner.log
	f, err = os.OpenFile(MinerLog, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write([]byte{}); err != nil {
		return err
	}
	util.Info.Println("miner.log文件创建成功")

	return nil
}

//创建新的挖矿脚本
func GethVbs() error {
	//不需要更新数据的不用重复创建
	if util.Exists(gethVbs) {
		util.Info.Println("geth.vbs已存在")
		return nil
	}
	f, err := os.OpenFile(gethVbs, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	defer f.Close()
	bat := `
Set ws = CreateObject("Wscript.Shell")
ws.run "cmd /c .\\bin\\geth.bat",0
`
	if _, err = f.Write([]byte(bat)); err != nil {
		return err
	}
	util.Info.Println("geth.vbs文件创建成功")
	return nil
}

//创建新的挖矿脚本
func MinerVbs() error {
	//不需要更新数据的不用重复创建
	if util.Exists(minerVbs) {
		util.Info.Println("miner.vbs已存在")
		return nil
	}
	f, err := os.OpenFile(minerVbs, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	defer f.Close()
	bat := `
Set ws = CreateObject("Wscript.Shell")
ws.run "cmd /c .\\bin\\miner.bat",0
`
	if _, err = f.Write([]byte(bat)); err != nil {
		return err
	}
	util.Info.Println("miner.vbs文件创建成功")
	return nil
}
