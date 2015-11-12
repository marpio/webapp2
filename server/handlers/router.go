package handlers

import (
	"github.com/gorilla/mux"
	"github.com/marpio/webapp2/server/middleware"
	"net/http"
)

func NewConfiguredRouter(env HandlerContext) *mux.Router {
	ctx := &handlersCtx{env}

	router := mux.NewRouter()
	router.HandleFunc("/", ctx.GetHome).Methods("GET")
	router.HandleFunc("/login", ctx.GetLogin).Methods("GET")
	router.HandleFunc("/login", ctx.PostLogin).Methods("POST")
	router.HandleFunc("/signup", ctx.GetSignup).Methods("GET")
	router.HandleFunc("/signup", ctx.PostSignup).Methods("POST")
	router.HandleFunc("/products", ctx.GetProducts).Methods("GET")
	router.HandleFunc("/orders/product", ctx.CreateProductOrder).Methods("POST")
	router.HandleFunc("/shoppingcart", ctx.GetShoppingCart).Methods("GET")
	router.HandleFunc("/shoppingcart/items", ctx.GetShoppingCartItems).Methods("GET")
	router.HandleFunc("/shoppingcart/items", ctx.DeleteShoppingCartItem).Methods("DELETE")

	router.Handle("/shoppingcart", middleware.RequireLogin(http.HandlerFunc(ctx.ConfirmOrder))).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	return router
}
