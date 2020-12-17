package main

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/axiaoxin-com/logging"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	db, err = gorm.Open(sqlite.Open("./sqlite3.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})
}

// G 模拟一次请求处理
func G(traceID string) {

	// 模拟一个 ctx ，并将 logger 和 traceID 设置到 ctx 中
	// 这里使用 Options 设置为打印 caller 字段
	ctx, _ := logging.NewCtxLogger(context.Background(), logging.CloneLogger("gorm"), traceID)

	// 新建会话模式设置 logger，也可以在 Open 时 使用 Config 设置
	db = db.Session(&gorm.Session{Logger: logging.NewGormLogger(zap.InfoLevel, zap.InfoLevel, time.Millisecond*200)})

	// 打印带 trace id 的 gorm 日志
	// 必须先对 db 对象设置带有 trace id 的 ctxlogger 作为 sql 日志打印的 logger
	// 后续的 gorm 操作使用新的 db 对象即可
	// 第三个参数为指定使用哪个级别的方法打印 sql 日志
	// Create
	db.WithContext(ctx).Create(&Product{Code: traceID, Price: 1000})
	// Query
	var products []Product
	db.WithContext(ctx).Find(&products)

	wg.Done()
}

func main() {
	// defer clear
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
// {"level":"INFO","time":"2020-10-14 15:15:32.914308","logger":"logging.ctx_logger","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"logging_bu3ab518d3b11hmdhac0","latency":0.000065205,"sql":"SELECT count(*) FROM sqlite_master WHERE type='table' AND name=\"products\"","rows":-1}
// {"level":"INFO","time":"2020-10-14 15:15:32.915582","logger":"logging.ctx_logger","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"logging_bu3ab518d3b11hmdhacg","latency":0.000999756,"sql":"CREATE TABLE `products` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`code` text,`price` integer,PRIMARY KEY (`id`))","rows":0}
// {"level":"INFO","time":"2020-10-14 15:15:32.916391","logger":"logging.ctx_logger","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"logging_bu3ab518d3b11hmdhad0","latency":0.000751082,"sql":"CREATE INDEX `idx_products_deleted_at` ON `products`(`deleted_at`)","rows":0}
// {"level":"INFO","time":"2020-10-14 15:15:32.918683","logger":"logging.example","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"trace-id-g4","latency":0.002183961,"sql":"INSERT INTO `products` (`created_at`,`updated_at`,`deleted_at`,`code`,`price`) VALUES (\"2020-10-14 15:15:32.916\",\"2020-10-14 15:15:32.916\",NULL,\"trace-id-g4\",1000)","rows":1}
// {"level":"INFO","time":"2020-10-14 15:15:32.918896","logger":"logging.example","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"trace-id-g4","latency":0.000176945,"sql":"SELECT * FROM `products` WHERE `products`.`deleted_at` IS NULL","rows":1}
// {"level":"INFO","time":"2020-10-14 15:15:32.922212","logger":"logging.example","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"trace-id-g1","latency":0.005698028,"sql":"INSERT INTO `products` (`created_at`,`updated_at`,`deleted_at`,`code`,`price`) VALUES (\"2020-10-14 15:15:32.917\",\"2020-10-14 15:15:32.917\",NULL,\"trace-id-g1\",1000)","rows":1}
// {"level":"INFO","time":"2020-10-14 15:15:32.922490","logger":"logging.example","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"trace-id-g1","latency":0.000233759,"sql":"SELECT * FROM `products` WHERE `products`.`deleted_at` IS NULL","rows":2}
// {"level":"INFO","time":"2020-10-14 15:15:32.925197","logger":"logging.example","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"trace-id-g3","latency":0.008652237,"sql":"INSERT INTO `products` (`created_at`,`updated_at`,`deleted_at`,`code`,`price`) VALUES (\"2020-10-14 15:15:32.92\",\"2020-10-14 15:15:32.92\",NULL,\"trace-id-g3\",1000)","rows":1}
// {"level":"INFO","time":"2020-10-14 15:15:32.925412","logger":"logging.example","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"trace-id-g3","latency":0.000175952,"sql":"SELECT * FROM `products` WHERE `products`.`deleted_at` IS NULL","rows":3}
// {"level":"INFO","time":"2020-10-14 15:15:32.926266","logger":"logging.example","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"trace-id-g2","latency":0.009718568,"sql":"INSERT INTO `products` (`created_at`,`updated_at`,`deleted_at`,`code`,`price`) VALUES (\"2020-10-14 15:15:32.917\",\"2020-10-14 15:15:32.917\",NULL,\"trace-id-g2\",1000)","rows":1}
// {"level":"INFO","time":"2020-10-14 15:15:32.926451","logger":"logging.example","caller":"logging/gorm.go:Trace:82","msg":"gorm trace","pid":4294,"server_ip":"10.66.41.115","trace_id":"trace-id-g2","latency":0.000150784,"sql":"SELECT * FROM `products` WHERE `products`.`deleted_at` IS NULL","rows":4}
