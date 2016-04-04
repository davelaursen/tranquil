package tranquil

import (
	"fmt"
	"reflect"
	"time"
)

// Assertion provides functions for performing assertions when testing.
type Assertion struct {
	Test  *Test
	Value interface{}
}

// NewAssertion creates and returns a new Assertion instance.
func NewAssertion(t *Test, value interface{}) *Assertion {
	return &Assertion{t, value}
}

// ToEqual asserts that the assertion value is equal to the given value.
func (as *Assertion) ToEqual(value interface{}) {
	if !as.areEqual(as.Value, value) {
		as.Test.PrintError(fmt.Sprintf("Expected `%v` to equal `%v`", as.Value, value))
	}
}

// ToBe asserts that the assertion value is equal to the given value (equivalent to ToEqual function).
func (as *Assertion) ToBe(value interface{}) {
	as.ToEqual(value)
}

// ToNotEqual asserts that the assertion value is not equal to the given value.
func (as *Assertion) ToNotEqual(value interface{}) {
	if as.areEqual(as.Value, value) {
		as.Test.PrintError(fmt.Sprintf("Expected `%v` to not equal `%v`", as.Value, value))
	}
}

// ToNotBe asserts that the assertion value is not equal to the given value (equivalent to ToNotEqual function).
func (as *Assertion) ToNotBe(value interface{}) {
	as.ToNotEqual(value)
}

// ToBeTrue asserts that the assertion value is true.
func (as *Assertion) ToBeTrue() {
	as.ToBe(true)
}

// ToBeFalse asserts that the assertion value is false.
func (as *Assertion) ToBeFalse() {
	as.ToBe(false)
}

// ToBeTheSame asserts that the assertion value references the same object as the given value.
func (as *Assertion) ToBeTheSame(value interface{}) {
	if !as.areTheSame(as.Value, value) {
		as.Test.PrintError(fmt.Sprintf("Expected `%v` to be the same instance as `%v`", as.Value, value))
	}
}

// ToNotBeTheSame asserts that the assertion value references a different object then the given value.
func (as *Assertion) ToNotBeTheSame(value interface{}) {
	if as.areTheSame(as.Value, value) {
		as.Test.PrintError(fmt.Sprintf("Expected `%v` to not be the same instance as `%v`", as.Value, value))
	}
}

// ToBeNil asserts that the assertion value is nil.
func (as *Assertion) ToBeNil() {
	if !as.isNil(as.Value) {
		as.Test.PrintError(fmt.Sprintf("Expected `%v` to not exist", as.Value))
	}
}

// ToNotBeNil asserts that the assertion value is not nil.
func (as *Assertion) ToNotBeNil() {
	if as.isNil(as.Value) {
		as.Test.PrintError(fmt.Sprintf("Expected `%v` to exist", as.Value))
	}
}

// ToBeEmpty asserts that the assertion value is empty. An 'empty' value includes:
//   nil, "", false, 0, an map/slice/chan with length 0, a zero time value
func (as *Assertion) ToBeEmpty() {
	if !as.isEmpty(as.Value) {
		as.Test.PrintError(fmt.Sprintf("Expected `%v` to be empty", as.Value))
	}
}

// ToNotBeEmpty asserts that the assertion value is not empty. An 'empty' value includes:
//   nil, "", false, 0, an map/slice/chan with length 0, a zero time value
func (as *Assertion) ToNotBeEmpty() {
	if as.isEmpty(as.Value) {
		as.Test.PrintError(fmt.Sprintf("Expected `%v` to not be empty", as.Value))
	}
}

// ToPanic asserts that the assertion value is a function that when executed 'throws' a panic.
func (as *Assertion) ToPanic() {
	var panicObj interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicObj = r
			}
		}()
		as.Value.(func())()
	}()

	if panicObj == nil {
		as.Test.PrintError("Expected panic")
	}
}

// ToNotPanic asserts that the assertion value is a function that when executed does not 'throw' a panic.
func (as *Assertion) ToNotPanic() {
	var panicObj interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicObj = r
			}
		}()
		as.Value.(func())()
	}()

	if panicObj != nil {
		as.Test.PrintError("Expected not to panic")
	}
}

func (*Assertion) isNil(obj interface{}) bool {
	if obj == nil {
		return true
	}

	value := reflect.ValueOf(obj)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}

func (*Assertion) areEqual(act, exp interface{}) bool {
	if act == nil || exp == nil {
		return act == exp
	}

	if reflect.DeepEqual(act, exp) {
		return true
	}

	actValue := reflect.ValueOf(act)
	expValue := reflect.ValueOf(exp)
	if actValue == expValue {
		return true
	}

	t := expValue.Type()
	if actValue.Type().ConvertibleTo(t) && actValue.Convert(t) == expValue {
		return true
	}

	return false
}

func (as *Assertion) areTheSame(act, exp interface{}) bool {
	actType := reflect.TypeOf(act)
	expType := reflect.TypeOf(exp)
	if actType != expType {
		return false
	}
	return as.areEqual(act, exp)
}

func (as *Assertion) isEmpty(obj interface{}) bool {
	if obj == nil {
		return true
	}
	if obj == "" {
		return true
	}
	if obj == false {
		return true
	}

	if f, err := as.getFloat(obj); err == nil {
		if f == float64(0) {
			return true
		}
	}

	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Map, reflect.Slice, reflect.Chan:
		return v.Len() == 0
	case reflect.Ptr:
		switch obj.(type) {
		case *time.Time:
			return obj.(*time.Time).IsZero()
		default:
			return false
		}
	}
	return false
}

func (as *Assertion) getFloat(obj interface{}) (float64, error) {
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v)
	floatType := reflect.TypeOf(float64(0))
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert to float64")
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}
