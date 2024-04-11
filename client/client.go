package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	// conn.Write([]byte("account=2&merchant=Zara&amount=20.05"))
	for {
		fmt.Print("masukkan amount : ")

		input, _ := reader.ReadString('\n')
		fmt.Println("mengirim payload...")
		conn.Write([]byte("account=2&merchant=Zara&amount=" + strings.Trim(input, "\n")))
		message, _ := bufio.NewReader(conn).ReadString('\n')

		log.Print("Server relay:", message)
	}
}
