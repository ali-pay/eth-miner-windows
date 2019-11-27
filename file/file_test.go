package file

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"
)

//测试文件获取算力
func TestGetPower(t *testing.T) {
	f := &logReader{"./log/miner.log", 0}
	br := bufio.NewReader(f)
	var count int
	//一直读取文件
	for {
		//读取数据
		log, _, err := br.ReadLine()
		if err == io.EOF {
			continue
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		count++
		if count == 600 {
			return
		}
		arr := strings.Split(string(log), " ")

		if len(arr) == 11 {

			if len(arr[7])>2 {
				continue
			}
			p := arr[6] + arr[7]
			fmt.Println("算力：" + p)
		}
	}
}

//找到算力的索引
func TestFindPower(t *testing.T) {
	arr := strings.Split(" m 16:06:57 <unknown> 0:00 A146:R13 10.70 Mh - cl0 10.70", " ")
	fmt.Println("len: ",len(arr))
	fmt.Println("arr: ",arr)
	for i,v:=range arr{
		fmt.Printf("%d：%s\r\n" , i,v)
	}
}

//找到数字的索引
func TestFindBlockNumber(t *testing.T) {
	arr := strings.Split("INFO [11-12|10:37:54.448] Loaded most recent local header          number=33053 hash=82dfe2…723969 td=180364710254 age=4m13s", " ")
	fmt.Println("len: ",len(arr))
	fmt.Println("arr: ",arr)
	for i,v:=range arr{
		if v==""{
			continue
		}
		fmt.Printf("%d：%s\r\n" , i,v)
	}
	fmt.Println("-------------------------------------------------")
	arr = strings.Split("INFO [11-12|10:37:54.448] Loaded most recent local full block      number=33043 hash=7b3406…315c61 td=179928711532 age=29m8s", " ")
	fmt.Println("len: ",len(arr))
	fmt.Println("arr: ",arr)
	for i,v:=range arr{
		if v==""{
			continue
		}
		fmt.Printf("%d：%s\r\n" , i,v)
	}
	fmt.Println("-------------------------------------------------")
	arr = strings.Split("INFO [11-12|10:37:54.448] Loaded most recent local fast block      number=33053 hash=82dfe2…723969 td=180364710254 age=4m13s", " ")
	fmt.Println("len: ",len(arr))
	fmt.Println("arr: ",arr)
	for i,v:=range arr{
		if v==""{
			continue
		}
		fmt.Printf("%d：%s\r\n" , i,v)
	}
}

//测试数字在哪行，事实证明，每次都不同
func TestFindBlockNumberInLine(t *testing.T) {
	f := &logReader{"./log/gpu.log", 0}
	br := bufio.NewReader(f)
	var count int
	//一直读取文件
	for {
		//读取数据
		log, _, err := br.ReadLine()
		if err == io.EOF {
			continue
		}
		if err != nil {
			fmt.Println(err.Error())
		}

		str := string(log)
		count++
		if count == 20 {
			return
		}
		if strings.Contains(str, "Loaded most recent local") {
			fmt.Println(count,str)
			arr := strings.Split(str, " ")
			//fmt.Println("len: ",len(arr))
			//fmt.Println("arr: ",arr)
			for i,v:=range arr{
				if v==""{
					continue
				}
				fmt.Printf("%d：%s\r\n" , i,v)
			}
			fmt.Println("-------------------------------------------------")
		}
	}
}
