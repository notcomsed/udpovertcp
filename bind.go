package main

import (
	"io"
	"net"
	"os"
	"fmt"
	"runtime"
	"time"
)

	var ipset_connect string

	func 处理新连接()(net.Conn) {
		for{
		remoteConn, err := net.DialTimeout("tcp", ipset_connect, time.Duration(time.Second*18))
		if err == nil {
			return remoteConn
		}
		fmt.Println("Connect remote :", err)
		time.Sleep(time.Second)
	}
	}


func main() {
	arg_num:=len(os.Args)
	if arg_num <= 2 {
			os.Exit(0)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	ipset_bind := os.Args[1]
	ipset_connect = os.Args[2]
	udpAddr,err:=net.ResolveUDPAddr("udp",ipset_bind)

	if err != nil {
		panic(err)
	}

	ipset, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening udp %s,to %s[tcp]\n", ipset_bind,ipset_connect)
	buffer:=make([]byte,1536)
	buffer_tcp:=make([]byte,1536)
	R_tcp_head:=make([]byte,2)
	S_tcp_head:=make([]byte,2)
	var hn uint16
	for {
	remoteConn:=处理新连接()

	n,conn, err := ipset.ReadFromUDP(buffer)
	go func() {
		for {
		_, err := io.ReadFull(remoteConn, R_tcp_head);
		if err != nil {
			fmt.Println(err)
			remoteConn.Close()
			ipset.Close()
			break
		}
		hn=((uint16(R_tcp_head[0]) << 8) | uint16(R_tcp_head[1]))

		_, err = io.ReadFull(remoteConn, buffer_tcp[:hn]);
		if err != nil {
			fmt.Println(err)
			remoteConn.Close()
			ipset.Close()
			break
		} else {
			ipset.WriteToUDP(buffer_tcp[:hn],conn)
		}
	}
	}()

	for {
		if err != nil {
			fmt.Println(err)
			remoteConn.Close()
			ipset.Close()
			break
		} else {
		S_tcp_head[0]=byte((n >> 8) & 0xff)
		S_tcp_head[1]=byte(n & 0xff)
		remoteConn.Write(S_tcp_head)
		remoteConn.Write(buffer[:n])
		}
		n,_, err = ipset.ReadFromUDP(buffer)
	}
	fmt.Printf("connect closed\n")
	for{
		ipset, err = net.ListenUDP("udp", udpAddr)
		if err == nil {
			break
		}
		fmt.Println("Error: fail to Listen udp", ipset_bind)
		time.Sleep(time.Second)
	}
	}
}
