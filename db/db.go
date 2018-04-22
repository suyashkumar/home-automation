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

// Handler abstracts away common persistence operations needed for this package
type Handler interface {
	// GetUser gets a user from the database that matches constraints on the input user
	GetUser(u auth.User) (auth.User, error)
	// UpsertUser updates a user (if input user UUID matches one in the db) or inserts a user
	UpsertUser(u auth.User) error
	// GetAccountSecret gets a user's device secret
	GetAccountSecret(uuid uuid.UUID) (entities.AccountSecret, error)
	// InsertAccountSecret updates or inserts a device secret for the User
	InsertAccountSecret(uuid uuid.UUID, ds entities.AccountSecret) error
	// GetDB returns the Handler's underlying *gorm.DB
	GetDB() *gorm.DB
}

type handler struct {
	db            *gorm.DB
	authDBHandler auth.DatabaseHandler
}

// NewHandler initializes and returns a new Handler
func NewHandler(dbConnection string) (Handler, error) {
	db, err := getDB(dbConnection)
	if err != nil {
		return nil, err
	}
	// AutoMigrate relevant schemas
	db.AutoMigrate(&entities.AccountSecret{})
	ah, err := auth.NewDatabaseHandlerFromGORM(db)
	if err != nil {
		return nil, err
	}
	return &handler{
		db:            db,
		authDBHandler: ah,
	}, nil
}

func (d *handler) GetUser(u auth.User) (auth.User, error) {
	return d.authDBHandler.GetUser(u)
}

func (d *handler) UpsertUser(u auth.User) error {
	return d.authDBHandler.UpsertUser(u)
}

func (d *handler) GetAccountSecret(uuid uuid.UUID) (entities.AccountSecret, error) {
	var foundDeviceSecret entities.AccountSecret
	// this could return multiple, but convention right now is one secret per user. May change in future
	err := d.db.Where(entities.AccountSecret{UserUUID: uuid}).Order("created_at desc").First(&foundDeviceSecret).Error
	if err != nil {
		return foundDeviceSecret, err
	}
	return foundDeviceSecret, nil
}

func (d *handler) InsertAccountSecret(uuid uuid.UUID, secret entities.AccountSecret) error {
	err := d.db.Create(&secret).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *handler) GetDB() *gorm.DB {
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
