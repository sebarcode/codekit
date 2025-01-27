package codekit_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/sebarcode/codekit"
)

func TestCast(t *testing.T) {
	dataCount := 10000000
	datas := make([]interface{}, dataCount)
	for i := 0; i < dataCount; i++ {
		data := time.Now().Add(-time.Duration(codekit.RandInt(1000)) * time.Second)
		datas[i] = &data
	}

	// cast normal
	t0 := time.Now()
	for i := 0; i < dataCount; i++ {
		_, ok := datas[i].(*time.Time)
		if !ok {
			t.Fatalf("cast failed: %v, index: %d", datas[i], i)
		}
		/*
			if cast.UnixMicro() != datas[i].(*time.Time).UnixMicro() {
				t.Fatalf("cast not equal: %v, index: %d", datas[i], i)
			}
		*/
	}
	t.Logf("cast normal: %s", time.Since(t0).String())

	// cast full reflect
	t0 = time.Now()
	for i := 0; i < dataCount; i++ {
		dataGen := datas[i]
		cast := new(time.Time)
		reflect.ValueOf(cast).Elem().Set(reflect.ValueOf(dataGen).Elem())
		if cast.UnixMicro() != datas[i].(*time.Time).UnixMicro() {
			t.Fatalf("cast not equal: %v, index: %d, value: %v", datas[i], i, cast)
		}
	}
	t.Logf("cast full reflect: %s", time.Since(t0).String())
}
