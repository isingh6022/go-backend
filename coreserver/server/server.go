package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/isingh6022/go-fullstack-ai/coreserver/types"
)

type App struct {
	server *http.Server
	mux    *http.ServeMux
	Logger *logger
}

func (a *App) Start() error {
	listener, err := net.Listen("tcp", a.server.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", a.server.Addr, err)
	}

	watchCloseSignalCh := make(chan os.Signal, 1)
	signal.Notify(watchCloseSignalCh, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(watchCloseSignalCh)

	watchServerErrCh := make(chan error, 1)

	startServer := func() {
		if err := a.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			watchServerErrCh <- err
			return
		}
		watchServerErrCh <- nil
	}

	go startServer()
	a.Logger.Info("server listening", "addr", listener.Addr().String())

	select {
	case err := <-watchServerErrCh:
		if err != nil {
			return errors.Join(err, a.Shutdown())
		}
		_ = a.Shutdown()
		return nil
	case sig := <-watchCloseSignalCh:
		signal.Stop(watchCloseSignalCh)
		a.Logger.Info("shutting down", "signal", sig.String())

		if err := a.Shutdown(); err != nil {
			return err
		}
		a.Logger.Info("server shut down gracefully")
		return nil
	}
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), types.DefaultShutdownTimeout)
	defer cancel()

	return a.ShutdownWithCtx(ctx)
}
func (a *App) ShutdownWithCtx(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		a.server.Close()
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}

func (a *App) Handle(pattern string, handler http.Handler) {
	a.mux.Handle(pattern, handler)
}

func (a *App) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	a.mux.HandleFunc(pattern, handler)
}

func CreateApp(port int, logLvl logLevel) *App {
	mux := http.NewServeMux()
	logger := CreateLogger(logLvl)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           getRecoverMiddleware(logger)(mux),
		ReadHeaderTimeout: types.DefaultReadHeaderTimeout,
		ReadTimeout:       types.DefaultReadTimeout,
		WriteTimeout:      types.DefaultWriteTimeout,
		IdleTimeout:       types.DefaultIdleTimeout,
	}

	return &App{
		server: server,
		mux:    mux,
		Logger: logger,
	}
}
