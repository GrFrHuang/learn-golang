package erupt_simultane

import (
	"github.com/gin-gonic/gin"
	"github.com/GrFrHuang/gox/log"
	"strings"
	"net/http"
	"sync"
	"time"
)

// Test golang frame work gin erupt simultaneously model

const (
	MaxWorker = 1000
	MaxJob    = 10000
)

var wg = &sync.WaitGroup{}

var Dispatchers map[string]map[string]func(ctx *gin.Context)

type Worker struct {
	quit     chan bool
	handlers map[string]func(ctx *gin.Context)
}

type Dispatcher struct {
	jobQueue chan *gin.Context
}

// New a request ctx(context) dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		jobQueue: make(chan *gin.Context, MaxJob),
	}
}

// Register api and handler to worker pool.
func (d *Dispatcher) RegisterMethodDispatcher(method, api string, handler func(ctx *gin.Context)) {
	bean := make(map[string]func(ctx *gin.Context))
	bean[api] = handler
	Dispatchers[method] = bean
}

func (d *Dispatcher) PushContextQueue(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Job queue has been closed : ", err)
		}
	}()
	c := ctx.Copy()
	d.jobQueue <- c
}

func (d *Dispatcher) GetJonNum() int {
	return len(d.jobQueue)
}

// New a gin handlers collection.
func NewWorkerPool() *Worker {
	return &Worker{
		quit:     make(chan bool),
		handlers: make(map[string]func(ctx *gin.Context)),
	}
}

// Distribute gin requests to the corresponding deal handler.
func (w *Worker) Work(ctx *gin.Context) {
	uri := strings.Split(ctx.Request.RequestURI, "?")[0]
	m := &sync.Mutex{}
	m.Lock()
	if _, ok := w.handlers[uri]; ok {
		w.handlers[uri](ctx)
	}
	m.Unlock()
	log.Error("api forget register ?")
}

// Register api and handler to worker pool.
func (w *Worker) RegisterWorker(api string, handler func(ctx *gin.Context)) {
	w.handlers[api] = handler
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func (d *Dispatcher) Run(w *Worker) {
	for i := 0; i < 20000; i++ {
		go func(queue chan *gin.Context, worker *Worker) {
			for {
				select {
				case ctx := <-queue:
					// we have received a work request.
					if ctx != nil {
						worker.Work(ctx)
					}
				case <-worker.quit:
					// we have received a signal to stop
					return
				}
			}
		}(d.jobQueue, w)
	}
}

func Get(ctx *gin.Context) {
	Dispa.PushContextQueue(ctx)
}

func Post(ctx *gin.Context) {
	Dispa.PushContextQueue(ctx)
}

func Option(ctx *gin.Context) {
	Dispa.PushContextQueue(ctx)
}

func getsome(ctx *gin.Context) {
	time.Sleep(time.Millisecond * 500)
	ctx.JSON(http.StatusOK, "get say get !")
	log.Info("get say get !")
}
func postsome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "post say post !")
}
func deletesome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "options say options !")
}

// Distinct listener api real handler.
func ListenApi(method, api string, handler func(ctx *gin.Context), router *gin.Engine) {
	switch method {
	case "GET":
		router.GET(api, handler)
	case "PUT":
		router.PUT(api, handler)
	case "POST":
		router.POST(api, handler)
	case "DELETE":
		router.DELETE(api, handler)
	}
}

var Dispa *Dispatcher
var WorkerPool *Worker

func init() {
	WorkerPool = NewWorkerPool()
	Dispa = NewDispatcher()
	Dispatchers = make(map[string]map[string]func(ctx *gin.Context))

	Dispa.RegisterMethodDispatcher("GET", "/getsome", Get)
	Dispa.RegisterMethodDispatcher("POST", "/postsome", Post)
	Dispa.RegisterMethodDispatcher("DELETE", "/deletesome", deletesome)

	WorkerPool.RegisterWorker("/getsome", getsome)
	WorkerPool.RegisterWorker("/postsome", postsome)
	WorkerPool.RegisterWorker("/deletesome", deletesome)
}
