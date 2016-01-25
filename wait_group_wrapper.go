package skip_list

import (
    "sync"
)

type WaitGroupWrapper struct {
    sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func(args[] interface{}), args... interface{}) {
    w.Add(1)
    go func(args[] interface{}) {
        cb(args)
        w.Done()
    } (args)
}
