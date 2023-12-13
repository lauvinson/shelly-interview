> The following demonstration uses netcat to connect endpoints and http api to get historical statistics

# Building and Running
### Requirements:
Go 1.12+

### Building:
```go build -o server .```
> cross-compile guide: https://golangcookbook.com/chapters/running/cross-compiling/


### Running:
```./server```


> This will start the server listening on port 8080 for TCP connections and 8081 for HTTP requests.


# Usage
### Connect to the server:
```nc localhost 8080```

### Send messages to the server:
```Hello world!```

### Get server statistics:
```curl http://localhost:8081/stats```

> This will return a JSON object containing the server statistics.