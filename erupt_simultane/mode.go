package erupt_simultane

import (
	"github.com/gin-gonic/gin"
	"github.com/GrFrHuang/gox/log"
	"strings"
	"net/http"
	"sync"
	"time"
	"strconv"
)

// Test golang frame work gin erupt simultaneously model

const (
	MaxProcess = 2     // Limit maximum cpu process numbers.
	MaxWorker  = 1000  // Limit maximum work pool worker numbers.
	MaxJob     = 10000 // Limit maximum request job numbers.
)

type Worker struct {
	rw       *sync.RWMutex
	quit     chan bool
	handlers map[string]func(ctx *gin.Context)
}

type Dispatcher struct {
	jobQueue chan *gin.Context
}

var HandlerCluster map[string]map[string]func(ctx *gin.Context)
var Dpt *Dispatcher
var WorkerPool *Worker

func init() {
	WorkerPool = NewWorkerPool()
	Dpt = NewDispatcher()
	HandlerCluster = make(map[string]map[string]func(ctx *gin.Context))

	WorkerPool.RegisterWorker(Dpt, "GET", "/getsome", getsome)
	WorkerPool.RegisterWorker(Dpt, "POST", "/postsome", Postsome)
	WorkerPool.RegisterWorker(Dpt, "DELETE", "/deletesome", deletesome)
}

// New a request ctx(context) dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		jobQueue: make(chan *gin.Context, MaxJob),
	}
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
		rw:       new(sync.RWMutex),
		quit:     make(chan bool),
		handlers: make(map[string]func(ctx *gin.Context)),
	}
}

// Distribute gin requests to the corresponding deal handler.
func (w *Worker) Work(ctx *gin.Context) {
	uri := strings.Split(ctx.Request.RequestURI, "?")[0]
	//w.rw.RLock()
	//defer w.rw.RUnlock()
	if _, ok := w.handlers[uri]; ok {
		w.handlers[uri](ctx)
		return
	}
	log.Error("api forget register ?")
}

// Register api and handler to worker pool.
func (w *Worker) RegisterWorker(dispatcher *Dispatcher, method, api string, handler func(ctx *gin.Context)) {
	// Register api dispatcher func.
	bean := make(map[string]func(ctx *gin.Context))
	bean[api] = handler
	HandlerCluster[method] = bean

	w.handlers[api] = handler
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func (d *Dispatcher) Run(w *Worker) {
	for i := 0; i < MaxWorker; i++ {
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

func getsome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "geter say get !")
}
func Postsome(ctx *gin.Context) {
	time.Sleep(time.Millisecond * 200)
	role := &Role{
		CpRoleId:   strconv.Itoa(time.Now().Nanosecond()),
		UserId:     1,
		GameId:     1,
		RoleName:   "战士",
		RoleGrade:  "79",
		GameRegion: "五行山",
	}
	Ol.Beans <- role
	ctx.JSON(http.StatusOK, "poster say post !")

}
func deletesome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "optioner say options !")
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
	case "OPTIONS":
		router.DELETE(api, handler)
	}
}
