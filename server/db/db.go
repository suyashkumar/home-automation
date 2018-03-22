package db

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
	"github.com/suyashkumar/auth"
	"github.com/suyashkumar/conduit/server/entities"
)

const DefaultMaxIdleConns = 5

var ErrorNoConnectionString = errors.New("A connection string must be specified on the first call to Get")

// DatabaseHandler abstracts away common persistence operations needed for this package
type DatabaseHandler interface {
	// GetUser gets a user from the database that matches constraints on the input user
	GetUser(u auth.User) (auth.User, error)
	// UpsertUser updates a user (if input user UUID matches one in the db) or inserts a user
	UpsertUser(u auth.User) error
	// GetDeviceSecret gets a user's device secret
	GetDeviceSecret(uuid uuid.UUID) (entities.DeviceSecret, error)
	// InsertDeviceSecret updates or inserts a device secret for the User TODO: make this non-descructive in future
	InsertDeviceSecret(uuid uuid.UUID, ds entities.DeviceSecret) error
	// GetDB returns the DatabaseHandler's underlying *gorm.DB
	GetDB() *gorm.DB
}

type databaseHandler struct {
	db            *gorm.DB
	authDBHandler auth.DatabaseHandler
}

// NewDatabaseHandler initializes and returns a new DatabaseHandler
func NewDatabaseHandler(dbConnection string) (DatabaseHandler, error) {
	db, err := getDB(dbConnection)
	if err != nil {
		return nil, err
	}
	// AutoMigrate relevant schemas
	db.AutoMigrate(&entities.DeviceSecret{})
	ah, err := auth.NewDatabaseHandlerFromGORM(db)
	if err != nil {
		return nil, err
	}
	return &databaseHandler{
		db:            db,
		authDBHandler: ah,
	}, nil
}

func (d *databaseHandler) GetUser(u auth.User) (auth.User, error) {
	return d.authDBHandler.GetUser(u)
}

func (d *databaseHandler) UpsertUser(u auth.User) error {
	return d.authDBHandler.UpsertUser(u)
}

func (d *databaseHandler) GetDeviceSecret(uuid uuid.UUID) (entities.DeviceSecret, error) {
	var foundDeviceSecret entities.DeviceSecret
	// this could return multiple, but convention right now is one secret per user. May change in future
	err := d.db.Where(entities.DeviceSecret{UserUUID: uuid}).Order("created_at desc").First(&foundDeviceSecret).Error
	if err != nil {
		return foundDeviceSecret, err
	}
	return foundDeviceSecret, nil
}

func (d *databaseHandler) InsertDeviceSecret(uuid uuid.UUID, secret entities.DeviceSecret) error {
	err := d.db.Create(&secret).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *databaseHandler) GetDB() *gorm.DB {
	return d.db
}

func getDB(dbConnection string) (*gorm.DB, error) {
	if dbConnection == "" {
		return nil, ErrorNoConnectionString
	}

	d, err := gorm.Open("postgres", dbConnection)
	if err != nil {
		return nil, err
	}

	d.DB().SetMaxIdleConns(DefaultMaxIdleConns)

	return d, nil

}
