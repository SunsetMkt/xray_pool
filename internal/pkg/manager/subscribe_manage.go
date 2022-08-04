package manager

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/core/subscribe"
	"strings"
)

func (m *Manager) SubForEach(funC func(int, *subscribe.Subscribe)) {
	for i, n := range m.Subscribes {
		funC(i+1, n)
	}
}

func (m *Manager) AddSubscribe(subscribe *subscribe.Subscribe) {
	if m.HasSub(subscribe.ID()) {
		logger.Warn("该订阅链接已存在")
	} else {
		m.Subscribes = append(m.Subscribes, subscribe)
		m.Save()

	}
}

func (m *Manager) SubLen() int {
	return len(m.Subscribes)
}

func (m *Manager) getSub(i int) *subscribe.Subscribe {
	if i >= 0 && i < m.SubLen() {
		return m.Subscribes[i]
	}
	return nil
}

func (m *Manager) GetSub(i int) *subscribe.Subscribe {
	return m.getSub(i - 1)
}

func (m *Manager) UpdataNode(opt *subscribe.UpdateOption) {
	if opt.Key == "" {
		m.SubForEach(func(i int, subscribe *subscribe.Subscribe) {
			if subscribe.Using {
				m.updateNode(subscribe, opt)
			}
		})
	} else {
		for _, index := range core.IndexList(opt.Key, m.SubLen()) {
			m.updateNode(m.GetSub(index), opt)
		}
	}
}

func (m *Manager) updateNode(subscribe *subscribe.Subscribe, opt *subscribe.UpdateOption) {
	links := subscribe.UpdateNode(opt)
	if len(links) == 0 {
		return
	}
	count := 0
	m.DelNodeById(subscribe.ID())
	for _, link := range links {
		if ok := m.AddNode(node.NewNode(link, subscribe.ID(), &m.wg)); ok {
			count += 1
		}
	}
	logger.Infof("从订阅 [%s] 获取了 '%d' 个节点", subscribe.Url, count)
}

func (m *Manager) HasSub(id string) bool {
	ok := false
	m.SubForEach(func(i int, subs *subscribe.Subscribe) {
		if subs.ID() == id {
			ok = true
		}
	})
	return ok
}

func (m *Manager) DelSub(key string) {
	indexList := core.IndexList(key, m.SubLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	newSubList := make([]*subscribe.Subscribe, 0)
	m.SubForEach(func(i int, subscirbe *subscribe.Subscribe) {
		if pkg.HasIn(i, indexList) == false {
			newSubList = append(newSubList, subscirbe)
		}
	})
	m.Subscribes = newSubList
}

func (m *Manager) SetSub(key string, using, url, name string) {
	indexList := core.IndexList(key, m.SubLen())
	if len(indexList) == 0 {
		return
	}
	if len(indexList) != 1 && url != "" {
		logger.Warn("订阅链接不可以批量更改")
		return
	}
	defer m.Save()
	for _, index := range indexList {
		sub := m.GetSub(index)
		switch strings.ToLower(using) {
		case "true", "yes", "y":
			sub.Using = true
		case "false", "no", "n":
			sub.Using = false
		}
		if url != "" {
			sub.Url = url
		}
		if name != "" {
			sub.Name = name
		}
	}
}
