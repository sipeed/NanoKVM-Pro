package stream

import (
	"bytes"
	"sync"
)

var BufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}
