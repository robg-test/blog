package internal

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
)

var UserSessionManager *scs.SessionManager

func SetupSessionManager() error {
	log.Println("Setting up session manager")
	opt, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)

	UserSessionManager = scs.New()
	UserSessionManager.Lifetime = 24 * time.Hour
	UserSessionManager.Cookie.Secure = true
	UserSessionManager.Cookie.SameSite = http.SameSiteStrictMode
	UserSessionManager.Store = goredisstore.New(client)

	return nil
}

func PutMessage(key string, value string, r *http.Request) {
	UserSessionManager.Put(r.Context(), key, value)
}

func GetMessage(key string, r *http.Request) string {
	msg := UserSessionManager.GetString(r.Context(), key)
	return msg
}

func ClearSession(r *http.Request) {
	UserSessionManager.Destroy(r.Context())
}
