package main

import (
	"flag"
	"fmt"
	"os"

	"dazzlefm/route"
	"dazzlefm/service"
	"dazzlefm/util"

	"github.com/go-ini/ini"
	"github.com/gorilla/sessions"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

// AppName : the applications name
const AppName string = "dazzlefm"

func main() {
	hash := flag.String("bcrypt", "", "string to bcrypt")
	flag.Parse()

	if len(*hash) > 0 {
		hashPassword(*hash)
		return
	}

	// configure the global logger
	logger := util.InitLogger()
	defer logger.Sync()

	// get initial config
	cfg, _, err := util.InitConfig(AppName, logger.Sugar())
	if err != nil {
		os.Exit(1)
	}

	// init db connection
	dbc, err := util.InitDb(cfg, logger.Sugar())
	if err != nil {
		return
	}

	// initialize echo
	rtr := initServer(cfg)

	// create env
	env := util.MakeEnv(AppName, rtr, cfg, dbc, logger)

	// start services
	if err := initServices(env); err != nil {
		log.Error(err)

		return
	}

	// mount handlers
	if err := initHandlers(env); err != nil {
		log.Error(err)

		return
	}

	bind := fmt.Sprintf(
		"%s:%s",
		cfg.Section("").Key("host").String(),
		cfg.Section("").Key("port").String(),
	)

	// rtr.Logger.Fatal(rtr.Start(bind))
	if err := rtr.Start(bind); err != nil {
		rtr.Logger.Fatal(err)
	}
}

func initServer(cfg *ini.File) (rtr *echo.Echo) {
	rtr = echo.New()

	// middleware
	lc := middleware.LoggerConfig{
		Format: `[${method}] ${status} - ${uri}` +
			` - ${latency_human}, rx:${bytes_in}, tx:${bytes_out}` + "\n",
	}
	rtr.Use(middleware.LoggerWithConfig(lc))

	origins := cfg.Section("url").Key("origin").Strings(",")
	// fmt.Println(origins)
	rtr.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     origins,
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	rtr.Use(middleware.Recover())

	// rtr.Use(session.Middleware(sessions.NewCookieStore([]byte("!chidinmaisafinegirl!"))))
	rtr.Use(session.Middleware(sessions.NewFilesystemStore("./tmp", []byte("!olusolaisafinegirl!"))))

	return
}

func initServices(env *util.Environment) (err error) {

	env.ServiceList["office"], err = service.NewOffice(env)
	if err != nil {
		return
	}

	env.ServiceList["content"], err = service.NewContent(env)
	if err != nil {
		return
	}

	env.ServiceList["auth"], err = service.NewAuthenticate(env)
	if err != nil {
		return
	}

	return
}

func initHandlers(env *util.Environment) (err error) {

	handlers := []util.IHandler{
		&route.Presenter{},
		&route.News{},
		&route.Gallery{},
		&route.Page{},
		&route.Music{},
		&route.Schedule{},
		&route.Event{},
		&route.User{},
		&route.Contact{},
		&route.Script{},
		&route.Settings{},
	}

	for _, i := range handlers {
		i.Init("/api", env)
	}

	// static folder
	env.Rtr.Static(env.URL("upload"), env.Path("upload"))
	env.Rtr.Static(env.URL("static"), env.Path("static"))
	env.Rtr.Static(env.URL("downloads"), env.Path("downloads"))

	return
}

func hashPassword(hash string) {
	h, _ := util.HashPassword(hash)
	fmt.Println("hash: ", h)
}
