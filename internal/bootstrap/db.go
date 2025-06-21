package bootstrap

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"interview-backend/internal/conf"
	"interview-backend/internal/db"
	stdlog "log"
	"strings"
	"time"
)

func InitDB() {
	logLevel := logger.Silent
	newLogger := logger.New(
		stdlog.New(log.StandardLogger().Out, "\r\n", stdlog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	var dB *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.Conf.Database.TablePrefix,
		},
		Logger: newLogger,
	}

	dbConf := conf.Conf.Database
	switch dbConf.Type {
	case "sqlite3":
		{
			if !(strings.HasSuffix(dbConf.DBFile, ".db") && len(dbConf.DBFile) > 3) {
				log.Fatalf("invalid database file name: %s", dbConf.DBFile)
			}
			// _journal=WAL设置为写前日志，提升并发性能，特别适合读多写少
			// _vacuum=incremental增量清理，使空间回收更加平缓，避免性能波动
			dB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental",
				dbConf.DBFile)), gormConfig)
		}
	case "mysql":
		{
			dsn := dbConf.DSN
			if dsn == "" {
				dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s",
					dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Name, dbConf.SSLMode)
			}
			dB, err = gorm.Open(mysql.Open(dsn), gormConfig)
		}
	case "postgres":
		{
			dsn := dbConf.DSN
			if dsn == "" {
				if dbConf.Password != "" {
					dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
						dbConf.Host, dbConf.User, dbConf.Password, dbConf.Name, dbConf.Port, dbConf.SSLMode)
				} else {
					dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
						dbConf.Host, dbConf.User, dbConf.Name, dbConf.Port, dbConf.SSLMode)
				}
			}
			dB, err = gorm.Open(postgres.Open(dsn), gormConfig)
		}
	default:
		log.Fatalf("Invalid database type: %s", dbConf.Type)
	}
	
	if err != nil {
		log.Fatalf("failed to connect database:%s", err.Error())
	}
	db.Init(dB)
}
