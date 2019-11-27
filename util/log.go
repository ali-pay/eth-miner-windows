package util

import (
	"io"
	"log"
	"os"
)

//使用: util.Info.Println("hello world") or util.Error.Println("hello world")
//参考: https://www.jianshu.com/p/73ae6dc4d16a

var (
	Info    *log.Logger // 运行信息
	Error   *log.Logger // 错误问题
	logPath = "./log"
	infoLog = "./log/info.log"
	errLog  = "./log/error.log"

	debug = false //设置info输出到文件
)

func init() {
	//创建log目录
	if !Exists(logPath) {
		if err := os.Mkdir("log", os.ModePerm); err != nil {
			log.Fatalln("Failed to make log dir:", err)
		}
	}

	//创建error.log
	errFile, err := os.OpenFile(errLog,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	//输出至文件
	Error = log.New(io.MultiWriter(errFile, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	//输出至控制台
	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	//debug模式，输出info.log文件
	if debug {
		//创建info.log
		infoFile, err := os.OpenFile(infoLog,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open info log file:", err)
		}

		//输出至文件
		Info = log.New(io.MultiWriter(infoFile, os.Stderr),
			"INFO: ",
			log.Ldate|log.Ltime|log.Lshortfile)
	}
}
