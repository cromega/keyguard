# PowerMux

[![Build Status](https://travis-ci.org/AndrewBurian/powermux.svg?branch=master)](https://travis-ci.org/AndrewBurian/powermux)
[![codecov](https://codecov.io/gh/AndrewBurian/powermux/branch/master/graph/badge.svg)](https://codecov.io/gh/AndrewBurian/powermux)
[![Release](https://img.shields.io/github/release/AndrewBurian/powermux.svg)](https://github.com/AndrewBurian/powermux/releases)


A drop-in replacement for Go's `http.ServeMux` with all the missing features

PowerMux stores routes in radix trees for fast route matching and lookup on large numbers of routes.

## Dependencies

PowerMux requires at least Go version 1.8.

## Setting up PowerMux

In all cases, PowerMux does not support routes with a trailing slash `/` other than the root node.
Requests to paths that end in a slash are automatically redirected using a permanent redirection.

### Using `http.ServeMux` syntax

You can use PowerMux exactly as you would use Go's server mux.

```go
// Golang default
mux := http.NewServeMux()
mux.Handle("/", myHandler)
  
// PowerMux
mux := powermux.NewServeMux()
mux.Handle("/", myHandler)
```

### Using the Route syntax

PowerMux also has a cleaner way to declare routes, using the `Route` function.

Each call to `Route()` returns a pointer to that particular path on the radix tree, creating it if necessary.
At each route, you can add middleware, set handlers, or descend further into the route:

```go
mux := powermux.NewServeMux()
 
// Set a GET handler for "/"
mux.Route("/").Get(myHandler)
 
// Set POST/DELETE handlers for "/"
mux.Route("/").
    Post(myPostHandler).
    Delete(myDeleteHandler)
```

Sequential calls to route have the same effect as a single call with a longer path:

```go
mux.Route("/a").Route("/b").Route("/c") == mux.Route("/a/b/c")
```

Since Handler methods also return the route, the syntax can also be chained like so:

```go
mux.Route("/").
    Get(rootHandler).
    
    Route("/a").
    Get(aGetHandler).
    Post(aPostHandler).
    
    Route("/b").
    Get(abGetHandler)
```

## Middleware

PowerMux has support for any kind of middleware that uses the common `func(res, req, next)` syntax.  
Middleware handler objects must implement the `ServeHTTPMiddleware` interface.

Middleware will **always** be executed before any handlers, including default or generated not found handlers.

Middleware can be added to any route:

```go
mux.Route("/users").
    Middleware(authMiddleware).
    Get(sensitiveInfoHandler)
    
// or
mux.Route("/books").MiddleWare(loggingMiddleware)
mux.Route("/books").Get(booksHandler)
```

Middleware will be run if it's attached to any part of the route above and including the final path:

```go
mux.Route("/").Middleware(midRoot)
mux.Route("/a").Middleware(midA)
mux.Route("/a/b").Middleware(midB)
mux.Route("/c").Middleware(midC)
 
// requests to /a/b will run midRoot, midA, midB, 
// then any handlers on Route("/a/b")
```

Middleware can also be set up to selectively execute based on the HTTP method of the request.

The middleware function variants `MiddlewareFor` and `MiddlewareExceptFor` either set middleware to execute on only
specified methods, on all methods except the specified ones respectively.

```go
// don't run this middleware on OPTIONS requests
mux.Route("/a").MiddlewareExceptFor(ignoreCorsMid, http.MethodOptions)
```

## Host specific routes

Unlike the Go default multiplexer, host specific routes need to be handled separately. Use the `*Host` variants of
common functions to achieve this.
```go
mux.Route("/test")
mux.RouteHost("example.com", "/text")

// request to any host other than example.com will go to the first handler
```

## Not Found and OPTIONS handlers

`Options` and `NotFound` handlers are treated specially. If one is not found on the Route node requested, 
the latest one above that node will be used. This allows whole sections of routes to be covered under custom CORS
responses or Not Found handlers

## Path Parameters

Routes may include path parameters, specified with `/:name`:

```go
mux.Route("/users/:id/info").
    Get(userInfoHander)
```

This will make the variable `id` available to the get handler and any middleware.  

To retrieve path parameters, use `PathParam()`:

```go
// called with /users/andrew/info
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
        id := powermux.PathParam(r, "id")
        // id == "andrew"
}
```

Path parameters that aren't found return an empty string.  
Path parameters are unescaped with `url.PathUnescape`.

## Wildcard patterns
Routes may be declared with a wildcard indicator `*` at the end. 
This will match any path that does not have a more specific handler registered.

```go
mux.Route("/static/*").Get(staticContentHandler)
mux.Route("/static/favicon").Get(faviconGenerator)
 
// requests to /static will all be mapped to static content handler
// EXCEPT for requests to /static/favicon
```

Declaring a wildcard route at the same level as a path parameter route will never be executed as the path parameter takes greater precedence.

```go
mux.Route("/users/:id") // valid
mux.Route("/users/*")   // never matched
```

More routes may be specified after a wildcard, but they will never be executed:

```go
r1 := mux.Route("/users/*") // valid
r1.Route("/further/paths")  // never matched
```

## Route precedence

If multiple routes are declared that could match a given path, they are selected in this order:

  1. A literal path `/users/andrew/info`
  2. A path with parameters `/users/:id/info`
  3. A wildcard path `/users/*`

## Retrieving the original route path

Handlers and Middleware may access the route pattern that was used by powermux to route any particular 
request with the `RequestPath` function.

```go
servemux.Route("/users/:id/info").Get(userHandler)
...
// envoked with /users/andrew/info
func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
        originalPath := powermux.RequestPath(req)
        
        originalPath == "/users/:id/info"
        req.URL.Path == "/users/andrew/info"
}
```

## Handler precedence

When multiple handlers are declared on a single route for different methods, they are selected in this order:

  1. An exact method match
  2. HEAD requests can use GET handlers
  3. The ANY handler
  4. A generated Method Not Allowed handler
