package manager

import "github.com/allanpk716/xray_pool/internal/pkg/core/routing"

func (m *Manager) AddRule(rt routing.Type, list ...string) {

	m.routing.AddRule(rt, list...)
}

func (m *Manager) GetRule(rt routing.Type, key string) [][]string {

	return m.routing.GetRule(rt, key)
}

func (m *Manager) GetRulesGroupData(rt routing.Type) ([]string, []string) {
	return m.routing.GetRulesGroupData(rt)
}

func (m *Manager) DelRule(rt routing.Type, key string) {
	m.routing.DelRule(rt, key)
}

func (m *Manager) RuleLen(rt routing.Type) int {
	return m.routing.RuleLen(rt)
}
