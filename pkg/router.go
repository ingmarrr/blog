package pkg

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

type Router struct {
	tree       *node
	middleware []MiddlewareFunc
}

func NewRouter() Router {
	return Router{tree: &node{}}
}

func (r *Router) Handle(path string, handler http.HandlerFunc) {
	parts, err := extractParts(path)
	if err != nil {
		panic("Invalid Path")
	}
	for _, mw := range r.middleware {
		handler = mw(handler)
	}
	r.tree.insert(parts, handler)
}

func (r *Router) Use(mw MiddlewareFunc) {
	r.middleware = append(r.middleware, mw)
}


func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	params := make(map[string]string)
	parts, err := extractParts(req.URL.Path)
	if err != nil {
    http.NotFound(w, req)
		return
	}

	handler, err := r.tree.search(parts, params)
	if err != nil {
    http.NotFound(w, req)
		return
	}

	req = req.WithContext(context.WithValue(req.Context(), "params", params))
	handler(w, req)
}

type node struct {
	children  map[string]*node
	handler   http.HandlerFunc
	isParam   bool
	paramName string
}

func (n *node) insert(parts []string, handler http.HandlerFunc) {
	if len(parts) == 0 {
		n.handler = handler
		return
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}

	path := parts[0]
	child, ok := n.children[path]
	if !ok {
		child = &node{}
		if strings.HasPrefix(path, ":") {
			child.isParam = true
			child.paramName = path[1:]
			path = ":"
		}
		n.children[path] = child
	}

	if child.isParam {
		child.insert(parts[1:], handler)
		return
	}
	child.insert(parts[1:], handler)
}

func (n *node) search(parts []string, params map[string]string) (http.HandlerFunc, error) {
	if len(parts) == 0 {
		return n.handler, nil
	}

	path := parts[0]
	if child, ok := n.children[path]; !ok {
		if param, ok := n.children[":"]; ok {
			params[param.paramName] = path
			return param.search(parts[1:], params)
		}
		return http.NotFound, errors.New("Not found")
	} else {
		return child.search(parts[1:], params)
	}
}

func extractParts(s string) ([]string, error) {
	if !strings.HasPrefix(s, "/") {
		return []string{}, errors.New("Path has start with '/'.")
	}
	return strings.Split(s[1:], "/"), nil
}
