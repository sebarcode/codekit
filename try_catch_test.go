package codekit_test

import (
	"errors"
	"math"
	"testing"

	"github.com/sebarcode/codekit"
)

func TestTryCatchBase(t *testing.T) {
	c := 0
	err := codekit.Try(func() {
		codekit.CheckError("test error")
	}).
		Catch(func(err error) {
			if err.Error() != "test error" {
				t.Error("error not match")
			}
		}).
		Finally(func() {
			c = 1
		}).
		Do()

	if err == nil {
		t.Error(errors.New("error not catched"))
	}

	if c != 1 {
		t.Error("finally not executed")
	}

	t.Log("TestTryCatch done")
}

func TestTryCatchOK(t *testing.T) {
	c := 0
	err := codekit.Try(func() {
		c = 5
	}).
		Catch(func(err error) {
			if err.Error() != "test error" {
				t.Error("error not match")
			}
		}).
		Finally(func() {
			c = c * 2
		}).
		Do()

	if err != nil {
		t.Error(errors.New("error not catched"))
	}

	if c != 10 {
		t.Error("finally not executed")
	}

	t.Log("TestTryCatchOK done")
}

func TestTryCatchDynamic(t *testing.T) {
	raw := 0
	c := 0
	d := 0
	err := codekit.Try(func() {
		a := codekit.RandInt(100)
		b := codekit.RandInt(5)
		c = a * b

		if math.Mod(float64(a), 2) == 1 {
			codekit.CheckError("odd value")
		}
		raw = c
		c = 50
	}).
		Catch(func(err error) {
			c = 101
		}).
		Finally(func() {
			d = c * 3
		}).
		Do()

	if err == nil {
		if c != 50 {
			t.Error("try doesnt executed")
		}
	} else if c != 101 {
		t.Error("catch doesnt executed")
	}

	if d != c*3 {
		t.Error("finally not executed")
	}

	t.Logf("raw: %d c: %d, d: %d", raw, c, d)
	t.Log("TestTryCatchDynamic done")
}
