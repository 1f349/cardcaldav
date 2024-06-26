package cardcaldav

import (
	"context"
	"database/sql"
	"errors"
	"github.com/1f349/cardcaldav/database"
	"github.com/charmbracelet/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type contextKey int

var authCtxKey contextKey = 0

type Context struct {
	UserName string
}

func NewContext(ctx context.Context, a *Context) context.Context {
	return context.WithValue(ctx, authCtxKey, a)
}

func FromContext(ctx context.Context) (*Context, bool) {
	a, ok := ctx.Value(authCtxKey).(*Context)
	return a, ok
}

type ProviderMiddleware interface {
	Middleware(next http.Handler) http.Handler
}

var authError = errors.New("auth context error")

func NewAuth(dbStr string, logger *log.Logger) *Auth {
	dbOpen, err := sql.Open("mysql", dbStr)
	if err != nil {
		logger.Fatal("sql.Open()", "err", err)
	}
	err = dbOpen.Ping()
	if err != nil {
		logger.Fatal("db.Ping()", "err", err)
	}

	dbQueries := database.New(dbOpen)

	return &Auth{DB: dbQueries}
}

type Auth struct {
	DB DbQueries
}

type DbQueries interface {
	GetPasswordHash(ctx context.Context, username string) (string, error)
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.Header().Add("WWW-Authenticate", `Basic realm="1f349/cardcaldav", charset="UTF-8"`)
			http.Error(w, "HTTP auth is required", http.StatusUnauthorized)
			return
		}
		username, accessToken, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Authorization invalid", http.StatusUnauthorized)
			return
		}

		// validate username and password
		if a.ValidateCredentials(r.Context(), username, accessToken) != nil {
			http.Error(w, "Authorization invalid", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(NewContext(r.Context(), &Context{UserName: username}))
		next.ServeHTTP(w, r)
	})
}

func (a *Auth) CurrentUserPrincipal(ctx context.Context) (string, error) {
	authCtx, ok := FromContext(ctx)
	if !ok {
		return "", authError
	}
	return "/" + authCtx.UserName + "/", nil
}

const blfCryptPrefix = "{BLF-CRYPT}"

var errNotBlfCrypt = errors.New("not BLF crypt")

func (a *Auth) ValidateCredentials(ctx context.Context, un, pw string) error {
	hash, err := a.DB.GetPasswordHash(ctx, un)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(hash, blfCryptPrefix) {
		return errNotBlfCrypt
	}
	hash = hash[len(blfCryptPrefix):]
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
