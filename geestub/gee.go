package gee

import (
	"net/http"
	"strings"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 中间件支持
	parent      *RouterGroup
	engine      *Engine
}

type Engine struct {
	*RouterGroup
	// router 内存储的是所有的路由
	router *router
	groups []*RouterGroup // 存储所有的分组
}

// New 实际上创建了一个默认的分组路由(根路由)
// 第15行 已经对*RouterGroup 实现了封装
// Engine struct 某种意义实现了 RouteGroup struct (鸭子模型)
func New() *Engine {
	engine := &Engine{router: newRoute()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup {engine.RouterGroup}
	return engine
}

// Group 用来创建新的分组
// g 的子分组 同时将这个子分组注册到engine 全局Group上
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newRouteGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, newRouteGroup)
	return newRouteGroup
}

// Use 给每个路由组添加中间件
func (g *RouterGroup) Use(middlewares ...HandlerFunc)  {
	g.middlewares = append(g.middlewares, middlewares...)
}

//给最后的节点添加路由
func (g *RouterGroup) addRoute(method string, subPath string, handler HandlerFunc) {
	pattern := g.prefix + subPath
	//注册到 engine 上的路由集合
	g.engine.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) Get(subPath string, handler HandlerFunc) {
	// 调用的第44行的方法
	g.addRoute("GET", subPath, handler)
}
func (g *RouterGroup) Post(subPath string, handler HandlerFunc) {
	// 调用的第44行的方法
	g.addRoute("POST", subPath, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c.handlers = middlewares
	e.router.handle(c)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
