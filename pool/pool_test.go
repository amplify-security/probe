package pool

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/amplify-security/probe"
	"github.com/stretchr/testify/assert"
)

var (
	logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
)

func waitForRunning(p *probe.Probe) {
	for i := 0; i < 10; i++ {
		if p.Running() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func waitForNotRunning(p *probe.Probe) {
	for i := 0; i < 10; i++ {
		if !p.Running() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func waitForIdle(p *probe.Probe) {
	for i := 0; i < 10; i++ {
		// wait for goroutine to become idle after finishing work
		if p.Idle() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestPool_Start(t *testing.T) {
	cases := []struct {
		size int
		e    int
		msg  string
	}{
		{
			size: 16,
			e:    16,
			msg:  "Start(dynamic == false) -> p.probes == 16",
		},
	}
	for _, c := range cases {
		p := NewPool(&PoolConfig{
			LogHandler: logHandler,
			Size:       c.size,
		})
		assert.Equal(t, c.e, len(p.probes), c.msg)
		// second start attempt tests to ensure we are not creating additional probes
		p.Start()
		assert.Equal(t, c.e, len(p.probes), c.msg)
		p.Stop(true)
		// third start attempt tests to ensure we are not creating additional probes on restart
		p.Start()
		assert.Equal(t, c.e, len(p.probes), c.msg)
	}
}

func TestPool_Stop(t *testing.T) {
	p := NewPool(&PoolConfig{
		LogHandler: logHandler,
		Size:       16,
	})
	assert.Equal(t, 16, len(p.probes), "Start -> p.probes == 16")
	p.Stop(true)
	assert.Equal(t, 16, len(p.probes), "Stop(wait == true) -> p.probes == 16")
	for i, probe := range p.probes {
		assert.False(t, probe.Running(), fmt.Sprintf("Stop(wait == true) -> probe[%d].Running == false", i))
	}
	p.Start()
	p.Stop(false)
	assert.Equal(t, 16, len(p.probes), "Stop(wait == false) -> p.probes == 16")
	for i, probe := range p.probes {
		waitForNotRunning(probe)
		assert.False(t, probe.Running(), fmt.Sprintf("Stop(wait == false) -> probe[%d].Running == false", i))
	}
}

func TestPool_Run(t *testing.T) {
	ctr := new(atomic.Int32)
	wg := new(sync.WaitGroup)
	wg.Add(16)
	r1 := func() {
		ctr.Add(1)
		wg.Done()
	}
	p := NewPool(&PoolConfig{
		LogHandler: logHandler,
		Size:       16,
	})
	for range 16 {
		p.Run(r1)
	}
	wg.Wait()
	assert.Equal(t, 16, int(ctr.Load()), "Run(r1) -> ctr == 16")
	p.Stop(true)
}

func TestPool_Idle(t *testing.T) {
	p := NewPool(&PoolConfig{
		LogHandler: logHandler,
		Size:       16,
	})
	for _, probe := range p.probes {
		waitForIdle(probe)
	}
	assert.Equal(t, 16, p.Idle(), "NewPool(16) -> p.Idle == 16")
}

func TestPool_Running(t *testing.T) {
	p := NewPool(&PoolConfig{
		LogHandler: logHandler,
		Size:       16,
	})
	for _, probe := range p.probes {
		waitForRunning(probe)
	}
	assert.Equal(t, 16, p.Running(), "NewPool(16) -> p.Running == 16")
}
