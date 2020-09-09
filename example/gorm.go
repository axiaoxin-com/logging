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
	ctx, _ := logging.NewCtxLogger(context.Background(), logging.CloneLogger("example").WithOptions(zap.AddCaller()), traceID)

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
// {"level":"DEBUG","time":"2020-09-09 08:50:33.159863","logger":"logging.example.ctx_logger.gorm","caller":"example/gorm.go:G:55","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":80928,"server_ip":"10.64.35.39","trace_id":"trace-id-g4","vars":["2020-09-09T08:50:33.158544+08:00","2020-09-09T08:50:33.158544+08:00",null,"trace-id-g4",1000],"affected":1,"duration":0.001300618}
// {"level":"DEBUG","time":"2020-09-09 08:50:33.160672","logger":"logging.example.ctx_logger.gorm","caller":"example/gorm.go:G:58","msg":"SELECT * FROM \"products\"  WHERE \"products\".\"deleted_at\" IS NULL","pid":80928,"server_ip":"10.64.35.39","trace_id":"trace-id-g4","vars":null,"affected":1,"duration":0.000169642}
// {"level":"DEBUG","time":"2020-09-09 08:50:33.161038","logger":"logging.example.ctx_logger.gorm","caller":"example/gorm.go:G:55","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":80928,"server_ip":"10.64.35.39","trace_id":"trace-id-g2","vars":["2020-09-09T08:50:33.159492+08:00","2020-09-09T08:50:33.159492+08:00",null,"trace-id-g2",1000],"affected":1,"duration":0.001508938}
// {"level":"DEBUG","time":"2020-09-09 08:50:33.161666","logger":"logging.example.ctx_logger.gorm","caller":"example/gorm.go:G:58","msg":"SELECT * FROM \"products\"  WHERE \"products\".\"deleted_at\" IS NULL","pid":80928,"server_ip":"10.64.35.39","trace_id":"trace-id-g2","vars":null,"affected":2,"duration":0.000141228}
// {"level":"DEBUG","time":"2020-09-09 08:50:33.163229","logger":"logging.example.ctx_logger.gorm","caller":"example/gorm.go:G:55","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":80928,"server_ip":"10.64.35.39","trace_id":"trace-id-g3","vars":["2020-09-09T08:50:33.159803+08:00","2020-09-09T08:50:33.159803+08:00",null,"trace-id-g3",1000],"affected":1,"duration":0.003394786}
// {"level":"DEBUG","time":"2020-09-09 08:50:33.163880","logger":"logging.example.ctx_logger.gorm","caller":"example/gorm.go:G:58","msg":"SELECT * FROM \"products\"  WHERE \"products\".\"deleted_at\" IS NULL","pid":80928,"server_ip":"10.64.35.39","trace_id":"trace-id-g3","vars":null,"affected":3,"duration":0.00015003}
// {"level":"DEBUG","time":"2020-09-09 08:50:33.169651","logger":"logging.example.ctx_logger.gorm","caller":"example/gorm.go:G:55","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":80928,"server_ip":"10.64.35.39","trace_id":"trace-id-g1","vars":["2020-09-09T08:50:33.159904+08:00","2020-09-09T08:50:33.159904+08:00",null,"trace-id-g1",1000],"affected":1,"duration":0.009720455}
// {"level":"DEBUG","time":"2020-09-09 08:50:33.170298","logger":"logging.example.ctx_logger.gorm","caller":"example/gorm.go:G:58","msg":"SELECT * FROM \"products\"  WHERE \"products\".\"deleted_at\" IS NULL","pid":80928,"server_ip":"10.64.35.39","trace_id":"trace-id-g1","vars":null,"affected":4,"duration":0.000159303}
