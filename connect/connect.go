package main

import (
	"io"
	"net"
	"os"
	"fmt"
	"runtime"
)

	var ipset_connect string

	func 处理新连接(conn net.Conn) {
	udpAddr,err:=net.ResolveUDPAddr("udp",ipset_connect)
		if err != nil {
			conn.Close()
			fmt.Println("Resolve remote :", err)
			return
		}

		remoteConn, err := net.DialUDP("udp",nil,udpAddr)
		if err != nil {
			conn.Close()
			fmt.Println("Connect remote :", err)
			return
		}
		buffer:=make([]byte,1536)
		buffer_tcp:=make([]byte,1536)
		R_tcp_head:=make([]byte,2)
		S_tcp_head:=make([]byte,2)
		var hn uint16
		go func() {
			for {
			n,_, err := remoteConn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println(err)
				remoteConn.Close()
				conn.Close()
				break
			} else {
				S_tcp_head[0]=byte((n >> 8) & 0xff)
				S_tcp_head[1]=byte(n & 0xff)
				conn.Write(S_tcp_head)
				conn.Write(buffer[:n])
			}
		}
		}()

		for {
		_, err := io.ReadFull(conn, R_tcp_head);

		if err != nil {
			fmt.Println(err)
			remoteConn.Close()
			conn.Close()
			break
		}
		hn=((uint16(R_tcp_head[0]) << 8) | uint16(R_tcp_head[1]))
		_, err = io.ReadFull(conn, buffer_tcp[:hn]);
		if err != nil {
			remoteConn.Close()
			conn.Close()
			fmt.Println(err)
			break
		} else {
			remoteConn.Write(buffer_tcp[:hn])
		}
	}
	}


func main() {
	arg_num:=len(os.Args)
	if arg_num <= 5 {
			os.Exit(0)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd := os.Args[1]
	if (cmd != "-s"){
	fmt.Println("command error")
	os.Exit(-1)
	}
	ipset_bind := os.Args[2]
	ipset_bindPort := os.Args[3]
	ipset_connect = os.Args[4]
	ipset_connectPort := os.Args[5]
	ipset_bind = ipset_bind +":"+ipset_bindPort
	ipset_connect=ipset_connect+":"+ipset_connectPort
	ipset, err := net.Listen("tcp", ipset_bind)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening %s[tcp] ,connect %s[udp]\n", ipset_bind,ipset_connect)
	for {
		conn, err := ipset.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
		
	 fmt.Println("new client:", conn.RemoteAddr())
		
		go 处理新连接(conn)
		}
		
	}

}
