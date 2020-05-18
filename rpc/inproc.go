// Authored and revised by YOC team, 2016-2018
// License placeholder #1

package rpc

import (
	"context"
	"net"
)

// NewInProcClient attaches an in-process connection to the given RPC server.
func DialInProc(handler *Server) *Client {
	initctx := context.Background()
	c, _ := newClient(initctx, func(context.Context) (net.Conn, error) {
		p1, p2 := net.Pipe()
		go handler.ServeCodec(NewJSONCodec(p1), OptionMethodInvocation|OptionSubscriptions)
		return p2, nil
	})
	return c
}
