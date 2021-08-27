# skid

## What was asked
Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds. The server should continue to the return the correct numbers after restarting it, by persisting data to a file.

## Running the project
### Parameters
I added a few parameters to give a few knobs and buttons to test the server.
* **size** integer, sets the size of the window, default value is 60 seconds
* **path** string, sets the path where to persist the window, default value is *.window.json*
* **restore** boolean, sets if must restore previous state from file, default value is `false`
* **port** integer, sets the port where to run the project, default value is 8080

The default values are set to what was asked in the description, so, by running
```
go run main.go
```
you have a server, listening to port 8080, with a window of size 60 seconds.

If you want to persist the state to a file, just stop the execution (`^C`) and the window will be persisted to what was passed to `--path`. To resume the execution from what was persisted, just run:
```
go run main.go --restore
```
