package handlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/handlers/ClientHandlers"
	"github.com/shipply-io/shipply-io-backend/handlers/OrganizationHandlers"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

// ** USER ROUTES ** //
func UserGet(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserGet(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UserGet(w, r)
	case util.ClientAdmin:
		ClientHandlers.UserGet(w, r)
	case util.ClientUser:
		ClientHandlers.UserGet(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UserUpdatePassword(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserUpdatePassword(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UserUpdatePassword(w, r)
	case util.ClientAdmin:
		ClientHandlers.UserUpdatePassword(w, r)
	case util.ClientUser:
		ClientHandlers.UserUpdatePassword(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UserGetByID(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserGetByID(w, r)
	case util.OrganizationUser:
		util.ErrorResponse(w, "user does not have access to other users", http.StatusForbidden)
	case util.ClientAdmin:
		ClientHandlers.UserGetByID(w, r)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access to other users", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UserCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserCreate(w, r)
	case util.OrganizationUser:
		util.ErrorResponse(w, "user does not have access to create users", http.StatusForbidden)
	case util.ClientAdmin:
		ClientHandlers.UserCreate(w, r)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access to create users", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func UserUpdate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserUpdate(w, r)
	case util.OrganizationUser:
		util.ErrorResponse(w, "user does not have access to update users", http.StatusForbidden)
	case util.ClientAdmin:
		ClientHandlers.UserUpdate(w, r)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access to update users", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func UserUpdateAvatar(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserUpdateAvatar(w, r)
	case util.OrganizationUser:
		util.ErrorResponse(w, "user does not have access to update avatar", http.StatusForbidden)
	case util.ClientAdmin:
		ClientHandlers.UserUpdateAvatar(w, r)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access to update avatar", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func UserDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserDelete(w, r)
	case util.OrganizationUser:
		util.ErrorResponse(w, "user does not have access to delete users", http.StatusForbidden)
	case util.ClientAdmin:
		ClientHandlers.UserDelete(w, r)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access to delete users", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END USER ROUTES ** //

// ** USER SAVED FILTERS ** //

func UserSavedFilterCreate(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserSavedFilterCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UserSavedFilterCreate(w, r)
	case util.ClientAdmin:
		ClientHandlers.UserSavedFilterCreate(w, r)
	case util.ClientUser:
		ClientHandlers.UserSavedFilterCreate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UserSavedFilterList(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserSavedFilterList(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UserSavedFilterList(w, r)
	case util.ClientAdmin:
		ClientHandlers.UserSavedFilterList(w, r)
	case util.ClientUser:
		ClientHandlers.UserSavedFilterList(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UserSavedFilterUpdate(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserSavedFilterUpdate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UserSavedFilterUpdate(w, r)
	case util.ClientAdmin:
		ClientHandlers.UserSavedFilterUpdate(w, r)
	case util.ClientUser:
		ClientHandlers.UserSavedFilterUpdate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UserSavedFilterDelete(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UserSavedFilterDelete(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UserSavedFilterDelete(w, r)
	case util.ClientAdmin:
		ClientHandlers.UserSavedFilterDelete(w, r)
	case util.ClientUser:
		ClientHandlers.UserSavedFilterDelete(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

// ** END USER SAVED FILTERS ** //

// ** PRODUCT ROUTES ** //

func SearchProducts(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.SearchProducts(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.SearchProducts(w, r)
	case util.ClientAdmin:
		ClientHandlers.SearchProducts(w, r)
	case util.ClientUser:
		ClientHandlers.SearchProducts(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func ListProducts(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListProducts(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListProducts(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListProducts(w, r)
	case util.ClientUser:
		ClientHandlers.ListProducts(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func GetProduct(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetProduct(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetProduct(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetProduct(w, r)
	case util.ClientUser:
		ClientHandlers.GetProduct(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func GetProductOrders(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetProductOrders(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetProductOrders(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetProductOrders(w, r)
	case util.ClientUser:
		ClientHandlers.GetProductOrders(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func GetProductInventory(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetProductInventory(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetProductInventory(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetProductInventory(w, r)
	case util.ClientUser:
		ClientHandlers.GetProductInventory(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func GetProductBundles(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetProductBundles(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetProductBundles(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetProductBundles(w, r)
	case util.ClientUser:
		ClientHandlers.GetProductBundles(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func GetProductBundleComponents(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetProductBundleComponents(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetProductBundleComponents(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetProductBundleComponents(w, r)
	case util.ClientUser:
		ClientHandlers.GetProductBundleComponents(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func GetProductStores(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetProductStores(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetProductStores(w, r)
	case util.ClientAdmin:
		// ClientHandlers.GetProductStores(w, r)
	case util.ClientUser:
		// ClientHandlers.GetProductStores(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

// ** END PRODUCT ROUTES ** //

// ** PURCHASE ORDER ROUTES ** //
func PurchaseOrderList(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderList(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderList(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderList(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderList(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderCreate(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderCreate(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderCreate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func PurchaseOrderGet(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderGet(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderGet(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderGet(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderGet(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func PurchaseOrderUpdate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderUpdate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderUpdate(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderUpdate(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderUpdate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderDelete(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderDelete(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderDelete(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderDelete(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END PURCHASE ORDER ROUTES ** //

// ** PURCHASE ORDER ITEM ROUTES ** //
func PurchaseOrderItemUpdateBulk(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderItemBulkUpdate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderItemBulkUpdate(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderItemBulkUpdate(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderItemBulkUpdate(w, r)
	default:
		util.ErrorResponse(w, "user role does not exist", http.StatusBadRequest)
	}

}

func PurchaseOrderItemCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderItemCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderItemCreate(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderItemCreate(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderItemCreate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderItemGet(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderItemGet(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderItemGet(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderItemGet(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderItemGet(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderItemUpdate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderItemUpdate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderItemUpdate(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderItemUpdate(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderItemUpdate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderItemDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderItemDelete(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderItemDelete(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderItemDelete(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderItemDelete(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderItemReceive(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderItemReceive(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderItemReceive(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "client users cannot receive purchase order items", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "client users cannot receive purchase order items", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderItemReject(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderItemReject(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderItemReject(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "client users cannot receive purchase order items", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "client users cannot receive purchase order items", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderItemUpdateIPAInfo(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderItemUpdateIPAInfo(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderItemUpdateIPAInfo(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "client users cannot receive purchase order items", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "client users cannot receive purchase order items", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderItemScanInput(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderItemScanInput(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderItemScanInput(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "client users cannot receive purchase order items", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "client users cannot receive purchase order items", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END PURCHASE ORDER ITEM ROUTES ** //

// ** PURCHASE ORDER STATUS ROUTES ** //
func PurchaseOrderStatusCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderStatusCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderStatusCreate(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderStatusCreate(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderStatusCreate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func PurchaseOrderStatusList(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderStatusList(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderStatusList(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderStatusList(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderStatusList(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func PurchaseOrderStatusUpdate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderStatusUpdate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderStatusUpdate(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderStatusUpdate(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderStatusUpdate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func PurchaseOrderStatusDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderStatusDelete(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderStatusDelete(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderStatusDelete(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderStatusDelete(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

// ** END PURCHASE ORDER STATUS ROUTES ** //

// ** PURCHASE ORDER NOTES ROUTES ** //
func PurchaseOrderHistoryCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderHistoryCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderHistoryCreate(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderHistoryCreate(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderHistoryCreate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END PURCHASE ORDER NOTES ROUTES ** //

// ** PURCHASE ORDER ATTACHMENTS ROUTES ** //
func PurchaseOrderAttachmentCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderAttachmentCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderAttachmentCreate(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderAttachmentCreate(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderAttachmentCreate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderAttachmentList(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderAttachmentList(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderAttachmentList(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderAttachmentList(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderAttachmentList(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderAttachmentDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderAttachmentDelete(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderAttachmentDelete(w, r)
	case util.ClientAdmin:
		ClientHandlers.PurchaseOrderAttachmentDelete(w, r)
	case util.ClientUser:
		ClientHandlers.PurchaseOrderAttachmentDelete(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END PURCHASE ORDER ATTACHMENTS ROUTES ** //

// ** VENDOR ROUTES **//
func VendorList(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.VendorList(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.VendorList(w, r)
	case util.ClientAdmin:
		ClientHandlers.VendorList(w, r)
	case util.ClientUser:
		ClientHandlers.VendorList(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func VendorCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.VendorCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.VendorCreate(w, r)
	case util.ClientAdmin:
		ClientHandlers.VendorCreate(w, r)
	case util.ClientUser:
		ClientHandlers.VendorCreate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func VendorGet(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.VendorGet(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.VendorGet(w, r)
	case util.ClientAdmin:
		ClientHandlers.VendorGet(w, r)
	case util.ClientUser:
		ClientHandlers.VendorGet(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func VendorUpdate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.VendorUpdate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.VendorUpdate(w, r)
	case util.ClientAdmin:
		ClientHandlers.VendorUpdate(w, r)
	case util.ClientUser:
		ClientHandlers.VendorUpdate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func VendorDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.VendorDelete(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.VendorDelete(w, r)
	case util.ClientAdmin:
		ClientHandlers.VendorDelete(w, r)
	case util.ClientUser:
		ClientHandlers.VendorDelete(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

//** END VENDOR ROUTES **//

// ** WAREHOUSE ROUTES ** //
func WarehouseList(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.WarehouseList(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.WarehouseList(w, r)
	case util.ClientAdmin:
		ClientHandlers.WarehouseList(w, r)
	case util.ClientUser:
		ClientHandlers.WarehouseList(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func WarehouseCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.WarehouseCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.WarehouseCreate(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot create warehouse", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot create warehouse", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func WarehouseGet(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.WarehouseGet(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.WarehouseGet(w, r)
	case util.ClientAdmin:
		ClientHandlers.WarehouseGet(w, r)
	case util.ClientUser:
		ClientHandlers.WarehouseGet(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func WarehouseUpdate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.WarehouseUpdate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.WarehouseUpdate(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot update warehouse", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot update warehouse", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func WarehouseDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.WarehouseDelete(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.WarehouseDelete(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot delete warehouse", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot delete warehouse", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END WAREHOUSE ROUTES ** //

// ** LOCATION ROUTES ** //
func LocationList(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.LocationList(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.LocationList(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot list locations", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot list locations", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func LocationCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.LocationCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.LocationCreate(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot create location", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot create location", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func LocationGet(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.LocationGet(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.LocationGet(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot get location", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot get location", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func LocationUpdate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.LocationUpdate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.LocationUpdate(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot update location", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot update location", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func LocationDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.LocationDelete(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.LocationDelete(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot delete location", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot delete location", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END LOCATION ROUTES ** //

// ** LOCATION TYPE ROUTES ** //
func LocationTypeList(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.LocationTypeList(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.LocationTypeList(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot list location types", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot list location types", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func LocationTypeCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.LocationTypeCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.LocationTypeCreate(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot create location type", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot create location type", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func LocationTypeGet(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.LocationTypeGet(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.LocationTypeGet(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot get location type", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot get location type", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func LocationTypeUpdate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.LocationTypeUpdate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.LocationTypeUpdate(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot update location type", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot update location type", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func LocationTypeDelete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.LocationTypeDelete(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.LocationTypeDelete(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot delete location type", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot delete location type", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END LOCATION TYPE ROUTES ** //

// ** CARRIER ROUTES ** //

func ListCarriers(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListCarriers(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListCarriers(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListCarriers(w, r)
	case util.ClientUser:
		ClientHandlers.ListCarriers(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END CARRIER ROUTES ** //

// ** CARRIER CONNECTION ROUTES ** //

func CreateCarrierConnection(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreateCarrierConnection(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.CreateCarrierConnection(w, r)
	case util.ClientAdmin:
		ClientHandlers.CreateCarrierConnection(w, r)
	case util.ClientUser:
		ClientHandlers.CreateCarrierConnection(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func GetCarrierConnection(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	//TODO need to remove carrier from all possible shipping method

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetCarrierConnection(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetCarrierConnection(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetCarrierConnection(w, r)
	case util.ClientUser:
		ClientHandlers.GetCarrierConnection(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func DisconnectCarrierConnection(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DisconnectCarrierConnection(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.DisconnectCarrierConnection(w, r)
	case util.ClientAdmin:
		ClientHandlers.DisconnectCarrierConnection(w, r)
	case util.ClientUser:
		ClientHandlers.DisconnectCarrierConnection(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func ListCarrierConnections(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListCarrierConnections(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListCarrierConnections(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListCarrierConnections(w, r)
	case util.ClientUser:
		ClientHandlers.ListCarrierConnections(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

// ** END CARRIER CONNECTION ROUTES ** //

// ** BOX ROUTES ** //

func BoxList(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.BoxList(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.BoxList(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot list boxes", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot list boxes", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func BoxCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.BoxCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.BoxCreate(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot create box", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot create box", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func GetBox(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetBox(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetBox(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot get box", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot get box", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func DeleteBox(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DeleteBox(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.DeleteBox(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot delete box", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot delete box", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func UpdateBox(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateBox(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateBox(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot update box", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot update box", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END BOX ROUTES ** //

// ** PICK SESSION ROUTES ** //
func CreatePickSession(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreatePickSession(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.CreatePickSession(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot create pick session", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot create pick session", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func GetActivePickSession(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetActivePickSession(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetActivePickSession(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot get active pick session", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot get active pick session", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func PickSessionSelectItem(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PickSessionSelectItem(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PickSessionSelectItem(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot select item", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot select item", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func PickSessionAssignTote(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PickSessionAssignTote(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PickSessionAssignTote(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot assign tote", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot assign tote", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func PickSessionConfirmTote(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PickSessionConfirmTote(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PickSessionConfirmTote(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot confirm tote", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot confirm tote", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func PickSessionPick(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PickSessionPick(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PickSessionPick(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot pick", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot pick", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func PickSessionComplete(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PickSessionComplete(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PickSessionComplete(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot complete pick session", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot complete pick session", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

// ** END PICK SESSION ROUTES ** //

//** SHIPPING METHOD ROUTES **//

func ListShippingMethods(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListShippingMethods(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListShippingMethods(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListShippingMethods(w, r)
	case util.ClientUser:
		ClientHandlers.ListShippingMethods(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}
}

func GetShippingMethod(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetShippingMethod(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetShippingMethod(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetShippingMethod(w, r)
	case util.ClientUser:
		ClientHandlers.GetShippingMethod(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func UpdateShippingMethod(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateShippingMethod(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateShippingMethod(w, r)
	case util.ClientAdmin:
		ClientHandlers.UpdateShippingMethod(w, r)
	case util.ClientUser:
		ClientHandlers.UpdateShippingMethod(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

//** END SHIPPING METHOD ROUTES **//

// ** STORE ROUTES ** //

func ListStores(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListStores(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListStores(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListStores(w, r)
	case util.ClientUser:
		ClientHandlers.ListStores(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func GetStore(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetStore(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetStore(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetStore(w, r)
	case util.ClientUser:
		ClientHandlers.GetStore(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func ActivateStore(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ActivateStore(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ActivateStore(w, r)
	case util.ClientAdmin:
		ClientHandlers.ActivateStore(w, r)
	case util.ClientUser:
		ClientHandlers.ActivateStore(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func DeactivateStore(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DeactivateStore(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.DeactivateStore(w, r)
	case util.ClientAdmin:
		ClientHandlers.DeactivateStore(w, r)
	case util.ClientUser:
		ClientHandlers.DeactivateStore(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func UpdateStore(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateStore(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateStore(w, r)
	case util.ClientAdmin:
		ClientHandlers.UpdateStore(w, r)
	case util.ClientUser:
		ClientHandlers.UpdateStore(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func DeleteStore(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DeleteStore(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.DeleteStore(w, r)
	case util.ClientAdmin:
		ClientHandlers.DeleteStore(w, r)
	case util.ClientUser:
		ClientHandlers.DeleteStore(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

// ** END STORE ROUTES ** //

// ** ORDER ROUTES ** //
func ListOrders(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListOrders(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListOrders(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListOrders(w, r)
	case util.ClientUser:
		ClientHandlers.ListOrders(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func GetOrder(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetOrder(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetOrder(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetOrder(w, r)
	case util.ClientUser:
		ClientHandlers.GetOrder(w, r)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

// ** END ORDER ROUTES ** //

// ** CLIENT ROUTES ** //

func ClientCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ClientCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ClientCreate(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user does not have access", http.StatusUnauthorized)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access", http.StatusUnauthorized)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func GetClient(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetClient(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetClient(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user does not have access", http.StatusUnauthorized)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access", http.StatusUnauthorized)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func ListClients(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)

	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListClients(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListClients(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user does not have access", http.StatusUnauthorized)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access", http.StatusUnauthorized)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func UpdateClient(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)

	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateClient(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateClient(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user does not have access", http.StatusUnauthorized)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access", http.StatusUnauthorized)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

func UpdateClientAvatar(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)

	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateClientAvatar(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateClientAvatar(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user does not have access", http.StatusUnauthorized)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access", http.StatusUnauthorized)
	default:
		util.ErrorResponse(w, "invalid User", http.StatusUnauthorized)
	}

}

// ** END CLIENT ROUTES ** //

// ** SHIPPING ROUTES ** //

func ShippingScanTote(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ShippingScanTote(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ShippingScanTote(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func ShippingGetPickSessionOrder(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ShippingGetPickSessionOrder(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ShippingGetPickSessionOrder(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func ShippingShopRates(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ShippingShopRates(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ShippingShopRates(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func ShippingSelectRate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ShippingSelectRate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ShippingSelectRate(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func ShippingPurchaseLabel(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ShippingPurchaseLabel(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ShippingPurchaseLabel(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END SHIPPING ROUTES ** //

// ** RECEIVING ROUTES ** //

func ReceivingListItems(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ReceivingListItems(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ReceivingListItems(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func PurchaseOrderItemGetReceivingDetails(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.PurchaseOrderItemGetReceivingDetails(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.PurchaseOrderItemGetReceivingDetails(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "User does not have access to this action", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END RECEIVING ROUTES ** //

//** PRODUCT LOT ROUTES **//

func ProductLotListByProduct(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ProductLotListByProduct(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ProductLotListByProduct(w, r)
	case util.ClientAdmin:
		ClientHandlers.ProductLotListByProduct(w, r)
	case util.ClientUser:
		ClientHandlers.ProductLotListByProduct(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func ProductLotCreate(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ProductLotCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ProductLotCreate(w, r)
	// case util.ClientAdmin:
	// 	ClientHandlers.ProductLotCreate(w, r)
	// case util.ClientUser:
	// 	ClientHandlers.ProductLotCreate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

//** END PRODUCT LOT ROUTES **//

// ** PRODUCT ALIAS ROUTES **//
func ProductAliasCreate(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ProductAliasCreate(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ProductAliasCreate(w, r)
	case util.ClientAdmin:
		ClientHandlers.ProductAliasCreate(w, r)
	case util.ClientUser:
		ClientHandlers.ProductAliasCreate(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func ProductAliasGetByBarcode(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ProductAliasGetByBarcode(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ProductAliasGetByBarcode(w, r)
	case util.ClientAdmin:
		ClientHandlers.ProductAliasGetByBarcode(w, r)
	case util.ClientUser:
		ClientHandlers.ProductAliasGetByBarcode(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func ProductAliasUpdateByBarcode(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ProductAliasUpdateByBarcode(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ProductAliasUpdateByBarcode(w, r)
	case util.ClientAdmin:
		ClientHandlers.ProductAliasUpdateByBarcode(w, r)
	case util.ClientUser:
		ClientHandlers.ProductAliasUpdateByBarcode(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func ProductAliasDeleteByBarcode(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ProductAliasDeleteByBarcode(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ProductAliasDeleteByBarcode(w, r)
	case util.ClientAdmin:
		ClientHandlers.ProductAliasDeleteByBarcode(w, r)
	case util.ClientUser:
		ClientHandlers.ProductAliasDeleteByBarcode(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

//** END PRODUCT ALIAS ROUTES **//
