package main

import (
	"bytes"
	"flag"
	"github.com/0xThiebaut/ctfd-koth/conf"
	"github.com/0xThiebaut/ctfd-koth/logger"
	"github.com/0xThiebaut/ctfd-koth/monitor"
	"gopkg.in/yaml.v3"
	"os"
	"os/signal"
)

func main() {
	// Optionally retrieve non-default options
	path := flag.String("config", "koth.yml", "The King of the Hill configuration file's path.")
	flag.Parse()
	// Open the configuration file
	file, err := os.Open(*path)
	if err != nil {
		logger.Critical.Println(err)
		return
	}
	// Read the configuration
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		logger.Critical.Println(err)
		return
	}
	// Parse the configuration
	c := &conf.Configuration{}
	if err := yaml.Unmarshal(buf.Bytes(), c); err != nil {
		logger.Critical.Println(err)
		return
	}
	// Start the application
	if err := run(c); err != nil {
		logger.Critical.Println(err)
	}
}

func run(c *conf.Configuration) error {
	// Some basic sanity checks
	if err := c.Check(); err != nil {
		return err
	}
	// Create a slice of monitoring objects to close
	ms := make([]*monitor.Monitor, len(c.Flags))
	// Monitor each flag
	for i, f := range c.Flags {
		m := monitor.New(f, c.API)
		if err := m.Start(); err != nil {
			logger.Warn.Println(err)
		}
		ms[i] = m
	}
	// Wait for and intercept an interruption signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// Close each monitoring agent
	for _, m := range ms {
		if m != nil {
			if err := m.Close(); err != nil {
				logger.Warn.Println(err)
			}
		}
	}
	// Bye bye
	logger.Info.Println("gracefully shut down on interruption signal")
	return nil
}
