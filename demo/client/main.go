package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "39.156.66.18:80")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	buffer := make([]byte, 1024)
	binary.BigEndian.PutUint32(buffer[0:4], 32)
	binary.BigEndian.PutUint32(buffer[4:8], 1)
	for i := 8; i < 40; i++ {
		buffer[i] = 66
	}
	n, _ := conn.Write(buffer[0:40])
	fmt.Printf("%v", n)
	conn.Read(buffer)
	fmt.Println(string(buffer))
	fmt.Scanf("%d\n")
}
