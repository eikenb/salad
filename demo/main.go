package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

// standalone server to demo process
// starts up, prints out it's port and waits for a connection
// it then returns a couple test packets and exits
// should show off ENV config + backoff working

func main() {
	ln, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	fmt.Println("demo server addr:", ln.Addr().String())

	defer ln.Close()
	for _, msg := range testMessages {
		conn, err := ln.Accept()
		switch {
		case errors.Is(err, context.Canceled):
			return
		case err != nil:
			panic(err)
		}
		// simulate waiting for more data
		time.Sleep(time.Second)
		bin, err := msg.MarshalBinary()
		io.Copy(conn, bytes.NewReader(bin))
		conn.Close()
	}
}

var testMessages []message = []message{
	{
		tail_number:  "N20904",
		engine_count: 2,
		engine_name:  "GEnx-1B",
		latitude:     39.11593389482025,
		longitude:    -67.32425341289998,
		altitude:     36895.5,
		temperature:  -53.2,
	},
	{
		tail_number:  "N20906",
		engine_count: 2,
		engine_name:  "GEnx-1B",
		latitude:     83.31593389482026,
		longitude:    -7.12425341290001,
		altitude:     16895.5,
		temperature:  -13.2,
	},
	{
		tail_number:  "N20907",
		engine_count: 4,
		engine_name:  "GEnx-1C",
		latitude:     -3.31593389482026,
		longitude:    17.12425341290001,
		altitude:     7032.5,
		temperature:  0.2,
	},
}

type message struct {
	tail_number                                string
	engine_count                               int
	engine_name                                string
	latitude, longitude, altitude, temperature float64
}

// MarshalBinary returns the struct in binary form in a byte array
// Wordy on purpose to make the layout very obvious.
func (m message) MarshalBinary() ([]byte, error) {
	var err error
	// header
	msg := []byte{0x41, 0x49, 0x52}

	tail_number_size := uint32(len(m.tail_number))
	if msg, err = extendMsg(msg, tail_number_size); err != nil {
		return nil, fmt.Errorf("tail_number_size: %w", err)
	}
	tail_number_value := []byte(m.tail_number)
	if msg, err = extendMsg(msg, tail_number_value); err != nil {
		return nil, fmt.Errorf("tail_number_value: %w", err)
	}

	engine_count := uint32(m.engine_count)
	if msg, err = extendMsg(msg, engine_count); err != nil {
		return nil, fmt.Errorf("engine_count: %w", err)
	}

	engine_name_size := uint32(len(m.engine_name))
	if msg, err = extendMsg(msg, engine_name_size); err != nil {
		return nil, fmt.Errorf("engine_name_size: %w", err)
	}
	engine_name_value := []byte(m.engine_name)
	if msg, err = extendMsg(msg, engine_name_value); err != nil {
		return nil, fmt.Errorf("engine_name_value: %w", err)
	}

	latitude := m.latitude
	if msg, err = extendMsg(msg, latitude); err != nil {
		return nil, fmt.Errorf("latitude: %w", err)
	}

	longitude := m.longitude
	if msg, err = extendMsg(msg, longitude); err != nil {
		return nil, fmt.Errorf("longitude: %w", err)
	}

	altitude := m.altitude
	if msg, err = extendMsg(msg, altitude); err != nil {
		return nil, fmt.Errorf("altitude: %w", err)
	}

	temperature := m.temperature
	if msg, err = extendMsg(msg, temperature); err != nil {
		return nil, fmt.Errorf("temperature: %w", err)
	}

	return msg, nil
}

var order = binary.BigEndian

// extendMsg returns a `[]byte` slice extended with the `encode`d data
func extendMsg(msg []byte, data any) ([]byte, error) {
	bs, err := encode(data)
	if err != nil {
		return nil, err
	}
	return append(msg, bs...), nil

}

// encode's the data into a `[]byte` slice
func encode(data any) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, order, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
