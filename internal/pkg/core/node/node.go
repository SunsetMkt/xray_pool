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
	SubID              string  `json:"sub_id"`
	Data               string  `json:"data"`
	TestResult         float64 `json:"-"`
	wg                 *sync.WaitGroup
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
func NewNode(link, id string, wg *sync.WaitGroup) *Node {
	if data := protocols.ParseLink(link); data != nil {
		return &Node{Protocol: data, SubID: id, wg: wg}
	}
	return nil
}

// NewNodeByData New一个节点
func NewNodeByData(protocol protocols.Protocol, wg *sync.WaitGroup) *Node {
	return &Node{Protocol: protocol, wg: wg}
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
func (n *Node) TCPing(wg *sync.WaitGroup) {
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
	if n.wg != nil {
		n.wg.Done()
	} else {
		wg.Done()
	}
}

func (n *Node) Show() {
	ShowTopBottomSepLine('=', strings.Split(n.GetInfo(), "\n")...)
}
