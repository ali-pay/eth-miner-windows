package main

import (
	"errors"
	"fmt"
	"github.com/gen2brain/dlgs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"ieth/acc"
	"ieth/conf"
	"ieth/file"
	"ieth/http"
	"ieth/miner"
	"ieth/util"
	"strconv"
	"time"
)

/** 编译说明

生成syso文件，执行一次就行
rsrc -arch="amd64" -ico ICON.ico -manifest main.manifest -o main.syso

开启黑窗输出调试信息
go build

关闭黑窗输出
go build -ldflags="-H windowsgui -w -s"

错误处理: 注释掉go-ethereum@v1.9.6\ethdb\leveldb\leveldb.go:103:3: DisableSeeksCompaction

参考: https://blog.csdn.net/kgjn__/article/details/89280235
 */

var (
	cpuBtn, gpuBtn, consoleBtn, stopBtn, accountBtn *walk.PushButton //按钮，控制它们的可点击状态

	logPrint            *walk.TextEdit //信息输出框
	password, etherbase *walk.LineEdit //密码输入框，地址输入框
	money, power, block *walk.Label    //钱包余额，算力，区块高度

	logCh   chan string //接受打印的数据
	powerCh chan string //接受算力的通道，发送数据到这里，会替换算力的显示文字
	state   chan string //按钮的点击事件，发送挖矿状态到这里做操作的分发

	errLog  chan error  //接受error数据，为了不让文件看起来太臃肿，用了这2个通道记录数据
	infoLog chan string //接受debug数据

	stopLog   chan bool //停止读取log文件
	stopPower chan bool //停止读取miner文件
	stopInfo  chan bool //挖矿数据的显示，true:显示算力 false:关闭显示，默认不发数据过去就只显示区块高度和金额，停止这个前先把stopPower算力读取给停了
	joinPeer  chan bool //连接节点的耗时显示，true:显示连接成功，开始挖矿 false:关闭耗时显示
	peer      = "0"     //节点的连接数量
)

func init() {

	//杀死残留挖矿进程
	miner.KillMiner()

	//初始化通道
	infoLog = make(chan string)
	errLog = make(chan error)
	state = make(chan string)
	logCh = make(chan string)
	powerCh = make(chan string)
	stopLog = make(chan bool)
	stopPower = make(chan bool)
	stopInfo = make(chan bool)
	joinPeer = make(chan bool)

	//读取配置文件
	if err := conf.ReadConfig(); err != nil {
		util.Error.Println(err.Error())
		return
	}
}

func main() {

	//开启log记录
	go doLog()

	//开启日志打印
	go printLog()

	//开启按钮监听
	go gogogo()

	//GUI窗口
	_, _ = MainWindow{
		Title:  "单机挖矿",                        //标题
		Size:   Size{Width: 600, Height: 600}, //窗口尺寸
		Layout: VBox{},                        //垂直排列
		Children: []Widget{
			//生成账号栏
			Composite{
				Layout: Grid{Columns: 3}, //3列布局
				Children: []Widget{
					Label{Text: "输入密码"},
					LineEdit{AssignTo: &password},
					PushButton{
						Text:     "生成账号",
						AssignTo: &accountBtn,
						OnClicked: func() {
							infoLog <- "点击生成账号"
							state <- "account"
						},
					},
				},
			},
			//日志打印窗口
			TextEdit{
				AssignTo:  &logPrint,
				ReadOnly:  true,
				Text:      "使用说明：\r\n\r\n\t1.请配置好钱包地址以及启动节点后再开始挖矿\r\n\r\n\t2.显卡挖矿如果报错，请添加信任或者关闭杀毒软件\r\n\r\n\t3.如遇到未知错误，建议重启软件即可\r\n\r\n",
				TextColor: walk.RGB(123, 1, 1),
			},
			//钱包地址栏
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "钱包地址"},
					LineEdit{
						AssignTo: &etherbase,
						Text:     conf.Coinbase, //默认显示配置文件中的钱包地址
					},
				},
			},
			//挖矿信息栏
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						AssignTo: &block,
						//Text:     "区块高度：N/N",
						Font: Font{PointSize: 12, Bold: true},
					},
					Label{
						AssignTo: &power,
						//Text:     "算力：N/N",
						Font: Font{PointSize: 12, Bold: true},
					},
					Label{
						AssignTo: &money,
						//Text:     "账户金额：N/N",
						Font: Font{PointSize: 12, Bold: true},
					},
				},
			},
			//4个按钮栏
			Composite{
				Layout: Grid{Columns: 4},
				Children: []Widget{
					//CPU挖矿
					PushButton{
						Text:     "CPU挖矿",
						MinSize:  Size{Width: 100, Height: 50},
						AssignTo: &cpuBtn,
						OnClicked: func() {
							infoLog <- "点击CPU挖矿"
							state <- "cpu"
						},
					},
					//显卡挖矿
					PushButton{
						Text:     "显卡挖矿",
						MinSize:  Size{Width: 100, Height: 50},
						AssignTo: &gpuBtn,
						OnClicked: func() {
							infoLog <- "点击显卡挖矿"
							state <- "gpu"
						},
					},
					//打开控制台
					PushButton{
						Text:     "打开控制台",
						MinSize:  Size{Width: 100, Height: 50},
						AssignTo: &consoleBtn,
						OnClicked: func() {
							infoLog <- "点击打开控制台"
							state <- "console"
						},
					},
					//停止挖矿
					PushButton{
						Text:     "停止挖矿",
						MinSize:  Size{Width: 100, Height: 50},
						AssignTo: &stopBtn,
						Enabled:  false,
						OnClicked: func() {
							infoLog <- "点击停止挖矿, 挖矿状态:" + miner.Mining
							state <- "stop"
						},
					},
				},
			},
		},
	}.Run()
}

//处理debug和error通道数据
func doLog() {
	for {
		select {
		case s := <-infoLog:
			util.Info.Println(s)

		case s := <-errLog:
			util.Error.Println(s.Error())
			logCh <- s.Error()
		}
	}
}

//log打印
func printLog() {
	infoLog <- "开启日志打印: printLog()"
	//在print打印的字符串
	logs := make([]string, 10)
	for {
		s := <-logCh
		//只显示15条数据
		if len(logs) == 15 {
			logs = logs[1:]
		}
		logs = append(logs, s)
		log := ""
		//构造显示的字符串
		for _, str := range logs {
			log += str + "\r\n"
		}
		//设置显示
		if err := logPrint.SetText(log); err != nil {
			errLog <- err
			//return //出错不退出
		}
	}
}

//按钮事件
func gogogo() {
	for {
		s := <-state //分发操作
		switch s {
		//启动链
		case "cpu", "gpu", "console":
			//前置操作
			if err := doCommon(); err != nil {
				errLog <- err
				break
			}

			//分发操作
			switch s {
			case "cpu", "gpu":
				//后台开启Geth
				miner.CpuMiner()

				//等待连接上节点
				go showJoinTime(s)

				//加载log显示
				go file.GetLog(logCh, stopLog)

			case "console":
				miner.OpenConsole(conf.Networkid, conf.Port, conf.Coinbase, conf.Bootnodes)
			}

		//停止挖矿
		case "stop":
			stopMiner()

		//创建账户
		case "account":
			if err := doAccount(); err != nil {
				errLog <- err
				break
			}
		}
	}
}

//停止挖矿
func stopMiner() {
	infoLog <- "停止挖矿: stopMiner()"

	str, err := http.Exit() //保存区块数据
	if err != nil {
		errLog <- err
		//return //手动关掉geth后，这会有连接不上的错误，还是注释了吧，不然停止不了
	}
	infoLog <- "保存区块数据: " + str

	time.Sleep(time.Second) //休眠1秒，等保存好数据

	//cpu和gpu有关闭其他的操作
	if miner.Mining != "console" {
		//没有连上其他节点时
		if peer == "0" {
			joinPeer <- false //关连接线程
		} else {
			stopInfo <- false //连接成功在挖矿时，关挖矿信息获取
		}
		stopLog <- true //关日志读取
	}

	peer = "0"            //重置节点数
	miner.Mining = "stop" //挖矿停止后修改状态

	miner.KillMiner() //杀挖矿进程

	//设置挖矿按钮可点击
	btnEnabled(true)
}

//显示连接耗时
func showJoinTime(state string) {
	infoLog <- "显示连接耗时: showJoinTime()"
	//1秒的定时器
	tick := time.NewTicker(1 * time.Second)
	var count int //记录时间
	for {
		select {
		//显示连接的耗时，其他2个数据隐藏
		case <-tick.C:
			//检查连接
			p, err := http.GetPeer()
			if err != nil {
				errLog <- err
				break
			}
			infoLog <- "节点数：" + p

			//如果连接成功
			if p != "0" {
				infoLog <- "节点连接成功"
				peer = p //记录连接节点数，方便停止挖矿时判断

				//担心挖矿和joinPeer退出冲突，所以单独在一个线程操作
				go startMiner(state)

				//需要用单独的线程给joinPeer发数据，这样才能退出当前的case
				go func() {
					joinPeer <- true //连接成功，关闭耗时显示
				}()
				break
			}

			//一直没连接上，每30秒提示一次要不要继续等待或者直接挖矿
			if count != 0 && count%30 == 0 {
				//询问要不用直接开始挖矿
				msg := "未发现可连接节点，请检查配置文件中的启动节点是否有误！\r\n是否直接开始挖矿？"
				//如果不独立一个线程运行的话，会连续弹出2个提示框
				go func() {
					yes, err := dlgs.Question("提示", msg, false)
					if err != nil {
						errLog <- err
						return
					}
					infoLog <- "按下了：" + strconv.FormatBool(yes)
					//确定要挖
					if yes {
						//设置连接数
						peer = "-1"
						//开始挖矿
						go startMiner(state)
						//停止连接等待
						joinPeer <- false //取消连接，关闭耗时显示
					}
				}()
			}

			//计时显示
			count++
			if err := block.SetText(fmt.Sprintf("正在连接节点。%d", count)); err != nil {
				errLog <- err
				break
			}
			if err := power.SetText(""); err != nil {
				errLog <- err
				break
			}
			if err := money.SetText(""); err != nil {
				errLog <- err
				break
			}

		//true: 连接成功，显示正在同步区块  false: 关闭显示  最后都会结束这个计时线程
		case b := <-joinPeer:
			//停止定时器
			tick.Stop() //放在第一行，尽快停掉它，否则又会发送新的数据

			infoLog <- "关闭连接耗时: showJoinTime()"
			if b {
				if err := block.SetText("节点连接成功"); err != nil {
					errLog <- err
					break
				}
			} else {
				if err := block.SetText(""); err != nil {
					errLog <- err
					break
				}
			}
			return
		}
	}
}

//开始挖矿
func startMiner(state string) {
	infoLog <- "开始挖矿：startMiner(): " + state
	time.Sleep(1 * time.Second) //休眠1秒，等joinPeer退出后再运行
	//发送挖矿请求
	_, err := http.Miner(true)
	if err != nil {
		errLog <- err
		return
	}
	//显示挖矿数据
	go showInfo()
	//使用gpu挖矿
	if state == "gpu" {
		miner.GpuMiner()
		stopInfo <- true //显示算力
	}
}

//显示区块数、账户金额、挖矿算力
func showInfo() {
	infoLog <- "显示挖矿信息: showInfo()"
	//5秒的定时器
	tick2 := time.NewTicker(5 * time.Second)
	var m2, b2 string
	for {
		select {
		//每5秒更新金额和区块高度
		case <-tick2.C:
			//设置账户金额
			m, err := http.GetMoney(conf.Coinbase)
			if err != nil {
				errLog <- err
				break
			}
			//修改金额显示，加个小数点
			if len(m) > 18 {
				point := len(m) - 18
				m = m[:point] + "." + m[point:]
			}
			//不相同才设置新的数据
			if m2 != m {
				m2 = m
				infoLog <- "账户金额：" + m
				if err := money.SetText("账户金额：" + m); err != nil {
					errLog <- err
					break
				}
			}

			//设置区块高度
			b, err := http.GetBlock()
			if err != nil {
				errLog <- err
				break
			}
			//不相同才设置新的数据
			if b2 != b {
				b2 = b
				infoLog <- "区块高度：" + b
				if err := block.SetText("区块高度：" + b); err != nil {
					errLog <- err
					break
				}
			}

		//是否显示算力 true读取power false退出关闭显示数据
		case b := <-stopInfo:
			if !b {
				infoLog <- "关闭挖矿信息：showInfo()"
				if miner.Mining == "gpu" {
					stopPower <- true //关算力读取
				}

				//停止定时器
				tick2.Stop() //不能放在第一行，否则已启动就停止了
				return
			}
			go file.GetPower(logCh, stopPower, powerCh)

		//设置算力
		case p := <-powerCh:
			infoLog <- p
			if err := power.SetText(p); err != nil {
				errLog <- err
				break
			}
		}
	}
}

//挖矿的前置操作
func doCommon() (err error) {
	infoLog <- "挖矿前置操作: doCommon()"
	//钱包地址是42位的
	coinbase := etherbase.Text()
	if len(coinbase) != 42 {
		var tmp walk.Form
		walk.MsgBox(tmp, "错误", "钱包地址有误", walk.MsgBoxIconInformation)
		err = errors.New("钱包地址有误")
		return
	}
	//设置按钮不可点击
	btnEnabled(false)
	//删除log
	if err = conf.NewLog(); err != nil {
		return
	}
	//生成新钱包地址的配置
	conf.Coinbase = coinbase
	if err = conf.NewFile(); err != nil {
		return
	}
	return
}

//按钮的可点击状态转换，true: 可点击 false: 不可点击
func btnEnabled(state bool) {
	if state {
		cpuBtn.SetEnabled(true)
		gpuBtn.SetEnabled(true)
		consoleBtn.SetEnabled(true)
		etherbase.SetEnabled(true)
		password.SetEnabled(true)
		accountBtn.SetEnabled(true)
		stopBtn.SetEnabled(false)
	} else {
		cpuBtn.SetEnabled(false)
		gpuBtn.SetEnabled(false)
		consoleBtn.SetEnabled(false)
		etherbase.SetEnabled(false)
		password.SetEnabled(false)
		accountBtn.SetEnabled(false)
		stopBtn.SetEnabled(true)
	}
}

//创建账户
func doAccount() (err error) {
	infoLog <- "创建账户: doAccount()"
	ps := password.Text() //获取密码
	if ps == "" {
		//错误弹窗
		var tmp walk.Form
		walk.MsgBox(tmp, "错误", "请输入密码", walk.MsgBoxIconInformation)
		err = errors.New("请输入密码")
		return
	}
	//创建账户
	path, err := acc.NewAccount(ps)
	if err != nil {
		return
	}
	//解析地址与私钥
	address, privateKey, err := acc.DecryptKeystore(path, ps)
	if err != nil {
		return
	}
	//在信息输出框显示地址与私钥
	if err = logPrint.SetText(fmt.Sprintf("账户: %s\r\n私钥: %s\r\n密钥存放位置: %s\r\n", address, privateKey, path)); err != nil {
		return
	}
	//在账户框中显示地址
	if err = etherbase.SetText(address); err != nil {
		return
	}
	return
}
