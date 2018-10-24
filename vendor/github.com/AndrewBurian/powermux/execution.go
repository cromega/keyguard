package powermux

import (
	"net/http"
	"sync"
)

// routeExecution is the complete instructions for running serve on a route
type routeExecution struct {
	pattern    string
	params     map[string]string
	notFound   http.Handler
	middleware []Middleware
	handler    http.Handler
}

func newExecution() *routeExecution {
	return &routeExecution{
		middleware: make([]Middleware, 0),
		params:     make(map[string]string),
	}
}

func (ex *routeExecution) Reset() {
	ex.middleware = ex.middleware[0:0]
	for key := range ex.params {
		delete(ex.params, key)
	}
	ex.handler = nil
	ex.notFound = nil
}

type executionPool struct {
	p *sync.Pool
}

func (ep *executionPool) Get() *routeExecution {
	return ep.p.Get().(*routeExecution)
}

func (ep *executionPool) Put(ex *routeExecution) {
	ex.Reset()
	ep.p.Put(ex)
}

func createExecution() interface{} {
	return newExecution()
}

func newExecutionPool() *executionPool {
	return &executionPool{
		p: &sync.Pool{
			New: createExecution,
		},
	}
}
