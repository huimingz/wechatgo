package msg

import (
	"bytes"
	"strconv"
)

func intSliceToString(v []int, step string) string {
	var buffer bytes.Buffer
	var length = len(v)
	for i := 0; i < length; i++ {
		buffer.WriteString(strconv.Itoa(v[i]))
		if i != length-1 {
			buffer.WriteString(step)
		}
	}

	return buffer.String()
}
