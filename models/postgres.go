package models

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(k *koanf.Koanf, debug bool) (*gorm.DB, error) {
	//set postgres dsn
	dsn := "host=" + k.String("postgres.hostname") +
		" user=" + k.String("postgres.username") +
		" password=" + k.String("postgres.password") +
		" dbname=" + k.String("postgres.database") +
		" port=5432 sslmode=disable TimeZone=UTC"

	newLogger := logger.Default
	if debug {
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
		return nil, err
	}

	db.AutoMigrate(&Order{},
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
		&ProductBundle{},
		&OrderTag{},
		&OrderStatus{},
		&OrderHistory{},
	)
	// TODO: ADD CHECKS FOR THE ABOVE MIGRATION, AND RETURN AN ERROR IF IT FAILS
	// Right now it just fails silently, and the app will crash later on when it tries to access a table that doesn't exist
	// It's actually already failing as can be seen by the first log message you get when the app is started.
	// Also the failure here is kind of complicated as the pick_session_order stuff seems to have a circular dependency with pick_session_order_{item,error} tables.
	// I'm not sure how to fix this, but I think it's a good idea to fix it before we go to production...
	// The only reason this probably doesn't happen in your local environment is because you have the tables already created from the previous version of the app.

	return db, nil
}

type NullInt64 struct {
	sql.NullInt64
}

type StringSlice []string

func (o *StringSlice) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		*o = strings.Split(string(src), ",")
		return nil
	case string:
		*o = strings.Split(src, ",")
		return nil
	default:
		return errors.New("incompatible type for StringSlice")
	}
}
func (o StringSlice) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	return strings.Join(o, ","), nil
}
