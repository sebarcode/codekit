package codekit

import (
	"fmt"
	"strings"
)

func Printfn(pattern string, parm ...interface{}) string {
	if strings.HasSuffix(pattern, "\n") {
		return fmt.Sprintf(pattern, parm...)
	}
	return fmt.Sprintf(pattern+"\n", parm...)
}
