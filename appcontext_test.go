package appcontext

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSyncContext(t *testing.T) {
	t.Run("AppContext will wait with no cleanup", testSyncContextWaitNoCleanup())
	t.Run("AppContext will wait with cleanup", testSyncContextWaitCleanup())
}

func testSyncContextWaitCleanup() func(*testing.T) {
	return func(t *testing.T) {
		hasRun := false
		ctx, cancel := NewSyncContext(context.Background())
		ctx.RegisterCleanup(func() {
			hasRun = true
		})
		go cancel()

		ctx.CleanupDone()
		assert.Equal(t, hasRun, true, "AppContext should wait for cleanup before finishing")
	}
}

func testSyncContextWaitNoCleanup() func(*testing.T) {
	return func(t *testing.T) {
		ctx, cancel := NewSyncContext(context.Background())

		hasRun := false
		go func() {
			hasRun = true
			go cancel()
		}()

		ctx.CleanupDone()
		assert.Equal(t, hasRun, true, "AppContext should wait for cancel before finishing")
	}
}

func TestWithValue(t *testing.T) {
	ctx, cancel := NewSyncContext(context.Background())
	vCtx := WithValue(ctx, "test", 150)
	assert.Equal(t, vCtx.Value("test").(int) == 150, true, "valueSyncContext should have a key test with a value 150")
	assert.Nil(t, ctx.Value("test"), "Original AppContext should not have a value for the key test")

	cancel()
	assert.EqualError(t, ctx.Err(), context.Canceled.Error(), "valueSyncContext should error when the original context is cancelled")
}