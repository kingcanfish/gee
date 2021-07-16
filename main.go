package main

import (
	"github.com/kingcanfish/gee"
	"log"
	"net/http"
	"time"
)
func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		//c.String(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Recovery())
	r.Use(gee.Logger()) // global midlleware
	r.Get("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	//r.Use(gee.Logger()) // global midlleware

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.Get("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
