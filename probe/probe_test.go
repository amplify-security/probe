package probe

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestProbe_Working(t *testing.T) {
	log := zerolog.New(os.Stdout)
	ctx := context.Background()
	w := make(chan Runner)
	p := NewProbe(ctx, w, &log)
	assert.False(t, p.Working(), "p.Working == false")
	p.Work()
	assert.True(t, p.Working(), "p.Working == true")
	p.Stop(true)
	assert.False(t, p.Working(), "p.Working == false")
}
