package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ribice/twisk/pkg/zerolog"

	"github.com/ribice/twisk/internal/openapi"

	"github.com/ribice/twisk/internal/iam"

	"github.com/justinas/alice"

	"github.com/ribice/twisk/pkg/mw"

	"github.com/ribice/twisk/pkg/hooks"

	"github.com/ribice/twisk/pkg/context"

	"github.com/gorilla/mux"

	"github.com/ribice/twisk/internal/iam/rbac"
	"github.com/ribice/twisk/internal/iam/secure"

	"github.com/ribice/twisk/pkg/config"
	"github.com/ribice/twisk/pkg/jwt"
	"github.com/ribice/twisk/pkg/postgres"

	iamdb "github.com/ribice/twisk/internal/iam/platform/postgres"
	iampb "github.com/ribice/twisk/rpc/iam"

	userdb "github.com/ribice/twisk/internal/user/platform/postgres"
	userpb "github.com/ribice/twisk/rpc/user"

	"github.com/ribice/twisk/internal/user"
)

func main() {

	cfgPath := flag.String("p", "./cmd/api/conf.local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	router := mux.NewRouter().StrictSlash(true)
	registerRoutes(router, cfg)

	openapi.New(router, cfg.OpenAPI.Username, cfg.OpenAPI.Password)

	mws := alice.New(mw.CORS, mw.AuthContext)

	http.ListenAndServe(cfg.Server.Port, mws.Then(router))
}

func registerRoutes(r *mux.Router, cfg *config.Configuration) {
	db, err := pgsql.New(cfg.DB.PSN, cfg.DB.LogQueries, cfg.DB.TimeoutSeconds)
	checkErr(err)

	rbacSvc := new(rbac.Service)
	ctxSvc := new(context.Service)
	log := zerolog.New()

	secureSvc := secure.New(cfg.App.MinPasswordStrength)

	j := jwt.New(cfg.JWT.Secret, cfg.JWT.Duration, cfg.JWT.Algorithm)

	userSvc := user.NewLoggingService(
		user.New(db, userdb.NewUser(), rbacSvc, secureSvc, ctxSvc), log)

	r.PathPrefix(userpb.UserPathPrefix).Handler(
		userpb.NewUserServer(userSvc, hooks.WithJWTAuth(j)))

	iamSvc := iam.NewLoggingService(iam.New(db, j, iamdb.NewUser(), secureSvc), log)

	r.PathPrefix(iampb.IAMPathPrefix).Handler(
		iampb.NewIAMServer(iamSvc, nil))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
