package manager

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
)

func (m *Manager) AddFilter(key string) {
	m.Filter = append(m.Filter, node.NewNodeFilter(key))
	m.Save()
}

func (m *Manager) RunFilter(key string) {
	defer m.Save()
	newNodeList := make([]*node.Node, 0)
	if key == "" {
		m.NodeForEach(func(i int, n *node.Node) {
			if f := m.IsCanFilter(n); f != nil {
				logger.Infof("规则 [%s] 过滤节点==> %s", f.String(), n.GetName())
			} else {
				newNodeList = append(newNodeList, n)
			}
		})

	} else if f := node.NewNodeFilter(key); f != nil {
		m.NodeForEach(func(i int, n *node.Node) {
			if f.IsMatch(n) {
				logger.Infof("规则 [%s] 过滤节点==> %s", f.String(), n.GetName())
			} else {
				newNodeList = append(newNodeList, n)
			}
		})
	}
	m.NodeList = newNodeList
}

func (m *Manager) FilterForEach(funC func(int, *node.Filter)) {
	for i, f := range m.Filter {
		funC(i+1, f)
	}
}

func (m *Manager) getFilter(i int) *node.Filter {
	if i >= 0 && i < m.FilterLen() {
		return m.Filter[i]
	}
	return nil
}

func (m *Manager) GetFilter(index int) *node.Filter {
	return m.getFilter(index - 1)
}

func (m *Manager) DelFilter(key string) {
	indexList := core.IndexList(key, m.FilterLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	newFilterList := make([]*node.Filter, 0)
	m.FilterForEach(func(i int, filter *node.Filter) {
		if pkg.HasIn(i, indexList) == false {
			newFilterList = append(newFilterList, filter)
		}
	})
	m.Filter = newFilterList
}

func (m *Manager) SetFilter(key string, isOpen bool) {
	indexList := core.IndexList(key, m.FilterLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	for _, index := range indexList {
		if filter := m.GetFilter(index); filter != nil {
			filter.IsUse = isOpen
		}
	}
}

func (m *Manager) FilterLen() int {
	return len(m.Filter)
}

func (m *Manager) IsCanFilter(n *node.Node) *node.Filter {
	if n == nil {
		return nil
	}
	var f *node.Filter = nil
	m.FilterForEach(func(i int, filter *node.Filter) {
		if filter.IsUse && filter.IsMatch(n) {
			f = filter
		}
	})
	return f
}
