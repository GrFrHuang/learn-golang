package test

import "github.com/gin-gonic/gin"

// Test golang frame work gin erupt simultaneously model

const (
	MaxWorker = 5000
	MaxJob    = 100000
)

type Dispatcher struct {
	JobQueue <-chan *gin.Context
	JobNum   int
}

// New a request ctx(context) dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		JobQueue: make(chan *gin.Context, MaxJob),
	}
}

// New a gin handlers collection.
func NewWorkerPool() {

}

// Register api and handler to worker pool.
func RegisterHandler() {

}

func Get() {

}

func Post() {

}

func Option() {

}

func main() {

}
