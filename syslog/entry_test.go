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

package syslog_test

import (
	"net"
	"strings"
	"testing"
	"time"

	syslog "github.com/m-z-b/syslogqd/syslog"
)

func TestNewEntry1(t *testing.T) {
	addr, _ := net.ResolveUDPAddr("udp", "192.168.1.99:5000")
	raw := []byte("hello")

	e := syslog.NewEntry(raw, addr)
	if e.HasSeverity() {
		t.Error("entry should not have a severity")
	}
	s := e.String()
	if strings.Index(s, "192.168.1.99") == -1 {
		t.Error("entry did not have expected IP address")
	}

	// Yes, there is a fraction of a second each year when this test is incorrect
	if strings.Index(s, time.Now().Format("2017")) == 0 {
		t.Error("entry did not include current time")
	}

	if strings.Index(s, "hello") == -1 {
		t.Error("entry did not include message")
	}
}

// We try and check that fundamental information is extracted and displayed, but not how
// exactly it is displayed
func TestNewEntry2(t *testing.T) {
	addr, _ := net.ResolveUDPAddr("udp", "192.168.1.99:5000")
	// Example from RFC5424
	raw := []byte("<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - ID47" +
		" - BOM'su root' failed for lonvick on /dev/pts/8")

	e := syslog.NewEntry(raw, addr)
	if !e.HasSeverity() {
		t.Error("entry should have a severity")
	}
	s := e.String()
	if strings.Index(s, "192.168.1.99") == -1 {
		t.Error("entry did not have expected IP address")
	}
	if strings.Index(s, "critical/auth") == -1 {
		t.Error("entry did not parse priority correctly")
	}
	// We reduce the resolution to the second
	if strings.Index(s, "2003-10-11T22:14:15Z") != 0 {
		t.Error("entry did not parse time correctly")
	}
	// The "1" at the beginning is an (optional) priority
	if strings.Index(s, "1 mymachine.example.com su - ID47 - BOM'su root' failed for lonvick on /dev/pts/8") == -1 {
		t.Error("entry did not preserve message correctly")
	}
}

// We try and check that fundamental information is extracted and displayed, but not how
// exactly it is displayed
func TestNewEntry3(t *testing.T) {
	addr, _ := net.ResolveUDPAddr("udp", "192.168.1.99:5000")
	// Example from RFC5424 - corrupted time stamp
	raw := []byte("<34> 2019-10-12X14:20:50.52+07:00 mymachine.example.com su - ID47 etc")

	e := syslog.NewEntry(raw, addr)
	if !e.HasSeverity() {
		t.Error("entry should have a severity")
	}
	s := e.String()
	if strings.Index(s, "192.168.1.99") == -1 {
		t.Error("entry did not have expected IP address")
	}
	if strings.Index(s, "critical/auth") == -1 {
		t.Error("entry did not parse priority correctly")
	}
	// We should ignore the given time and use the current time
	if strings.Index(s, "2019") == 0 {
		t.Error("entry parsed invalid time")
	}

	// Yes, there is a fraction of a second each year when this test is incorrect
	if strings.Index(s, time.Now().Format("2017")) == 0 {
		t.Error("entry did not include time")
	}

	// The "1" at the beginning is an (optional) priority
	if strings.Index(s, "2019-10-12X14:20:50.52+07:00 mymachine.example.com su - ID47 etc") == -1 {
		t.Error("entry did not preserve message correctly")
	}
}

func TestNewEntry4(t *testing.T) {
	addr, _ := net.ResolveUDPAddr("udp", "192.168.1.99:5000")
	// Example from RFC5424 - corrupted time stamp
	raw := []byte("<34> 2019-10-12T14:20:50.52+07:00 mymachine.example.com su - ID47 etc")

	e := syslog.NewEntry(raw, addr)
	if !e.HasSeverity() {
		t.Error("entry should have a severity")
	}
	s := e.String()
	if strings.Index(s, "192.168.1.99") == -1 {
		t.Error("entry did not have expected IP address")
	}
	if strings.Index(s, "critical/auth") == -1 {
		t.Error("entry did not parse priority correctly")
	}
	// We should ignore the given time and use the current time
	if strings.Index(s, "2019-10-12T07:20:50Z") != 0 {
		t.Error("time was not converted to UTC")
	}
	// The "1" at the beginning is an (optional) priority
	if strings.Index(s, "mymachine.example.com su - ID47 etc") == -1 {
		t.Error("entry did not preserve message correctly")
	}
}
