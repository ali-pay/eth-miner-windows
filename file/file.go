package file

import (
	"bufio"
	"ieth/conf"
	"ieth/util"
	"io"
	"os"
	"strings"
	"time"
)

type logReader struct {
	file   string
	offset int64
}

func (f *logReader) Read(p []byte) (n int, err error) {
	reader, err := os.Open(f.file)
	if err != nil {
		return 0, err
	}
	defer reader.Close()
	if _, err = reader.Seek(f.offset, 0);
		err != nil {
		return 0, err
	}
	n, err = reader.Read(p)

	if err == io.EOF {
		time.Sleep(1 * time.Second)
	}
	f.offset += int64(n)

	return n, err
}

//读取日志 geth.log
func GetLog(logCh chan string, stopLog chan bool) {
	util.Info.Println("开启日志读取：readLog()")
	time.Sleep(2 * time.Second)
	f := &logReader{conf.GethLog, 0}
	br := bufio.NewReader(f)

	//一直读取文件
	for {
		select {
		//读取数据
		default:
			log, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				util.Error.Println(err.Error())
				logCh <- err.Error()
				break
			}

			//将新的数据发送到logCh
			logCh <- string(log)

		//关闭读取
		case <-stopLog:
			util.Info.Println("关闭日志读取：readLog()")
			return
		}
	}
}

//获取算力 miner.log
func GetPower(logCh chan string, stopPower chan bool, powerCh chan string) {
	util.Info.Println("开启算力读取：getPower()")
	time.Sleep(5 * time.Second)
	f := &logReader{conf.MinerLog, 0}
	br := bufio.NewReader(f)
	var power string
	//一直读取文件
	for {
		select {
		//读取数据
		default:
			log, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				util.Error.Println(err.Error())
				logCh <- err.Error()
				break
			}

			//截取算力字符
			arr := strings.Split(string(log), " ")
			if len(arr) == 11 {
				if len(arr[7]) > 2 || len(arr[7]) == 0 { //判断单位的长度 单位有 h Mh Gh
					break
				}
				p := arr[6] + arr[7]
				if p != power { //不相同才设置新的数据
					power = p
					powerCh <- "算力：" + p //算力：10.70Mh  算力：0.00h
				}
			}

		//关闭读取
		case <-stopPower:
			util.Info.Println("关闭算力读取：getPower()")
			return
		}
	}
}
