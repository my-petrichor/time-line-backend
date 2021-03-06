package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	account "github.com/my-Sakura/time-line-backend/pkg/account/controller"
	timeline "github.com/my-Sakura/time-line-backend/pkg/timeline/controller"
)

const (
	timeLineRouterGroup = "/api/v1/timeLine"
	accountRouterGroup  = "/api/v1/account"
)

func main() {
	fmt.Println(time.Now())
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:123456@tcp(mysql:3306)/mysql?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		panic(err)
	}

	timeLineConn := timeline.New(dbConn)
	accountConn := account.New(dbConn)

	timeLineConn.RegistRouter(router.Group(timeLineRouterGroup))
	accountConn.RegistRouter(router.Group(accountRouterGroup))

	log.Fatal(router.Run("0.0.0.0:10002"))
}
