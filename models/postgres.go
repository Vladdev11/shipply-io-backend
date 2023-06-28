package models

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var PGDB *gorm.DB

func PostgresInit() {
	//set postgres dsn
	//if development, use local postgres
	dsn := "host=" + util.ConfigPgAddr + " user=" + util.ConfigPgUsername + " password=" + util.ConfigPgPassword + " dbname=" + util.ConfigPgDatabase + " port=5432 sslmode=disable TimeZone=UTC"
	if *util.DevelopmentMode {
		dsn = "host=" + util.ConfigPgDevelopmentAddr + " user=" + util.ConfigPgDevelopmentUsername + " password=" + util.ConfigPgDevelopmentPassword + " dbname=" + util.ConfigPgDevelopmentDatabase + " port=5432 sslmode=disable TimeZone=UTC"
	}

	newLogger := logger.Default
	if *util.PrintSQL {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,       // Disable color
			},
		)
	}

	//init postgres with GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(errors.New("failed to connect to postgres database"))
	}

	//set global db
	PGDB = db

	// Migrate the schema
	PGDB.AutoMigrate(&Order{},
		&OrderItem{},
		&Inventory{},
		&Product{},
		&ProductAlias{},
		&Kit{},
		&ProductKit{},
		&OrderTags{},
		&PurchaseOrder{},
		&PurchaseOrderItem{},
		&Client{},
		&Organization{},
		&User{},
		&Warehouse{},
		&Vendor{},
		&Location{},
		&PurchaseOrderStatus{},
		&PurchaseOrderTag{},
		&PurchaseOrderHistory{},
		&Attachment{},
		&PurchaseOrderAttachment{},
		&Address{},
		&LocationType{},
		&Carrier{},
		&CarrierConnection{},
		&Store{},
		&ShopifyProduct{},
		&ShopifyApiLog{},
		&ShopifyLocation{},
		&SystemTask{},
		&SystemError{},
		&ShipEngineLog{},
		&Shipment{},
		&AutomationRule{},
		&PickSession{},
		&PickSessionOrder{},
		&PickSessionOrderItem{},
		&PickSessionOrderError{},
		&ShippingRate{},
		&ShippingMethod{},
		&PasswordResetToken{},
		&InventoryAuditLog{},
		&ProductLot{},
		&PurchaseOrderItemRejection{},
		&PurchaseOrderItemRejectionAttachment{},
		&PurchaseOrderItemHistory{},
	)

}

type NullInt64 struct {
	sql.NullInt64
}
