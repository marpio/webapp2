package handlers

import (
	"encoding/json"
	"github.com/marpio/webapp2/server/models"
	"github.com/marpio/webapp2/server/services/orderService"
	vm "github.com/marpio/webapp2/server/viewmodels"
	"net/http"
	"strconv"
)

func (env handlersCtx) GetHome(w http.ResponseWriter, r *http.Request) {
	data := struct {
		CurrentUser             *models.UserRow
		ShoppingCartItemsNumber int
	}{
		env.getCurrentUser(r),
		env.getShoppingCartItemsCount(r),
	}
	env.HandlerContext.GetRenderer().HTML(w, http.StatusOK, "home", data)
}

func (env handlersCtx) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := env.HandlerContext.GetDatastore().AllProducts()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		CurrentUser             *models.UserRow
		Products                []*models.ProductRow
		ShoppingCartItemsNumber int
	}{
		env.getCurrentUser(r),
		products,
		env.getShoppingCartItemsCount(r),
	}
	env.HandlerContext.GetRenderer().HTML(w, http.StatusOK, "products", data)
}

func (env handlersCtx) GetServiceProviders(w http.ResponseWriter, r *http.Request) {
	session, _ := env.HandlerContext.GetCookieStore().Get(r, "holzrepublic-session")
	var user *models.UserRow
	if session.Values["user"] != nil {
		user = session.Values["user"].(*models.UserRow)
	}
	var order *vm.OrderDto
	if session.Values["order"] != nil {
		order = session.Values["order"].(*vm.OrderDto)
	}
	products, err := env.HandlerContext.GetDatastore().AllProducts()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		CurrentUser             *models.UserRow
		Products                []*models.ProductRow
		ShoppingCartItemsNumber int
	}{
		user,
		products,
		orderService.GetItemsCount(order),
	}
	env.HandlerContext.GetRenderer().HTML(w, http.StatusOK, "products", data)
}

func (env handlersCtx) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	// session, _ := env.HandlerContext.GetCookieStore().Get(r, "holzrepublic-session")
	// order := session.Values["order"].(*vm.OrderDto)

	http.Redirect(w, r, "/products", 302)
}

func (env handlersCtx) GetShoppingCart(w http.ResponseWriter, r *http.Request) {
	order := env.getOrder(r)
	data := struct {
		CurrentUser             *models.UserRow
		ShoppingCartItemsNumber int
		Order                   *vm.OrderDto
		TotalPrice              string
	}{
		env.getCurrentUser(r),
		env.getShoppingCartItemsCount(r),
		order,
		strconv.FormatFloat(orderService.GetTotalPrice(order), 'f', 2, 32),
	}
	env.HandlerContext.GetRenderer().HTML(w, http.StatusOK, "shoppingcart", data)
}

/*
func (env *Env) SearchBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	bookmarks, err := env.db.SearchBookmarks(q)
	if checkErrorAndWriteToRes(err, w) {
		return
	}
	bs, err := json.Marshal(bookmarks)
	if checkErrorAndWriteToRes(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

*/
func (env handlersCtx) GetShoppingCartItems(w http.ResponseWriter, r *http.Request) {
	order := env.getOrder(r)
	if order == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(nil)
		return
	}
	orderJson, _ := json.Marshal(order)

	w.Header().Set("Content-Type", "application/json")
	w.Write(orderJson)
}

func (env handlersCtx) DeleteShoppingCartItem(w http.ResponseWriter, r *http.Request) {
	session, _ := env.HandlerContext.GetCookieStore().Get(r, "holzrepublic-session")
	order := env.getOrder(r)
	if order == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(nil)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var m map[string]string
	err := decoder.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	orderService.RemoveProductOrder(order, m["productID"])
	session.Values["order"] = order
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	orderJson, _ := json.Marshal(order)
	w.Write(orderJson)
}

func (env handlersCtx) CreateProductOrder(w http.ResponseWriter, r *http.Request) {
	session, _ := env.HandlerContext.GetCookieStore().Get(r, "holzrepublic-session")
	order := env.getOrder(r)

	productID := r.FormValue("ID")
	productName := r.FormValue("Name")
	price, _ := strconv.ParseFloat(r.FormValue("Price"), 2)
	orderedAmount := r.FormValue("Amount")
	amount, err := strconv.ParseInt(orderedAmount, 10, 0)
	po := vm.ProductOrderDto{
		ProductID:    productID,
		ProductName:  productName,
		ProductPrice: price,
		Amount:       int(amount),
	}
	if order == nil {
		order = vm.NewOrderDto()
	}
	orderService.AddProductOrder(order, po)

	session.Values["order"] = order
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/products", 302)
}
