package main

import "github.com/allanpk716/xray_pool/internal/backend"

func main() {

	restartSignal := make(chan interface{}, 1)
	exitSignal := make(chan interface{}, 1)
	defer close(restartSignal)
	defer close(exitSignal)
	bend := backend.NewBackEnd(restartSignal, exitSignal)
	go bend.Restart()
	restartSignal <- 1
	// 阻塞
	select {}
}
