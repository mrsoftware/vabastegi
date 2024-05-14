# :unicorn: Vabastegi [![GoDoc](https://pkg.go.dev/badge/github.com/mrsoftware/vabastegi)](https://pkg.go.dev/github.com/mrsoftware/vabastegi) [![Github release](https://img.shields.io/github/release/mrsoftware/vabastegi.svg)](https://github.com/mrsoftware/vabastegi/releases)
> [!WARNING]  
> Be aware this package is not ready for production usage and is in the development phase.

Vabastegi (in farsi) is a dependency management system for Go.

**Benefits**

- Eliminate globals: the Vabastegi helps you remove global state from your application.
  No more `init()` or global variables.
- Code reuse: the Vabastegi lets teams within your organization build loosely-coupled
  and well-integrated shareable components.
- Battle tested: the Vabastegi is the backbone of all Go services at Snapp! Cab Pricing Team.

See our [docs](#) to get started and/or
learn more about Vabastegi.

## Installation

Use Go modules to install the Vabastegi in your application.

```shell
go get github.com/mrsoftware/vabastegi
```

## Getting started

```go
package main

import (
    "context"
  
    "github.com/gofiber/fiber/v2"
    "github.com/mrsoftware/vabastegi"
)

type Config struct {
    Port        int64
    LoggerLevel int
}

type Logger struct {
    Level int
}

func (receiver *Logger) Sync() error {
    // sync buffered logs into output
    return nil
}

type Hub struct {
    Config *Config
    Logger *Logger
    Fiber  *fiber.App
}

// Providers normally go to `internal/app` directory and a file for each group of providers.
// The main.go file will only contain the main function.

// ProvideConfig create and load config and store it into Hub.
func ProvideConfig(ctx context.Context, application *vabastegi.App[Hub]) error {
    application.Hub.Config = &Config{}
  
    return nil
}

// ProvideLogger create logger and store it into Hub.
// This provider is depending on Config.
func ProvideLogger(ctx context.Context, application *vabastegi.App[Hub]) error {
    application.Hub.Logger = &Logger{Level: application.Hub.Config.LoggerLevel}
  
    // you can register a callback and will call when application is getting shutdown.
    application.OnShutdown(func(ctx context.Context) error {
      return application.Hub.Logger.Sync()
    })
  
    return nil
}

// ProvideFiber create Fiber app (server) and store it into Hub.
// We break http server into (Create/Register/Serve) to each part work independently. 
func ProvideFiber(ctx context.Context, application *vabastegi.App[Hub]) error {
    application.Hub.Fiber = fiber.New(fiber.Config{})
  
    return nil
}

// RegisterObservabilityRoutes we create a passage provider so each handler can have its own Provider.
func RegisterObservabilityRoutes(ctx context.Context, application *vabastegi.App[Hub]) error {
    application.Hub.Fiber.Get("/ping", func(ctx *fiber.Ctx) error { return ctx.SendString("pong") })
  
    return nil
}

// StartServer is normally the final step of a server and all routes/middlewares registered before the ServerStart.
func StartServer(ctx context.Context, application *vabastegi.App[Hub]) error {
    // register a callback for application shutdown, to shut down fiber.
    application.OnShutdown(application.Hub.Fiber.ShutdownWithContext)
  
    // go run background tack, you need to call RunTask, and this will let the application be aware of background tasks and will wait for them.
    application.RunTask(func() {
      if err := application.Hub.Fiber.Listen(application.Hub.Config.Port); err != nil {
        application.Shutdown(ctx, err.Error())
      }
    })
  
    return nil
}

func main() {
	// Creating list or application provider.
    providers := []vabastegi.Provider[Hub]{
      // list of your application dependency provider like config and logger, http server, repository, service creation, for each you create a provider
      ProvideConfig, ProvideLogger, ProvideFiber, RegisterObservabilityRoutes, StartServer,
    }
  
    // creating an application with graceful shutdown, after received an interrupt signal, the application will call all shutdown callbacks (check logger)
    application := vabastegi.New[app.Hub](vabastegi.WithGraceFullShutdown(true))
  
    // pass all providers and to build the application.
    if err := application.Builds(context.TODO(), providers...); err != nil {
      application.Logger().Errorf("creating application: %s", err)
  
      return
    }
  
    // we will wait for all background tasks to complete.
    if err := application.Wait(); err != nil {
      application.Logger().Errorf("application is finished: %s", err)
  
      return
    }
  
    application.Logger().Infof("application is finished successfully!")
}

```

for mode details, check the [documentation](https://godoc.org/github.com/mrsoftware/vabastegi)


## Roadmap
- [ ] Complete README.md file
- [ ] Handle routine errors
- [ ] Unit test