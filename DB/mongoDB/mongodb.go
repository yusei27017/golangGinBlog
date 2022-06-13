package MongoDB

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//-----------------------------------------------------------------------------
//	資料結構
//-----------------------------------------------------------------------------

// MongoDBType : MongoDB 主結構
type MongoDBType struct {
	sync.Mutex // 繼承

	sURLBase string // 資料庫基本 sURLBase
	sDBName  string // 資料庫名稱

	// 連線物件
	pDBConn   *mongo.Client
	pDBConnDB *mongo.Database

	// control flag
	bIsLink bool

	// Quit chan
	bFlagQuit chan bool
}

//-----------------------------------------------------------------------------
//	變數
//-----------------------------------------------------------------------------
var pMongoDBObjs []*MongoDBType

var pMongoDBMapLock sync.Mutex
var pMongoDBMap map[string]*MongoDBType

// ----------------------------------------------------------------------------
// 函式
// ----------------------------------------------------------------------------
func init() {
	fmt.Println("init function was called.")
	pMongoDBObjs = make([]*MongoDBType, 0)
	pMongoDBMap = make(map[string]*MongoDBType)
}

func New(sName string) *MongoDBType {
	var pErr error

	if sName == "" {
		sName = "Base"
	}

	pMongoDBMapLock.Lock()
	pMongoDB, bExist := pMongoDBMap[sName]
	pMongoDBMapLock.Unlock()
	if bExist {
		return pMongoDB
	}

	pMongoDB = new(MongoDBType)

	pMongoDB.bFlagQuit = make(chan bool)
	pMongoDB.bIsLink = false

	pMongoDB.sURLBase = "mongodb://127.0.0.1:27017/"

	pMongoDB.sDBName = "golangGinBlog"

	pMongoDB.pDBConn, pErr = mongo.NewClient(options.Client().ApplyURI(pMongoDB.sURLBase))
	if pErr != nil {
		fmt.Println(pErr)
		pMongoDB = nil
		return pMongoDB
	}

	pContext, pContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pContextCancel()

	pErr = pMongoDB.pDBConn.Connect(pContext)
	if pErr == nil {
		pErr = pMongoDB.pDBConn.Ping(pContext, nil)
		if pErr != nil {
			fmt.Println(pErr)
			pMongoDB.bIsLink = false
			pErr = pMongoDB.pDBConn.Disconnect(pContext)
			if pErr != nil {
				fmt.Println(pErr)
			}
		} else {
			pMongoDB.bIsLink = true
		}
	}

	if pMongoDB.IsAlive() {
		pMongoDB.pDBConnDB = pMongoDB.pDBConn.Database(pMongoDB.sDBName)
	}

	go pMongoDB.pingServer()

	pMongoDBObjs = append(pMongoDBObjs, pMongoDB)

	pMongoDBMapLock.Lock()
	pMongoDBMap[sName] = pMongoDB
	pMongoDBMapLock.Unlock()

	return pMongoDB
}

func (pObj *MongoDBType) pingServer() {
	// 建立 1 分鐘 的 Ticker
	pTicker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-pObj.bFlagQuit:
			return
		// Ticket 時間
		case <-pTicker.C:
			pObj.IsAlive()
		}
	}
}

// IsAlive :
func (pObj *MongoDBType) IsAlive() bool {
	var pErr error
	var times int

	pContext, pContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pContextCancel()

	for times = 0; times < 2; times++ {
		pObj.Lock()

		if !pObj.bIsLink {
			pErr = pObj.pDBConn.Connect(pContext)
			if pErr != nil {
				pObj.pDBConnDB = nil
				pObj.Unlock()
				return false
			}

			pObj.pDBConnDB = pObj.pDBConn.Database(pObj.sDBName)

			pObj.bIsLink = true
		}

		// ping
		pErr = pObj.pDBConn.Ping(pContext, nil)
		if pErr != nil {
			pObj.pDBConnDB = nil
			pErr = pObj.pDBConn.Disconnect(pContext)
			if pErr != nil {
				fmt.Println(pErr)
			}
			pObj.Unlock()
			pObj.bIsLink = false
			continue
		}

		pObj.Unlock()
		break
	}

	return pObj.bIsLink
}

// InsertOne :
func (pObj *MongoDBType) InsertOne(sTableName string, pData interface{}) bool {
	var pErr error

	if !pObj.IsAlive() {
		fmt.Println("Can't connect to MongoDB server.")
		return false
	}
	fmt.Println(sTableName)
	pContext, pContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pContextCancel()
	fmt.Println(pData)

	fmt.Println("<<Insert>> Table :", sTableName, "Data : ", pData)

	pObj.Lock()

	pCollection := pObj.pDBConnDB.Collection(sTableName)
	if pCollection == nil {
		fmt.Println("Can't get the", sTableName)
		pObj.Unlock()
		return false
	}

	_, pErr = pCollection.InsertOne(pContext, pData)
	if pErr != nil {
		fmt.Println(pErr)
		pObj.Unlock()
		return false
	}

	pObj.Unlock()

	return true
}

// GetMany :
func (pObj *MongoDBType) GetMany(sTableName string, pFilter interface{}, pTmpData interface{}, pAppend func()) bool {
	var pErr error
	if !pObj.IsAlive() {
		fmt.Println("Can't connect to MongoDB server.")
		return false
	}

	pContext, pContextCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pContextCancel()

	fmt.Println("<<Get>> Table :", sTableName, "Filter : ", pFilter)

	pObj.Lock()

	pCollection := pObj.pDBConnDB.Collection(sTableName)
	if pCollection == nil {
		fmt.Println("Can't get the", sTableName)
		pObj.Unlock()
		return false
	}

	pCursor, pErr := pCollection.Aggregate(pContext, pFilter)
	if pErr != nil {
		fmt.Println(pErr)
		pObj.Unlock()
		return false
	}

	for pCursor.Next(pContext) {
		// create a value into which the single document can be decoded
		pErr = pCursor.Decode(pTmpData)

		if pErr != nil {
			fmt.Println(pErr)
			pObj.Unlock()
			return false
		}
		pAppend()
	}

	pErr = pCursor.Err()
	if pErr != nil {
		fmt.Println(pErr)
		pObj.Unlock()
		return false
	}

	// Close the cursor once finished
	err := pCursor.Close(pContext)
	if err != nil {
		return false
	}
	pObj.Unlock()

	return true
}
