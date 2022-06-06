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

package severity

import (
	"fmt"
	"strconv"
	"strings"
)

// Severity of a syslog message as defined in RFC 5424
type Severity uint8

// names are taken from RFC 5424 Table 2 - we use debug(7) for entries with no defined severity
var names []string = []string{"emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"}

func (s Severity) String() string {
	if int(s) < len(names) {
		return names[s]
	} else {
		return fmt.Sprintf("severity(%d)!", int(s))
	}
}

// AsOrMoreSevereThan(other) returns T if self is as or more severe than other
func (s Severity) AsOrMoreSevereThan(other Severity) bool {
	// 0 is more severe than 1
	return s <= other
}

// Parse converts a string into a severity
//
// Values which are not recognized are given an UnknownSeverity() with a non-nil err
func Parse(s string) (Severity, error) {
	for i, v := range names {
		if strings.EqualFold(s, v) {
			return Severity(i), nil
		}
	}
	i, err := strconv.ParseUint(s, 10, 3)
	if err == nil && int(i) < len(names) {
		return Severity(i), nil
	}
	return Default(), fmt.Errorf("Unknown severity value \"%s\"", s)
}

// Default (the lowest) is used if no recognizable severity is specified
//
// Note that we can't use the Go convention of the zero value being the default
// as the number/value mapping is defined in the RFC
func Default() Severity {
	return Severity(len(names) - 1)
}
