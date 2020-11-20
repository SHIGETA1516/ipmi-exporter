package test

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	_ "io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

type GPUInfos struct{
	Pid string
	Uid string
	Cmd string
	Utilization float64 //GPU 利用率
	Mem float64			//显存利用率
	Idx int
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}


func parseGPU(contents string){
	var gpuInfos GPUInfos
	if contents!=""{
		fields:=strings.Fields(contents)
		if len(fields)>0 && IsNum(fields[0]){
			fmt.Println("Idx: ",fields[0])
			idx,err:=strconv.Atoi(fields[0])
			if err != nil {
				log.Errorf("error occured: %s", err.Error())
			}
			pid:=fields[1]
			sm,err:=strconv.ParseFloat(fields[3],64)
			if err != nil {
				log.Errorf("error occured: %s", err.Error())
			}
			mem,err:=strconv.ParseFloat(fields[4],64)
			if err != nil {
				log.Errorf("error occured: %s", err.Error())
			}
			enc,err:=strconv.ParseFloat(fields[5],64)
			if err != nil {
				log.Errorf("error occured: %s", err.Error())
			}
			dec,err:=strconv.ParseFloat(fields[6],64)
			if err != nil {
				log.Errorf("error occured: %s", err.Error())
			}
			cmd:=fields[7]
			util:=sm+enc+dec
			gpuInfos.Idx=idx
			gpuInfos.Pid=pid
			gpuInfos.Utilization=util
			gpuInfos.Mem=mem
			gpuInfos.Cmd=cmd
			fmt.Printf("Idx : %d  Pid: %s sm: %f mem:%f enc:%f dec:%f cmd:%s  \n",idx,pid,sm,mem,enc,dec,cmd)
			fmt.Printf("GPUInfos: Pid:%s Cmd:%s GPU利用率:%f 显存:%f Idx:%d \n",gpuInfos.Pid,gpuInfos.Cmd,gpuInfos.Utilization,gpuInfos.Mem,gpuInfos.Idx)
		}
	}
}

func syncGPU(stdout io.ReadCloser) {
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		//fmt.Println(line)
		parseGPU(line)
	}
}


func doSomething(){
	for i:=0;i<100;i++{
		fmt.Println(i*i)
	}
}

//测试这个方法就行
func TestGPU(t *testing.T) {
	cmdStr:=`nvidia-smi pmon -d 5`
	cmd := exec.Command("bash", "-c", cmdStr)
	//这里得到标准输出和标准错误输出的两个管道，此处获取了错误处理
	cmdStdoutPipe, _ := cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	go syncGPU(cmdStdoutPipe)

	doSomething()
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
