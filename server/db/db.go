package db

import (
	"fmt"

	"item/server/config"
	"item/server/log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var connObj *gorm.DB

func init() {
	dsn := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=%s TimeZone=%s",
		config.DbHost, config.DbPort, config.DbUser, config.DbName, config.DbPassword, "disable", config.DbTz)

	connObj, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{
		FullSaveAssociations: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			DatabaseLocation, tzErr := time.LoadLocation(config.DbTz)
			if tzErr != nil {
				log.Error.Println("Cannot load database timezone info from .env, defaulting to UTC...")
				return time.Now().UTC().Truncate(time.Microsecond)
			}
			return time.Now().In(DatabaseLocation).Truncate(time.Microsecond)
		},
		PrepareStmt: true,
		Logger:      log.GetQueryLogger(),
	})

	if err := Ping(); err != nil {
		log.Error.Println(err)
		return
	} else {
		log.Info.Println("Successfully connected to database..")
	}
}

// Conn is globally accessible database handle
func Conn() *gorm.DB {
	if Ping() != nil {
		log.Error.Panicln("Cannot connect to database..")
	}

	if config.Debug {
		return connObj.Debug()
	}

	return connObj
}

// ConnWithoutDebug is globally accessible database handle
func ConnWithoutDebug() *gorm.DB {
	if Ping() != nil {
		log.Error.Panicln("Cannot connect to database..")
	}

	return connObj
}

// Ping will return error if database connection was unsuccessful
func Ping() error {
	db, err := connObj.DB()
	if err != nil {
		return err
	}

	return db.Ping()
}

// PreloadUnscoped can be used as a callback function
// for Preload() to preload soft deleted records
func PreloadUnscoped(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}
