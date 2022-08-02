package main

import "github.com/allanpk716/xray_pool/internal/backend"

func main() {

	restartSignal := make(chan interface{}, 1)
	defer close(restartSignal)
	bend := backend.NewBackEnd(19035, restartSignal)
	go bend.Restart()
	restartSignal <- 1
	// 阻塞
	select {}
}
