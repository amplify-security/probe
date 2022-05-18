package probe

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	// MockRandom is a mock Random stream for testing.
	MockRandom struct {
		mock.Mock
	}
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

func TestProbe_Working(t *testing.T) {
	var test bool
	runner := func() {
		test = true
	}
	log := zerolog.New(os.Stdout)
	ctx, cancel := context.WithCancel(context.Background())
	w := make(chan Runner)
	p := NewProbe(ctx, w, &log)
	assert.True(t, p.Working(), "NewProbe -> p.Working == true")
	w <- runner
	p.Stop(true)
	assert.False(t, p.Working(), "p.Stop -> p.Working == false")
	assert.True(t, test, "test == true")
	p.Work()
	cancel()
	for i := 0; i < 10; i++ {
		// wait for goroutine to return after cancel
		if !p.Working() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	assert.False(t, p.Working(), "ctx cancel -> p.Working == false")
}

func TestProbe_Idle(t *testing.T) {
	log := zerolog.New(os.Stdout)
	ctx := context.Background()
	w := make(chan Runner)
	done := make(chan struct{})
	runner := func() {
		<-done
	}
	p := NewProbe(ctx, w, &log)
	assert.True(t, p.Idle(), "NewProbe -> p.Idle == true")
	w <- runner
	assert.False(t, p.Idle(), "runner -> p.Idle == false")
	close(done)
	for i := 0; i < 10; i++ {
		// wait for goroutine to become idle after finishing work
		if p.Idle() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	assert.True(t, p.Idle(), "close(done) -> p.Idle == true")
	p.Stop(true)
}

func TestProbe_ID(t *testing.T) {
	log := zerolog.New(os.Stdout)
	ctx := context.Background()
	w := make(chan Runner)
	p := NewProbe(ctx, w, &log)
	assert.NotEmpty(t, p.ID(), "p.ID -> !empty")
}

func TestProbe_Stop(t *testing.T) {
	log := zerolog.New(os.Stdout)
	ctx := context.Background()
	w := make(chan Runner)
	done := make(chan struct{})
	runner := func() {
		<-done
	}
	p := NewProbe(ctx, w, &log)
	assert.True(t, p.Working(), "NewProbe -> p.Working == true")
	p.Stop(true)
	assert.False(t, p.Working(), "p.Stop(true) -> p.Working == false")
	p.Work()
	assert.True(t, p.Working(), "p.Work -> p.Working == true")
	w <- runner
	p.Stop(false)
	assert.True(t, p.Working(), "p.Stop(false) -> p.Working == true")
	close(done)
	p.Stop(true)
	assert.False(t, p.Working(), "p.Stop(false) -> p.Working == false")
	p.Stop(false)
}

func TestProbe_Work(t *testing.T) {
	log := zerolog.New(os.Stdout)
	ctx := context.Background()
	w := make(chan Runner)
	p := NewProbe(ctx, w, &log)
	old := p.done
	p.Work()
	assert.Equal(t, old, p.done, "p.Work && p.Working -> old == new")
	p.Stop(true)
	p.Work()
	assert.NotEqual(t, old, p.done, "p.Work && !p.Working -> old != new")
	p.Stop(true)
}

func TestGetID(t *testing.T) {
	log := zerolog.New(os.Stdout)
	id := getID(&log)
	assert.NotEmpty(t, id, "id -> !empty")
	m := &MockRandom{}
	m.On("Read", mock.Anything).Return(3, nil)
	randStream = m
	id = getID(&log)
	assert.Equal(t, "000000", id, "id -> 000000")
	m = &MockRandom{}
	m.On("Read", mock.Anything).Return(0, errors.New("unexpected error"))
	randStream = m
	assert.Panics(t, func() {
		id = getID(&log)
	}, "getID -> panic")
	randStream = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func BenchmarkGetID(b *testing.B) {
	log := zerolog.New(os.Stdout)
	for i := 0; i < b.N; i++ {
		getID(&log)
	}
}
