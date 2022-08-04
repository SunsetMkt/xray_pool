package core

import (
	"github.com/allanpk716/xray_pool/internal/pkg"
	"path/filepath"
)

var (
	AppSettings = filepath.Join(pkg.GetConfigRootDirFPath(), "xray_pool_config.json")
	RoutingFile = filepath.Join(pkg.GetConfigRootDirFPath(), "routing.json")
)
