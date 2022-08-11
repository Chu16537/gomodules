package json

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

func Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return jsoniter.MarshalIndent(v, prefix, indent)
}

func MarshalToString(v interface{}) (string, error) {
	return jsoniter.MarshalToString(v)
}

func UnmarshalFromString(str string, v interface{}) error {
	return jsoniter.UnmarshalFromString(str, v)
}

func Get(data []byte, path ...interface{}) jsoniter.Any {
	return jsoniter.Get(data, path...)
}

func NewEncoder(writer io.Writer) *jsoniter.Encoder {
	return jsoniter.NewEncoder(writer)
}

func NewDecoder(reader io.Reader) *jsoniter.Decoder {
	return jsoniter.NewDecoder(reader)
}

func Valid(data []byte) bool {
	return jsoniter.Valid(data)
}

func RegisterExtension(extension jsoniter.Extension) {
	jsoniter.RegisterExtension(extension)
}
