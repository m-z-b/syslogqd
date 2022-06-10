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
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/m-z-b/syslogqd/internal/syslog"
)

// Maximum number of bytes we will accept in a UDP message - although in theory this
// is 65527 (+ 8 bytes header ) - in practice a packet that size will be split
const maxUDPmessage = 512

// UDPListener listens for incoming syslog messages
//
// Create with NewUDPListener(), then
// call the Listen() function to handle syslog messages
type UDPListener struct {
	port      int
	sock      *net.UDPConn
	reporting syslog.Channel
}

// NewUDPListener returns a new UDPListener on the given port
// If the listener can't be created, an error is returned as
// the second return value
func NewUDPListener(port int, reporting syslog.Channel) (*UDPListener, error) {
	u := UDPListener{port: port, reporting: reporting}

	var err error
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to resolve UDP port %d: %s", port, err.Error()))
	}
	u.sock, err = net.ListenUDP("udp", addr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to listen to UDP port %d: %s\n", port, err.Error()))
	}
	return &u, nil
}

func (self *UDPListener) Listen() {
	for {
		buf := make([]byte, 512) // TODO Check max UDP. Reuse of buffers
		nBytes, remoteAddress, err := self.sock.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Socket Read Error: %s", err.Error())
			continue
		}
		for ; nBytes > 0 && buf[nBytes-1] == '\n'; nBytes-- {
		}
		if nBytes > 0 {
			self.reporting <- syslog.NewEntry(buf[0:nBytes], remoteAddress)
		}
	}
}
