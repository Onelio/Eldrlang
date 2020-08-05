package object

import (
	"fmt"
	"io/ioutil"
	"os"
)

type BuiltFun func(args ...Object) Object

var Builtins = map[string]*Builtin{
	"len":    {Size: 1, Fun: builtLen},
	"print":  {Size: -1, Fun: builtPrint},
	"scan":   {Size: 1, Fun: builtScan},
	"fopen":  {Size: 1, Fun: builtFOpen},
	"fclose": {Size: 1, Fun: builtFClose},
	"fread":  {Size: 1, Fun: builtFRead},
	"fwrite": {Size: 2, Fun: builtFWrite},
}

func builtLen(args ...Object) Object {
	switch arg := args[0].(type) {
	case *String:
		return &Integer{Value: int64(len(arg.Value))}
	default:
		return nil
	}
}

func builtPrint(args ...Object) Object {
	for _, arg := range args {
		switch elem := arg.(type) {
		case *String:
			fmt.Print(elem.Value)
		case *Integer:
			fmt.Print(elem.Value)
		case *Boolean:
			fmt.Print(elem.Value)
		}
	}
	return nil
}

func builtScan(args ...Object) Object {
	switch arg := args[0].(type) {
	case *String:
		_, _ = fmt.Scanln(&arg.Value)
	case *Integer:
		_, _ = fmt.Scanln(&arg.Value)
	case *Boolean:
		_, _ = fmt.Scanln(&arg.Value)
	}
	return nil
}

func builtFOpen(args ...Object) Object {
	switch arg := args[0].(type) {
	case *String:
		file, err := os.Open(arg.Value)
		if err != nil {
			return &Integer{Value: 0}
		}
		return &Integer{Value: int64(file.Fd())}
	default:
		return nil
	}
}

func builtFClose(args ...Object) Object {
	switch arg := args[0].(type) {
	case *Integer:
		file := os.NewFile(uintptr(arg.Value), "pipe")
		_ = file.Close()
	}
	return nil
}

func builtFRead(args ...Object) Object {
	switch arg := args[0].(type) {
	case *Integer:
		file := os.NewFile(uintptr(arg.Value), "pipe")
		data, err := ioutil.ReadAll(file)
		if err == nil {
			return &String{Value: string(data)}
		}
	}
	return nil
}

func builtFWrite(args ...Object) Object {
	switch arg := args[0].(type) {
	case *Integer:
		file := os.NewFile(uintptr(arg.Value), "pipe")
		switch arg1 := args[1].(type) {
		case *String:
			_, _ = file.Write([]byte(arg1.Value))
		}
	}
	return nil
}
