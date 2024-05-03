package vabastegi

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-runewidth"
)

// Version of current Vabastegi package
const Version = "0.1.0"

// Provider is what a pkgs use to create dependency.
type Provider[T any] func(context.Context, *App[T]) error

// App is the dependency injection manger.
type App[T any] struct {
	waitGroup            sync.WaitGroup
	errList              []error
	onShutdown           []func(ctx context.Context) error
	Hub                  T
	options              Options
	backgroundTasksCount int
}

// New instance of App Dependency management.
func New[T any](options ...Option) *App[T] {
	app := App[T]{options: Options{}}
	app.UpdateOptions(options...)

	return &app
}

// UpdateOptions is used if you want to change any options for App.
func (a *App[T]) UpdateOptions(options ...Option) {
	for _, option := range options {
		option(&a.options)
	}

	if a.options.GracefulShutdown {
		a.registerGracefulShutdown()
	}

	// app required a logger.
	if a.options.Logger == nil {
		a.options.Logger = NewIOLogger(os.Stdout, InfoLogLevel)
	}
}

// Builds the dependency structure of your app.
func (a *App[T]) Builds(ctx context.Context, providers ...Provider[T]) error {
	startAt := time.Now()
	for _, provider := range providers {
		if err := a.Build(ctx, provider); err != nil {
			return err
		}
	}

	a.banner(len(providers), time.Since(startAt))

	return nil
}

// Build use the provider to set a dependency.
func (a *App[T]) Build(ctx context.Context, provider Provider[T]) (err error) {
	logMessage := a.getProviderName(provider)

	defer func() {
		if err != nil {
			logMessage = logMessage + " ✕"
		} else {
			logMessage = logMessage + " ✓"
		}

		a.options.Logger.Infof(logMessage)
	}()

	return provider(ctx, a)
}

// Logger of application.
func (a *App[T]) Logger() Logger {
	return a.options.Logger
}

// RunTask in background.
func (a *App[T]) RunTask(fn func()) {
	go fn()

	a.backgroundTasksCount++

	a.waitGroup.Add(1)
}

// Wait for background task to done or any shutdown signal.
func (a *App[T]) Wait() error {
	a.waitGroup.Wait()

	// todo: handle merging and returning all errors.
	if len(a.errList) != 0 {
		return a.errList[0]
	}

	return nil
}

// Shutdown ths application.
func (a *App[T]) Shutdown(ctx context.Context, reason string) {
	a.options.Logger.Infof("Shutting down( %s ) ...", reason)

	for i := len(a.onShutdown) - 1; i >= 0; i-- {
		if err := a.onShutdown[i](ctx); err != nil {
			a.errList = append(a.errList, err)
		}
	}

	for i := 0; i < a.backgroundTasksCount; i++ {
		a.waitGroup.Done()
	}
}

// OnShutdown register any method for Shutdown method.
func (a *App[T]) OnShutdown(fn func(ctx context.Context) error) {
	a.onShutdown = append(a.onShutdown, fn)
}

func (a *App[T]) getProviderName(creator interface{}) string {
	parts := strings.Split(runtime.FuncForPC(reflect.ValueOf(creator).Pointer()).Name(), ".")

	return parts[len(parts)-1]
}

func (a *App[T]) registerGracefulShutdown() {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	go func() {
		appSignal := <-interruptChan
		a.Shutdown(context.Background(), appSignal.String())
	}()
}

func (a *App[T]) banner(providersCount int, duration time.Duration) {
	const (
		colorBlack = "\u001b[90m"
		colorCyan  = "\u001b[96m"
		colorReset = "\u001b[0m"
	)

	centerValue := func(s string, width int) string {
		const padDiv = 2
		pad := strconv.Itoa((width - runewidth.StringWidth(s)) / padDiv)
		str := fmt.Sprintf("%"+pad+"s", " ")
		str += fmt.Sprintf("%s%s%s", colorCyan, s, colorBlack)
		str += fmt.Sprintf("%"+pad+"s", " ")
		if runewidth.StringWidth(s)-10 < width && runewidth.StringWidth(s)%2 == 0 {
			// add an ending space if the length of str is even and str is not too long
			str += " "
		}
		return str
	}

	value := func(s string, width int) string {
		pad := width - len(s)
		str := ""
		for i := 0; i < pad; i++ {
			str += "."
		}
		if s == "Disabled" {
			str += " " + s
		} else {
			str += fmt.Sprintf(" %s%s%s", colorCyan, s, colorBlack)
		}
		return str
	}

	const lineLen = 49
	mainLogo := colorBlack + " ┌───────────────────────────────────────────────────┐\n"
	if a.options.AppName != "" {
		mainLogo += " │ " + centerValue(a.options.AppName, lineLen) + " │\n"
	}
	mainLogo += " │ " + centerValue("Vabastegi v"+Version, lineLen) + " │\n"
	mainLogo += " │ " + centerValue("Dependency Manager", lineLen) + " │\n"

	mainLogo += fmt.Sprintf(
		" │                                                   │\n"+
			" │ Providers %s │\n"+
			" │ Duration %s  │\n"+
			" └───────────────────────────────────────────────────┘"+
			colorReset, value(strconv.Itoa(providersCount), 38),
		value(duration.String(), 38),
	)

	fmt.Fprint(os.Stdout, mainLogo)
}
