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

func main() {
	// Mock a context with a trace id and logger
	traceID := "logging-fake-trace-id"
	ctx := logging.Context(context.Background(), logging.DefaultLogger(), traceID)

	// Create gorm db instance
	db, err := gorm.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	defer os.Remove("./sqlite3.db")

	// Enable Logger, show detailed log
	db.LogMode(true)

	// Set logging GormLogger for gorm
	logging.SetGormLogger(ctx, db)

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})
}

// log:
// {"level":"DEBUG","time":"2020-04-20 17:27:55.915805","logger":"root.ctxLogger","msg":"Running AtomicLevel HTTP server on :1903","pid":79239}
// {"level":"DEBUG","time":"2020-04-20 17:27:55.916801","logger":"root.gorm","msg":"logging create and set GormLogger successful","pid":79239,"traceID":"logging-fake-trace-id"}
// {"level":"DEBUG","time":"2020-04-20 17:27:55.918425","logger":"root.gorm","msg":"CREATE TABLE \"products\" (\"id\" integer primary key autoincrement,\"created_at\" datetime,\"updated_at\" datetime,\"deleted_at\" datetime,\"code\" varchar(255),\"price\" integer )","pid":79239,"traceID":"logging-fake-trace-id","vars":null,"rowsAffected":0,"duration":0.001237568}
// {"level":"DEBUG","time":"2020-04-20 17:27:55.919377","logger":"root.gorm","msg":"CREATE INDEX idx_products_deleted_at ON \"products\"(deleted_at) ","pid":79239,"traceID":"logging-fake-trace-id","vars":null,"rowsAffected":0,"duration":0.000748753}
// {"level":"DEBUG","time":"2020-04-20 17:27:55.919790","logger":"root.gorm","msg":"INSERT INTO \"products\" (\"created_at\",\"updated_at\",\"deleted_at\",\"code\",\"price\") VALUES (?,?,?,?,?)","pid":79239,"traceID":"logging-fake-trace-id","vars":["2020-04-20T17:27:55.919448+08:00","2020-04-20T17:27:55.919448+08:00",null,"L1212",1000],"rowsAffected":1,"duration":0.000332846}
