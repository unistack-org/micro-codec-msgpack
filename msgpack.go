// Package msgpack provides a msgpack codec
package msgpack // import "go.unistack.org/micro-codec-msgpack/v3"

import (
	"io"

	"github.com/vmihailenco/msgpack/v5"
	"go.unistack.org/micro/v3/codec"
	rutil "go.unistack.org/micro/v3/util/reflect"
)

type msgpackCodec struct {
	opts codec.Options
}

const (
	flattenTag = "flatten"
)

func (c *msgpackCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, err := rutil.StructFieldByTag(v, options.TagName, flattenTag); err == nil {
		v = nv
	}

	if m, ok := v.(*codec.Frame); ok {
		return m.Data, nil
	}

	return msgpack.Marshal(v)
}

func (c *msgpackCodec) Unmarshal(b []byte, v interface{}, opts ...codec.Option) error {
	if len(b) == 0 || v == nil {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, err := rutil.StructFieldByTag(v, options.TagName, flattenTag); err == nil {
		v = nv
	}

	if m, ok := v.(*codec.Frame); ok {
		m.Data = b
		return nil
	}

	return msgpack.Unmarshal(b, v)
}

func (c *msgpackCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *msgpackCodec) ReadBody(conn io.Reader, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := io.ReadAll(conn)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return nil
	}

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
	}

	return c.Unmarshal(buf, v)
}

func (c *msgpackCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := c.Marshal(v)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf)
	return err
}

func (c *msgpackCodec) String() string {
	return "msgpack"
}

func NewCodec(opts ...codec.Option) codec.Codec {
	return &msgpackCodec{opts: codec.NewOptions(opts...)}
}
