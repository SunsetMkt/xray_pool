package core

import (
	"github.com/allanpk716/xray_pool/internal/pkg"
	"path/filepath"
)

var (
	DataFile    = filepath.Join(pkg.GetConfigRootDirFPath(), "data.json")
	SettingFile = filepath.Join(pkg.GetConfigRootDirFPath(), "setting.toml")
	RoutingFile = filepath.Join(pkg.GetConfigRootDirFPath(), "routing.json")
)
