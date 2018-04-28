package erupt_simultane

import (
	"time"
	"fmt"
	"github.com/GrFrHuang/gox/log"
	"github.com/go-xorm/xorm"
	"reflect"
)

type Role struct {
	Id         int    `xorm:"not null pk autoincr INT(10)" json:"id"`
	CpRoleId   string `xorm:"comment('cp的角色id') VARCHAR(255)" json:"cp_role_id"`
	GameId     int    `xorm:"comment('关联的game_id') INT(10)" json:"game_id"`
	UserId     int    `xorm:"comment('关联的用户id') INT(10)" json:"user_id"`
	RoleName   string `xorm:"default '' comment('游戏角色名') VARCHAR(255)" json:"role_name"`
	RoleGrade  string `xorm:"default '' comment('角色等级') VARCHAR(255)" json:"role_grade"`
	GameRegion string `xorm:"default '' comment('游戏区服') VARCHAR(255)" json:"game_region"`
	CreateTime int    `xorm:"INT(11) created" json:"create_time"`
	UpdateTime int    `xorm:"INT(11) updated" json:"update_time"`
}

const (
	MaxSql = 1000
)

type OrmListener struct {
	Beans chan interface{}
}

var Count int
var Engine *xorm.Engine
var Ol *OrmListener

func init() {
	Ol = NewOrmListener()
	UrlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		"root",
		"baiying",
		"127.0.0.1",
		"3306",
		"sdk")
	Engine, err := xorm.NewEngine("mysql", UrlStr)
	if err != nil {
		log.Panic("Initialize the xorm err : ", err)
	}
	Engine.SetMaxOpenConns(1000)
	Engine.SetMaxIdleConns(500)
}

func NewOrmListener() *OrmListener {
	return &OrmListener{
		Beans: make(chan interface{}, MaxSql),
	}
}

func (ol *OrmListener) ListenOrm() {
	go func(beans chan interface{}) {
		for {
			if len(beans) >= 2 {
				fmt.Println("=======", reflect.TypeOf(<-beans))
				for i := 0; i < len(beans); i++ {
					_, err := Engine.Table("role").Insert(<-beans)
					if err != nil {
						log.Error(err)
					}
				}
			}
			time.Sleep(time.Millisecond * 10)
		}
	}(ol.Beans)
}

//func Handler(w http.ResponseWriter, r *http.Request) {
//	Count ++
//	role := &Role{
//		CpRoleId:   strconv.Itoa(Count),
//		UserId:     1,
//		GameId:     1,
//		RoleName:   "战士",
//		RoleGrade:  "79",
//		GameRegion: "五行山",
//	}
//	Beans <- role
//	fmt.Fprintf(w, "Hello, GrFrHuang !")
//}
