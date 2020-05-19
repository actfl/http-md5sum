# http-md5sum 

### How to build
```
go build 
```

### How to execute the file
You could specify the number of goroutine with `--parallel, -p, --p` arguments.

```
> ./http-md5sum --p 5 google.com naver.com jtbc.joins.com anaconda.com yahoo.com cnn.com nytimes.com
2020/05/19 14:47:14 creating 5 goroutine
google.com               c7be58e10eec940ee551789fef99fc40
anaconda.com             a4e6a108c9717a1e13f19a70b8319f31
cnn.com                  99849984bb65fb7a3ff0c9f795d812d7
nytimes.com              e69930f571bc19e2794cacdc432009b4
yahoo.com                baf7039b7035ce51a8bd29d87659c3b1
naver.com                timeout on `naver.com`, HTTP Timeout error
jtbc.joins.com           fe3c3871ade861556706a350ba5e3b97
```

### Unittests
```
go test ./pkg/httpsum/ -v
=== RUN   TestHttpSum_get
=== RUN   TestHttpSum_get/0:_successful_result
=== RUN   TestHttpSum_get/1:_url_not_found
=== RUN   TestHttpSum_get/2:_timeout
--- PASS: TestHttpSum_get (1.00s)
    --- PASS: TestHttpSum_get/0:_successful_result (0.00s)
    --- PASS: TestHttpSum_get/1:_url_not_found (0.00s)
    --- PASS: TestHttpSum_get/2:_timeout (1.00s)
PASS
ok      github.com/tomahawk28/http-md5sum/pkg/httpsum   1.018s

```