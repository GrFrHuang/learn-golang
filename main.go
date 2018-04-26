package main

import (
	"github.com/gin-gonic/gin"
	"learn-golang/erupt_simultane"
)

func main() {
	router := gin.Default()
	for k, v := range erupt_simultane.Dispatchers {
		for i, j := range v {
			erupt_simultane.ListenApi(k, i, j, router)
		}
	}
	erupt_simultane.Dispa.Run(erupt_simultane.WorkerPool)
	//router.GET("getsome", func(ctx *gin.Context) {
	//	time.Sleep(time.Millisecond * 500)
	//	ctx.JSON(http.StatusOK, "get say get !")
	//	log.Info("get say get !")
	//})
	router.Run("127.0.0.1:8098")
}
