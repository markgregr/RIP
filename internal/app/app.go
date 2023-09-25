package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markgregr/RIP/internal/app/config"
	"github.com/markgregr/RIP/internal/app/dsn"
	"github.com/markgregr/RIP/internal/app/repository"
)

// Application представляет основное приложение.
type Application struct {
    Config       *config.Config
    Router       *mux.Router
    Repository   *repository.Repository
    RequestLimit int
}

// New создает новый объект Application и настраивает его.
func New(ctx context.Context) (*Application, error) {
    // Инициализируйте конфигурацию
    cfg, err := config.NewConfig(ctx)
    if err != nil {
        return nil, err
    }

    // Инициализируйте подключение к базе данных (DB)
    repo, err := repository.New(dsn.FromEnv())
    if err != nil {
        return nil, err
    }

    // Создайте новый маршрутизатор (Router)
    router := mux.NewRouter()

    // Инициализируйте и настройте объект Application
    app := &Application{
        Config: cfg,
        Router: router,
        Repository: repo,
        // Установите другие параметры вашего приложения, если необходимо
    }

    return app, nil
}

// Run запускает приложение.
func (app *Application) Run() error {
    // Настройте обработчики маршрутов
    // Пример:
    // app.Router.HandleFunc("/api/someendpoint", app.handleSomeEndpoint).Methods("GET")

    // Запустите веб-сервер с вашим маршрутизатором
    http.Handle("/", app.Router)

    // Запустите сервер на указанном порту и хосте
    addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
    return http.ListenAndServe(addr, nil)
}
