package tranquil

import (
	"fmt"
	"testing"
)

// It is a function that tests a particular behavor.
type It func(desc string, expectFn func(Expect))

// Expect creates an assertion for a given value.
type Expect func(val interface{}) *Assertion

// Describe provides the foundation for testing behavior against a subject (e.g. object, function).
func Describe(desc string, t *testing.T, wrapper func(*Setup, It)) {
	s := &Setup{beforeEachFns: []func(){}, afterEachFns: []func(){}}
	wrapper(s, func(itDesc string, expectFn func(Expect)) {
		test := Test{
			Description: fmt.Sprintf("%s %s", desc, itDesc),
			T:           t,
			ExpectFn:    expectFn,
			Setup:       s,
		}
		test.Run()
	})

}
