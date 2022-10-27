package main

import (
	"github.com/gin-gonic/gin"
	"goTest/tmpCtrl"

	"goTest/ctrl"
)

func main() {

	mainRoute := gin.Default()
	mainRoute.LoadHTMLGlob("tmpl/*.tmpl")
	mainRoute.Static("/static", "./static")

	indexRouter := mainRoute.Group("/index")
	indexRouter.GET("/", tmpCtrl.MainView)
	indexRouter.GET("/aboutMe", tmpCtrl.AboutMe)
	indexRouter.GET("/ps5Page", tmpCtrl.Ps5Page)
	indexRouter.GET("/pixivPage", tmpCtrl.PixivPage)
	indexRouter.GET("/skillPage", tmpCtrl.SkillPage)
	indexRouter.GET("/aboutJp", tmpCtrl.AboutJp)
	indexRouter.GET("/linkPage", tmpCtrl.LinkPage)
	indexRouter.GET("/aboutSite", tmpCtrl.AboutSite)
	indexRouter.GET("/runTimer", tmpCtrl.RunTimer)

	mainRoute.POST("/apiIndex", ctrl.ApiIndex)

	mainRoute.Run()
}
