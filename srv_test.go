package wsrv

import (
	"net/http"
	"testing"

	"github.com/eyasliu/websocket"
	"github.com/go-eyas/toolkit/log"
)

func TestSrv(t *testing.T) {
	server := New(&websocket.Config{
		Logger: log.SugaredLogger,
	})
	// server.UseRequest(func(c *Context) {
	// 	log.Debugf("ws request middleware, sid=%d", c.SessionID)
	// })
	// server.UseResponse(func(c *Context) {
	// 	log.Debugf("ws response middleware, sid=%d", c.SessionID)
	// })
	// server.UseRequest(func(c *Context) {
	// 	// uid, ok := c.Get("uid").(int64)
	// 	// if !ok || uid == 0 {
	// 	// 	c.Abort()
	// 	// }
	// })
	server.Use(func(c *Context) {
		log.Debugf("ws request middleware, sid=%d, cmd=%s, data=%s", c.SessionID, c.CMD, string(c.Request.Data))
		c.Next()
		log.Debugf("ws response middleware, sid=%d cmd=%s, data=%v", c.SessionID, c.CMD, c.Response.Data)
	})
	server.Handle("register")
	server.Handle("register", func(c *Context) {
		c.Set("uid", int(123))
		for sid, vals := range server.Session {
			if uid, ok := vals["uid"]; ok {
				server.Push(sid, &WSResponse{
					CMD: "have_user_register",
					Data: map[string]interface{}{
						"uid": uid,
					},
				})
			}
		}
		c.OK()
	})

	t.Log("init ws srv ok ")
	http.HandleFunc("/ws", server.Engine.HTTPHandler)
	http.HandleFunc("/play", server.Engine.Playground)
	http.ListenAndServe("127.0.0.1:9000", nil)

}
