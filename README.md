# probe

[![Amplify Security](https://github.com/amplify-security/probe/actions/workflows/amplify.yml/badge.svg?branch=main)](https://github.com/amplify-security/probe/actions/workflows/amplify.yml)
[![test](https://github.com/github/docs/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/amplify-security/probe/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/amplify-security/probe)](https://goreportcard.com/report/github.com/amplify-security/probe)

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
fmt.Println(p.Idle()) // 14 (this is a race condition with the Pool work channel, your output may differ)
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
instead of using the default, blocking channel.

```go
work := make(chan probe.Runner, 1024)
p := probe.NewProbe(&probe.ProbeConfig{ WorkChan: work })
work <- func() {
    fmt.Println("Hello from a buffered Probe!")
}
```

## Why use a goroutine pool?

Go is excellent at parallelization and includes concurrency and sychronization mechanisms "out of the box."
Go also multiplexes goroutines onto system threads for threading efficiency. So, with this in mind,
why use a goroutine pool at all? One, they can keep code cleaner with lower cognitive load for developers.
It's easy to conceptualize a single pool of goroutines whereas starting multiple goroutines in different
places in code can be difficult to grok. Two, debugging synchronization issues with a goroutine Pool is
easier than many concurrent goroutines that may be started in different places in code. The latter point is
especially true when dealing with graceful shutdowns or context cancellation.

