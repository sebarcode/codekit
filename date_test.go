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
		{"03-05-2022 21:30:00+0700", "dd-MM-yyyy HH:mm:ssTH", dt0},
		{"03 May 2022 21:30:00+0700", "dd MMM yyyy HH:mm:ssTH", dt0},
		{"03/May/2022 09:30 PM +0700", "dd/MMM/yyyy hh:mm A TH", dt0},
		{"05/03/2022 09:30:00 -0500", "MM/dd/yyyy HH:mm:ss TH", dt0},
	}

	for _, td := range testData {
		res := codekit.String2Date(td.Txt, td.Format)
		if !res.UTC().Equal(td.Result.UTC()) {
			t.Errorf("expected %s got %s with format %s", dt0.UTC().String(), res.String(), codekit.GetFormatDate(td.Txt, td.Format))
		}
	}
}

func TestDateTrim(t *testing.T) {
	dtTest := time.Now()
	dtOnly := codekit.DateOnly(dtTest)
	tmOnly := codekit.TimeOnly(dtTest, true, false)

	if dtOnly.Format("15-04-05") != "00-00-00" {
		t.Fatalf("should return date only, got %s for %v", dtOnly.Format("15-04-05"), dtOnly)
	}

	if tmOnly.Format("2006-01-02") != "1900-01-01" {
		t.Fatalf("should return time only, got %s for %v", tmOnly.Format("2006-01-02"), tmOnly)
	}
}
