In Go, **channels** are a built-in concurrency primitive used for communication between goroutines.
The core idea is:

> **"Do not communicate by sharing memory; instead, share memory by communicating."** — Rob Pike

Channels let goroutines safely exchange data without manually managing locks in many cases.

---

# 1. What is a Channel?

A channel is a typed conduit through which you can send and receive values.

```go
ch := make(chan int)
```

This creates a channel that transports `int` values.

---

# 2. Basic Syntax

## Sending

```go
ch <- 10
```

Send value `10` into the channel.

---

## Receiving

```go
x := <-ch
```

Receive a value from the channel.

---

# 3. Why Channels Exist

Without channels, goroutines running concurrently would need:

* mutexes
* shared variables
* synchronization logic
* race-condition handling

Channels provide:

* synchronization
* communication
* coordination

all together.

---

# 4. Simple Example

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan string)

    go func() {
        ch <- "hello from goroutine"
    }()

    msg := <-ch

    fmt.Println(msg)
}
```

---

## What Happens Here

### Step 1

Main goroutine creates channel:

```go
ch := make(chan string)
```

---

### Step 2

A new goroutine starts:

```go
go func() {
    ch <- "hello"
}()
```

---

### Step 3

The goroutine tries to send:

```go
ch <- "hello"
```

But it blocks until someone receives.

---

### Step 4

Main goroutine waits:

```go
msg := <-ch
```

Once receive happens:

* send completes
* receive completes
* program continues

---

# 5. Channels are Blocking

This is extremely important.

---

## Send Blocks

```go
ch <- value
```

Blocks until receiver is ready.

---

## Receive Blocks

```go
value := <-ch
```

Blocks until sender sends something.

---

# 6. Buffered vs Unbuffered Channels

---

# Unbuffered Channel

```go
ch := make(chan int)
```

Capacity = 0.

Send and receive must happen simultaneously.

Think of it like a handshake.

---

## Example

```go
ch := make(chan int)

go func() {
    ch <- 5
}()

fmt.Println(<-ch)
```

Works because:

* sender waits
* receiver waits
* synchronization occurs

---

# Buffered Channel

```go
ch := make(chan int, 3)
```

Capacity = 3.

Can store values before receiver consumes them.

---

## Example

```go
ch := make(chan int, 2)

ch <- 1
ch <- 2

fmt.Println(<-ch)
fmt.Println(<-ch)
```

No goroutine needed because buffer has space.

---

## Buffer Full Behavior

```go
ch := make(chan int, 2)

ch <- 1
ch <- 2
ch <- 3 // BLOCKS
```

Third send blocks because buffer is full.

---

# 7. Channel Direction

Channels can be restricted.

---

## Send-only

```go
func send(ch chan<- int) {
    ch <- 10
}
```

Can only send.

---

## Receive-only

```go
func receive(ch <-chan int) {
    fmt.Println(<-ch)
}
```

Can only receive.

---

## Why Useful?

Improves API safety.

You can prevent misuse.

---

# 8. Closing Channels

Channels can be closed using:

```go
close(ch)
```

---

## Why Close?

Signals:

> "No more values will be sent."

---

# Example

```go
ch := make(chan int)

go func() {
    ch <- 1
    ch <- 2
    close(ch)
}()

for v := range ch {
    fmt.Println(v)
}
```

Output:

```text
1
2
```

Loop ends automatically after close.

---

# 9. Reading from Closed Channel

After channel closes:

```go
v, ok := <-ch
```

---

## Behavior

| State           | v            | ok    |
| --------------- | ------------ | ----- |
| value available | actual value | true  |
| channel closed  | zero value   | false |

---

## Example

```go
ch := make(chan int, 1)

ch <- 5
close(ch)

fmt.Println(<-ch)

v, ok := <-ch

fmt.Println(v, ok)
```

Output:

```text
5
0 false
```

---

# 10. Important Rules

---

## Only Sender Should Close

Bad:

```go
receiver closes channel
```

Good:

```go
sender closes channel
```

Because sender knows when work is done.

---

## Sending to Closed Channel Panics

```go
close(ch)
ch <- 1 // panic
```

---

## Closing Nil Channel Panics

```go
var ch chan int
close(ch) // panic
```

---

# 11. Nil Channels

```go
var ch chan int
```

Zero value of channel is `nil`.

---

## Behavior

| Operation | Result         |
| --------- | -------------- |
| send      | blocks forever |
| receive   | blocks forever |
| close     | panic          |

---

# 12. Select Statement

`select` lets you wait on multiple channel operations.

---

## Example

```go
select {
case msg := <-ch1:
    fmt.Println(msg)

case ch2 <- 10:
    fmt.Println("sent")

default:
    fmt.Println("nothing ready")
}
```

---

# How select Works

* waits for channel operation
* picks one ready case
* blocks if none ready
* `default` makes it non-blocking

---

# 13. Timeout Pattern

Very common in Go.

---

## Example

```go
select {
case msg := <-ch:
    fmt.Println(msg)

case <-time.After(2 * time.Second):
    fmt.Println("timeout")
}
```

If no message arrives within 2 seconds:

```text
timeout
```

---

# 14. Fan-Out Pattern

Multiple workers consume jobs.

---

## Example

```go
jobs := make(chan int)

for i := 0; i < 3; i++ {
    go worker(jobs)
}
```

Workers compete for jobs.

---

# 15. Fan-In Pattern

Multiple goroutines send into one channel.

---

## Example

```go
func producer(out chan<- int, start int) {
    for i := start; i < start+3; i++ {
        out <- i
    }
}
```

Many producers → one consumer.

---

# 16. Worker Pool Example

```go
package main

import (
    "fmt"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        fmt.Println("worker", id, "processing", job)
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 5)
    results := make(chan int, 5)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= 5; j++ {
        jobs <- j
    }

    close(jobs)

    for a := 1; a <= 5; a++ {
        fmt.Println(<-results)
    }
}
```

---

# 17. Deadlocks

Go detects deadlocks.

---

## Example

```go
ch := make(chan int)
ch <- 5
```

Panic:

```text
fatal error: all goroutines are asleep - deadlock!
```

Because nobody receives.

---

# 18. Channels vs Mutexes

---

## Use Channels When

* passing ownership
* coordinating goroutines
* pipelines
* worker systems
* streaming data

---

## Use Mutexes When

* protecting shared state
* high-performance shared memory access
* frequent reads/writes

---

# 19. Internal Mental Model

A channel internally has:

* queue/buffer
* send wait queue
* receive wait queue
* lock

When send/receive occurs:

1. runtime locks channel
2. matches sender/receiver
3. copies data
4. wakes goroutine

---

# 20. Pipelines

Channels are heavily used in pipelines.

---

## Example

```go
func generate(nums ...int) <-chan int {
    out := make(chan int)

    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()

    return out
}
```

---

# 21. Context Cancellation with Channels

`context.Context` internally uses channels.

---

## Example

```go
select {
case <-ctx.Done():
    return
}
```

This listens for cancellation signal.

---

# 22. Common Mistakes

---

## Goroutine Leak

```go
go func() {
    ch <- data
}()
```

If nobody receives, goroutine blocks forever.

---

## Forgetting close

Consumers waiting forever.

---

## Overusing Channels

Not every concurrent problem needs channels.

---

# 23. Advanced: Channel Ownership Pattern

Good architecture pattern.

---

## Rule

Function that creates channel:

* writes to it
* closes it

Consumers only read.

---

# 24. Buffered Channel as Semaphore

```go
sem := make(chan struct{}, 3)
```

Limit concurrency to 3 tasks.

---

## Example

```go
sem <- struct{}{} // acquire

go func() {
    defer func() { <-sem }() // release
}()
```

---

# 25. Range over Channels

```go
for v := range ch {
    fmt.Println(v)
}
```

Equivalent to:

```go
for {
    v, ok := <-ch
    if !ok {
        break
    }

    fmt.Println(v)
}
```

---

# 26. Best Practices

---

## Prefer channels for communication

But not blindly.

---

## Keep channels directional

Use:

```go
chan<- T
<-chan T
```

when possible.

---

## Close channels carefully

Only producers close.

---

## Avoid giant buffered channels

They can hide bugs.

---

# 27. Real-World Uses

Channels power:

* worker pools
* job queues
* event systems
* streaming systems
* request coordination
* cancellation
* rate limiting
* pub/sub systems

---

# 28. Key Insight

A channel is BOTH:

* a queue
* a synchronization mechanism

That dual nature is why Go concurrency feels elegant.

---

# 29. One-Line Summary

> Channels allow goroutines to safely communicate and synchronize by sending typed values through blocking queues.

---

# 30. Tiny Cheat Sheet

```go
ch := make(chan int)        // unbuffered
ch := make(chan int, 5)     // buffered

ch <- v                     // send
v := <-ch                   // receive

close(ch)                   // close

v, ok := <-ch               // closed check

for v := range ch {}        // iterate

select {}                   // multiplex channels
```

If you want, I can also explain:

* channel internals inside Go runtime
* scheduler + channels
* select fairness
* channel performance
* advanced concurrency patterns
* real production examples
* channels vs async/await
* CSP model behind Go
* how channels are implemented in runtime (`hchan`)
