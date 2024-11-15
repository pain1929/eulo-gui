package main

import (
	Eulogist "Eulogist/eulogist"
	"Eulogist/message"
	_ "embed"
	"github.com/pterm/pterm"
	"net"
	"os"
)

func main() {

	dial, err2 := net.Dial("tcp", "127.0.0.1:1930")
	if err2 != nil {
		println(err2.Error())
		return
	}
	err := Eulogist.Eulogist(os.Args[1], os.Args[2], os.Args[3], dial)
	if err != nil {
		pterm.Error.Println(err)
		message.SendMsg(false, err.Error(), dial)
	}
}
