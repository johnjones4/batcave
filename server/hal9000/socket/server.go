package socket

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/learning"
	"github.com/johnjones4/hal-9000/server/hal9000/runtime"
)

type Server struct {
	Host     string
	Runtime  *runtime.Runtime
	Location core.Coordinate
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.Host)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.newConnection(conn)
	}
}

func (s *Server) newConnection(conn net.Conn) {
	handleError := func(err error) {
		log.Println(err)
		conn.Write([]byte("Sorry. Something went wrong.\n\r"))
	}
	handle := func() error {
		conn.SetDeadline(time.Now().Add(time.Minute * 5))

		buffer := bufio.NewReader(conn)
		usernamePassword, err := buffer.ReadString('\n')
		if err != nil {
			return err
		}

		usernamePasswordSplit := strings.Split(strings.TrimSpace(usernamePassword), ":")
		if len(usernamePasswordSplit) != 2 {
			return errors.New("bad username password string")
		}

		user, err := s.Runtime.UserStore.Login(usernamePasswordSplit[0], usernamePasswordSplit[1])
		if err != nil {
			return err
		}

		conn.Write([]byte("ok\n"))

		for {
			input, err := buffer.ReadString('\n')
			if err != nil {
				return err
			}

			cleanInput := strings.TrimSpace(input)

			if cleanInput == "" {
				continue
			}

			state, err := s.Runtime.StateStore.GetStateForUser(user)
			if err != nil {
				handleError(err)
				continue
			}

			req, err := s.Runtime.Parse(core.InboundBody{
				Body:     cleanInput,
				Location: s.Location,
			}, state)
			if err != nil {
				handleError(err)
				continue
			}

			res, err := s.Runtime.Intents.ProcessRequest(req)
			if err != nil {
				handleError(err)
				continue
			}

			conn.Write([]byte(res.OutboundBody.Body + "\n\r"))

			err = s.Runtime.Logger.Log(learning.InteractionEvent{
				Request:  req,
				Response: res,
			})
			if err != nil {
				handleError(err)
				continue
			}

			err = s.Runtime.StateStore.SetStateForUser(state)
			if err != nil {
				handleError(err)
				continue
			}
		}
	}
	err := handle()
	if err != nil {
		log.Println(err)
	}
	conn.Close()
}
