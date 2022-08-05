package manager

import (
	"github.com/allanpk716/xray_pool/internal/pkg/core/subscribe"
	"testing"
	"time"
)

func TestManager_AddSubscribe(t *testing.T) {

	m := NewManager()
	sub := subscribe.NewSubscribe("V2RAY_SUBSCRIBE_URL", "")
	m.AddSubscribe(sub)
}

func TestManager_UpdateNode(t *testing.T) {

	m := NewManager()
	opt := subscribe.NewUpdateOption(subscribe.NONE, "", 0, 5*time.Second)
	m.UpdateNode(opt)
}
