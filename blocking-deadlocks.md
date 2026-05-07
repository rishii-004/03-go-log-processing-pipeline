In Go concurrency, **blocking** and **deadlocking** are closely related concepts, but they are not the same.

---

# 1. Blocking

A goroutine is **blocked** when it cannot continue execution until some event happens.

This is normal and expected in concurrent programs.

---

# Example: Blocking on Receive

```go id="md3r6t"
package main

import "fmt"

func main() {
    ch := make(chan int)

    fmt.Println("waiting...")

    x := <-ch

    fmt.Println(x)
}
```

---

## What Happens

```go id="e6f91e"
x := <-ch
```

The main goroutine pauses because:

* channel is empty
* no sender exists yet

So it waits.

This waiting state = **blocked**.

---

# Blocking is NOT Bad

Blocking is fundamental to Go concurrency.

Examples:

* waiting for data
* waiting for I/O
* waiting for lock
* waiting for network
* waiting for channel

---

# Analogy

Imagine:

* sender = person delivering package
* receiver = person waiting at door

Receiver blocks until package arrives.

---

# 2. Types of Blocking in Channels

---

# A. Receive Blocking

```go id="37d86x"
x := <-ch
```

Blocks when:

* no value available

---

# B. Send Blocking

```go id="tq02r4"
ch <- 5
```

Blocks when:

* unbuffered channel has no receiver
  OR
* buffered channel is full

---

# Example

```go id="prgqeu"
ch := make(chan int)

ch <- 10
```

Main goroutine blocks forever because nobody receives.

---

# 3. Buffered Channel Blocking

---

## Buffered Example

```go id="ryb3m8"
ch := make(chan int, 2)

ch <- 1
ch <- 2
ch <- 3
```

---

## What Happens

Buffer capacity = 2

State:

| Operation | Buffer        |
| --------- | ------------- |
| send 1    | [1]           |
| send 2    | [1,2]         |
| send 3    | FULL → blocks |

---

# 4. Deadlock

A **deadlock** happens when:

> goroutines are blocked forever and none can make progress.

This is a program failure.

---

# Example Deadlock

```go id="5jdr0o"
package main

func main() {
    ch := make(chan int)

    ch <- 10
}
```

---

## Why Deadlock Happens

```go id="t6ifh6"
ch <- 10
```

Needs a receiver.

But:

* only main goroutine exists
* no receiver exists
* main blocks forever

Go runtime detects:

> all goroutines are asleep

and crashes.

---

# Runtime Error

```text id="l2umqo"
fatal error: all goroutines are asleep - deadlock!
```

---

# 5. Important Difference

| Blocking                | Deadlock          |
| ----------------------- | ----------------- |
| temporary waiting       | permanent waiting |
| normal behavior         | fatal error       |
| goroutine may resume    | nobody can resume |
| expected in concurrency | bug in program    |

---

# 6. Blocking Example That Works

```go id="m6d3mo"
package main

import "fmt"

func main() {
    ch := make(chan int)

    go func() {
        ch <- 42
    }()

    x := <-ch

    fmt.Println(x)
}
```

---

## Timeline

### Main goroutine

Blocks here:

```go id="9v4lwr"
x := <-ch
```

---

### Worker goroutine

Eventually sends:

```go id="d6h6z8"
ch <- 42
```

---

### Result

* sender unblocks receiver
* receiver unblocks sender
* program continues

This is healthy blocking.

---

# 7. Deadlock with Multiple Goroutines

Deadlock doesn't require one goroutine.

---

## Example

```go id="bllf0i"
func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        <-ch1
        ch2 <- 1
    }()

    go func() {
        <-ch2
        ch1 <- 1
    }()

    select {}
}
```

---

## What Happens

Goroutine A waits for `ch1`

Goroutine B waits for `ch2`

Nobody sends first.

Both blocked forever.

Circular waiting = deadlock.

---

# 8. Circular Wait

Most deadlocks involve circular dependency.

---

# Example Mental Model

```text id="0e22go"
A waits for B
B waits for C
C waits for A
```

Nobody proceeds.

---

# 9. Deadlock with Mutexes

Channels aren't the only source.

---

## Example

```go id="6s4s2x"
mu.Lock()
mu.Lock()
```

Same goroutine tries locking twice.

Second lock blocks forever.

---

# 10. Goroutine Leak vs Deadlock

These are different.

---

# Goroutine Leak

Some goroutines blocked forever,
BUT program still runs.

---

## Example

```go id="gh3n6w"
go func() {
    <-ch
}()
```

If nobody sends:

* goroutine stuck forever
* main may still continue

This is a leak, not global deadlock.

---

# 11. How Go Detects Deadlock

Go runtime checks:

```text id="57brv5"
Are all goroutines blocked?
```

If yes:

```text id="a1n6q9"
panic: deadlock
```

---

# Important

If even ONE goroutine can still run:

* no deadlock panic

---

# 12. Nil Channel Blocking

Very important interview topic.

---

# Nil Channel

```go id="l4ckl2"
var ch chan int
```

`ch == nil`

---

# Behavior

```go id="x4mly6"
<-ch
```

Blocks forever.

---

## Why?

Nil channels have no internal channel structure.

No sender/receiver can ever wake it.

---

# 13. Select and Blocking

---

# Without Default

```go id="9ng70m"
select {
case msg := <-ch:
    fmt.Println(msg)
}
```

Blocks until channel ready.

---

# With Default

```go id="uk2uxv"
select {
case msg := <-ch:
    fmt.Println(msg)

default:
    fmt.Println("not ready")
}
```

Non-blocking.

---

# 14. Channel States and Blocking

---

## Unbuffered Channel

| Operation | Condition          |
| --------- | ------------------ |
| send      | waits for receiver |
| receive   | waits for sender   |

---

## Buffered Channel

| Operation | Condition            |
| --------- | -------------------- |
| send      | waits if buffer full |
| receive   | waits if empty       |

---

# 15. Visual Timeline

---

# Unbuffered

```text id="ey1l0t"
Sender ----waiting---->
                     Receiver arrives
Sender resumes
Receiver resumes
```

---

# Buffered

```text id="ukod8x"
Sender -> buffer -> receiver later
```

Blocking delayed until buffer fills.

---

# 16. Common Causes of Deadlock

---

## A. No Receiver

```go id="i7uxif"
ch <- 1
```

---

## B. No Sender

```go id="9tk18s"
<-ch
```

---

## C. Full Buffered Channel

```go id="af0w6k"
ch := make(chan int, 1)
ch <- 1
ch <- 2
```

---

## D. Range Without Close

```go id="m9r6yu"
for v := range ch {
}
```

If channel never closes:

* loop waits forever

---

## E. Circular Dependencies

Goroutines waiting on each other.

---

# 17. Debugging Deadlocks

Go gives stack traces.

Example:

```text id="0jol1v"
goroutine 1 [chan send]:
```

Meaning:

* goroutine 1 blocked on send

---

# 18. Prevention Strategies

---

## Use Timeouts

```go id="7eh98m"
select {
case msg := <-ch:
    fmt.Println(msg)

case <-time.After(time.Second):
    fmt.Println("timeout")
}
```

---

## Close Channels Properly

Consumers should know completion.

---

## Avoid Circular Waits

Design ownership carefully.

---

## Use Buffered Channels Carefully

Large buffers can hide problems.

---

# 19. Internal Runtime View

When goroutine blocks:

* runtime parks goroutine
* scheduler runs another goroutine

Blocked goroutines consume little CPU.

---

# 20. Key Insight

Blocking is a **feature**.

Deadlock is a **bug**.

---

# Tiny Summary

| Concept        | Meaning                                |
| -------------- | -------------------------------------- |
| blocking       | goroutine temporarily waits            |
| deadlock       | all goroutines wait forever            |
| send block     | no receiver / full buffer              |
| receive block  | no sender / empty buffer               |
| nil channel    | blocks forever                         |
| deadlock panic | runtime detects zero progress possible |

---

# Final Mental Model

Think of goroutines as people and channels as doors.

---

## Blocking

```text id="2f3gg8"
Person waits at door for package.
```

Normal.

---

## Deadlock

```text id="s0fj5f"
Everyone waits at doors,
but nobody delivers anything.
```

System frozen.
