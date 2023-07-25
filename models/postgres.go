package models

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/shipply-io/shipply-io-backend/util"
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
		&UserSavedFilter{},
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
		&AutomationRule{},
		&PickSession{},
		&PickSessionOrder{},
		&PickSessionOrderItem{},
		&ShippingRate{},
		&ShippingMethod{},
		&Shipment{},
		&PasswordResetToken{},
		&InventoryAuditLog{},
		&ProductLot{},
		&PurchaseOrderItemRejection{},
		&PurchaseOrderItemRejectionAttachment{},
		&PurchaseOrderItemHistory{},
		&ProductBundle{},
		&ProductImage{},
		&OrderTag{},
		&OrderStatus{},
		&OrderHistory{},
	)

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

// ContextWithTx runs a callback function with a transaction, and rolls back the transaction if the callback returns an error.
// If the callback returns nil, the transaction is committed.
// The transaction is attached to the context, so it can be accessed with util.DBFromContext(ctx) in all model functions.
func ContextWithTx(ctx context.Context, cb func(context context.Context) error) error {
    err := util.DBFromContext(ctx).Transaction(func(tx *gorm.DB) error {
        return cb(util.ContextWithDB(ctx, tx))
    })
    return err
}