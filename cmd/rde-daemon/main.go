package main

import (
	"log"
	"os"
	"os/signal"
	"rde-daemon/internal/daemon"
	"syscall"
)

func main() {
	var err chan error = make(chan error)
	serviceFS := daemon.FS{}
	go serviceFS.Run(&err)

	// MARK: - Streamer
	streamer := daemon.Streamer{}
	go streamer.Run(&err)

	var sig chan os.Signal = make(chan os.Signal)
	signal.Notify(
		sig,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGKILL,
		syscall.SIGINT,
	)

	select {
	case <-sig:
	case e := <-err:
		log.Printf("err catched: %v", e)
	}
	log.Println("Shutting down")
	serviceFS.Shutdown()
	streamer.Shutdown()

	close(sig)
	os.Exit(0)
}
