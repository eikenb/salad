package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// message is structure to hold the json and serialize to/from binary
//
// camelCase is usually the standard in Go but I feel it is better to stick
// close to the spec when there is one and it primarily uses snake_case.
type message struct {
	tail_number                                string
	engine_count                               int
	engine_name                                string
	latitude, longitude, altitude, temperature float64
}

func (m *message) UnmarshalBinary(data []byte) error {
	if len(data) < 3 {
		return fmt.Errorf("message unmarhsal binary: no data")
	}

	buf := bytes.NewReader(data)
	// trim off header
	bs := make([]byte, 3)
	if err := binary.Read(buf, order, &bs); err != nil {
		return fmt.Errorf("header: %w", err)
	}
	// tail_number
	var size uint32
	if err := binary.Read(buf, order, &size); err != nil {
		return fmt.Errorf("tail_number(size): %w", err)
	}
	bs = make([]byte, size)
	if err := binary.Read(buf, order, &bs); err != nil {
		return fmt.Errorf("tail_number: %w", err)
	}
	m.tail_number = string(bs)
	// engine_count
	var engine_count uint32
	if err := binary.Read(buf, order, &engine_count); err != nil {
		return fmt.Errorf("engine_count: %w", err)
	}
	m.engine_count = int(engine_count)
	// engine_name
	size = 0
	if err := binary.Read(buf, order, &size); err != nil {
		return fmt.Errorf("engine_name(size)): %w", err)
	}
	bs = make([]byte, size)
	if err := binary.Read(buf, order, &bs); err != nil {
		return fmt.Errorf("engine_name: %w", err)
	}
	m.engine_name = string(bs)
	// latitude
	if err := binary.Read(buf, order, &m.latitude); err != nil {
		return fmt.Errorf("latitude: %w", err)
	}
	// longitude
	if err := binary.Read(buf, order, &m.longitude); err != nil {
		return fmt.Errorf("longitude: %w", err)
	}
	// altitude
	if err := binary.Read(buf, order, &m.altitude); err != nil {
		return fmt.Errorf("altitude: %w", err)
	}
	// temperature
	if err := binary.Read(buf, order, &m.temperature); err != nil {
		return fmt.Errorf("temperature: %w", err)
	}
	return nil
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

func extendMsg(msg []byte, data any) ([]byte, error) {
	bs, err := encode(data)
	if err != nil {
		return nil, err
	}
	return append(msg, bs...), nil

}

func encode(data any) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, order, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
