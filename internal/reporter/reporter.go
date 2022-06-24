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

package reporter

import (
	"fmt"
	"os"
	"regexp"

	"github.com/m-z-b/syslogqd/internal/severity"
	"github.com/m-z-b/syslogqd/internal/syslog"
)

// A Reporter repeatedly receives a syslog.Entry and writes it to a set of output streams
//
//    newswire := make( syslog.Channel, 10 )
//    ...
//    r := NewReporter()
//    r.AddOutput(os.StdOut)
//    go r.Report(newswire, Severity.Default() )
//
type Reporter struct {
	files     []*os.File
	mustMatch *regexp.Regexp
}

// NewReporter constructs a new Reporter instance
func NewReporter(mustMatch *regexp.Regexp) (r *Reporter) {
	r = &Reporter{mustMatch: mustMatch}
	r.files = make([]*os.File, 0, 5)
	return
}

// AddOutput adds a stream the output should be sent to
func (self *Reporter) AddOutput(file *os.File) *Reporter {
	self.files = append(self.files, file)
	return self
}

// Write a syslog entry to all the file streams
func (self *Reporter) reportEntry(e *syslog.Entry) {
	for _, f := range self.files {
		fmt.Fprintln(f, e)
	}
}

// Report gets a new SyslogEntry from the newswire channel and reports it to all outputs
// if it has at least the given severity
func (self *Reporter) Report(newswire syslog.Channel, minSeverity severity.Severity) {
	for {
		var e = <-newswire
		if !e.HasSeverity() || e.Severity().AsOrMoreSevereThan(minSeverity) {
			if e.Matches(self.mustMatch) { // e.Matches handles nil as match any
				self.reportEntry(e)
			}
		}
	}
}
