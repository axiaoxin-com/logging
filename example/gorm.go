package main

import (
	"context"
	"os"
	"sync"

	"github.com/axiaoxin-com/logging"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.uber.org/zap"
)

// Product test model
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var (
	db  *gorm.DB
	err error
	wg  sync.WaitGroup
)

func init() {
	// Create gorm db instance
	db, err = gorm.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})
	// Enable Logger, show detailed log
	db.LogMode(true)

}

// G 模拟一次请求处理
func G(traceID string) {

	// Mock a context with a trace id and logger
	ctx := logging.Context(context.Background(), logging.DefaultLogger().WithOptions(zap.AddCaller()), traceID)

	// 打印带 trace id 的 gorm 日志
	// 必须先对 db 对象设置带有 trace id 的 ctxlogger 作为 sql 日志打印的 logger
	// 后续的 gorm 操作使用新的 db 对象即可
	db := logging.GormDBWithCtxLogger(ctx, db)
	// Create
	db.Create(&Product{Code: traceID, Price: 1000})
	wg.Done()
}

func main() {
	// defer clear
	defer db.Close()
	defer os.Remove("./sqlite3.db")

	// 模拟并发
	wg.Add(4)
	go G("g1")
	go G("g2")
	go G("g3")
	go G("g4")
	wg.Wait()
}

// log:
// {"level":"DEBUG","time":"2020-04-21 17:08:44.449254","logger":"root.gorm","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":29826,"traceID":"g4","vars":["2020-04-21T17:08:44.448622+08:00","2020-04-21T17:08:44.448622+08:00",null,"g4",1000],"rowsAffected":1,"duration":0.000613636}
// {"level":"DEBUG","time":"2020-04-21 17:08:44.452657","logger":"root.gorm","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":29826,"traceID":"g2","vars":["2020-04-21T17:08:44.44919+08:00","2020-04-21T17:08:44.44919+08:00",null,"g2",1000],"rowsAffected":1,"duration":0.0034358}
// {"level":"DEBUG","time":"2020-04-21 17:08:44.458721","logger":"root.gorm","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":29826,"traceID":"g1","vars":["2020-04-21T17:08:44.44946+08:00","2020-04-21T17:08:44.44946+08:00",null,"g1",1000],"rowsAffected":1,"duration":0.009227084}
// {"level":"DEBUG","time":"2020-04-21 17:08:44.471094","logger":"root.gorm","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":29826,"traceID":"g3","vars":["2020-04-21T17:08:44.449387+08:00","2020-04-21T17:08:44.449387+08:00",null,"g3",1000],"rowsAffected":1,"duration":0.021678226}
