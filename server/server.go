package server

import (
	//"bytes"
	"fmt"
	"github.com/tomyhero/ore_server/context"
	"io"
	"net"
	"os"
	"reflect"
)

type Server struct {
	Port       int
	conns      []net.Conn
	ln         net.Listener
	dispatcher *Dispatcher
}

func (s *Server) Setup(handlers []context.IHandler) {
	s.dispatcher = NewDispatcher(handlers)
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error Accepting", err.Error())
			os.Exit(1)
		}
		go handle(s.dispatcher, conn)
	}
}

// Handles incoming requests.
func handle(dispatcher *Dispatcher, conn net.Conn) {

	// when out of for loop, close the connection.
	defer conn.Close()

	for {
		fmt.Println("start")
		cdata := CData{SerializorType: SERIALIZOR_TYPE_MESSAGE_PACK}
		data, err := cdata.Receive(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client dissconected")
				break
			} else {
				fmt.Println("receive cdata", err)
				break
			}
		}

		c, err := context.NewContext(data)
		if err != nil {
			fmt.Println("create context", err)
			break
		}

		loginAction, onLogin := dispatcher.LoginActions[c.Req.GetCMD()]
		// do login
		if onLogin {
			fmt.Println(c.Req.GetCMD(), c.Req.Header, dispatcher.LoginActions, loginAction)
			ok := loginAction.Call([]reflect.Value{reflect.ValueOf(c)})[0].Bool()
			if ok {
				c.Res.Header["STATUS"] = context.STATUS_OK
			} else {
				c.Res.Header["STATUS"] = context.STATUS_NOT_OK
			}
			// do auth action 
		} else {

			action, find := dispatcher.Actions[c.Req.GetCMD()]
			fmt.Println(c.Req.GetCMD(), c.Req.Header, dispatcher.Actions, action, find)
			if find {
				if dispatcher.ExecAuth(c, c.Req.GetCMD()) {

					// BEFORE_EXECUTE
					dispatcher.BeforeExecute(c, c.Req.GetCMD())

					action.Call([]reflect.Value{reflect.ValueOf(c)})
					c.Res.Header["STATUS"] = context.STATUS_OK

					// AFTER_EXECUTE
					dispatcher.AfterExecute(c, c.Req.GetCMD())

				} else {
					c.Res.Header["STATUS"] = context.STATUS_FORBIDDEN
				}
			} else {
				c.Res.Header["STATUS"] = context.STATUS_COMMAND_NOT_FOUND
			}
		}

		err = cdata.Send(conn, c.Res.GetData())
		if err != nil {
			fmt.Println("send fail", err)
			break
		}
		fmt.Println("end")
	}
}
