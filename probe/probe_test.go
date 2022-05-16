package probe

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestProbe_Working(t *testing.T) {
	var test bool
	runner := func() {
		test = true
	}
	log := zerolog.New(os.Stdout)
	ctx := context.Background()
	w := make(chan Runner)
	p := NewProbe(ctx, w, &log)
	assert.True(t, p.Working(), "NewProbe -> p.Working == true")
	w <- runner
	p.Stop(true)
	assert.False(t, p.Working(), "p.Stop -> p.Working == false")
	assert.True(t, test, "test == true")
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
