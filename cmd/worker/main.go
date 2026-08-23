package worker

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	log.Println("industrial worker started")
	for {
		select {
		case <-ctx.Done():
			log.Println("industrial worker stopped")
			return
		case <-ticker.C:
			log.Println("heartbeat monitor tick")
		}
	}
}
