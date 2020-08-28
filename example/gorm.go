package main

import (
	"context"
	"os"
	"sync"

	"github.com/axiaoxin-com/logging"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	// 模拟一个 ctx ，并将 logger 和 traceID 设置到 ctx 中
	// 这里使用 Options 设置为打印 caller 字段
	ctx := logging.Context(context.Background(), logging.CloneDefaultLogger("example").WithOptions(zap.AddCaller()), traceID)

	// 打印带 trace id 的 gorm 日志
	// 必须先对 db 对象设置带有 trace id 的 ctxlogger 作为 sql 日志打印的 logger
	// 后续的 gorm 操作使用新的 db 对象即可
	// 第三个参数为指定使用哪个级别的方法打印 sql 日志
	db := logging.GormDBWithCtxLogger(ctx, db, zapcore.DebugLevel)
	// Create
	db.Create(&Product{Code: traceID, Price: 1000})
	// Query
	var products []Product
	db.Find(&products)
	wg.Done()
}

func main() {
	// defer clear
	defer db.Close()
	defer os.Remove("./sqlite3.db")

	// 模拟并发
	wg.Add(4)
	go G("trace-id-g1")
	go G("trace-id-g2")
	go G("trace-id-g3")
	go G("trace-id-g4")
	wg.Wait()
}

// log:
// {"level":"DEBUG","time":"2020-06-17 13:16:41.601297","logger":"logging.gorm","caller":"example/gorm.go:G:55","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":9748,"traceID":"trace-id-g4","vars":["2020-06-17T13:16:41.600679+08:00","2020-06-17T13:16:41.600679+08:00",null,"trace-id-g4",1000],"rowsAffected":1,"duration":0.000602749}
// {"level":"DEBUG","time":"2020-06-17 13:16:41.603107","logger":"logging.gorm","caller":"example/gorm.go:G:58","msg":"SELECT * FROM \"products\"  WHERE \"products\".\"deleted_at\" IS NULL","pid":9748,"traceID":"trace-id-g4","vars":null,"rowsAffected":1,"duration":0.000159561}
// {"level":"DEBUG","time":"2020-06-17 13:16:41.605189","logger":"logging.gorm","caller":"example/gorm.go:G:55","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":9748,"traceID":"trace-id-g2","vars":["2020-06-17T13:16:41.601395+08:00","2020-06-17T13:16:41.601395+08:00",null,"trace-id-g2",1000],"rowsAffected":1,"duration":0.003753052}
// {"level":"DEBUG","time":"2020-06-17 13:16:41.605765","logger":"logging.gorm","caller":"example/gorm.go:G:58","msg":"SELECT * FROM \"products\"  WHERE \"products\".\"deleted_at\" IS NULL","pid":9748,"traceID":"trace-id-g2","vars":null,"rowsAffected":2,"duration":0.000129308}
// {"level":"DEBUG","time":"2020-06-17 13:16:41.610385","logger":"logging.gorm","caller":"example/gorm.go:G:55","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":9748,"traceID":"trace-id-g1","vars":["2020-06-17T13:16:41.601498+08:00","2020-06-17T13:16:41.601498+08:00",null,"trace-id-g1",1000],"rowsAffected":1,"duration":0.008860571}
// {"level":"DEBUG","time":"2020-06-17 13:16:41.611072","logger":"logging.gorm","caller":"example/gorm.go:G:58","msg":"SELECT * FROM \"products\"  WHERE \"products\".\"deleted_at\" IS NULL","pid":9748,"traceID":"trace-id-g1","vars":null,"rowsAffected":3,"duration":0.000143793}
// {"level":"DEBUG","time":"2020-06-17 13:16:41.621077","logger":"logging.gorm","caller":"example/gorm.go:G:55","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":9748,"traceID":"trace-id-g3","vars":["2020-06-17T13:16:41.601322+08:00","2020-06-17T13:16:41.601322+08:00",null,"trace-id-g3",1000],"rowsAffected":1,"duration":0.019732596}
// {"level":"DEBUG","time":"2020-06-17 13:16:41.622074","logger":"logging.gorm","caller":"example/gorm.go:G:58","msg":"SELECT * FROM \"products\"  WHERE \"products\".\"deleted_at\" IS NULL","pid":9748,"traceID":"trace-id-g3","vars":null,"rowsAffected":4,"duration":0.000171964}
