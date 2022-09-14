package glider_helper

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/pkg/errors"
	"os/exec"
	"path/filepath"
	"strings"
)

type GliderHelper struct {
	gliderCmd  *exec.Cmd // 正向代理服务器实例
	gliderPath string    // glider 程序的路径
}

func NewGliderHelper() *GliderHelper {
	return &GliderHelper{}
}

// Check 检查 Xray 程序和需求的资源是否已经存在，不存在则需要提示用户去下载
func (g *GliderHelper) Check() bool {

	if g.gliderPath == "" {
		// 在这个目录下进行搜索是否存在 Xray 程序
		nowRootPath := pkg.GetBaseThingsFolderFPath()
		gliderExeName := pkg.GetGliderExeName()
		gliderExeFullPath := filepath.Join(nowRootPath, gliderExeName)
		if pkg.IsFile(gliderExeFullPath) == false {
			logger.Error(GliderDownloadInfo)
			return false
		}
		g.gliderPath = gliderExeFullPath
		return true
	} else {
		return true
	}
}

func (g *GliderHelper) Start(healthCheckUrl string, healthCheckInterval int, forwardServerHttpPort int, socksPorts []int, GliderStrategy string) error {

	// 构建正向代理服务器启动的命令
	runCommand := fmt.Sprintf("-listen :%d -strategy %s", forwardServerHttpPort, GliderStrategy)
	for _, socksPort := range socksPorts {
		runCommand += fmt.Sprintf(" -forward socks5://127.0.0.1:%d", socksPort)
	}
	// -check http://www.msftconnecttest.com/connecttest.txt#expect=200
	if healthCheckUrl != "" {
		runCommand += fmt.Sprintf(" -check %s#expect=200", healthCheckUrl)
		runCommand += fmt.Sprintf(" -checkinterval %d", healthCheckInterval)
	}
	cmdArgs := strings.Fields(runCommand)
	g.gliderCmd = exec.Command(g.gliderPath, cmdArgs...)
	err := g.gliderCmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (g *GliderHelper) Stop() error {
	defer func() {
		g.gliderCmd = nil
	}()

	if g.gliderCmd != nil {
		err := g.gliderCmd.Process.Kill()
		if err != nil {
			return nil
		}
	}
	return nil
}

var (
	GliderDownloadInfo = errors.New(fmt.Sprintf("缺少 Glider 可执行程序，请去 https://github.com/nadoo/glider/releases 下载对应平台的程序，解压放入 %v 文件夹中", pkg.GetBaseThingsFolderAbsFPath()))
)
