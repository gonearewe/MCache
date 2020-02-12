package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gonearewe/MCache/client"
)

func main() {
	var client_ = client.New("tcp", "localhost")
	var sc = bufio.NewScanner(os.Stdin)
	for {
		printPrompt()
		sc.Scan()
		cmd := strings.Split(sc.Text(), " ")

		if len(cmd) == 2 && cmd[0] == "get" {
			var req = &client.Request{
				Type:  "get",
				Key:   cmd[1],
				Val:   nil,
				Error: nil,
			}
			client_.Run(req)

			if req.Error != nil || len(req.Val) == 0 {
				fmt.Println("FAIL: ", req.Error)
			} else {
				fmt.Println("key: ", cmd[1], " value: ", string(req.Val))
			}
		}

		if len(cmd) == 3 && cmd[0] == "set" {
			var req = &client.Request{
				Type:  "set",
				Key:   cmd[1],
				Val:   []byte(cmd[2]),
				Error: nil,
			}
			client_.Run(req)

			if req.Error != nil {
				fmt.Println("FAIL: ", req.Error)
			} else {
				fmt.Println("OK")
			}
		}
	}
}

func printPrompt() {
	fmt.Print("MCache Client > ")
}
