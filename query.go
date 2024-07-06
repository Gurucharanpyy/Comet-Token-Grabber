package binder

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const defaultSeparator = ","

// FromQuery...
func FromQuery(dst interface{}, src url.Values) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("invalid interface type")
	}

	rValue := reflect.ValueOf(dst).Elem()
	rType := rValue.Type()

	for i := 0; i < rValue.NumField(); i++ {
		tag, ok := rType.Field(i).Tag.Lookup("json")
		if !ok || tag == "" || tag == "-" {
			continue
		}

		queryParamValue := src.Get(tag)
		if queryParamValue == "" {
			continue
		}

		fld := rValue.Field(i)
		if !fld.CanSet() {
			return fmt.Errorf("field %s is not addressable", rType.Field(i).Name)
		}

		switch fld.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			p, err := strconv.Atoi(queryParamValue)
			if err != nil {
				return err
			}
			fld.SetInt(int64(p))
		case reflect.String:
			fld.SetString(queryParamValue)
		case reflect.Bool:
			if strings.ToLower(queryParamValue) == "true" {
				fld.SetBool(true)
			}
		case reflect.Slice:
			var values []string
			for _, v := range strings.Split(queryParamValue, defaultSeparator) {
				if v == "" {
					continue
				}
				values = append(values, v)
			}

			if len(values) == 0 {
				continue
			}

			dst := reflect.MakeSlice(fld.Type(), len(values), len(values))
			switch fld.Interface().(type) {
			case []int:
				var vs []int
				for _, s := range values {
					n, err := strconv.Atoi(s)
					if err != nil {
						return fmt.Errorf("can't parse value: %s", s)
					}
					vs = append(vs, n)
				}

				reflect.Copy(dst, reflect.ValueOf(vs))
			case []string:
				reflect.Copy(dst, reflect.ValueOf(values))
			default:
				return fmt.Errorf("unsupported slice type")
			}
			fld.Set(dst)

		default:
			return fmt.Errorf("%s unsupported type", fld.Type())
		}
	}
	return nil
}
