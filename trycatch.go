package codekit

import (
	"errors"
	"fmt"
)

type tcf struct {
	TryFn     func()
	CatchFn   func(error)
	FinallyFn func()
}

func Try(fn func()) *tcf {
	tcf := &tcf{TryFn: fn}
	return tcf
}

func (tcf *tcf) Catch(fn func(error)) *tcf {
	tcf.CatchFn = fn
	return tcf
}

func (tcf *tcf) Finally(fn func()) *tcf {
	tcf.FinallyFn = fn
	return tcf
}

func (tcf *tcf) Do() error {
	var err error

	if tcf.FinallyFn != nil {
		defer tcf.FinallyFn()
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case error:
					err = x

				case string:
					err = errors.New(x)

				default:
					err = fmt.Errorf("%v", x)
				}

				if tcf.CatchFn != nil {
					tcf.CatchFn(err)
				}
			}
		}()

		tcf.TryFn()
	}()

	return err
}

func CheckError(err interface{}) {
	switch x := err.(type) {
	case error:
		panic(x)

	case string:
		panic(errors.New(err.(string)))

	case nil:
		return

	default:
		panic(fmt.Errorf("%v", x))
	}
}
