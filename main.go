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

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/m-z-b/syslogqd/internal/listener"
	"github.com/m-z-b/syslogqd/internal/reporter"
	"github.com/m-z-b/syslogqd/internal/severity"
	"github.com/m-z-b/syslogqd/internal/syslog"
)

// Tombstone information for the program
const (
	NAME    = "syslogqd"
	VERSION = "1.0"
)

// Command line arguments
// (note that these are displayed alphabetically)
var (
	optPort     = flag.Int("port", 514, "port to listen on (UDP-only)")
	optFilename = flag.String("file", "", "write output to file")
	optQuiet    = flag.Bool("quiet", false, "do not write to standard output")
	optSeverity = flag.String("severity", "debug", "minimum severity of events to report")
)

var (
	output      *os.File                               // File to write to (if not null)
	minSeverity severity.Severity = severity.Default() // Minimum severity to display
)

// FatalError prints a message followed by a newline to stderr and exits the program
//
// If no args are supplied, the format string is written directly. If args are supplied,
// the format string is used along with the arguments.
func FatalError(format string, args ...any) {
	format += "\n" // Feels wrong
	if len(args) == 0 {
		os.Stderr.WriteString(format)
	} else {
		fmt.Fprintf(os.Stderr, format, args...)
	}
	os.Exit(1)
}

// CheckForFatalError calls FatalError using the remaining arguments if err is not nil
//
// Arguments which support the error interface will be converted to strings using an Error()
// call before being interpolated in the format string.
func CheckForFatalErrorF(err error, format string, args ...any) {
	if err == nil {
		return
	}
	for i := range args {
		if arg, ok := args[i].(error); ok {
			args[i] = arg.Error() // Replace an error with a error.Error()
		}
	}
	FatalError(format, args...)
}

// If err is not nil, print its error message and exit
func CheckForFatalError(err error) {
	if err != nil {
		FatalError(err.Error())
	}
}

// The Entry point...
//
// This interprets the command line arguments and sets up the listeners and a reporter
func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: %s V%s [options]\n", NAME, VERSION)
		flag.PrintDefaults()
		fmt.Fprintf(w, "\n Severity Values: %s\n", severity.PossibleValues())
	}

	flag.Parse()

	if *optPort <= 0 || *optPort > 65535 {
		FatalError("-port must be in the range 1..65535")
	}

	if *optFilename == "" && *optQuiet {
		FatalError("Can only specify -quiet if -file is specified")
	}

	var err error
	if *optSeverity != "" {
		minSeverity, err = severity.Parse((*optSeverity))
		CheckForFatalError(err)
	}

	reporter := reporter.NewReporter()

	if *optFilename != "" {
		output, err = os.OpenFile(*optFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		CheckForFatalErrorF(err, "Could not open %s: %s", *optFilename, err)
		defer output.Close()
		reporter.AddOutput(output)
	}

	if !*optQuiet {
		fmt.Printf("%s V%s listening on port %d for severity >= %s\nUse Ctrl-C to exit\n", NAME, VERSION, *optPort, minSeverity)
		reporter.AddOutput(os.Stdout)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTERM)

	newswire := make(syslog.Channel, 10)
	go reporter.Report(newswire, minSeverity)

	udpListener, err := listener.NewUDPListener(*optPort, newswire)
	CheckForFatalError(err)
	go udpListener.Listen()

	tcpListener, err := listener.NewTCPListener(*optPort, newswire)
	CheckForFatalError(err)
	go tcpListener.Listen()

	<-done // Wait

	if !*optQuiet {
		fmt.Println(NAME, "Normal exit")
	}

}
