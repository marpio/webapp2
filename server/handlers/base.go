package handlers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/marpio/webapp2/server/models"
	"github.com/marpio/webapp2/server/services/orderService"
	vm "github.com/marpio/webapp2/server/viewmodels"
	"github.com/unrolled/render"
	"net/http"
)

type HandlerContext interface {
	GetLogger() *log.Logger
	GetName() string
	GetCookieStore() *sessions.CookieStore
	GetRenderer() *render.Render
	GetDatastore() models.Datastore
}

type handlersCtx struct {
	HandlerContext
}

func (env handlersCtx) getOrder(r *http.Request) *vm.OrderDto {
	session, _ := env.HandlerContext.GetCookieStore().Get(r, "holzrepublic-session")

	var order *vm.OrderDto
	if session.Values["order"] != nil {
		order = session.Values["order"].(*vm.OrderDto)
	}
	return order
}

func (env handlersCtx) getShoppingCartItemsCount(r *http.Request) int {
	var order = env.getOrder(r)
	return orderService.GetItemsCount(order)
}

func (env handlersCtx) getCurrentUser(r *http.Request) *models.UserRow {
	session, _ := env.HandlerContext.GetCookieStore().Get(r, "holzrepublic-session")
	var user *models.UserRow
	if session.Values["user"] != nil {
		user = session.Values["user"].(*models.UserRow)
	}
	return user
}
