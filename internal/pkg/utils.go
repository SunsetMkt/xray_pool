package pkg

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/WQGroup/logger"
	browser "github.com/allanpk716/fake-useragent"
	detector "github.com/allanpk716/go-protocol-detector/pkg"
	"github.com/go-resty/resty/v2"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 存在且是文件
func IsFile(filePath string) bool {
	s, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// WriteJSON 将对象写入json文件
func WriteJSON(v interface{}, path string) error {
	file, e := os.Create(path)
	if e != nil {
		return e
	}
	defer file.Close()
	jsonEncode := json.NewEncoder(file)
	jsonEncode.SetIndent("", "\t")
	return jsonEncode.Encode(v)
}

// CopyDir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.DirEntry
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// CopyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer func() {
		_ = srcFd.Close()
	}()

	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer func() {
		_ = dstFd.Close()
	}()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

func HasIn(index int, indexList []int) bool {
	for _, i := range indexList {
		if i == index {
			return true
		}
	}
	return false
}

func GetXrayExeName() string {

	xrayExeName := XrayName
	sysType := runtime.GOOS
	if sysType == "windows" {
		xrayExeName += ".exe"
	}

	return xrayExeName
}

func GetGliderExeName() string {

	xrayExeName := GliderName
	sysType := runtime.GOOS
	if sysType == "windows" {
		xrayExeName += ".exe"
	}

	return xrayExeName
}

// ScanAlivePortList 扫描本地空闲的端口
func ScanAlivePortList(portRange string) []int {

	logger.Info("扫描本地可用端口，开始...")
	defer logger.Info("扫描本地可用端口，完成")
	scan := detector.NewScanTools(10, 50*time.Millisecond)

	outInfo, err := scan.Scan(detector.Common,
		detector.InputInfo{Host: "127.0.0.1", Port: portRange},
		false)
	if err != nil {
		logger.Errorf("scan alive port list error: %s", err.Error())
		return nil
	}

	outPort := make([]int, 0)
	for _, ports := range outInfo.FailedMapString {
		for _, port := range ports {
			port, err := strconv.Atoi(port)
			if err != nil {
				continue
			}
			outPort = append(outPort, port)
		}
	}

	return outPort
}

// TimeCost 耗时统计
func TimeCost() func(funcName string) {
	start := time.Now()
	return func(funcName string) {
		tc := time.Since(start)
		logger.Infof("-------------------------")
		logger.Infof("%v Time Cost = %v", funcName, tc)
		logger.Infof("-------------------------")
	}
}

func RandomSecondDuration(min, max int32) time.Duration {
	tmp := src.Int31n(max-min) + min
	return time.Duration(tmp) * time.Second
}

func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// NewHttpClient 新建一个 resty 的对象
func NewHttpClient(opt *HttpClientOptions) (*resty.Client, error) {
	//const defUserAgent = "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50"
	//const defUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36 Edg/91.0.864.41"

	var UserAgent string
	// ------------------------------------------------
	// 随机的 Browser
	UserAgent = browser.Random()
	// ------------------------------------------------
	httpClient := resty.New().SetTransport(&http.Transport{
		DisableKeepAlives:   true,
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
	})
	httpClient.SetTimeout(opt.HTMLTimeOut)
	httpClient.SetRetryCount(1)
	// ------------------------------------------------
	// 设置 Referer
	if len(opt.Referer) > 0 {
		httpClient.SetHeader("Referer", opt.Referer)
	}
	// ------------------------------------------------
	// 设置 Header
	httpClient.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   UserAgent,
	})
	// ------------------------------------------------
	// 不要求安全链接
	httpClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	// ------------------------------------------------
	// http 代理
	if opt.HttpProxyUrl != "" {
		httpClient.SetProxy(opt.HttpProxyUrl)
	} else {
		httpClient.RemoveProxy()
	}

	return httpClient, nil
}

type HttpClientOptions struct {
	HTMLTimeOut  time.Duration
	HttpProxyUrl string
	Referer      string
}

func NewHttpClientOptions(HTMLTimeOut time.Duration, httpProxyUrl string, referer string) *HttpClientOptions {
	return &HttpClientOptions{HTMLTimeOut: HTMLTimeOut, HttpProxyUrl: httpProxyUrl, Referer: referer}
}

var (
	src = rand.New(rand.NewSource(time.Now().UnixNano()))
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

const (
	XrayName   = "xray"
	GliderName = "glider"
)
