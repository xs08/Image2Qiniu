package utils

import "bytes"

// JoinStrs join strs
func JoinStrs(strs ...string) string {
	var buf bytes.Buffer
	for _, str := range strs {
		buf.WriteString(str)
	}
	return buf.String()
}
