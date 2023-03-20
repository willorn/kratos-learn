package data

// 用于连接和管理数据库，并定义了一些数据库操作的函数和结构体

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kratos-realworld/internal/conf"
)

// ProviderSet 是 Wire 注册数据层组件所需的依赖集合。它包含了若干个函数，用于创建数据库连接、数据操作等。
var ProviderSet = wire.NewSet(NewData, NewDB, NewUserRepo, NewProfileRepo, NewArticleRepo, NewCommentRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .用于创建 Data 类型的实例。
// 它需要接收三个参数：一个 conf.Data 类型的配置对象，一个 Kratos 的日志对象，以及一个 GORM 数据库连接对象。
// 它返回 Data 实例、一个函数（用于清理资源）和一个错误对象。
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db}, cleanup, nil
}

// NewDB  在数据层(NewData最后一层)，注入config
// 用于创建 GORM 数据库连接对象。它接收一个 conf.Data 类型的配置对象，并使用其 DSN 字段来连接 MySQL 数据库。
func NewDB(c *conf.Data) *gorm.DB {
	// DSN:Data Source Name
	//dsn := "root:@tcp(localhost:3306)/sql_test?charset=utf8mb4&parseTime=True"
	//db, _ = sql.Open("mysql", dsn)
	//db, err := gorm.Open(mysql.Open(dsn))

	db, err := gorm.Open(mysql.Open(c.Database.Dsn), &gorm.Config{
		//当进行迁移时，是否禁用外键约束。
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	// 如果连接失败，则会抛出 panic 异常。
	if err != nil {
		panic("failed to connect database")
	}
	// 连接成功后，它会调用 InitDB 函数来进行数据库迁移，并返回连接对象。
	InitDB(db)
	return db
}

// InitDB 它接收一个 GORM 数据库连接对象，并使用 AutoMigrate 函数来自动创建或修改数据库表结构。它会创建名为 User、Article、Comment、ArticleFavorite 和 Following 的表。
func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(
		&User{},
		&Article{},
		&Comment{},
		&ArticleFavorite{},
		&Following{},
	); err != nil {
		panic(err)
	}
}
