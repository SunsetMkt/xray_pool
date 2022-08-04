package node

import (
	"fmt"
	"github.com/allanpk716/xray_pool/internal/pkg/protocols"
	"net"
	"strings"
	"sync"
	"time"
)

type Node struct {
	protocols.Protocol `json:"-"`
	SubID              string         `json:"sub_id"`
	Data               string         `json:"data"`
	TestResult         float64        `json:"-"`
	Wg                 sync.WaitGroup `json:"-"`
}

func (n *Node) TestResultStr() string {
	if n.TestResult == 0 {
		return ""
	} else if n.TestResult >= 99998 {
		return "-1ms"
	} else {
		return fmt.Sprintf("%.4vms", n.TestResult)
	}
}

// NewNode New一个节点
func NewNode(link, id string) *Node {
	if data := protocols.ParseLink(link); data != nil {
		return &Node{Protocol: data, SubID: id}
	}
	return nil
}

func NewNodeByData(protocol protocols.Protocol) *Node {
	return &Node{Protocol: protocol}
}

// ParseData 反序列化Data
func (n *Node) ParseData() {
	n.Protocol = protocols.ParseLink(n.Data)
}

// Serialize2Data 序列化数据-->Data
func (n *Node) Serialize2Data() {
	n.Data = n.GetLink()
}

// TCPing 测试 Node TCP 连接
func (n *Node) TCPing() {
	count := 3
	var sum float64
	var timeout = 3 * time.Second
	isTimeout := false
	for i := 0; i < count; i++ {
		start := time.Now()
		d := net.Dialer{Timeout: timeout}
		conn, err := d.Dial("tcp", fmt.Sprintf("%s:%d", n.GetAddr(), n.GetPort()))
		if err != nil {
			isTimeout = true
			break
		}
		_ = conn.Close()
		elapsed := time.Since(start)
		sum += float64(elapsed.Nanoseconds()) / 1e6
	}
	if isTimeout {
		n.TestResult = 99999
	} else {
		n.TestResult = sum / float64(count)
	}
	n.Wg.Done()
}

func (n *Node) Show() {
	ShowTopBottomSepLine('=', strings.Split(n.GetInfo(), "\n")...)
}
