package node

import (
	"github.com/WQGroup/logger"
	"regexp"
	"strconv"
	"strings"
)

type Filter struct {
	Mode  FilterMode `json:"mode"`
	Re    string     `json:"re"`
	IsUse bool       `json:"is_use"`
}

type FilterMode string

const (
	NameFilter     FilterMode = "name"
	AddrFilter     FilterMode = "addr"
	PortFilter     FilterMode = "port"
	ProtocolFilter FilterMode = "proto"
)

func NewNodeFilter(key string) *Filter {
	if strings.HasPrefix(key, "name:") {
		return &Filter{Mode: NameFilter, Re: key[5:], IsUse: true}
	} else if strings.HasPrefix(key, "addr:") {
		return &Filter{Mode: AddrFilter, Re: key[5:], IsUse: true}
	} else if strings.HasPrefix(key, "port:") {
		return &Filter{Mode: PortFilter, Re: key[5:], IsUse: true}
	} else if strings.HasPrefix(key, "proto:") {
		return &Filter{Mode: ProtocolFilter, Re: key[6:], IsUse: true}
	} else {
		return &Filter{Mode: NameFilter, Re: key, IsUse: true}
	}

}

func (nf *Filter) IsMatch(n *Node) bool {
	reg, err := regexp.Compile(nf.Re)
	if err != nil {
		logger.Errorf("IsMatch.Compile Error: %v", err.Error())
	}
	if n != nil {
		return reg.MatchString(nf.getData(n))
	}
	return false
}

func (nf *Filter) String() string {
	switch nf.Mode {
	case AddrFilter:
		return "addr:" + nf.Re
	case PortFilter:
		return "port:" + nf.Re
	case ProtocolFilter:
		return "proto:" + nf.Re
	case NameFilter:
		return "name:" + nf.Re
	default:
		return "name:" + nf.Re
	}
}

func (nf *Filter) getData(n *Node) string {
	if n == nil {
		return ""
	}
	switch nf.Mode {
	case AddrFilter:
		return n.GetAddr()
	case PortFilter:
		return strconv.Itoa(n.GetPort())
	case ProtocolFilter:
		return string(n.GetProtocolMode())
	case NameFilter:
		return n.GetName()
	default:
		return n.GetName()
	}
}
