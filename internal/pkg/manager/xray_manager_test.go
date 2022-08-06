package manager

import (
	"testing"
)

func TestManager_StartXray(t *testing.T) {

	m := NewManager()
	bok, aliveNodeIndexList, alivePorts := m.GetsValidNodesAndAlivePorts()
	if bok == false {
		t.Error("获取有效的节点失败")
		return
	}
	bok = m.StartXray(aliveNodeIndexList, alivePorts)
	if bok == false {
		t.Error("启动 xray 失败")
		return
	}
	select {}
}
