# TODO: update doc for request, response

# Marine-Snow
A web framework



# URL Route

Base on trie tree, and support *named parameters*, *registered route match* and *bulk route registration*.

## route registration

**Demo for *registered route match***:

```go
package main

import (
	"MarineSnow/framework"
	"fmt"
	"net/http"
)

func nilHandler(ctx *framework.Context) error {
	return nil
}

const SERVER_ADDR = "127.0.0.1:80"

func main() {
	fmt.Printf("welcome to MarineSnow, start now. Listen on [%s]\n", SERVER_ADDR)
	core := framework.NewCore()
	server := &http.Server{
		Addr:    SERVER_ADDR,
		Handler: core,
	}

	/* register HTTP handler and route */
	core.GetRegisterFunc("/hello", nilHandler)
	core.GetRegisterFunc("/timeout", nilHandler)
	core.GetRegisterFunc("/timeout/demo", nilHandler)
	core.GetRegisterFunc("/hello/demo", nilHandler)
	//core.GetRegisterFunc("/hello/demo", nilHandler) // conflict with /hello/demo
	core.PostRegisterFunc("/hello/demo", nilHandler) // not conflict with GET /hello/demo

	core.DumpRoutes()

	_ = server.ListenAndServe()

	fmt.Println("bye~, exit Marine-Snow now.")
}
```

After run above code, you can get the routes have been registered from output:

```bash
[root@LC Marine_Snow]# make run
====== make debug ======
go build -gcflags="all=-N -l" -o server .
====== start to run ======
welcome to MarineSnow, start now. Listen on [127.0.0.1:80]
===== GET tire tree dump =====
+ [/hello] --> [main.nilHandler()]
+ [/hello/demo] --> [main.nilHandler()]
+ [/timeout] --> [main.nilHandler()]
+ [/timeout/demo] --> [main.nilHandler()]

===== POST tire tree dump =====
+ [/hello/demo] --> [main.nilHandler()]
```

Try to access the URL by:

```bash
[root@LC Marine_Snow]# curl http://127.0.0.1:80/hello
[root@LC Marine_Snow]# curl http://127.0.0.1:80/hello/demo
[root@LC Marine_Snow]# curl http://127.0.0.1:80/hello/timeout
[root@LC Marine_Snow]# curl http://127.0.0.1:80/timeout
[root@LC Marine_Snow]# curl http://127.0.0.1:80/timeout/demo

log output:

```

log output:

```bash
==> request[GET:/hello], match [/hello], forwarding to [main.nilHandler()]
==> request[GET:/hello/demo], match [/hello/demo], forwarding to [main.nilHandler()]
can not find any handler for [GET:/hello/timeout]
==> request[GET:/timeout], match [/timeout], forwarding to [main.nilHandler()]
==> request[GET:/timeout/demo], match [/timeout/demo], forwarding to [main.nilHandler()]
```





## named parameters

**Demo for *named parameters***:

```go
package main

import (
	"MarineSnow/framework"
	"fmt"
	"net/http"
)

func nilHandler(ctx *framework.Context) error {
	return nil
}

const SERVER_ADDR = "127.0.0.1:80"

func main() {
	fmt.Printf("welcome to MarineSnow, start now. Listen on [%s]\n", SERVER_ADDR)
	core := framework.NewCore()
	server := &http.Server{
		Addr:    SERVER_ADDR,
		Handler: core,
	}

	/* register HTTP handler and route */
	// named parameter
	core.PostRegisterFunc("/parameter/:id", nilHandler)
	//core.PostRegisterFunc("/parameter/:name", nilHandler) // conflict with /parameter/:id
	core.PostRegisterFunc("/parameter/:id/demo", nilHandler)
	//core.PostRegisterFunc("/parameter/:id/:name", nilHandler) // conflict with /parameter/:id/demo
	core.PostRegisterFunc("/parameter/:id/:name/end", nilHandler)
	core.PostRegisterFunc("/parameter/:age/:name/new-end", nilHandler)
    
	core.DumpRoutes()

	_ = server.ListenAndServe()

	fmt.Println("bye~, exit Marine-Snow now.")
}
```

After run above code, you can get the routes have been registered from output:

```bash
[root@LC Marine_Snow]# make run
====== make debug ======
go build -gcflags="all=-N -l" -o server .
====== start to run ======
welcome to MarineSnow, start now. Listen on [127.0.0.1:80]
===== POST tire tree dump =====
+ [/parameter/:id] --> [main.nilHandler()]
+ [/parameter/:id/demo] --> [main.nilHandler()]
+ [/parameter/:id/:name/end] --> [main.nilHandler()]
+ [/parameter/:age/:name/new-end] --> [main.nilHandler()]
```

Try to access the URL by:

```bash
[root@LC Marine_Snow]# curl -X POST http://127.0.0.1:80/parameter/123
[root@LC Marine_Snow]# curl -X POST http://127.0.0.1:80/parameter/123/demo
[root@LC Marine_Snow]# curl -X POST http://127.0.0.1:80/parameter/123/name/end
[root@LC Marine_Snow]# curl -X POST http://127.0.0.1:80/parameter/age-123/name/new-end
```

log output:

```bash
==> request[POST:/parameter/123], match [/parameter/:id], forwarding to [main.nilHandler()]
==> request[POST:/parameter/123/demo], match [/parameter/:id/demo], forwarding to [main.nilHandler()]
==> request[POST:/parameter/123/name/end], match [/parameter/:id/:name/end], forwarding to [main.nilHandler()]
==> request[POST:/parameter/age-123/name/new-end], match [/parameter/:age/:name/new-end], forwarding to [main.nilHandler()]
```







## bulk route registration

**Demo for *bulk route registration***:

```go
package main

import (
	"MarineSnow/framework"
	"fmt"
	"net/http"
)

func nilHandler(ctx *framework.Context) error {
	return nil
}

const SERVER_ADDR = "127.0.0.1:80"

func main() {
	fmt.Printf("welcome to MarineSnow, start now. Listen on [%s]\n", SERVER_ADDR)
	core := framework.NewCore()
	server := &http.Server{
		Addr:    SERVER_ADDR,
		Handler: core,
	}

	/* register HTTP handler and route */
	// group route register
	group := core.NewRouteGroup("/group/route")
	group.GetRegisterFunc("/name", nilHandler)
	group.GetRegisterFunc("/time", nilHandler)
	group.GetRegisterFunc("/id/:name", nilHandler)
	//core.GetRegisterFunc("/group/route/name", nilHandler) // conflict with group.GetRegisterFunc("/name", nilHandler)

	core.GetRegisterFunc("/group/route/dup", nilHandler)
	//group.GetRegisterFunc("/dup", nilHandler) // conflict with core.GetRegisterFunc("/group/route/dup", nilHandler)

	// inner RouteGroup
	upGroup := core.NewRouteGroup("/up/group")
	upGroup.GetRegisterFunc("/route-1", nilHandler)
	innerGroup := upGroup.Group("/inner")
	innerGroup.GetRegisterFunc("/route-2", nilHandler)
    
	core.DumpRoutes()

	_ = server.ListenAndServe()

	fmt.Println("bye~, exit Marine-Snow now.")
}
```

After run above code, you can get the routes have been registered from output:

```bash
[root@LC Marine_Snow]# make run
====== make debug ======
go build -gcflags="all=-N -l" -o server .
====== start to run ======
welcome to MarineSnow, start now. Listen on [127.0.0.1:80]
===== GET tire tree dump =====
+ [/group/route/name] --> [main.nilHandler()]
+ [/group/route/time] --> [main.nilHandler()]
+ [/group/route/id/:name] --> [main.nilHandler()]
+ [/group/route/dup] --> [main.nilHandler()]
+ [/up/group/route-1] --> [main.nilHandler()]
+ [/up/group/inner/route-2] --> [main.nilHandler()]
```

Try to access the URL by:

```bash
[root@LC Marine_Snow]# curl http://127.0.0.1:80/group/route/name
[root@LC Marine_Snow]# curl http://127.0.0.1:80/group/route/time
[root@LC Marine_Snow]# curl http://127.0.0.1:80/group/route/id/name-parameter
[root@LC Marine_Snow]# curl http://127.0.0.1:80/group/route/dup
[root@LC Marine_Snow]# curl http://127.0.0.1:80/up/group/route-1
[root@LC Marine_Snow]# curl http://127.0.0.1:80/up/group/up/group/inner/route-2
```

log output:

```bash
==> request[GET:/group/route/name], match [/group/route/name], forwarding to [main.nilHandler()]
==> request[GET:/group/route/time], match [/group/route/time], forwarding to [main.nilHandler()]
==> request[GET:/group/route/id/name-parameter], match [/group/route/id/:name], forwarding to [main.nilHandler()]
==> request[GET:/group/route/dup], match [/group/route/dup], forwarding to [main.nilHandler()]
==> request[GET:/up/group/route-1], match [/up/group/route-1], forwarding to [main.nilHandler()]
can not find any handler for [GET:/up/group/up/group/inner/route-2]
```








## Note:
Before commit file to main branch, please perform `make checkout`.