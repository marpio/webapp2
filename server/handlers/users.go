package handlers

import (
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"net/http"
)

func (env handlersCtx) GetSignup(w http.ResponseWriter, r *http.Request) {
	env.HandlerContext.GetRenderer().HTML(w, http.StatusOK, "users/signup", nil, render.HTMLOptions{Layout: "users/login-signup-parent"})
}

func (env handlersCtx) PostSignup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("Email")
	password := r.FormValue("Password")
	passwordAgain := r.FormValue("PasswordAgain")

	_, err := env.HandlerContext.GetDatastore().Signup(email, password, passwordAgain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Perform login
	env.PostLogin(w, r)
}

func (env handlersCtx) GetLoginWithoutSession(w http.ResponseWriter, r *http.Request) {
	env.HandlerContext.GetRenderer().HTML(w, http.StatusOK, "users/login", nil, render.HTMLOptions{Layout: "users/login-signup-parent"})
}

func (env handlersCtx) GetLogin(w http.ResponseWriter, r *http.Request) {
	cookieStore := context.Get(r, "cookieStore").(*sessions.CookieStore)
	session, _ := cookieStore.Get(r, "holzrepublic-session")

	currentUserInterface := session.Values["user"]
	if currentUserInterface != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	env.GetLoginWithoutSession(w, r)
}

// PostLogin performs login.
func (env handlersCtx) PostLogin(w http.ResponseWriter, r *http.Request) {
	cookieStore := context.Get(r, "cookieStore").(*sessions.CookieStore)

	email := r.FormValue("Email")
	password := r.FormValue("Password")

	user, err := env.HandlerContext.GetDatastore().GetUserByEmailAndPassword(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := cookieStore.Get(r, "holzrepublic-session")
	session.Values["user"] = user

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func (env handlersCtx) GetLogout(w http.ResponseWriter, r *http.Request) {
	cookieStore := context.Get(r, "cookieStore").(*sessions.CookieStore)

	session, _ := cookieStore.Get(r, "holzrepublic-session")

	delete(session.Values, "user")
	session.Save(r, w)

	http.Redirect(w, r, "/login", 302)
}

// func (env *Env) PostPutDeleteUsersID(w http.ResponseWriter, r *http.Request) {
// 	method := r.FormValue("_method")
// 	if method == "" || strings.ToLower(method) == "post" || strings.ToLower(method) == "put" {
// 		env.PutUsersID(w, r)
// 	} else if strings.ToLower(method) == "delete" {
// 		env.DeleteUsersID(w, r)
// 	}
// }
