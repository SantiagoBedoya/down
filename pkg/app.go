package down

import "sync"

type App struct {
	Concurrency int
	URI         string
	Chunks      map[int][]byte
	Err         error
	Destination string
	*sync.Mutex
}
