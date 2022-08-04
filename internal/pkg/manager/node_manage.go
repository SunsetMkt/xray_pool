package manager

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
)

// NodeLen 节点数量
func (m *Manager) NodeLen() int {
	return len(m.NodeList)
}

// GetNodeByIndex 获取节点
func (m *Manager) GetNode(index int) *node.Node {
	return m.getNode(index - 1)
}

// getNode 获取节点
func (m *Manager) getNode(i int) *node.Node {
	if i < m.NodeLen() && i >= 0 {
		return m.NodeList[i]
	}
	return nil
}

func (m *Manager) addNode(n *node.Node) bool {
	if n == nil {
		return false
	}
	if f := m.IsCanFilter(n); f != nil {
		logger.Infof("规则 [%s] 过滤节点==> %s", f.String(), n.GetName())
		return false
	}
	n.Serialize2Data()
	m.NodeList = append(m.NodeList, n)
	return true
}

func (m *Manager) AddNode(n *node.Node) bool {
	ok := false
	if ok = m.addNode(n); ok {
		m.Save()
	}
	return ok
}

func (m *Manager) NodeForEach(funC func(int, *node.Node)) {
	for i, n := range m.NodeList {
		funC(i+1, n)
	}
}

func (m *Manager) TCPing() {
	m.NodeForEach(func(i int, n *node.Node) {
		m.wg.Add(1)
		go n.TCPing()
	})
	m.wg.Wait()
	defer m.Save()
	m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
		return n1.TestResult < n2.TestResult
	})
}

func (m *Manager) NodeSort(less func(*node.Node, *node.Node) bool) {
	if m.NodeLen() <= 1 {
		return
	}
	for i := 1; i < m.NodeLen(); i++ {
		preIndex := i - 1
		current := m.getNode(i)
		for preIndex >= 0 && !less(m.getNode(preIndex), current) {
			m.NodeList[preIndex+1] = m.NodeList[preIndex]
			preIndex -= 1
		}
		m.NodeList[preIndex+1] = current
	}
}

func (m *Manager) Sort(mode int) {

	switch mode {
	case 0:
		for i := 0; i < m.NodeLen()/2; i++ {
			j := m.NodeLen() - i - 1
			m.NodeList[i], m.NodeList[j] = m.NodeList[j], m.NodeList[i]
		}
	case 1:
		m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
			return n1.GetProtocolMode() < n2.Protocol.GetProtocolMode()
		})
	case 2:
		m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
			return n1.GetName() < n2.GetName()
		})
	case 3:
		m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
			return n1.GetAddr() < n2.GetAddr()
		})
	case 4:
		m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
			return n1.GetPort() < n2.GetPort()
		})
	case 5:
		m.NodeSort(func(n1 *node.Node, n2 *node.Node) bool {
			return n1.TestResult < n2.TestResult
		})
	default:
		return
	}
	defer m.Save()
}

func (m *Manager) DelNode(key string) {
	indexList := core.IndexList(key, m.NodeLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	newNodeList := make([]*node.Node, 0)
	m.NodeForEach(func(i int, n *node.Node) {
		if pkg.HasIn(i, indexList) {
		} else {
			newNodeList = append(newNodeList, n)
		}
	})
	m.NodeList = newNodeList
}

func (m *Manager) DelNodeById(id string) {
	defer m.Save()
	newNodeList := make([]*node.Node, 0)
	m.NodeForEach(func(i int, n *node.Node) {
		if n.SubID == id {
		} else {
			newNodeList = append(newNodeList, n)
		}
	})
	m.NodeList = newNodeList
}

func (m *Manager) GetNodeLink(key string) []string {
	links := make([]string, 0)
	for _, index := range core.IndexList(key, m.NodeLen()) {
		links = append(links, m.GetNode(index).GetLink())
	}
	return links
}
