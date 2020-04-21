package main

import (
	"context"
	"os"

	"github.com/axiaoxin-com/logging"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
)

func init() {
	// Create gorm db instance
	db, err = gorm.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		panic(err)
	}
	// Enable Logger, show detailed log
	db.LogMode(true)
}

func main() {
	// defer clear
	defer db.Close()
	defer os.Remove("./sqlite3.db")

	// Mock a context with a trace id and logger
	traceID := "logging-fake-trace-id"
	ctx := logging.Context(context.Background(), logging.DefaultLogger(), traceID)

	// 打印带 trace id 的 gorm 日志
	// 必须先对 db 对象设置带有 trace id 的 ctxlogger 作为 sql 日志打印的 logger
	// 后续的 gorm 操作使用新的 db 对象即可
	ctxLoggerDB := logging.GormDBWithCtxLogger(ctx, db)

	// Migrate the schema
	ctxLoggerDB.AutoMigrate(&Product{})

	// Create
	ctxLoggerDB.Create(&Product{Code: "L1212", Price: 1000})
}

// log:
// {"level":"DEBUG","time":"2020-04-20 17:27:55.918425","logger":"root.gorm","msg":"CREATE TABLE \"products\" (\"id\" integer primary key autoincrement,\"created_at\" datetime,\"updated_at\" datetime,\"deleted_at\" datetime,\"code\" varchar(255),\"price\" integer )","pid":79239,"traceID":"logging-fake-trace-id","vars":null,"rowsAffected":0,"duration":0.001237568}
// {"level":"DEBUG","time":"2020-04-20 17:27:55.919377","logger":"root.gorm","msg":"CREATE INDEX idx_products_deleted_at ON \"products\"(deleted_at) ","pid":79239,"traceID":"logging-fake-trace-id","vars":null,"rowsAffected":0,"duration":0.000748753}
// {"level":"DEBUG","time":"2020-04-20 17:27:55.919790","logger":"root.gorm","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":79239,"traceID":"logging-fake-trace-id","vars":["2020-04-20T17:27:55.919448+08:00","2020-04-20T17:27:55.919448+08:00",null,"L1212",1000],"rowsAffected":1,"duration":0.000332846}
