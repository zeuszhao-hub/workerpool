package worker

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestWorker normal
func TestWorker(t *testing.T) {
	w := NewWorker()
	w.HandleWork(1, 1, 2*time.Second, func(ctx context.Context, data interface{}) {
		t.Log(ctx.Deadline())
		t.Log(data)
	})
	w.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	w.Process(ctx, "hello worker pool")

	time.Sleep(1 * time.Second)
	w.Shutdown()
}

// TestPipeFull timeout control
func TestPipeFull(t *testing.T) {
	w := NewWorker()
	w.HandleWork(0, 1, 2*time.Second, func(ctx context.Context, data interface{}) {
		// delay 5s
		time.Sleep(5 * time.Second)
		t.Log(ctx.Deadline())
		t.Log(data)
	})
	w.Run()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := w.Process(ctx, "hello worker pool 1")
		if err != nil {
			t.Log(err.Error())
		}
	}()

	// 3s timeout
	wg.Add(1)
	go func() {
		wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := w.Process(ctx, "hello worker pool 2")
		if err != nil {
			t.Log(err.Error())
		}
	}()

	wg.Wait()
	w.Shutdown()
}

// TestWorkerPanic worker panic control
func TestWorkerPanic(t *testing.T) {
	w := NewWorker()
	w.HandleWork(0, 1, 2*time.Second, func(ctx context.Context, data interface{}) {
		panic("panic err")
		t.Log(ctx.Deadline())
		t.Log(data)
	})
	w.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	w.Process(ctx, "test worker panic")
	w.Shutdown()
}

// TestWorkerPanic Woker waits for processing to finish and closes
func TestWaitWorker(t *testing.T) {
	w := NewWorker()
	w.HandleWork(0, 1, 2*time.Second, func(ctx context.Context, data interface{}) {
		time.Sleep(5 * time.Second)
		t.Log(ctx.Deadline())
		t.Log(data)
	})
	w.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	w.Process(ctx, "send message")

	// shutdown worker
	w.Shutdown()

	// 5s after print
	t.Log("worker shutdown")
}
