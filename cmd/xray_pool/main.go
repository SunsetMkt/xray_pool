package main

import (
	"github.com/allanpk716/xray_pool/internal/backend"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"os"
)

func main() {

	defer func() {
		_ = os.RemoveAll(pkg.GetTmpFolderFPath())
	}()

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
