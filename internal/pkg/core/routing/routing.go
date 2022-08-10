package routing

import (
	"encoding/json"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/allanpk716/xray_pool/internal/pkg/core"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Routing struct {
	Proxy  []*OneRouting `json:"proxy"`
	Direct []*OneRouting `json:"direct"`
	Block  []*OneRouting `json:"block"`
}

func NewRouting() *Routing {

	route := &Routing{
		Proxy:  make([]*OneRouting, 0),
		Direct: make([]*OneRouting, 0),
		Block:  make([]*OneRouting, 0),
	}
	if _, err := os.Stat(core.RoutingFile); os.IsNotExist(err) {
		route.save()
	} else {
		file, _ := os.Open(core.RoutingFile)
		defer func() {
			_ = file.Close()
		}()
		err = json.NewDecoder(file).Decode(route)
		if err != nil {
			logger.Panic(err)
		}
	}

	return route
}

func (r *Routing) save() {
	err := pkg.WriteJSON(r, core.RoutingFile)
	if err != nil {
		logger.Error(err)
	}
}

// GetRuleMode 判断是IP规则还是域名规则
// IP|Domain
func GetRuleMode(str string) Mode {
	if strings.HasPrefix(str, "geoip:") {
		return ModeIP
	}
	if strings.Contains(str, "ip.dat:") {
		return ModeIP
	}
	pattern := `(?:^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})(?:/(?:[1-9]|[1-2][0-9]|3[0-2]){1})?$`
	re, _ := regexp.Compile(pattern)
	if re.MatchString(str) {
		return ModeIP
	}
	return ModeDomain
}

// AddRule 添加规则
func (r *Routing) AddRule(rt Type, list ...string) int {
	defer r.save()
	count := 0
	for _, rule := range list {
		if rule != "" {
			oneRouting := &OneRouting{
				Data: rule,
				Mode: GetRuleMode(rule),
			}
			count += 1
			switch rt {
			case TypeBlock:
				r.Block = append(r.Block, oneRouting)
			case TypeDirect:
				r.Direct = append(r.Direct, oneRouting)
			case TypeProxy:
				r.Proxy = append(r.Proxy, oneRouting)
			}
		}
	}
	return count
}

// GetRule 获取规则
func (r *Routing) GetRule(rt Type, key string) [][]string {
	var rules []*OneRouting
	switch rt {
	case TypeDirect:
		rules = r.Direct
	case TypeProxy:
		rules = r.Proxy
	case TypeBlock:
		rules = r.Block
	}
	indexList := core.IndexList(key, len(rules))
	result := make([][]string, 0, len(indexList))
	for _, x := range indexList {
		r := rules[x-1]
		result = append(result, []string{
			strconv.Itoa(x),
			string(r.Mode),
			r.Data,
		})
	}
	return result
}

// GetRulesGroupData 对路由数据进行分组
func (r *Routing) GetRulesGroupData(rt Type) ([]string, []string) {
	ips := make([]string, 0)
	domains := make([]string, 0)
	var rules []*OneRouting
	switch rt {
	case TypeDirect:
		rules = r.Direct
	case TypeProxy:
		rules = r.Proxy
	case TypeBlock:
		rules = r.Block
	}
	for _, x := range rules {
		if x.Mode == "Domain" {
			domains = append(domains, x.Data)
		} else {
			ips = append(ips, x.Data)
		}
	}
	return ips, domains
}

// DelRule 删除规则
func (r *Routing) DelRule(rt Type, key string) {
	var rules []*OneRouting
	switch rt {
	case TypeDirect:
		rules = r.Direct
	case TypeProxy:
		rules = r.Proxy
	case TypeBlock:
		rules = r.Block
	}
	indexList := core.IndexList(key, len(rules))
	if len(indexList) == 0 {
		return
	}
	defer r.save()
	result := make([]*OneRouting, 0)
	for i, rule := range rules {
		if pkg.HasIn(i+1, indexList) == false {
			result = append(result, rule)
		}
	}
	switch rt {
	case TypeDirect:
		r.Direct = result
	case TypeProxy:
		r.Proxy = result
	case TypeBlock:
		r.Block = result
	}
	logger.Info("删除了 [", len(indexList), "] 条规则")
}

// RuleLen 有多少条规则
func (r *Routing) RuleLen(rt Type) int {
	switch rt {
	case TypeDirect:
		return len(r.Direct)
	case TypeProxy:
		return len(r.Proxy)
	case TypeBlock:
		return len(r.Block)
	default:
		return 0
	}
}
