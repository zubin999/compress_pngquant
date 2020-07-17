package main

import (
	"net/rpc"
	"net"
	"github.com/labstack/gommon/log"
	"net/rpc/jsonrpc"
	"os/exec"
	"bytes"
	"strconv"
)

type PngCompressor struct {

}

type Args struct {
	SrcImg  string
	SaveImg string
	Min     int
	Max     int
}

type Success struct {
	Code int8
	Msg  string
}

func (p *PngCompressor) Compress(args Args, reply *Success) error {
	min := strconv.Itoa(args.Min)
	max := strconv.Itoa(args.Max)
	cmd := exec.Command("/usr/bin/pngquant", "--force", "--quality", min + "-" + max, "-o", args.SaveImg, args.SrcImg)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		reply.Code = 1
		reply.Msg  = err.Error() + stderr.String()
	}else {
		reply.Code = 0
		reply.Msg  = "success"
	}

	return nil
}


func main() {
	rpc.Register(new(PngCompressor))
	listener, err := net.Listen("tcp", "127.0.0.1:1215")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}