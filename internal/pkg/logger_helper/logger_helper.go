package logger_helper

import (
	"fmt"
	"github.com/WQGroup/logger"
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/nxadm/tail"
	cmap "github.com/orcaman/concurrent-map/v2"
	"sync"
	"time"
)

// Listen 需要在 logger 使用一次后再调用这个函数
func Listen() {

	var err error
	logger.Infoln("Start Listen Log File...")

	logQueue = make([]string, 0)
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
				ClearLogQueue()
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

	// Create the new MQTT Server.
	mqttServer = mqtt.NewServer(nil)
	// Create a TCP listener on a standard port.
	ws := listeners.NewWebsocket("t1", ":19039")
	//tcp := listeners.NewTCP("t1", ":19039")
	// Add the listener to the mqttServer with default options (nil).
	err = mqttServer.AddListener(ws, nil)
	if err != nil {
		logger.Panic(err)
	}
	go func() {
		err = mqttServer.Serve()
		if err != nil {
			logger.Panic(err)
		}
	}()
	// 连接
	mqttServer.Events.OnConnect = func(cl events.Client, pk events.Packet) {
		fmt.Printf("<< OnConnect client connected %s\n", cl.ID)
		AddClientStatus(cl.ID)
		go sendLog(cl.ID)
	}
	// 取消订阅主题
	mqttServer.Events.OnDisconnect = func(cl events.Client, err error) {
		fmt.Printf("<< OnDisconnect client disconnected %s: %v\n", cl.ID, err)
		DelClientStatus(cl.ID)
	}

	select {}
}

func readLogOut() {
	for line := range tailInstance.Lines {
		Add2LogQueue(line.Text)
	}
}

func initMQTTServer() {

	select {}
}

func Add2LogQueue(message string) {
	logQueueLock.Lock()
	defer logQueueLock.Unlock()
	logQueue = append(logQueue, message)
}

func GetOneLogQueue(index int) string {
	logQueueLock.Lock()
	defer logQueueLock.Unlock()
	if index < len(logQueue) {
		return logQueue[index]
	} else {
		return ""
	}
}

func ClearLogQueue() {
	logQueueLock.Lock()
	defer logQueueLock.Unlock()
	logQueue = make([]string, 0)
}

func sendLog(clientId string) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 为每一个客户端开启一个发送日志的线程
	for true {
		select {
		case <-GetClientStatus(clientId).ExitSignal:
			return
		default:
			// 发送日志
			time.Sleep(time.Millisecond * 100)
			nowClient := GetClientStatus(clientId)
			if nowClient == nil {
				return
			}
			oneLineLogMessage := GetOneLogQueue(nowClient.LogIndex)
			if oneLineLogMessage == "" {
				continue
			}
			nowClient.LogIndex++

			err := mqttServer.Publish(fmt.Sprintf("log/%s", clientId), []byte(oneLineLogMessage), false)
			if err != nil {
				logger.Errorln("Publish Error:", err)
			}
		}
	}
}

var (
	nowLogFileName = ""                        // 当前日志的缓存名称，需要循环对比，因为每天会新建一个日志名称
	tailInstance   *tail.Tail                  // 读取本地日志缓存的 tail 实例
	logQueue       = make([]string, 0)         // 本地缓存日志的队列
	logQueueLock   sync.Mutex                  // 本地缓存日志的队列的锁
	mqttServer     *mqtt.Server                // MQTT 服务器实例
	clientIds      = cmap.New[*ClientStatus]() // 每个客户端的状态
)

type ClientStatus struct {
	ID         string
	LogIndex   int
	ExitSignal chan interface{}
}

func AddClientStatus(clientId string) {
	clientIds.Set(clientId, &ClientStatus{ID: clientId, LogIndex: 0, ExitSignal: make(chan interface{}, 1)})
}

func GetClientStatus(clientId string) *ClientStatus {
	client, found := clientIds.Get(clientId)
	if found == false {
		return nil
	}
	return client
}

func DelClientStatus(clientId string) {

	client, found := clientIds.Get(clientId)
	if found == false {
		return
	}
	client.ExitSignal <- true
	close(client.ExitSignal)
	clientIds.Remove(clientId)
}
