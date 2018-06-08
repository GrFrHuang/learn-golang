package main

import (
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"learn-golang/erupt_simultane"
	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(erupt_simultane.MaxProcess)
	router := gin.Default()
	for k, v := range erupt_simultane.HandlerCluster {
		for i, j := range v {
			erupt_simultane.ListenApi(k, i, j, router)
		}
	}
	//erupt_simultane.Ol.ListenOrm()
	erupt_simultane.Dpt.Run(erupt_simultane.WorkerPool)
	router.Run("127.0.0.1:8098")
}
