// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package xflags

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/tommika/gorilla/util"
)

const (
	flagTag = "flag"
)

type xFlagSet struct {
	flagSet     *flag.FlagSet
	remArgs     *reflect.Value
	remArgsHelp string
}

// ParseArgs parses the given arguments based on the given value
func ParseArgs(name string, a any, args []string) error {
	xfs, err := NewFlagSet(name, a)
	if err != nil {
		return err
	}
	if err = xfs.Parse(args); err != nil {
		return err
	}
	return err
}

func Usage(out io.Writer, name string, a any) error {
	xfs, err := NewFlagSet(name, a)
	if err != nil {
		return err
	}
	xfs.Usage(out)
	return nil
}

func (xfs *xFlagSet) Usage(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s [options] %s\n", xfs.flagSet.Name(), xfs.remArgsHelp)
	xfs.flagSet.SetOutput(out)
	xfs.flagSet.PrintDefaults()
}

func (xfs *xFlagSet) Parse(args []string) error {
	if err := xfs.flagSet.Parse(args); err != nil {
		return err
	}
	if xfs.flagSet.NArg() > 0 && xfs.remArgs != nil {
		xfs.remArgs.Set(reflect.ValueOf(xfs.flagSet.Args()))
	}
	return nil
}

// NewFlagSet creates a FlagSet from the given value.
// Must be a pointer to a structure.
func NewFlagSet(name string, a any) (*xFlagSet, error) {
	xfs := xFlagSet{
		flagSet: flag.NewFlagSet(name, flag.ContinueOnError),
	}
	xfs.flagSet.Usage = func() {
		xfs.Usage(os.Stderr)
	}
	if err := xfs.InitFlagSet(a); err != nil {
		return nil, err
	}
	return &xfs, nil
}

func (xfs *xFlagSet) InitFlagSet(a any) error {
	v := reflect.ValueOf(a)
	if !(v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Struct) {
		return fmt.Errorf("parameter must be a pointer to a struct")
	}
	// deference the pointer
	v = v.Elem()
	// get the actual type
	t := v.Type()
	var err error
	for i := 0; err == nil && i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		var nameList []string
		var usage, def string
		if tag, found := field.Tag.Lookup(flagTag); !found {
			if field.IsExported() {
				nameList = []string{defFlagName(field)}
			}
		} else {
			// flag:"name,usage,default"
			tagSplit := util.SplitAndTrim(tag, ",")
			if l := len(tagSplit); l > 0 {
				names := tagSplit[0]
				if names == "-" {
					// ignore this field
					continue
				}
				if names == "*" {
					// field that will hold the remaining args
					if !isSliceOfKind(fieldValue, reflect.String) {
						// has to be a slice of string
						return fmt.Errorf("inappropriate type for remaining args field: %s", field.Name)
					}
					xfs.remArgs = &fieldValue
					if len(tagSplit) > 1 {
						xfs.remArgsHelp = tagSplit[1]
					}
					continue
				}
				if len(names) == 0 {
					// use default name
					nameList = []string{defFlagName(field)}
				} else {
					nameList = util.SplitAndTrim(names, "|")
				}
				if l > 1 {
					usage = tagSplit[1]
					if l > 2 {
						def = tagSplit[2]
					}
				}
			}
		}
		fs := xfs.flagSet
		if nameList != nil {
			if len(def) == 0 {
				def = fmt.Sprintf("%v", fieldValue)
			}
			switch fieldValue.Kind() {
			default:
				err = fmt.Errorf("flag type is not supported: %v", field)
			case reflect.String:
				var ptr *string
				reflect.ValueOf(&ptr).Elem().Set(fieldValue.Addr())
				for _, name := range nameList {
					fs.StringVar(ptr, name, def, usage)
				}
			case reflect.Bool:
				defBool := util.IsTrue(def)
				var ptr *bool
				reflect.ValueOf(&ptr).Elem().Set(fieldValue.Addr())
				for _, name := range nameList {
					fs.BoolVar(ptr, name, defBool, usage)
				}
			case reflect.Int:
				defInt, _ := strconv.Atoi(def)
				var ptr *int
				reflect.ValueOf(&ptr).Elem().Set(fieldValue.Addr())
				for _, name := range nameList {
					fs.IntVar(ptr, name, defInt, usage)
				}
			case reflect.Uint:
				n, _ := strconv.Atoi(def)
				defUint := uint(n)
				var ptr *uint
				reflect.ValueOf(&ptr).Elem().Set(fieldValue.Addr())
				for _, name := range nameList {
					fs.UintVar(ptr, name, defUint, usage)
				}
			case reflect.Int64:
				defInt64, _ := strconv.ParseInt(def, 10, 64)
				var ptr *int64
				reflect.ValueOf(&ptr).Elem().Set(fieldValue.Addr())
				for _, name := range nameList {
					fs.Int64Var(ptr, name, defInt64, usage)
				}
			case reflect.Uint64:
				defUint64, _ := strconv.ParseUint(def, 10, 64)
				var ptr *uint64
				reflect.ValueOf(&ptr).Elem().Set(fieldValue.Addr())
				for _, name := range nameList {
					fs.Uint64Var(ptr, name, defUint64, usage)
				}
			}
		}
	}
	return err
}

func defFlagName(field reflect.StructField) string {
	return strings.ToLower(field.Name)
}

func isSliceOfKind(val reflect.Value, kind reflect.Kind) bool {
	return val.Kind() == reflect.Slice && val.Type().Elem().Kind() == kind
}
