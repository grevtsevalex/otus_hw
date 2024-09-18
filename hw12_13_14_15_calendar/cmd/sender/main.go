package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/logger"
	rabbitqueue "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/queue/rabbitQueue"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/sender"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/sender/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := NewConfig(configFile)
	if err != nil {
		err = fmt.Errorf("config initialization: %w", err)
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	logg := logger.New(config.Logger.Level, os.Stdout)

	// create queue.
	q := rabbitqueue.NewQueue(config.Queue.Address, logg)

	// create sender.
	daemon := sender.NewSender(logg, q, time.Second*time.Duration(config.Cron.Period))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		logg.Info("sender starting...")
		daemon.Start(ctx)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	daemon.Stop(ctx)
}
