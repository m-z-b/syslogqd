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

package severity_test

import (
	"testing"

	"github.com/m-z-b/syslogqd/severity"
)

type _parseExample struct {
	s     string
	value uint8
}

// These should return no error AND the supplied value
var happyParseExamples = []_parseExample{
	{"0", 0},
	{"emerGency", 0},
	{"6", 6},
	{"debug", 7},
}

// These should return an error AND the supplied value
var sadParseExamples = []_parseExample{
	{"8", 7},
	{"emerg", 7},
	{"", 7},
	{"-1", 7},
	{"6hello", 7},
	{"spline", 7},
}

func TestHappyParseExamples(t *testing.T) {
	for _, ex := range happyParseExamples {
		got, err := severity.Parse(ex.s)
		if err != nil {
			t.Error(err)
		}
		wanted := severity.Severity(ex.value)
		if got != wanted {
			t.Errorf("got %q, wanted %q", got, wanted)
		}
	}
}

func TestSadParseExamples(t *testing.T) {
	for _, ex := range sadParseExamples {
		got, err := severity.Parse(ex.s)
		wanted := severity.Severity(ex.value)
		if got != wanted {
			t.Errorf("expected value of %q to be %s", ex.s, wanted)
		}
		if err == nil {
			t.Errorf("expected error parsing %q", ex.s)
		}
	}
}

// a should be greater than or equal to b
var comparisonExamples = []struct {
	more string
	less string
}{
	{"0", "0"},
	{"0", "1"},
	{"2", "7"},
	{"emergency", "info"},
	{"info", "debug"},
	{"unknown", "debug"},
}

func TestAsOrMoreSevereThan(t *testing.T) {
	for _, ex := range comparisonExamples {
		more, _ := severity.Parse(ex.more)
		less, _ := severity.Parse(ex.less)
		if !more.AsOrMoreSevereThan(less) {
			t.Errorf("expected %q to be as or more severe than %q", ex.more, ex.less)
		}
	}
}
