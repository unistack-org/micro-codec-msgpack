// Package msgpack provides a msgpack codec
package msgpack // import "go.unistack.org/micro-codec-msgpack/v3"

import (
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
	pb "go.unistack.org/micro-proto/v3/codec"
	"go.unistack.org/micro/v3/codec"
	rutil "go.unistack.org/micro/v3/util/reflect"
)

type msgpackCodec struct {
	opts codec.Options
}

func (c *msgpackCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, err := rutil.StructFieldByTag(v, options.TagName, "flatten"); err == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	case codec.RawMessage:
		return []byte(m), nil
	case *codec.RawMessage:
		return []byte(*m), nil
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

	if options.Flatten {
		if nv, err := rutil.StructFieldByTag(v, options.TagName, "flatten"); err == nil {
			v = nv
			rv := reflect.ValueOf(v)
			if rv.Kind() != reflect.Pointer &&
				rv.Kind() != reflect.Map {
				v = reflect.New(rv.Type()).Interface()
			}
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = b
		return nil
	case *pb.Frame:
		m.Data = b
		return nil
	case *codec.RawMessage:
		*m = append((*m)[0:0], b...)
		return nil
	case codec.RawMessage:
		copy(m, b)
		return nil
	}

	return msgpack.Unmarshal(b, v)
}

func (c *msgpackCodec) String() string {
	return "msgpack"
}

func NewCodec(opts ...codec.Option) codec.Codec {
	return &msgpackCodec{opts: codec.NewOptions(opts...)}
}
