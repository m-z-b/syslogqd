// Copyright 2022 Mike Bell, Albion Research Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !windows

/*
syslogc - A simple test log program

Writes a sequence of syslog messages to the given host/port using the specified protocol then exits

Only builds / runs on Unix hosts

Usage:

  syslogc -host example.com -protocol udp -severity debug -delay 250ms -count 10

All options have defaults. Option -help will list them.


Mike Bell, May 2022
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"
	"time"

	"github.com/m-z-b/syslogqd/internal/severity"
)

const (
	NAME    = "syslogc"
	VERSION = "1.0"
)

// Command line options
var (
	optHost     = flag.String("host", "localhost", "syslog host or ip")
	optPort     = flag.Int("port", 514, "syslog port")
	optCount    = flag.Int("count", 5, "Number of log entries to send")
	optDelay    = flag.Duration("delay", 1*time.Second, "Delay between each log entry (e.g. 200ms)")
	optProtocol = flag.String("protocol", "udp", "Either udp or tcp")
	optSeverity = flag.String("severity", "debug", "severity name or number")
)

var minSeverity = severity.Default()

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: %s V%s [options] [message text]\n", NAME, VERSION)
		flag.PrintDefaults()
		fmt.Fprintf(w, "\n Severity Values: %s\n", severity.PossibleValues())
	}
	flag.Parse()

	if len(*optHost) < 1 {
		log.Fatal("Missing host or ip")
	}

	if *optPort < 1 || *optPort > 65535 {
		log.Fatal("port number must be in range 1..65535")
	}

	if *optCount < 1 {
		log.Fatal("-count must be greater than 1")
	}

	if *optDelay < 0 {
		log.Fatal("-delay must be greater than or equal to zero")
	}

	*optProtocol = strings.ToLower(*optProtocol)
	if *optProtocol != "udp" && *optProtocol != "tcp" {
		log.Fatal("-protocol must be udp or tcp")
	}

	msgSeverity, err := severity.Parse(*optSeverity)
	if err != nil {
		log.Fatal(err.Error())
	}

	sysLog, err := syslog.Dial(*optProtocol, fmt.Sprintf("%s:%d", *optHost, *optPort),
		syslog.LOG_WARNING|syslog.LOG_DAEMON, "syslogc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sending %d packets to %s:%d using %s with severity %s\n", *optCount, *optHost, *optPort, *optProtocol, msgSeverity)

	args := flag.Args()
	msg := strings.Join(args, " ")

	for i := 1; i <= *optCount; i++ {
		if i > 1 {
			time.Sleep(*optDelay)
		}
		os.Stdout.WriteString("*")

		s := msg
		if len(s) == 0 {
			s = fmt.Sprintf("This is message #%d from syslogc", i)
		}
		switch msgSeverity {
		case 0:
			sysLog.Emerg(s)
		case 1:
			sysLog.Alert(s)
		case 2:
			sysLog.Crit(s)
		case 3:
			sysLog.Err(s)
		case 4:
			sysLog.Warning(s)
		case 5:
			sysLog.Notice(s)
		case 6:
			sysLog.Info(s)
		case 7:
			sysLog.Debug(s)
		default:
			log.Fatalf("Illegal severity value %d", int(msgSeverity))
		}
	}
	os.Stdout.WriteString("\n")
}
