/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package libnan

import (
	"errors"
	"fmt"
	//"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/ipv4"

	"time"
)

// ===================================================================================================

var ()

// Returns nil in case of success
func Ping(_targetip string) error {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return err
	}
	defer c.Close()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		return err
	}

	addr := net.ParseIP(_targetip)

	if len(addr) == 0 {
		return errors.New("Expected an IP address as parameter (DNS resolution not done)")
	}

	if _, err := c.WriteTo(wb, &net.IPAddr{IP: addr}); err != nil {
		return err
	}

	c.SetReadDeadline(time.Now().Add(2 * time.Second))

	rb := make([]byte, 1500)
	n, _ /*peer*/, err := c.ReadFrom(rb)
	if err != nil {
		return err
	}

	rm, err := icmp.ParseMessage(iana.ProtocolICMP, rb[:n])
	if err != nil {
		return err
	}

	if rm.Type == ipv4.ICMPTypeEchoReply {
		return nil
	}

	msg := fmt.Sprintf("Received a response to our ping request but it's not an ICMP message: %s", rm.Type)
	return errors.New(msg)
}
