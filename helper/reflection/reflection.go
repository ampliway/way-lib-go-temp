package reflection

import (
	"reflect"
	"runtime"
	"strings"
)

func TypeName(value interface{}) string {
	if value == nil {
		return ""
	}

	t := reflect.TypeOf(value)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

const appNameMax = 3

func AppName(msg interface{}) string {
	typeOfEvent := reflect.TypeOf(msg)
	if typeOfEvent.Kind() == reflect.Ptr {
		typeOfEvent = typeOfEvent.Elem()
	}

	pkgPath := typeOfEvent.PkgPath()

	pkgPathParts := strings.Split(pkgPath, "/")
	if len(pkgPathParts) >= appNameMax {
		pkgPath = pkgPathParts[appNameMax-1]
	}

	return pkgPath
}

func AppNamePkg() string {
	var pkgName string

	for i := 0; i <= 10; i++ {
		pc, _, _, _ := runtime.Caller(i)
		funcName := runtime.FuncForPC(pc).Name()

		lastSlash := strings.LastIndexByte(funcName, '/')
		if lastSlash < 0 {
			lastSlash = 0
		}

		lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

		pkgName = funcName[:lastDot]
		pkgNameParts := strings.Split(pkgName, "/")

		if len(pkgNameParts) >= appNameMax {
			pkgName = pkgNameParts[appNameMax-1]
		}

		if pkgName != "way-lib-go" {
			break
		}
	}

	return pkgName
}
