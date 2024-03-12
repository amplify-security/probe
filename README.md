# probe

[![Amplify Security](https://github.com/amplify-security/probe/actions/workflows/amplify.yml/badge.svg?branch=main)](https://github.com/amplify-security/probe/actions/workflows/amplify.yml)
[![test](https://github.com/github/docs/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/amplify-security/probe/actions/workflows/test.yml)
![coverage](https://raw.githubusercontent.com/amplify-security/probe/badges/.badges/main/coverage.svg)

[![Go Reference](https://pkg.go.dev/badge/github.com/amplify-security/probe.svg)](https://pkg.go.dev/github.com/amplify-security/probe)

A modern, zero-dependency goroutine pool. Probe is designed to abstract goroutine synchronization
and control flow to enable cleaner concurrent code.

## What can you do with Probe?

You can easily create reusable goroutines called Probes to do side channel work with less synchronization
code. For example:

```go
p := probe.NewProbe(&probe.ProbeConfig{})
p.WorkChan() <-func() {
    fmt.Println("Hello from Probe!")
}
p.Stop()
```

You can also create a Probe Pool and run functions on a configurable pool of goroutines:

```go
p := pool.NewPool(&pool.PoolConfig{ Size: 16 })
p.Run(func() {
    fmt.Println("Hello from Probe Pool!")
})
p.Stop()
```

You can check to see how many Probes in the Pool are idle for monitoring and tuning Pool sizes:

```go
ctrlChan := make(chan struct{})
f := func() {
    <-ctrlChan
}
p := pool.NewPool(&pool.PoolConfig{ Size: 16 })
p.Run(f)
p.Run(f)
// this is a race condition with the Pool work channel, your output may differ
fmt.Println(p.Idle()) // 14
```

Channels can be used to get results with type safety back from your functions:

```go
f, _ := os.Open("/tmp/test")
b := new(bytes.Buffer)
returnChan := make(chan struct{int; error})
p := probe.NewProbe(&probe.ProbeConfig{})
p.WorkChan() <-func() {
    n, err := f.Read(b)
    returnChan <- struct{int; error}{n, err}
}
r := <-returnChan // access with r.int, r.error
```

## Configuration

Some common configuration scenarios for an individual Probe may be passing in a buffered channel
instead of using the default, blocking channel, or passing in a shared context so that all probes
are stopped if the context is canceled.

```go
ctx, cancel := context.WithCancel(context.Background())
work := make(chan probe.Runner, 1024)
p := probe.NewProbe(&probe.ProbeConfig{
    Ctx:      ctx,
    WorkChan: work, 
})
work <- func() {
    fmt.Println("Hello from a buffered Probe!")
}
cancel()
```

Pools may likewise be configured with a `Size`, `Ctx`, and `BufferSize`.

```go
ctx, cancel := context.WithCancel(context.Background())
p := pool.NewPool(&pool.PoolConfig{
    Ctx:        ctx,
    Size:       1024,
    BufferSize: 2048,
})
p.Run(func() {
    fmt.Println("Hello from a custom Pool!")
})
// note that a pool canceled like this cannot be restarted
// useful if unique requests within a larger system each create a new pool
cancel()
```

## Logging

Probe uses the `slog.Handler` interface for logging to maximize logging compatibility. By default,
`Probe` and `Pool` use the `logging.NoopLogHandler` which does not log. For information on building
a `slog.Handler` for your logger of choice, see the [slog Handler guide](https://github.com/golang/example/blob/master/slog-handler-guide/README.md).

## Why use a goroutine pool?

Go is excellent at parallelization and includes concurrency and sychronization mechanisms "out of the box."
Go also multiplexes goroutines onto system threads for threading efficiency. So, with this in mind,
why use a goroutine pool at all? 

One, they can keep code cleaner with lower cognitive load for developers. It's easy to conceptualize a
single pool of goroutines whereas starting multiple goroutines in different places in code can be
difficult to grok. 

Two, debugging synchronization issues with a goroutine Pool is easier than many concurrent goroutines 
that may be started in different places in code. The latter point is especially true when dealing with
graceful shutdowns or context cancellation.

## Why Probe?

Probe is named after the Protoss Probe unit from the popular video game series Starcraft. These
faithful workers are helpful for collecting resources. The name was chosen to be both fitting and fun!

![Probe](https://raw.githubusercontent.com/amplify-security/probe/main/doc/Protoss_Probe.jpg)

