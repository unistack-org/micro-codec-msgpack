package msgpack

import (
	"bytes"
	"testing"

	"go.unistack.org/micro/v3/codec"
)

func TestFrame(t *testing.T) {
	s := &codec.Frame{Data: []byte("test")}

	buf, err := NewCodec().Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte(`test`)) {
		t.Fatalf("bytes not equal %s != %s", buf, `test`)
	}
}

func TestFrameFlatten(t *testing.T) {
	s := &struct {
		One  string
		Name *codec.Frame `json:"name" codec:"flatten"`
	}{
		One:  "xx",
		Name: &codec.Frame{Data: []byte("test")},
	}

	buf, err := NewCodec().Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte(`test`)) {
		t.Fatalf("bytes not equal %s != %s", buf, `test`)
	}
}

func TestStructByTag(t *testing.T) {
	type Str struct {
		Name []string `msgpack:"name" codec:"flatten"`
	}

	val := &Str{Name: []string{"first", "second"}}

	c := NewCodec()
	buf, err := c.Marshal(val)
	if err != nil {
		t.Fatal(err)
	}

	nval := &Str{}
	err = c.Unmarshal(buf, nval)
	if err != nil {
		t.Fatal(err)
	}

	if len(nval.Name) != 2 {
		t.Fatalf("invalid unmarshal: %v", nval)
	}
}

func TestReadBody(t *testing.T) {
	s := &struct {
		Name string
	}{}
	c := NewCodec()
	b := bytes.NewReader(nil)
	err := c.ReadBody(b, s)
	if err != nil {
		t.Fatal(err)
	}
}
