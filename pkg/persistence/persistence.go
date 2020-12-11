package persistence

import (
	"fmt"

	"github.com/spoonrocker/cart-go-sonalys/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Persistence interface {
	Create(interface{}) error
	Find(interface{}) error
	FindAll(interface{}, ...interface{}) error
	Update(interface{}) error
	Delete(interface{}) error
	LoadRelationships() Persistence
}

type persistence struct {
	db *gorm.DB
}

func (p persistence) LoadRelationships() Persistence {
	return persistence{
		db: p.db.Preload(clause.Associations),
	}
}

func (p persistence) Create(item interface{}) error {
	return p.db.Create(item).Error
}

func (p persistence) Find(item interface{}) error {
	return p.db.First(item).Error
}

func (p persistence) FindAll(items interface{}, conds ...interface{}) error {
	return p.db.Find(items, conds...).Error
}

func (p persistence) Update(item interface{}) error {
	return p.db.Save(item).Error
}

func (p persistence) Delete(item interface{}) error {
	return p.db.Delete(item).Error
}

func BuildDatabaseConnString(user, password, host, port, dbName string, ssl bool) string {
	var sslFlag string
	if ssl {
		sslFlag = "enabled"
	} else {
		sslFlag = "disable"
	}
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user,
		password,
		host,
		port,
		dbName,
		sslFlag)
}

func NewPersistence(connURL string) (Persistence, error) {
	db, err := gorm.Open(postgres.Open(connURL), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(models.Cart{}, models.CartItem{})
	db.Logger = db.Logger.LogMode(logger.Silent)

	return persistence{
		db: db,
	}, nil
}
