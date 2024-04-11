package app

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

// Depend is the dependency injection manger.
type Depend[T any] struct {
	waitGroup  sync.WaitGroup
	errList    []error
	onShutdown []func() error
	hub        T
}

// AddWaiting used to add waiting for application.
func (h *Depend[T]) AddWaiting(delta int) {
	h.waitGroup.Add(delta)
}

// Done reduce waite counter, on error store is in error group.
func (h *Depend[T]) Done(err error) {
	h.waitGroup.Done()
}

// Wait for all registered waite group to call Done.
func (h *Depend[T]) Wait() error {
	h.waitGroup.Wait()

	return nil // h.waitGroup.Err()
}

// Shutdown ths application.
func (h *Depend[T]) Shutdown(reason string) {
	log.Default().Printf("Shutting down( %s ) ...", reason)

	for i := len(h.onShutdown) - 1; i >= 0; i-- {
		h.Done(h.onShutdown[i]())
	}
}

// OnShutdown register any method for Shutdown method.
func (h *Depend[T]) OnShutdown(fn func() error) {
	h.AddWaiting(1)

	h.onShutdown = append(h.onShutdown, fn)
}

// ShutdownOnSignal listen on passed channel and shutdown application.
func (h *Depend[T]) ShutdownOnSignal(c <-chan os.Signal) {
	go func() {
		appSignal := <-c
		h.Shutdown(appSignal.String())
	}()
}

func (h *Depend[T]) registerGracefulShutdown() {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	h.ShutdownOnSignal(interruptChan)
}
