package logger_helper

import (
	"fmt"
	"github.com/WQGroup/logger"
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/nxadm/tail"
	"time"
)

// Listen 需要在 logger 使用一次后再调用这个函数
func Listen() {

	logger.Infoln("Start Listen Log File...")

	go initMQTTServer()

	var err error
	tailInstance, err = tail.TailFile(logger.CurrentFileName(), tail.Config{Follow: true, Poll: true})
	if err != nil {
		logger.Panic(err)
	}
	defer func() {
		_ = tailInstance.Stop()
		tailInstance.Cleanup()
	}()
	if nowLogFileName == "" {
		nowLogFileName = logger.CurrentFileName()
	}
	go func() {
		// 第一次从这里，读取更新日志，一旦隔天，那么这里就会退出
		readLogOut()
	}()

	go func() {
		// 因为在这里使用的 Log 文件存储分两个，一个是 Link，一个才是具体的文件
		// 那么在跨天的时候，需要重新打开文件
		for {
			time.Sleep(time.Second * 10)
			if nowLogFileName != logger.CurrentFileName() && logger.CurrentFileName() != "" {
				nowLogFileName = logger.CurrentFileName()
				err = tailInstance.Stop()
				if err != nil {
					logger.Panic(err)
				}
				tailInstance, err = tail.TailFile(nowLogFileName, tail.Config{Follow: true, Poll: true})
				if err != nil {
					logger.Panic(err)
				}

				go func() {
					// 读取更新日志
					readLogOut()
				}()
			}
		}
	}()

	select {}
}

func readLogOut() {
	for line := range tailInstance.Lines {
		fmt.Println("Tail --", line.Text)
	}
}

func initMQTTServer() {

	// Create the new MQTT Server.
	server := mqtt.NewServer(nil)
	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", ":19039")
	// Add the listener to the server with default options (nil).
	err := server.AddListener(tcp, nil)
	if err != nil {
		logger.Panic(err)
	}
	go func() {
		err = server.Serve()
		if err != nil {
			logger.Panic(err)
		}
	}()

	go func() {

		index := 0
		for range time.Tick(time.Second * 5) {

			showMessage := fmt.Sprintf("Hello World %d", index)
			err := server.Publish("test/02", []byte(showMessage), false)
			if err != nil {
				logger.Errorln("Publish Error:", err)
			}
			println(showMessage)
			index++
		}
	}()

	select {}
}

var (
	nowLogFileName = ""
	tailInstance   *tail.Tail
)
