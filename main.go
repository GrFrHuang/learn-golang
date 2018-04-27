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
	erupt_simultane.Dpt.Run(erupt_simultane.WorkerPool)
	router.Run("127.0.0.1:8098")

	//ExecSql(channels)
	//router.POST("postsome", erupt_simultane.Postsome)
	//router.Run("127.0.0.1:8098")
}

//var channels = make(chan *Role, 10000)
//
//func ExecSql(cls chan *Role) {
//	go func(cs chan *Role) {
//		for {
//			if len(cs) >= 2 {
//				for j := 0; j < len(cs); j++ {
//					_, err := Engine.InsertOne(<-cs)
//					if err != nil {
//						log.Error(err)
//					}
//				}
//			}
//			time.Sleep(time.Millisecond * 10)
//		}
//	}(cls)
//}
//
//func handler(w http.ResponseWriter, r *http.Request) {
//
//	i ++
//	role := &Role{
//		CpRoleId:   strconv.Itoa(i),
//		UserId:     1,
//		GameId:     1,
//		RoleName:   "战士",
//		RoleGrade:  "79",
//		GameRegion: "五行山",
//	}
//	channels <- role
//	log.Info("run here")
//	fmt.Println(" length is : ", len(channels))
//	fmt.Fprintf(w, "Hello, GrFrHuang !")
//}
//
//var i int = 0
//var Engine *xorm.Engine
//
//func init() {
//	UrlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
//		"root",
//		"baiying",
//		"127.0.0.1",
//		"3306",
//		"sdk")
//	log.Info(UrlStr)
//	engine, err := xorm.NewEngine("mysql", UrlStr)
//	if err != nil {
//		log.Panic("Initialize the xorm err : ", err)
//	}
//	Engine = engine
//	Engine.SetMaxOpenConns(2000)
//	Engine.SetMaxIdleConns(1000)
//}
//
//type Role struct {
//	Id         int    `xorm:"not null pk autoincr INT(10)" json:"id"`
//	CpRoleId   string `xorm:"comment('cp的角色id') VARCHAR(255)" json:"cp_role_id"`
//	GameId     int    `xorm:"comment('关联的game_id') INT(10)" json:"game_id"`
//	UserId     int    `xorm:"comment('关联的用户id') INT(10)" json:"user_id"`
//	RoleName   string `xorm:"default '' comment('游戏角色名') VARCHAR(255)" json:"role_name"`
//	RoleGrade  string `xorm:"default '' comment('角色等级') VARCHAR(255)" json:"role_grade"`
//	GameRegion string `xorm:"default '' comment('游戏区服') VARCHAR(255)" json:"game_region"`
//	CreateTime int    `xorm:"INT(11) created" json:"create_time"`
//	UpdateTime int    `xorm:"INT(11) updated" json:"update_time"`
//}
