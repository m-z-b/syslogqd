// Copyright 2022 Mike Bell, Albion Research Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package listener

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/m-z-b/syslogqd/syslog"
)

// ReadTimeout is the timeout for a TCP read
const ReadTimeout = 1000 * time.Millisecond

// Regex used to identify the start of messages
var entryStart = regexp.MustCompile(`<[0-9]{2,3}>`)

type TCPListener struct {
	listener  net.Listener
	reporting syslog.Channel
}

func NewTCPListener(port int, reporting syslog.Channel) (*TCPListener, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	r := &TCPListener{
		listener:  l,
		reporting: reporting,
	}
	return r, nil
}

// If we can find a syslog message start after the beginning of the buffer
// we can treat all bytes up to that point as a message and remove it
// from the buffer
func (self *TCPListener) removeComplete(msg []byte, c net.Conn) []byte {
	if len(msg) > 0 {
		for loc := entryStart.FindIndex(msg[1:]); loc != nil; loc = entryStart.FindIndex(msg[1:]) {
			// We have found a <99> at location loc[0]+1
			self.reporting <- syslog.NewEntry(msg[0:loc[0]+1], c.RemoteAddr())
			remaining := copy(msg, msg[loc[0]+1:])
			msg = msg[:remaining] // Note remaining > 0 since matches entryStart
		}
	}
	return msg
}

// Read connection and break into bytes based on timeout / or start of syslog header
// If we knew there would only be a few messages, we could just buffer everything and parse
// on timeout. As it is, we need to deal with a single read producing a message and possibly
// a message fragment.
func (self *TCPListener) Accept(c net.Conn) {
	msg := make([]byte, 0, 2048)
	for {
		buffer := make([]byte, 1024, 1024)
		c.SetReadDeadline(time.Now().Add(ReadTimeout))
		n, err := c.Read(buffer)
		if n > 0 {
			msg = append(msg, buffer[:n]...)
			msg = self.removeComplete(msg, c)
		}
		if err != nil {
			msg = self.removeComplete(msg, c)
			if len(msg) > 0 {
				self.reporting <- syslog.NewEntry(msg, c.RemoteAddr())
				msg = msg[:0]
			}
			switch {
			case err == io.EOF:
				c.Close()
				return
			case os.IsTimeout(err):
				continue
			default:
				log.Println(err.Error())
				return
			}
		}
	}
}

func (self *TCPListener) Listen() {
	for {
		conn, err := self.listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		go self.Accept(conn)
	}
}
