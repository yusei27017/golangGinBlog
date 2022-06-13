package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	MongoDB "goTest/DB/mongoDB"
)

//new一個字定義的 mongo實例
var mongo = MongoDB.New("")

//map映射
var apiIndexCallBackMap = map[string]func(*gin.Context, map[string]interface{}){
	"getData": getData,
}

func ApiIndex(pCont *gin.Context) {

	// 從 body 獲取資料
	pBodyData, pErr := pCont.GetRawData()

	// 如果出錯
	if pErr != nil {
		fmt.Println(pErr)
		return
	}

	var pParseData map[string]interface{}
	// json 結構解析
	pErr = json.Unmarshal(pBodyData, &pParseData)
	if pErr != nil {
		fmt.Println(pErr)
		return
	}

	resFun, _ := apiIndexCallBackMap[pParseData["apitype"].(string)]
	resFun(pCont, pParseData["params"].(map[string]interface{}))
}

type LogInfo struct {
	Sort string            `bson:"sort,omitempty"`
	Data map[string]string `bson:"data,omitempty"`
}

func getData(pCont *gin.Context, params map[string]interface{}) {

	fmt.Println(params)
	var pTmpData = new(LogInfo)
	var dataSli = make([]LogInfo, 0)
	operations := []bson.M{
		//可以加些複合查詢條件
		{"$match": bson.M{"sort": params["sort"]}},
	}

	_ = mongo.GetMany("aboutMe", operations, pTmpData, func() {
		dataSli = append(dataSli, *pTmpData)
	})

	pCont.JSON(200, dataSli)

}
