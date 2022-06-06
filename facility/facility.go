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

package facility

import (
	"fmt"
)

// Facility - as defined in RFC 5424
type Facility uint8

// names - Names of facility values
// These are taken from Wikipedia rather than the RFC 5424
// The values vary between systems and implementations
var names []string = []string{
	"kernel",
	"user",
	"mail",
	"daemon",
	"auth",
	"syslog",
	"lpr",
	"news",
	"uucp",
	"cron",
	"authpriv",
	"ftp",
	"ntp",
	"security",
	"console",
	"solaris-cron",
	"local0",
	"local1",
	"local2",
	"local3",
	"local4",
	"local5",
	"local6",
	"local7",
}

func (s Facility) String() string {
	if int(s) < len(names) {
		return names[s]
	} else {
		return fmt.Sprintf("facility(%d)!", int(s))
	}
}

// Default() is used when a facility is missing from a message
//
// Note that we can't use the Go convention of using the zero value
// as the numeric mapping is defined in the RFC and we don't want to
// use "kernel" as the default
func Default() Facility {
	return Facility(1) // "user"
}
