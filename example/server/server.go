package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var t =time.Now()


func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 9981})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Local: <%s> \n", listener.LocalAddr().String())
	data := make([]byte, 1024)
	nums:=make([]uint32,0,100)
	time.Sleep(time.Second*3)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		// fmt.Println("length:",len(data[:n]))
		a,_:=strconv.Atoi(string(data[:n][4:]))
		nums=append(nums,uint32(a))
		fmt.Printf("len:%d;cap:%d\n",len(nums),cap(nums))
		if t.Add(time.Second*3).Before(time.Now()){
			// fmt.Println("begin")
			go testResult(&nums)
		}
		t =time.Now()
		_, err = listener.WriteToUDP([]byte(data[:n]), remoteAddr)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
	
}

func testResult(nums* []uint32){
	t =time.Now()
	for{
		time.Sleep(time.Millisecond)
		// fmt.Println("wait",time.Now())
		if t.Add(time.Second*3).Before(time.Now()){
			fmt.Println(len(*nums))
			fmt.Println(*nums)
			*nums=make([]uint32,0,100)
			break
		}
	}
}