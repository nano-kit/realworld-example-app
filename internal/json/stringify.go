package json

import jsoniter "github.com/json-iterator/go"

func Stringify(v interface{}) (s string) {
	s, _ = jsoniter.MarshalToString(v)
	return
}
