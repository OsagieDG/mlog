# mlog
- Logging Middleware for Go HTTP Server


## Installation
```
     go get github.com/OsagieDG/mlog
```

## Usage
import `github.com/OsagieDG/mlog/service/middleware`

```go

mlog := middleware.MLog(
    middleware.LogResponse,
    middleware.RecoverPanic,
)

listenAddr := ":6862"
log.Printf("Server is listening on %s", listenAddr)
if err := http.ListenAndServe(listenAddr, mlog(router)); err != nil {
    log.Fatal("HTTP server error:", err)
}

```

![mlog](https://github.com/OsagieDG/mlog/blob/main/blob/mlog.png)



