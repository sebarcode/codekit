package codekit_test

import (
	"testing"
	"time"

	"github.com/sebarcode/codekit"
)

func TestString2Date(t *testing.T) {
	locJkt, _ := time.LoadLocation("Asia/Jakarta")
	dt0 := time.Date(2022, 05, 03, 21, 30, 0, 0, locJkt)
	testData := []struct {
		Txt    string
		Format string
		Result time.Time
	}{
		{"03-05-2022 21:30:00+07:00", "dd-MM-yyyy HH:mm:ssTZ:00", dt0},
		{"03 May 2022 21:30:00+07:00", "dd MMM yyyy HH:mm:ssTZ:00", dt0},
		{"03/May/2022 09:30 PM +07:00", "dd/MMM/yyyy HH:mm A TZ:00", dt0},
	}

	for _, td := range testData {
		res := codekit.String2Date(td.Txt, td.Format)
		if !res.Equal(td.Result) {
			t.Errorf("expected %s got %s", dt0.String(), res.String())
		}
	}
}
