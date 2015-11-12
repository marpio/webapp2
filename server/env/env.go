package env

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/marpio/webapp2/server/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/unrolled/render"
	"html/template"
	"os"
	"strconv"
)

type Env struct {
	Logger      *log.Logger
	Datastore   models.Datastore
	Renderer    *render.Render
	CookieStore *sessions.CookieStore
}

func (env *Env) GetLogger() *log.Logger {
	return env.Logger
}
func (env *Env) GetName() string {
	return "webapp"
}
func (env *Env) GetCookieStore() *sessions.CookieStore {
	return env.CookieStore
}
func (env *Env) GetRenderer() *render.Render {
	return env.Renderer
}
func (env *Env) GetDatastore() models.Datastore {
	return env.Datastore
}

func NewLogger() *log.Logger {
	var log = &log.Logger{
		Out:       os.Stdout,
		Formatter: new(log.JSONFormatter),
		Level:     log.DebugLevel,
	}
	return log
}
func multiply(a int, b float64) string {
	return strconv.FormatFloat(float64(a)*b, 'f', 2, 32)
}

var funcMap = template.FuncMap{
	"multiply": multiply,
}

func NewEnv() *Env {
	logger := NewLogger()
	renderer := render.New(render.Options{
		Layout:        "layout",
		Funcs:         []template.FuncMap{funcMap},
		IsDevelopment: true,
	})
	db := sqlx.MustConnect("sqlite3", "./holzrepublic.db")
	cookieStore := sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(128)))
	env := &Env{
		Logger:      logger,
		Datastore:   models.NewDatastore(db),
		Renderer:    renderer,
		CookieStore: cookieStore,
	}
	return env
}
