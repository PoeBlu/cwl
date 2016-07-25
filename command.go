package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type command struct {
	profile      string
	region       string
	logGroupName string
	filter       string
	start        time.Time
	end          time.Time
	limit        int64
	tail         bool
	refresh      int64
	interleaved  bool
	help         bool
}

// ParseCommand parses the command line and creates a new command to run.
func parseCommand() *command {
	startParam := "1 minute ago"
	endParam := "now"

	command := &command{interleaved: true, limit: 50, tail: false}

	flag.StringVar(&command.profile, "profile", "", "AWS credential profile to use.")
	flag.StringVar(&command.region, "region", "", "AWS region to request logs from")
	flag.StringVar(&command.logGroupName, "group", "", "Log group name to read from")
	flag.StringVar(&command.filter, "filter", "", "Filter pattern to appy")
	flag.StringVar(&startParam, "start", "1 minute ago", "The RFC3339 time that log events should start from")
	flag.StringVar(&endParam, "end", "now", "The RFC3339 time that log events should end")

	flag.BoolVar(&command.tail, "tail", false, "Read log messages continuously")
	flag.Int64Var(&command.refresh, "refresh", 5, "Refresh rate for tailing logs, in seconds.")
	flag.BoolVar(&command.interleaved, "interleaved", true, "Interleave log messages between sources")
	flag.Parse()

	if command.help {
		usage()
	}

	if command.region == "" || command.logGroupName == "" {
		usage()
	}

	if startParam == "1 minute ago" {
		command.start = time.Now().Add(-1 * time.Minute)
	} else {
		startTime, err := time.Parse(time.RFC3339, startParam)

		if err != nil {
			fmt.Printf("Start time must be given in RFC 3339 time: %s\n", time.Now().Format(time.RFC3339))
			fmt.Println()
			usage()
		}
		command.start = startTime
	}

	if endParam == "now" {
		command.end = time.Now()
	} else {
		endTime, err := time.Parse(time.RFC3339, endParam)

		if err != nil {
			fmt.Printf("End time must be given in RFC 3339 time: %s\n", time.Now().Format(time.RFC3339))
			fmt.Println()
			usage()
		}
		command.end = endTime
	}

	return command
}

func usage() {
	fmt.Println("cwl - A command line tool for reviewing Amazon CloudWatch Logs")
	fmt.Println()
	fmt.Println("Parameters:")
	flag.PrintDefaults()
	os.Exit(-1)
}
