package dao

import (
	"fmt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/xray_pool/internal/models"
	"github.com/allanpk716/xray_pool/internal/pkg"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	dbLogger "gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"sync"
)

// Get 获取数据库实例
func Get() *gorm.DB {
	if db == nil {
		once.Do(func() {
			err := initDB()
			if err != nil {
				logger.Panicln(err)
			}
		})
	}
	return db
}

func initDB() error {

	var err error

	// sqlite3
	nowDBFName := filepath.Join(".", dbFileName)
	dbDir := filepath.Dir(nowDBFName)
	if pkg.IsDir(dbDir) == false {
		err = os.MkdirAll(dbDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	db, err = gorm.Open(sqlite.Open(nowDBFName), &gorm.Config{})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to connect database, %s", err.Error()))
	}
	// 降低 gorm 的日志级别
	db.Logger = dbLogger.Default.LogMode(dbLogger.Silent)
	// 迁移 schema
	err = db.AutoMigrate(
		&models.Subscribe{},
		&models.Node{},
	)
	if err != nil {
		return errors.New(fmt.Sprintf("db AutoMigrate error, %s", err.Error()))
	}

	return nil
}

const (
	dbFileName = "xray_pool.db"
)

var (
	db   *gorm.DB
	once sync.Once
)
