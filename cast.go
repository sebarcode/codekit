package codekit

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	RoundingAuto = "RoundAuto"
	RoundingUp   = "RoundUp"
	RoundingDown = "RoundDown"
)

func Value(o interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(o))
}

func Kind(o interface{}) reflect.Kind {
	return Value(o).Kind()
}

func ToString(o interface{}) string {
	v := Value(o)
	k := v.Kind()

	if k == reflect.Interface && v.IsNil() {
		return ""
	} else if k == reflect.String {
		if val, ok := o.(string); ok {
			return val
		}
		return fmt.Sprintf("%s", o)
	} else if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 {
		return fmt.Sprintf("%d", o)
	} else if k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 {
		return fmt.Sprintf("%d", o)
	} else if k == reflect.Float32 || k == reflect.Float64 {
		return fmt.Sprintf("%f", o)
	} else if k == reflect.Bool {
		return fmt.Sprintf("%t", o)
	} else if k == reflect.Array || k == reflect.Slice || k == reflect.Map || k == reflect.Struct {
		return fmt.Sprintf("%v", o)
	} else {
		return ""
	}
}

/*===============================
== LEGEND ========================
=================================
d 			= date
dd 			= date 2 digit
M 			= month
MM 			= month 2 digit
MMM 		= month in name, 3 chars
MMMM 		= month in name, full
YY 			= Year 2 digit
YYYY 		= Year 4 digit
h 			= hour
hh 			= hour 2 digit
H 			= hour in 24 cycle
HH 			= hour in 24 cycle 2 digit
m 			= minute
mm 			= minute 2 digits
s 			= Second
ss 			= second 2 digit
A 			= AMPM
TZ 			= Timezone using code
TH			= Timezone using +hhmm
=================================*/

func GetFormatDate(o interface{}, dateFormat string) string {

	var dateMap = map[string]string{"dd": "02", "d": "2", "MMMM": "January", "MMM": "Jan", "MM": "01", "M": "1",
		"YYYY": "2006", "YY": "06", "hh": "03", "h": "3", "HH": "15", "mm": "04", "m": "4", "ss": "05", "s": "5",
		"A": "PM", "TZ": "MST", "TH": "-0700",
	}

	var dateOrder = map[int]string{1: "dd", 2: "d", 3: "MMMM", 4: "MMM", 5: "MM", 6: "M", 7: "YYYY", 8: "YY",
		9: "hh", 10: "h", 11: "HH", 12: "mm", 13: "m", 14: "ss", 15: "s", 16: "A", 17: "TZ", 18: "TH",
	}

	var keys []int
	for k := range dateOrder {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	dateFormat = strings.Replace(dateFormat, "y", "Y", -1)
	for _, k := range keys {
		dateFormat = strings.Replace(dateFormat, dateOrder[k], dateMap[dateOrder[k]], -1)
	}

	if strings.Contains(dateFormat, "H") {
		if Value(o).Kind() == reflect.String {

			dateFormat = strings.Replace(dateFormat, "H", "15", -1)
		} else {
			if o.(time.Time).Hour() < 10 {
				dateFormat = strings.Replace(dateFormat, "H", "3", -1)
			} else {
				dateFormat = strings.Replace(dateFormat, "H", "15", -1)
			}
		}
	}

	return dateFormat

}

var _defaultDateFormat string

func SetDefaultDateFormat(f string) {
	_defaultDateFormat = f
}

func DefaultDateFormat() string {
	if _defaultDateFormat == "" {
		_defaultDateFormat = "dd-MMM-yyyy"
	}
	return _defaultDateFormat
}

func Date2String(t time.Time, dateFormat string) string {
	if dateFormat == "" {
		dateFormat = DefaultDateFormat()
	}
	dateFormat = GetFormatDate(t, dateFormat)
	return t.Format(dateFormat)
}

func String2Date(dateString string, dateFormat string) time.Time {
	if dateFormat == "" {
		dateFormat = DefaultDateFormat()
	}
	dateFormat = GetFormatDate(dateString, dateFormat)
	t, e := time.Parse(dateFormat, dateString)
	if e != nil {
		fmt.Println("codekit.Date2String error: " + e.Error())
	}
	return t
}

func ToInt(o interface{}, rounding string) int {
	switch o.(type) {
	case int:
		return o.(int)

	case int8, int16, int32, int64:
		return int(Value(o).Int())

	case uint, uint8, uint16, uint32, uint64:
		return int(Value(o).Uint())

	case float32, float64:
		return int(ToFloat64(o, 0, rounding))

	default:
		return 0
	}
}

func ToFloat32(o interface{}, decimalPoint int, rounding string) float32 {
	var f float64

	switch o.(type) {
	case int, int8, int16, int32, int64:
		f = ToFloat64(Value(o).Int(), decimalPoint, rounding)

	case uint, uint8, uint16, uint32, uint64:
		f = ToFloat64(Value(o).Uint(), decimalPoint, rounding)

	case float32, float64:
		var ok bool
		f, ok = o.(float64)
		if !ok {
			f = float64(o.(float32))
		}
		f = ToFloat64(f, decimalPoint, rounding)

	default:
		f = ToFloat64(Value(o).String(), decimalPoint, rounding)
	}

	if math.IsNaN(f) || math.IsInf(f, 0) {
		f = 0
	}

	return float32(f)
}

func ToFloat64(o interface{}, decimalPoint int, rounding string) float64 {
	var f float64
	var e error

	f, ok := o.(float64)
	if !ok {
		v := Value(o)
		t := strings.ToLower(v.Type().Name())

		if t != "interface{}" && strings.HasPrefix(t, "int") {
			f = float64(v.Int())
		} else if strings.HasPrefix(t, "uint") {
			f = float64(v.Uint())
		} else if strings.HasPrefix(t, "float") {
			f = float64(v.Float())
		} else {
			f, e = strconv.ParseFloat(v.String(), 64)
			if e != nil {
				return 0
			}
		}
	}

	if decimalPoint < 0 {
		return f
	}

	switch rounding {
	case RoundingAuto:
		return RoundingAuto64(f, decimalPoint)
	case RoundingDown:
		return RoundingDown64(f, decimalPoint)
	case RoundingUp:
		return RoundingUp64(f, decimalPoint)
	}

	if math.IsNaN(f) || math.IsInf(f, 0) {
		f = 0
	}

	return f
}

func RoundingAuto64(f float64, decimalPoint int) (retValue float64) {

	tempPow := math.Pow(10, float64(decimalPoint))
	f = f * tempPow

	if f < 0 {
		f = math.Ceil(f - 0.5)
	} else {
		f = math.Floor(f + 0.5)
	}

	retValue = f / tempPow
	return
}

func RoundingDown64(f float64, decimalPoint int) (retValue float64) {
	tempPow := math.Pow(10, float64(decimalPoint))
	f = f * tempPow
	f = math.Floor(f)
	retValue = f / tempPow
	return
}

func RoundingUp64(f float64, decimalPoint int) (retValue float64) {
	tempPow := math.Pow(10, float64(decimalPoint))
	f = f * tempPow
	f = math.Ceil(f)
	retValue = f / tempPow
	return
}

func ToDate(o interface{}, formatDate string) time.Time {
	v := reflect.Indirect(reflect.ValueOf(o))
	t := strings.ToLower(v.Type().String())
	if strings.Contains(t, "int") {
		intDate := v.Int()
		return time.Unix(intDate, 0)
	} else if strings.Contains(t, "string") {
		return String2Date(o.(string), formatDate)
	} else if t == "*time.Time" {
		dt, ok := o.(*time.Time)
		if !ok {
			return time.Now()
		}
		return *dt
	} else if strings.HasSuffix(t, "time.time") {
		return o.(time.Time)
	}
	return time.Now()
}

func ToDuration(o interface{}) time.Duration {
	return time.Duration(ToInt(o, RoundingAuto)) * time.Second
}
