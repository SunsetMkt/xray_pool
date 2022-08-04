package routing

type Type string

const (
	TypeProxy  Type = "Proxy"
	TypeDirect Type = "Direct"
	TypeBlock  Type = "Block"
)

type Mode string

const (
	ModeIP     Mode = "IP"
	ModeDomain      = "Domain"
)

type OneRouting struct {
	Data string `json:"data"`
	Mode Mode   `json:"mode"`
}
