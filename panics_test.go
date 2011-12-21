// Copyright 2011 Aaron Jacobs. All Rights Reserved.
// Author: aaronjjacobs@gmail.com (Aaron Jacobs)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oglematchers_test

import (
	"errors"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
)

////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////

type PanicsTest struct {
	matcherCalled bool
	suppliedCandidate interface{}
	wrappedResult MatchResult
	wrappedError error

	matcher Matcher
}

func init() { RegisterTestSuite(&PanicsTest{}) }

func (t *PanicsTest) SetUp() {
	wrapped := &fakeMatcher{
		func(c interface{}) (MatchResult, error) {
			t.matcherCalled = true
			t.suppliedCandidate = c
			return t.wrappedResult, t.wrappedError
		},
		"foo",
	}

	t.matcher = Panics(wrapped)
}

////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////

func (t *PanicsTest) Description() {
	ExpectThat(t.matcher.Description(), Equals("panics with: foo"))
}

func (t *PanicsTest) CandidateIsNil() {
	res, err := t.matcher.Matches(nil)

	ExpectThat(res, Equals(MATCH_UNDEFINED))
	ExpectThat(err, Error(Equals("which is not a zero-arg function")))
}

func (t *PanicsTest) CandidateIsString() {
	res, err := t.matcher.Matches("taco")

	ExpectThat(res, Equals(MATCH_UNDEFINED))
	ExpectThat(err, Error(Equals("which is not a zero-arg function")))
}

func (t *PanicsTest) CandidateTakesArgs() {
	res, err := t.matcher.Matches(func(i int) string { return "" })

	ExpectThat(res, Equals(MATCH_UNDEFINED))
	ExpectThat(err, Error(Equals("which is not a zero-arg function")))
}

func (t *PanicsTest) CallsFunction() {
	callCount := 0
	t.matcher.Matches(func(i int) string {
		callCount++
		return ""
	})

	ExpectThat(callCount, Equals(1))
}

func (t *PanicsTest) FunctionDoesntPanic() {
	res, err := t.matcher.Matches(func() {})

	ExpectThat(res, Equals(MATCH_FALSE))
	ExpectThat(err, Error(Equals("which didn't panic")))
}

func (t *PanicsTest) CallsWrappedMatcher() {
	expectedErr := 17
	t.matcher.Matches(func() { panic(expectedErr) })

	ExpectThat(t.suppliedCandidate, Equals(expectedErr))
}

func (t *PanicsTest) ReturnsWrappedMatcherResult() {
	t.wrappedResult = MatchResult(17)
	t.wrappedError = errors.New("taco")

	res, err := t.matcher.Matches(func() { panic(nil) })

	ExpectThat(res, Equals(t.wrappedResult))
	ExpectThat(err, Equals(t.wrappedError))
}
