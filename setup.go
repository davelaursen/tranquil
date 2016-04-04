package tranquil

// Setup provides the ability to setup tests with automatically executed functions.
type Setup struct {
	beforeEachFns []func()
	afterEachFns  []func()
}

// BeforeEach adds a function that gets executed before each behavior ('it' function) is tested.
func (s *Setup) BeforeEach(fn func()) {
	if fn != nil {
		s.beforeEachFns = append(s.beforeEachFns, fn)
	}
}

// AfterEach adds a function that gets executed after each behavior ('it' function) is tested.
func (s *Setup) AfterEach(fn func()) {
	if fn != nil {
		s.afterEachFns = append(s.afterEachFns, fn)
	}
}

func (s *Setup) execBeforeEachFns() {
	for _, f := range s.beforeEachFns {
		f()
	}
}

func (s *Setup) execAfterEachFns() {
	for _, f := range s.afterEachFns {
		f()
	}
}
