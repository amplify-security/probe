package probe

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	// MockRandom is a mock Random stream for testing.
	MockRandom struct {
		mock.Mock
	}
)

var (
	logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	log        = slog.New(logHandler)
)

// Read implementation of io.Reader for MockRandom.
func (m *MockRandom) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	for _, i := range p {
		p[i] = 0
	}
	if args.Get(1) != nil {
		err = args.Error(1)
	}
	return args.Int(0), err
}

func waitForRunning(p *Probe) {
	for i := 0; i < 10; i++ {
		if p.Running() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func waitForIdle(p *Probe) {
	for i := 0; i < 10; i++ {
		// wait for goroutine to become idle after finishing work
		if p.Idle() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestProbe_Working(t *testing.T) {
	var test bool
	runner := func() {
		test = true
	}
	ctx, cancel := context.WithCancel(context.Background())
	p := NewProbe(&ProbeConfig{
		Ctx:        ctx,
		LogHandler: logHandler,
	})
	waitForRunning(p)
	assert.True(t, p.Running(), "NewProbe -> p.Running == true")
	p.WorkChan() <- runner
	p.Stop(true)
	assert.False(t, p.Running(), "p.Stop -> p.Running == false")
	assert.True(t, test, "test == true")
	p.Run()
	cancel()
	waitForIdle(p)
	assert.False(t, p.Running(), "ctx cancel -> p.Running == false")
}

func TestProbe_Idle(t *testing.T) {
	ctx := context.Background()
	done := make(chan struct{})
	runner := func() {
		<-done
	}
	p := NewProbe(&ProbeConfig{
		Ctx:        ctx,
		LogHandler: logHandler,
	})
	assert.True(t, p.Idle(), "NewProbe -> p.Idle == true")
	p.WorkChan() <- runner
	assert.False(t, p.Idle(), "runner -> p.Idle == false")
	close(done)
	waitForIdle(p)
	assert.True(t, p.Idle(), "close(done) -> p.Idle == true")
	p.Stop(true)
}

func TestProbe_ID(t *testing.T) {
	ctx := context.Background()
	p := NewProbe(&ProbeConfig{
		Ctx:        ctx,
		LogHandler: logHandler,
	})
	assert.NotEmpty(t, p.ID(), "p.ID -> !empty")
}

func TestProbe_Stop(t *testing.T) {
	ctx := context.Background()
	done := make(chan struct{})
	runner := func() {
		<-done
	}
	p := NewProbe(&ProbeConfig{
		Ctx:        ctx,
		LogHandler: logHandler,
	})
	waitForRunning(p)
	assert.True(t, p.Running(), "NewProbe -> p.Running == true")
	p.Stop(true)
	assert.False(t, p.Running(), "p.Stop(true) -> p.Running == false")
	p.Run()
	waitForRunning(p)
	assert.True(t, p.Running(), "p.Run -> p.Running == true")
	p.WorkChan() <- runner
	p.Stop(false)
	assert.True(t, p.Running(), "p.Stop(false) -> p.Running == true")
	close(done)
	p.Stop(true)
	assert.False(t, p.Running(), "p.Stop(false) -> p.Running == false")
	p.Stop(false)
}

func TestProbe_Run(t *testing.T) {
	ctx := context.Background()
	p := NewProbe(&ProbeConfig{
		Ctx:        ctx,
		LogHandler: logHandler,
	})
	old := p.done
	waitForRunning(p)
	p.Run()
	assert.Equal(t, old, p.done, "p.Run && p.Running -> old == new")
	p.Stop(true)
	p.Run()
	waitForRunning(p)
	assert.NotEqual(t, old, p.done, "p.Run && !p.Running -> old != new")
	p.Stop(true)
}

func TestGetID(t *testing.T) {
	id := getID(log)
	assert.NotEmpty(t, id, "id -> !empty")
	m := &MockRandom{}
	m.On("Read", mock.Anything).Return(3, nil)
	randStream = m
	id = getID(log)
	assert.Equal(t, "000000", id, "id -> 000000")
	m = &MockRandom{}
	m.On("Read", mock.Anything).Return(0, errors.New("unexpected error"))
	randStream = m
	assert.Panics(t, func() {
		id = getID(log)
	}, "getID -> panic")
	randStream = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func BenchmarkGetID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getID(log)
	}
}
