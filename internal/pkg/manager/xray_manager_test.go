package manager

import (
	"testing"
)

func TestManager_StartXray(t *testing.T) {

	m := NewManager()
	m.StartXray()

	select {}
}
