package tmpCtrl

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func MainView(pCont *gin.Context) {
	pCont.HTML(http.StatusOK, "homePage.tmpl", nil)
}

func AboutMe(pCont *gin.Context) {
	pCont.HTML(http.StatusOK, "aboutMe.tmpl", nil)
}

func Ps5Page(pCont *gin.Context) {
	pCont.HTML(http.StatusOK, "ps5Page.tmpl", nil)
}

func PixivPage(pCont *gin.Context) {
	pCont.HTML(http.StatusOK, "pixivPage.tmpl", nil)
}

func SkillPage(pCont *gin.Context) {
	pCont.HTML(http.StatusOK, "skillPage.tmpl", nil)
}

func AboutJp(pCont *gin.Context) {
	pCont.HTML(http.StatusOK, "aboutJp.tmpl", nil)
}

func LinkPage(pCont *gin.Context) {
	pCont.HTML(http.StatusOK, "linkPage.tmpl", nil)
}

func AboutSite(pCont *gin.Context) {
	pCont.HTML(http.StatusOK, "aboutSite.tmpl", nil)
}
