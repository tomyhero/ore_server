package util

import (
	"github.com/stretchr/testify/assert"
	//	"reflect"
	"testing"
)

/*
func TestGetMethods(t *testing.T) {
	methods := GetMethods(&sample.EchoHandler{})
	c := &server.Context{}
	methods["prefix_Echo"].Call([]reflect.Value{reflect.ValueOf(c)})
	assert.NotNil(t, methods)
}
*/

func TestMessagePack(t *testing.T) {
	in := map[string]interface{}{"h": map[string]interface{}{"cmd": "prefix_Echo"}, "b": map[string]interface{}{"text": "Hello World\n"}}
	buf, err := PackMP(in)
	assert.Nil(t, err)
	out, err := UnpackMP(buf)
	assert.Nil(t, err)
	assert.Equal(t, "prefix_Echo", out["h"].(map[string]interface{})["cmd"])
}
