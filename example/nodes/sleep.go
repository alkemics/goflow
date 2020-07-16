package nodes

import (
	"time"
)

func Sleeper(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
