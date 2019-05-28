package appcontext

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type AppContext interface {
	context.Context
	CleanupDone()
	RegisterCleanup(cleanupFn func())
}

type appContext struct {
	context.Context
	wg *sync.WaitGroup
}

// Key is reserved for app-wide value access. Do not use it if you want private type access
// Ex: Key("thisKey")
type Key string

func (sc appContext) CleanupDone() {
	<- sc.Done()
	sc.wg.Wait()
}

func (sc appContext) RegisterCleanup(cleanupFn func()) {
	sc.wg.Add(1)
	go sc.cleanup(cleanupFn)
}

func (sc appContext) cleanup(cleanupFn func()) {
	<- sc.Done()
	if cleanupFn != nil {
		cleanupFn()
	}
	sc.wg.Done()
}

func WithValue(sc AppContext, key, val interface{}) AppContext {
	if key == nil {
		panic("nil key")
	}

	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	return &valueSyncContext{sc, key, val}
}

func NewSyncContext(parent context.Context) (ctx AppContext, cancel func()) {
	cancelCtx, cancel := context.WithCancel(parent)

	ctx = appContext{
		Context: cancelCtx,
		wg: &sync.WaitGroup{},
	}

	return
}

type valueSyncContext struct {
	AppContext
	key, val interface{}
}

func (sc *valueSyncContext) String() string {
	return fmt.Sprintf("%v.WithValue(%#v, %#v)", sc.AppContext, sc.key, sc.val)
}

func (sc *valueSyncContext) Value(key interface{}) interface{} {
	if sc.key == key {
		return sc.val
	}
	return sc.AppContext.Value(key)
}