package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mkideal/cli"
)

const DescriptionTemplate = `
usage: go run ping.go --Host <Host name or ip address> --TTL <the value TTL you want to set>

This CLI ping program is for Cloudflare Internship Application: Systems.
Author: Shi Tang
Email: tangshi6666@gmail.com
YOu can run the program as:
go run ping.go --Host 8.8.8.8 --TTL 55

Enjoy!
`

type CLIOpts struct {
	Help      bool   `cli:"!h,help" usage:"Show help."`
	Condensed bool   `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`
	Host      string `cli:"Host" usage:"Host name or ip address."`
	TTL       int    `cli:"TTL" usage:"Set TTL by yourself, this will launch time exceeded."`
}

type PingOption struct {
	Count      int
	Size       int
	Timeout    int64
	Nerverstop bool
}

func NewPingOption() *PingOption {
	return &PingOption{
		Count:      4,
		Size:       16,
		Timeout:    1000,
		Nerverstop: false,
	}
}

const (
	icmpv4EchoRequest = 8
	icmpv4EchoReply   = 0
	icmpv6EchoRequest = 128
	icmpv6EchoReply   = 129
)

type icmpMessage struct {
	Type     int             // type
	Code     int             // code
	Checksum int             // checksum
	Body     icmpMessageBody // body
}

type icmpMessageBody interface {
	Len() int
	Marshal() ([]byte, error)
}

func (m *icmpMessage) Marshal() ([]byte, error) {
	b := []byte{byte(m.Type), byte(m.Code), 0, 0}
	if m.Body != nil && m.Body.Len() != 0 {
		mb, err := m.Body.Marshal()
		if err != nil {
			return nil, err
		}
		b = append(b, mb...)
	}
	switch m.Type {
	case icmpv6EchoRequest, icmpv6EchoReply:
		return b, nil
	}
	csumcv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	// Place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero.
	b[2] ^= byte(^s & 0xff)
	b[3] ^= byte(^s >> 8)
	return b, nil
}

func parseICMPMessage(b []byte) (*icmpMessage, error) {
	msglen := len(b)
	if msglen < 4 {
		return nil, errors.New("message too short")
	}
	m := &icmpMessage{Type: int(b[0]), Code: int(b[1]), Checksum: int(b[2])<<8 | int(b[3])}
	//fmt.Print("\n Here is parce b:", b)
	if msglen > 4 {
		var err error
		switch m.Type {
		case icmpv4EchoRequest, icmpv4EchoReply, icmpv6EchoRequest, icmpv6EchoReply:
			m.Body, err = parseICMPEcho(b[4:])
			if err != nil {
				return nil, err
			}
		}
	}
	//fmt.Print("    ----------------test-------------m", m)
	return m, nil
}

// imcpEcho represenets an ICMP echo request or reply message body.
type icmpEcho struct {
	ID   int    // identifier
	Seq  int    // sequence number
	Data []byte // data
}

func (p *icmpEcho) Len() int {
	if p == nil {
		return 0
	}
	return 4 + len(p.Data)
}

func (p *icmpEcho) Marshal() ([]byte, error) {
	b := make([]byte, 4+len(p.Data))
	b[0], b[1] = byte(p.ID>>8), byte(p.ID&0xff)
	b[2], b[3] = byte(p.Seq>>8), byte(p.Seq&0xff)
	copy(b[4:], p.Data)
	return b, nil
}

func parseICMPEcho(b []byte) (*icmpEcho, error) {
	bodylen := len(b)
	p := &icmpEcho{ID: int(b[0])<<8 | int(b[1]), Seq: int(b[2])<<8 | int(b[3])}
	if bodylen > 4 {
		p.Data = make([]byte, bodylen-4)
		copy(p.Data, b[4:])
	}
	return p, nil
}

func main() {
	//argsmap:=map[string]interface{}{}
	cli.SetUsageStyle(cli.DenseManualStyle)
	cli.Run(new(CLIOpts), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CLIOpts)
		if argv.Help || len(argv.Host) == 0 {
			com := ctx.Command()
			com.Text = DescriptionTemplate
			ctx.String(com.Usage(ctx))
			return nil
		}
		ipaddress := argv.Host
		if argv.TTL == 0 {
			argv.TTL = 1000
		}
		vttl := argv.TTL
		//fmt.Print("TTL:", vttl)
		argsmap := map[string]interface{}{}
		p := NewPingOption()
		p.ping3(ipaddress, argsmap, vttl)
		return nil
	})

}

//ping连接用的协议是ICMP，原理：
//Ping的基本原理是发送和接受ICMP请求回显报文。接收方将报文原封不动的返回发送方，发送方校验报文，校验成功则表示ping通。
//一台主机向一个节点发送一个类型字段值为8的ICMP报文，如果途中没有异常（如果没有被路由丢弃，目标不回应ICMP或者传输失败），
//则目标返回类型字段值为0的ICMP报文，说明这台主机可达
func (p *PingOption) ping3(host string, args map[string]interface{}, vttl int) {


	var count int = 10000
	var timeout int64 = 1000
	var neverstop bool = false
	if len(args) != 0 {
		count = args["n"].(int)
		timeout = args["w"].(int64)
		neverstop = args["t"].(bool)
	}

	starttime := time.Now()
	c, err := net.DialTimeout("ip6:ipv6-icmp", host, time.Duration(timeout*1000*1000))
	if err != nil {
		fmt.Println("error ", err)
		return
	}
	sendN := 0
	recvN := 0
	lostN := 0
	shortT := -1
	longT := -1
	sumT := 0
	cflag := make(chan os.Signal, 2)
	signal.Notify(cflag, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cflag
		stat4(host, sendN-1, lostN, recvN, shortT, longT, sumT)
		os.Exit(0)
	}()
	for count > 0 || neverstop {

		sendN++
		starttime = time.Now()
		typ := icmpv6EchoRequest
		xid, xseq := os.Getpid()&0xffff, 1
		wb, err := (&icmpMessage{
			Type: typ, Code: 0,
			Body: &icmpEcho{
				ID: xid, Seq: xseq,
				Data: bytes.Repeat([]byte("Go Go Gadget Ping!!!"), 3),
			},
		}).Marshal()

		//fmt.Print("  Here is wb:", wb)
		if err != nil {
			fmt.Print("\n---------------test3----------")
			return
		}
		if _, err = c.Write(wb); err != nil {
			fmt.Print("\n---------------test4----------")
			return
		}
		var m *icmpMessage
		rb := make([]byte, 20+len(wb))

		for {

			if _, err = c.Read(rb); err != nil {
				return
			}
			if m, err = parseICMPMessage(rb); err != nil {
				return
			}
			switch m.Type {
			case icmpv4EchoRequest, icmpv6EchoRequest:
				fmt.Println("type ", m.Type)
				continue
			}
			break
		}
		var endduration int = int(int64(time.Since(starttime)) / (1000))

		sumT += endduration

		time.Sleep(1000 * 1000 * 1000)

		if shortT == -1 {
			shortT = endduration
		} else if shortT > endduration {
			shortT = endduration
		}
		if longT == -1 {
			longT = endduration
		} else if longT < endduration {
			longT = endduration
		}
		recvN++
		// ttl := int(receive[8])
		//          fmt.Println(ttl)
		fmt.Println("16 bytes come from " + "[" + host + "]" + " reponds time=" + strconv.Itoa(endduration) + "us") //+ "ms TTL=" + strconv.Itoa(ttl)
		// seq++
		count--
	}

}

func checkSum3(msg []byte) uint16 {
	sum := 0

	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256 // notice here, why *256?
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}

func checkError3(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func gensequence3(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genidentifier3(host string) (byte, byte) {
	return host[0], host[1]
}

func stat3(ip string, sendN int, lostN int, recvN int, shortT int, longT int, sumT int) {
	fmt.Println()
	fmt.Println(ip, "  Ping static message:")
	fmt.Printf("    datapacket: Already sent = %d，Already received = %d，lost = %d (%d%% lost)，\n", sendN, recvN, lostN, int(lostN*100/sendN))
	fmt.Println("Estimate time(ms):")
	if recvN != 0 {
		fmt.Printf("    Min = %dus，Max = %dus，Average = %dus\n", shortT, longT, sumT/sendN)
	}
}

func stat4(ip string, sendN int, lostN int, recvN int, shortT int, longT int, sumT int) {
	if sendN == 0 {
		fmt.Println(" Cannot ping ipv6 adress:" + ip)
		fmt.Printf("    Min = %dus，Max = %dus，Average = %dus\n", 0, 0, 0)
	} else {
		fmt.Println()
		fmt.Println(ip, "  Ping static message:")
		fmt.Printf("    datapacket: Already sent = %d，Already received = %d，lost = %d (%d%% lost)，\n", sendN, recvN, lostN, int(lostN*100/sendN))
		fmt.Println("Estimate time(ms):")
		if recvN != 0 {
			fmt.Printf("    Min = %dus，Max = %dus，Average = %dus\n", shortT, longT, sumT/sendN)
		}
	}

}
