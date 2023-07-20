package handlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/handlers/ClientHandlers"
	"github.com/shipply-io/shipply-io-backend/handlers/OrganizationHandlers"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

// ** USER ROUTES ** //
func GetUserSelf(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetUserSelf(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetUserSelf(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetUserSelf(w, r)
	case util.ClientUser:
		ClientHandlers.GetUserSelf(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetUser(w, r)
	case util.OrganizationUser:
		util.ErrorResponse(w, "user does not have access to other users", http.StatusForbidden)
	case util.ClientAdmin:
		ClientHandlers.GetUser(w, r)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access to other users", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreateUser(w, r)
	case util.OrganizationUser:
		util.ErrorResponse(w, "user does not have access to create users", http.StatusForbidden)
	case util.ClientAdmin:
		ClientHandlers.CreateUser(w, r)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access to create users", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateUser(w, r)
	case util.OrganizationUser:
		util.ErrorResponse(w, "user does not have access to update users", http.StatusForbidden)
	case util.ClientAdmin:
		ClientHandlers.UpdateUser(w, r)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access to update users", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DeleteUser(w, r)
	case util.OrganizationUser:
		util.ErrorResponse(w, "user does not have access to delete users", http.StatusForbidden)
	case util.ClientAdmin:
		ClientHandlers.DeleteUser(w, r)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access to delete users", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateUserPassword(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateUserPassword(w, r)
	case util.ClientAdmin:
		ClientHandlers.UpdateUserPassword(w, r)
	case util.ClientUser:
		ClientHandlers.UpdateUserPassword(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateUserAvatar(w, r)
	case util.OrganizationUser:
		util.ErrorResponse(w, "user does not have access to update avatar", http.StatusForbidden)
	case util.ClientAdmin:
		ClientHandlers.UpdateUserAvatar(w, r)
	case util.ClientUser:
		util.ErrorResponse(w, "user does not have access to update avatar", http.StatusForbidden)
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

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreateProduct(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.CreateProduct(w, r)
	case util.ClientAdmin:
		ClientHandlers.CreateProduct(w, r)
	case util.ClientUser:
		ClientHandlers.CreateProduct(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateProduct(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateProduct(w, r)
	case util.ClientAdmin:
		ClientHandlers.UpdateProduct(w, r)
	case util.ClientUser:
		ClientHandlers.UpdateProduct(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func AddProductImage(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.AddProductImage(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.AddProductImage(w, r)
	case util.ClientAdmin:
		ClientHandlers.AddProductImage(w, r)
	case util.ClientUser:
		ClientHandlers.AddProductImage(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func DeleteProductImage(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DeleteProductImage(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.DeleteProductImage(w, r)
	case util.ClientAdmin:
		ClientHandlers.DeleteProductImage(w, r)
	case util.ClientUser:
		ClientHandlers.DeleteProductImage(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UpdateProductImageOrder(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateProductImageOrder(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateProductImageOrder(w, r)
	case util.ClientAdmin:
		ClientHandlers.UpdateProductImageOrder(w, r)
	case util.ClientUser:
		ClientHandlers.UpdateProductImageOrder(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

// ** END PRODUCT ROUTES ** //

// ** PURCHASE ORDER ROUTES ** //

func GetPurchaseOrder(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetPurchaseOrder(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetPurchaseOrder(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetPurchaseOrder(w, r)
	case util.ClientUser:
		ClientHandlers.GetPurchaseOrder(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func ListPurchaseOrders(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListPurchaseOrders(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListPurchaseOrders(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListPurchaseOrders(w, r)
	case util.ClientUser:
		ClientHandlers.ListPurchaseOrders(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreatePurchaseOrder(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.CreatePurchaseOrder(w, r)
	case util.ClientAdmin:
		ClientHandlers.CreatePurchaseOrder(w, r)
	case util.ClientUser:
		ClientHandlers.CreatePurchaseOrder(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UpdatePurchaseOrder(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdatePurchaseOrder(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdatePurchaseOrder(w, r)
	case util.ClientAdmin:
		ClientHandlers.UpdatePurchaseOrder(w, r)
	case util.ClientUser:
		ClientHandlers.UpdatePurchaseOrder(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func DeletePurchaseOrder(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DeletePurchaseOrder(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.DeletePurchaseOrder(w, r)
	case util.ClientAdmin:
		ClientHandlers.DeletePurchaseOrder(w, r)
	case util.ClientUser:
		ClientHandlers.DeletePurchaseOrder(w, r)
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
func ListPurchaseOrderStatuses(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListPurchaseOrderStatuses(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListPurchaseOrderStatuses(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListPurchaseOrderStatuses(w, r)
	case util.ClientUser:
		ClientHandlers.ListPurchaseOrderStatuses(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

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

func ListPurchaseOrderAttachments(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListPurchaseOrderAttachments(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListPurchaseOrderAttachments(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListPurchaseOrderAttachments(w, r)
	case util.ClientUser:
		ClientHandlers.ListPurchaseOrderAttachments(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

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

func GetVendor(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetVendor(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetVendor(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetVendor(w, r)
	case util.ClientUser:
		ClientHandlers.GetVendor(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func ListVendors(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListVendors(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListVendors(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListVendors(w, r)
	case util.ClientUser:
		ClientHandlers.ListVendors(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func CreateVendor(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreateVendor(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.CreateVendor(w, r)
	case util.ClientAdmin:
		ClientHandlers.CreateVendor(w, r)
	case util.ClientUser:
		ClientHandlers.CreateVendor(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func UpdateVendor(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateVendor(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateVendor(w, r)
	case util.ClientAdmin:
		ClientHandlers.UpdateVendor(w, r)
	case util.ClientUser:
		ClientHandlers.UpdateVendor(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}
}

func DeleteVendor(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DeleteVendor(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.DeleteVendor(w, r)
	case util.ClientAdmin:
		ClientHandlers.DeleteVendor(w, r)
	case util.ClientUser:
		ClientHandlers.DeleteVendor(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

//** END VENDOR ROUTES **//

// ** WAREHOUSE ROUTES ** //
func ListWarehouses(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListWarehouses(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListWarehouses(w, r)
	case util.ClientAdmin:
		ClientHandlers.ListWarehouses(w, r)
	case util.ClientUser:
		ClientHandlers.ListWarehouses(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func GetWarehouse(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetWarehouse(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetWarehouse(w, r)
	case util.ClientAdmin:
		ClientHandlers.GetWarehouse(w, r)
	case util.ClientUser:
		ClientHandlers.GetWarehouse(w, r)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func CreateWarehouse(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreateWarehouse(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.CreateWarehouse(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot create warehouse", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot create warehouse", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func UpdateWarehouse(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateWarehouse(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateWarehouse(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot update warehouse", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot update warehouse", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func DeleteWarehouse(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DeleteWarehouse(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.DeleteWarehouse(w, r)
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
func ListLocations(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListLocations(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListLocations(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot list locations", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot list locations", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func GetLocation(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetLocation(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetLocation(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot get location", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot get location", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func CreateLocation(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreateLocation(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.CreateLocation(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot create location", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot create location", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func UpdateLocation(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateLocation(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateLocation(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot update location", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot update location", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func DeleteLocation(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DeleteLocation(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.DeleteLocation(w, r)
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
func ListLocationTypes(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListLocationTypes(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListLocationTypes(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot list location types", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot list location types", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func GetLocationType(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.GetLocationType(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.GetLocationType(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot get location type", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot get location type", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func CreateLocationType(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreateLocationType(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.CreateLocationType(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot create location type", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot create location type", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func UpdateLocationType(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.UpdateLocationType(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.UpdateLocationType(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot update location type", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot update location type", http.StatusForbidden)
	default:
		util.ErrorResponse(w, "Invalid User", http.StatusUnauthorized)
	}

}

func DeleteLocationType(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.DeleteLocationType(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.DeleteLocationType(w, r)
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

func ListBoxes(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.ListBoxes(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.ListBoxes(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot list boxes", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot list boxes", http.StatusForbidden)
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

func CreateBox(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
		return
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreateBox(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.CreateBox(w, r)
	case util.ClientAdmin:
		util.ErrorResponse(w, "user cannot create box", http.StatusForbidden)
	case util.ClientUser:
		util.ErrorResponse(w, "user cannot create box", http.StatusForbidden)
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

func CreateClient(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusUnauthorized)
	}

	switch user.GetRole() {
	case util.OrganizationAdmin:
		OrganizationHandlers.CreateClient(w, r)
	case util.OrganizationUser:
		OrganizationHandlers.CreateClient(w, r)
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
