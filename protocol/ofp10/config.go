package ofp10

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/jonstout/ogo/protocol/ofpxx"
)

func NewConfigRequest() *ofpxx.Header {
	h := ofpxx.NewOfp10Header()
	h.Type = T_GET_CONFIG_REQUEST
	return &h
}

// ofp_config_flags 1.0
const (
	C_FRAG_NORMAL = 0
	C_FRAG_DROP   = 1
	C_FRAG_REASM  = 2
	C_FRAG_MASK   = 3
)

// ofp_switch_config 1.0
type SwitchConfig struct {
	Header      ofpxx.Header
	Flags       uint16 // OFPC_* flags
	MissSendLen uint16
}

func NewSetConfig() *SwitchConfig {
	h := ofpxx.NewOfp10Header()
	h.Length = 12
	h.Type = T_SET_CONFIG

	c := new(SwitchConfig)
	c.Header = h
	c.Flags = 0
	c.MissSendLen = 0
	return c
}

func (c *SwitchConfig) Len() (n uint16) {
	return c.Header.Length
}

func (c *SwitchConfig) GetHeader() *ofpxx.Header {
	return &c.Header
}

func (c *SwitchConfig) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, c)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (c *SwitchConfig) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	err = c.Header.UnmarshelBinary(buf.Next(8))

	if err = binary.Read(buf, binary.BigEndian, &c.Flags); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &c.MissSendLen); err != nil {
		return
	}
	n += 2
	return
}
