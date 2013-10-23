package potato

import (
    "os"
    "log"
    "fmt"
    "net/http"
)

var (
    AppName = "a potato application"
    Version = "0.0.1"
    Env     = "prod"

    Host    = "localhost"
    Port    = 80
    Timeout = 30

    Dir     = &AppDirStruct{
        Config:     "./config/",
        Controller: "./controller/",
        Model:      "./model/",
        Static:     "./static/",
        Log:        "./log/",
    }

    NotFoundRoute = &Route{
        Controller: "Error",
        Action: "NotFound",
    }

    ServerErrorRoute = &Route{
        Controller: "Error",
        Action: "ServerError",
    }

    R *Router
    S *http.Server

    Logger *log.Logger
)

type AppDirStruct struct {
    Config string
    Controller string
    Model string
    Static string
    Log string
}

func Init() {
    //initialize config
    config := new(Tree)
    if e := LoadYaml(&config.Data, Dir.Config + "config.yml"); e != nil {
        log.Fatal(e)
    }

    if name, ok := config.String("name"); ok {
        AppName = name
    }

    if env, ok := config.String("env"); ok {
        Env = env
    }

    //http config
    if v, ok := config.String("http.host"); ok {
        Host = v
    }
    if v, ok := config.Int("http.port"); ok {
        Port = v
    }
    if v, ok := config.Int("http.timeout"); ok {
        Timeout = v
    }

    //dir config
    if dir, ok := config.String("static_dir"); ok {
        Dir.Static = dir
    }
    if dir, ok := config.String("log_dir"); ok {
        Dir.Log = dir
    }

    //error handlers
    if v, ok := config.String("error_handler.not_found.controller"); ok {
        NotFoundRoute.Controller = v
    }
    if v, ok := config.String("error_handler.not_found.action"); ok {
        NotFoundRoute.Action = v
    }
    if v, ok := config.String("error_handler.server_error.controller"); ok {
        ServerErrorRoute.Controller = v
    }
    if v, ok := config.String("error_handler.server_error.action"); ok {
        ServerErrorRoute.Action = v
    }

    //initialize logger
    file, e := os.OpenFile(Dir.Log + Env + ".log",
            os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
    if e != nil {
        log.Fatal("Error init log file:", e)
    }

    Logger = log.New(file, "", log.LstdFlags)

    //initialize router and load routes config file
    R = NewRouter()
    R.InitConfig(Dir.Config + "routes.yml")

    //initialize server
    S = &http.Server{
        Addr: fmt.Sprintf(":%d", Port),
        Handler: R,
    }

    Logger.Println(fmt.Sprintf("Server ready: %s:%d", Host, Port))
}
