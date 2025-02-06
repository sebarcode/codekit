package codekit_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/sebarcode/codekit"
)

type TestModel struct {
	Name          string
	NamePtr       *string
	Created       time.Time
	ActivatedDate *time.Time
}

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

func TestCastModel(t *testing.T) {
	dataCount := 100
	datas := make([]*TestModel, dataCount)
	for i := 0; i < dataCount; i++ {
		data := &TestModel{
			Name:          codekit.RandomString(10),
			Created:       time.Now().Add(-time.Duration(codekit.RandInt(1000)) * time.Second),
			ActivatedDate: nil,
		}
		datas[i] = data
	}

	t0 := time.Now()
	namePtr := "Hahaha"
	datas[10].ActivatedDate = &t0
	datas[10].NamePtr = &namePtr

	mData, err := codekit.ToM(datas[0])
	if err != nil {
		t.Fatalf("tom failed: %v", err)
	}

	if mData.Get("Created", time.Now()).(time.Time).UnixMicro() != datas[0].Created.UnixMicro() {
		t.Fatalf("get created failed: %v", mData.Get("Created", time.Now()).(time.Time))
	}
	if nilName, _ := mData.Get("NamePtr", nil).(*string); nilName != nil {
		t.Fatalf("get name ptr nil failed: %v", nilName)
	}

	mData, err = codekit.ToM(datas[10])
	if err != nil {
		t.Fatalf("tom 10 failed: %v", err)
	}
	t10 := datas[10].ActivatedDate.UnixMicro()
	if t10 != mData.Get("ActivatedDate", &t0).(*time.Time).UnixMicro() {
		t.Fatalf("get activated date failed: %v", mData.Get("ActivatedDate", &t0).(*time.Time))
	}
	if *(mData.Get("NamePtr", &namePtr).(*string)) != namePtr {
		t.Fatalf("get name ptr failed: %v", mData.Get("NamePtr", &namePtr).(*string))
	}
}
