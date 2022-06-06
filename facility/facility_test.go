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

package facility_test

import (
	"testing"

	"github.com/m-z-b/syslogqd/facility"
)

type _facilityTestExample struct {
	number uint8
	name   string
}

// Out of range values show a made up name
var _facilityExamples = []_facilityTestExample{
	{0, "kernel"},
	{1, "user"},
	{23, "local7"},
	{24, "facility(24)!"},
}

func TestFacilityExamples(t *testing.T) {
	for _, ex := range _facilityExamples {
		got := facility.Facility(ex.number).String()
		if got != ex.name {
			t.Errorf("got %q, wanted %q", got, ex.name)
		}
	}
}
