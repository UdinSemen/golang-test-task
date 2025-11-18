package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/UdinSemen/golang-test-task/internal/config"
	"github.com/UdinSemen/golang-test-task/pkg/closer"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	shutdownTimeout = time.Duration(5) * time.Second
	timeReadHeader  = time.Duration(5) * time.Second
)

type App struct {
	cfg *config.App

	appCtx       context.Context
	appCtxCancel context.CancelFunc

	serviceProvider *serviceProvider
	httpServer      *gin.Engine
}

func New(ctx context.Context) (*App, error) {
	const op = "app.New"

	ctx, cancel := context.WithCancel(ctx)
	app := &App{
		appCtx:       ctx,
		appCtxCancel: cancel,
	}
	if err := app.initDeps(ctx); err != nil {
		return nil, fmt.Errorf("%s.%s: %w", op, "init", err)
	}

	return app, nil
}

func (a *App) Run() error {
	const op = "app.Run"

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		defer func() {
			log.Println("Stopping http server...")
		}()
		if err := a.runHttpServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("%s.%s: %s", op, "http server", err)
		}
	})

	<-stop
	log.Println("Received shutdown signal")
	a.appCtxCancel()
	closer.CloseAll()
	time.Sleep(shutdownTimeout)

	closer.Wait()

	wg.Wait()
	log.Println("Application finished")

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	const op = "app.initDeps"

	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHttpServer,
		a.initRoutes,
	}

	for _, fun := range inits {
		if err := fun(ctx); err != nil {
			return fmt.Errorf("%s.%s: %w", op, "init", err)
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	const op = "app.initConfig"

	pathConfig := os.Getenv("PATH_CONFIG")
	if pathConfig == "" {
		return errors.New("PATH_CONFIG is not set")
	}

	if _, err := os.Stat(pathConfig); os.IsNotExist(err) {
		return fmt.Errorf("%s.%s: config file does not exist: %s", op, "config file", pathConfig)
	}

	var cfg config.App
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return fmt.Errorf("%s.%s: %w", op, "env", err)
	}

	if err := cleanenv.ReadConfig(pathConfig, &cfg); err != nil {
		return fmt.Errorf("%s.%s: %w", op, "file", err)
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("%s.%s: %w", op, "config.Validate", err)
	}

	a.cfg = &cfg

	return nil
}

// initServiceProvider initializes the serviceProvider field of the App struct using a new serviceProvider instance.
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.cfg)
	return nil
}

// initHttpServer initializes the HTTP server with middleware, hides banners, and sets up tracing, logging, and metrics.
func (a *App) initHttpServer(_ context.Context) error {
	a.httpServer = gin.New()

	a.httpServer.Use(gin.Recovery())
	a.httpServer.Use(gin.Logger())

	return nil
}

// initRoutes sets up application routes, registers controllers, and logs the registered routes for debugging purposes.
func (a *App) initRoutes(_ context.Context) error {
	a.httpServer.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	numeratorV1 := a.httpServer.Group("v1/numerator")
	{
		a.serviceProvider.NumeratorRouter().RegisterRoutes(numeratorV1)
	}

	return nil
}

func (a *App) runHttpServer() error {
	log.Println("starting http server",
		"address", fmt.Sprintf("http://%s", a.cfg.Connectors.HTTPServer.GetAddress()))
	srv := &http.Server{
		Addr:              a.cfg.Connectors.HTTPServer.GetAddress(),
		Handler:           a.httpServer,
		ReadHeaderTimeout: timeReadHeader,
	}

	closer.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return srv.Shutdown(ctx)
	})

	return srv.ListenAndServe()
}
