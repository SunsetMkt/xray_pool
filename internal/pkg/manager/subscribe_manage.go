package manager

import (
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core"
	"github.com/allanpk716/xray_pool/internal/pkg/core/node"
	"github.com/allanpk716/xray_pool/internal/pkg/core/subscribe"
	"strings"
)

func (m *Manager) SubscribeForEach(funC func(int, *subscribe.Subscribe)) {
	for i, n := range m.Subscribes {
		funC(i+1, n)
	}
}

func (m *Manager) AddSubscribe(subscribe *subscribe.Subscribe) {
	if m.HasSubscribe(subscribe.ID()) {
		logger.Warn("该订阅链接已存在")
	} else {
		m.Subscribes = append(m.Subscribes, subscribe)
		m.Save()

	}
}

func (m *Manager) SubscribeLen() int {
	return len(m.Subscribes)
}

func (m *Manager) getSub(i int) *subscribe.Subscribe {
	if i >= 0 && i < m.SubscribeLen() {
		return m.Subscribes[i]
	}
	return nil
}

func (m *Manager) GetSubscribe(i int) *subscribe.Subscribe {
	return m.getSub(i - 1)
}

func (m *Manager) UpdateNode(opt *subscribe.UpdateOption) {
	if opt.Key == "" {
		m.SubscribeForEach(func(i int, subscribe *subscribe.Subscribe) {
			if subscribe.Using {
				m.updateNode(subscribe, opt)
			}
		})
	} else {
		for _, index := range core.IndexList(opt.Key, m.SubscribeLen()) {
			m.updateNode(m.GetSubscribe(index), opt)
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

// HasSubscribe 是否已经订阅，ID 由 URL 的 MD5 值组成
func (m *Manager) HasSubscribe(id string) bool {
	ok := false
	m.SubscribeForEach(func(i int, subs *subscribe.Subscribe) {
		if subs.ID() == id {
			ok = true
		}
	})
	return ok
}

func (m *Manager) DelSubscribe(key string) {
	indexList := core.IndexList(key, m.SubscribeLen())
	if len(indexList) == 0 {
		return
	}
	defer m.Save()
	newSubList := make([]*subscribe.Subscribe, 0)
	m.SubscribeForEach(func(i int, subscirbe *subscribe.Subscribe) {
		if pkg.HasIn(i, indexList) == false {
			newSubList = append(newSubList, subscirbe)
		}
	})
	m.Subscribes = newSubList
}

func (m *Manager) SetSubscribe(key string, using, url, name string) {
	indexList := core.IndexList(key, m.SubscribeLen())
	if len(indexList) == 0 {
		return
	}
	if len(indexList) != 1 && url != "" {
		logger.Warn("订阅链接不可以批量更改")
		return
	}
	defer m.Save()
	for _, index := range indexList {
		sub := m.GetSubscribe(index)
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
