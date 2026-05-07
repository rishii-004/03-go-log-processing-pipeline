# 03-go-log-processing-pipeline
This teaches:  Pipeline pattern Fan-out workers Fan-in aggregator Channels design Clean concurrency architecture

Goroutines
Channels
Blocking
Deadlocks
Unbuffered vs Buffered channels
Channel closing
Concurrency synchronization

https://bwoff.medium.com/the-comprehensive-guide-to-concurrency-in-golang-aaa99f8bccf6


## blocking 
1. If there is no data in the channel, reading from the channel blocks the goroutine until data becomes available.
2. If there is no receiver ready to read from the channel, writing to the channel blocks the goroutine until a receiver is available.
3. different behavior for buffered and unbuffered channels 

## buffered vs unbuffered channels 
1. buffered channels have capacity, and only block when capacity is full 
    ch: make(channel int, capacity)
2. unbuffered channels, they have no capacity, sends are blocked to a channel, until a reciever receives the value
    ch: make(channel int)

## Fan-out 
1. single channel output is distributed among multiple goroutines to parallelize cpu usage 

## fan-in 
1. data from multiple channels is consolidated into a single channel 

