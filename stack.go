package logx

import (
	"runtime"
	"strconv"
)

const (
	maxStackDepth = 32
)

func getPcs(skip int) []uintptr {
	pcs := make([]uintptr, maxStackDepth)
	n := runtime.Callers(skip, pcs)
	// skip runtime.main()
	// skip runtime.goexit()
	if n -= 2; n < 0 {
		n = 0
	}
	return pcs[:n]
}

func appendStrings(buffer []byte, list ...string) (n int) {
	for _, str := range list {
		n += copy(buffer[n:], str)
		if n >= len(buffer) {
			break
		}
	}
	return
}

func getStack(pcs []uintptr, buffer []byte) (n int) {
	bufLen := len(buffer)
	if bufLen == 0 {
		return
	}
	for _, pc := range pcs {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}

		name := fn.Name()
		n += appendStrings(buffer[n:], name, "()\n")
		if n >= bufLen {
			break
		}

		file, line := fn.FileLine(pc)
		n += appendStrings(buffer[n:], "    ", file, ":")
		if n >= bufLen {
			break
		}

		n += appendStrings(buffer[n:], strconv.Itoa(line), "\n")
		if n >= bufLen {
			break
		}
	}
	return n
}
