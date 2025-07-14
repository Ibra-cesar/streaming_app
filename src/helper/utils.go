package helper

import "net/http"


//Middleware Chaining Helpers
type Middleware func(http.Handler) http.Handler

func ChainMiddleware(ms ...Middleware) Middleware{
	return func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i--{
			x := ms[i]
			next = x(next)
		}
		return  next
	}
}
