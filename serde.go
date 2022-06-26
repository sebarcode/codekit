package codekit

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

func Serde(source interface{}, dest interface{}, serdeType string) error {
	vt1 := reflect.TypeOf(source)
	vt2 := reflect.TypeOf(dest)
	vt1kind := vt1.Kind()
	if vt1kind == reflect.Ptr {
		vt1kind = reflect.ValueOf(source).Elem().Type().Kind()
	}

	vt1name := vt1.String()
	vt2name := vt2.String()
	if vt1name == vt2name && vt1kind != reflect.Map && vt1kind != reflect.Slice {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(source).Elem())
	}

	if serdeType == "" || serdeType == "reflect" {
		return SerdeUsingReflect(source, dest)
	}

	bs, e := ToBytesWithError(source, serdeType)
	if len(bs) == 0 {
		return errors.New("toolkit.Serde: Serialization Fail " + e.Error())
	}

	e = FromBytes(bs, serdeType, dest)
	if e != nil {
		return errors.New("toolkit.Serde: Deserialization fail " + e.Error())
	}

	return nil
}

func SerdeUsingReflect(source, dest interface{}) error {
	vSource := reflect.Indirect(reflect.ValueOf(source))
	isSourceSlice := vSource.Kind() == reflect.Slice

	vDest := reflect.ValueOf(dest)
	vDestKind := vDest.Kind()
	if vDestKind != reflect.Ptr {
		return errors.New("destination should be a pointer")
	}
	isDestSlice := vDestKind == reflect.Slice
	if isSourceSlice && !isDestSlice {
		return errors.New("destination should be a slice")
	}

	if isSourceSlice {
		return SerdeSlice(vSource, vDest)
	}

	return CopyValue(vSource, vDest.Elem())
}

func SerdeSlice(source, dest reflect.Value) error {
	if dest.Kind() == reflect.Ptr {
		return SerdeSlice(source, dest.Elem())
	}

	elemType := dest.Type().Elem()
	elemIsPtr := elemType.Kind() == reflect.Ptr
	sourceLen := source.Len()
	destBuffer := reflect.MakeSlice(elemType, sourceLen, sourceLen)
	for i := 0; i <= sourceLen; i++ {
		sourceItem := source.Index(i)
		destItem := CreatePtrFromType(elemType)
		if e := CopyValue(sourceItem, destItem.Elem()); e != nil {
			return fmt.Errorf("errors processing index %d, %s", i, e.Error())
		}
		if elemIsPtr {
			destBuffer.Index(i).Set(destItem)
		} else {
			destBuffer.Index(i).Set(destItem.Elem())
		}
	}
	dest.Set(destBuffer)
	return nil
}

func CopyValue(source, dest reflect.Value) error {
	if source.Kind() == reflect.Ptr {
		return CopyValue(source.Elem(), dest)
	}

	if dest.Kind() == reflect.Ptr {
		return CopyValue(source, dest.Elem())
	}

	sourceIsMap := source.Kind() == reflect.Map
	destIsMap := dest.Kind() == reflect.Map

	if destIsMap {
		return copyValueToMap(source, dest, sourceIsMap, true)
	}
	return copyValueToStruct(source, dest, sourceIsMap, true)
}

func copyValueToStruct(source, dest reflect.Value, sourceIsMap, ignoreError bool) error {
	destType := dest.Type()
	fieldCount := destType.NumField()
	for i := 0; i < fieldCount; i++ {
		fieldMeta := destType.Field(i)
		sourceField := getFieldFromRV(source, fieldMeta.Name, sourceIsMap)
		if !sourceField.IsZero() {
			var eSet error
			func() {
				defer RecoverToError(&eSet)
				sourceFieldTypeName := sourceField.Type().Name()
				destField := dest.FieldByName(fieldMeta.Name)
				destFieldTypeName := destField.Type().String()

				// check if source and dest has same type, if not then do cast
				if sourceFieldTypeName == destFieldTypeName {
					destField.Set(sourceField)
				} else {
					switch destFieldTypeName {

					case "time.Time":
						switch sourceFieldTypeName {
						case "string":
							if dt, err := time.Parse(time.RFC3339, destField.String()); err == nil {
								destField.Set(reflect.ValueOf(dt))
							} else {
								dt := ToDate(destField.String(), DefaultDateFormat())
								destField.Set(reflect.ValueOf(dt))
							}

						case "int", "int16", "int32", "int64", "uint", "uint16", "uint32", "uint64":
							dt := time.Unix(int64(destField.Int()), 0)
							destField.Set(reflect.ValueOf(dt))
						}

					default:
						/*
								destFieldKind := destField.Kind()
							switch destFieldKind {
											case reflect.Ptr:
												newPtr := CreatePtrFromType(destField.Type())
												newPtr.Elem().Set(sourceField)

											case reflect.Slice:
											default:
												destField.Set(sourceField)
											}
						*/
						destField.Set(sourceField)
					}
				}
			}()
			if eSet != nil {
				if !ignoreError {
					return fmt.Errorf("fail processing %s. %ss", fieldMeta.Name, eSet.Error())
				} else {
					fmt.Printf("fail processing %s. %ss", fieldMeta.Name, eSet.Error())
				}
			}
		}
	}
	return nil
}

func copyValueToMap(source, dest reflect.Value, sourceIsMap, ignoreError bool) error {
	keys := []reflect.Value{}
	destType := dest.Type()
	fieldCount := 0
	if sourceIsMap {
		keys = dest.MapKeys()
		fieldCount = len(keys)
	} else {
		fieldCount = dest.NumField()
	}

	for i := 0; i < fieldCount; i++ {
		var (
			sourceField reflect.Value
			eSet        error
			fieldName   = ""
		)

		func() {
			defer RecoverToError(&eSet)
			if sourceIsMap {
				key := keys[i]
				fieldName = fmt.Sprintf("%v", key.Interface())
				sourceField = source.MapIndex(key)
				if !sourceField.IsZero() {
					dest.SetMapIndex(key, sourceField)
				}
			} else {
				fieldName = destType.Field(i).Name
				sourceField = source.Field(i)
				if !sourceField.IsZero() {
					dest.Set(sourceField)
				}
			}
		}()

		if eSet != nil {
			if !ignoreError {
				return fmt.Errorf("fail processing %s. %ss", fieldName, eSet.Error())
			} else {
				fmt.Printf("fail processing %s. %ss", fieldName, eSet.Error())
			}
		}
	}

	return nil
}

func CreatePtrFromType(t reflect.Type) reflect.Value {
	isPtr := t.Kind() == reflect.Ptr
	elemType := t
	if isPtr {
		elemType = elemType.Elem()
	}

	return reflect.New(elemType)
}

func getFieldFromRV(rv reflect.Value, name string, isMap bool) reflect.Value {
	if isMap {
		return rv.MapIndex(reflect.ValueOf(name))
	}
	return rv.FieldByName(name)
}

func RecoverToError(e *error) string {
	if r := recover(); r != nil {
		switch r.(type) {
		case string:
			*e = errors.New(r.(string))
		default:
			*e = errors.New(fmt.Sprintf("%v", r))
		}
		txts := strings.Split(string(debug.Stack()), "\n")
		codeTraces := []string{}
		for _, txt := range txts {
			if strings.Contains(txt, ".go") {
				codeTraces = append(codeTraces, txt)
			}
		}
		return strings.Join(codeTraces, "\n")
	}
	return ""
}
