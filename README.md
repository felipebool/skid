# skid

## Task
Using only the standard library, create a Go HTTP server that on each request
responds with a counter of the total number of requests that it has received
during the previous 60 seconds (moving window). The server should continue to
the return the correct numbers after restarting it, by persisting data to a
file.

## Requirements
* Create an HTTP server which responds with a counter
* Use moving window algorithm (60 seconds)
* Persist current window to file

## Considerations
* since `time.Now()` returns `int64`, was adopted throughout the solution for simplicity

## *Most importantly... enjoy it :)*
As the task description said, I took this opportunity to visit a few topics I
had reserved to read.