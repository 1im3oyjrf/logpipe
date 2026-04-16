package enricher

import (
	"os"
	"sync"
)

var (
	hostnameCache string
	import_os_once sync.Once
)

func hostnameOnce() string {
	import_os_once.Do(func() {
		h, err := os.Hostname()
		if err == nil {
			hostnameCache = h
		}
	})
	return hostnameCache
}
