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

package syslog

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/m-z-b/syslogqd/internal/facility"
	"github.com/m-z-b/syslogqd/internal/severity"
)

var (
	rPriority        = regexp.MustCompile(`^(\d+ )?<\d{1,3}>`)
	rTimeStamp       = regexp.MustCompile(`\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d(\.\d+)?(Z|(\+|-)\d\d:\d\d)`)
	rMultiWhiteSpace = regexp.MustCompile(`\s+`)
)

// A syslog entry received from a remote client
type Entry struct {
	text        string // The entry
	remoteIP    string
	time        time.Time         // Time in UTC - either received time or time parsed from string
	severity    severity.Severity // 0..7
	facility    facility.Facility // 0..23 = kernel..local7
	hasSeverity bool              // Was severity/priority supplied?
}

// Create a syslog entry from a set of bytes
func NewEntry(bytes []byte, remoteAddress net.Addr) *Entry {
	r := &Entry{time: time.Now().UTC(), severity: severity.Default(), facility: facility.Default()}
	switch addr := remoteAddress.(type) {
	case *net.UDPAddr:
		r.remoteIP = addr.IP.String()
	case *net.TCPAddr:
		r.remoteIP = addr.IP.String()
	}

	// If this is a properly formatted string, it starts with
	// <priority>timestamp
	priority := rPriority.Find(bytes)
	if priority != nil {
		// priority is a slice <123> - we ignore the facility part
		i, err := strconv.ParseUint(string(priority[1:len(priority)-1]), 10, 8)
		if err == nil {
			s, f := uint8(i%8), uint8(i/8)
			if f <= 23 { // Legal facility index
				r.severity = severity.Severity(s)
				r.facility = facility.Facility(f)
				r.hasSeverity = true
				bytes = bytes[len(priority):]
			}
		}
	}
	// If available, the timestamp is extracted from the bytes and used as the
	// time of the entry
	ts := rTimeStamp.FindIndex(bytes)
	if ts != nil {
		t, err := time.Parse(time.RFC3339Nano, string(bytes[ts[0]:ts[1]]))
		if err == nil {
			r.time = t.UTC()
			bytes = append(bytes[0:ts[0]], bytes[ts[1]:]...)
		}
	}
	// Clean up the text by removing leading/trailing/multiple white space
	r.text = strings.TrimSpace(string(bytes))
	r.text = rMultiWhiteSpace.ReplaceAllLiteralString(r.text, " ")
	return r
}

func (self *Entry) Severity() severity.Severity {
	if !self.hasSeverity {
		panic("syslog.Entry: asked for severity when none supplied")
	}
	return self.severity
}

func (self *Entry) HasSeverity() bool {
	return self.hasSeverity
}

// A nil regex matches everything
func (self *Entry) Matches(regex *regexp.Regexp) bool {
	if regex == nil {
		return true
	} else {
		return regex.MatchString(self.text)
	}
}

func (self *Entry) String() string {
	if self.hasSeverity {
		return fmt.Sprintf("%s %s %s/%s: %s",
			self.time.Format(time.RFC3339),
			self.remoteIP,
			self.severity,
			self.facility,
			self.text)
	} else {
		return fmt.Sprintf("%s %s: %s",
			self.time.Format(time.RFC3339),
			self.remoteIP,
			self.text)
	}
}
