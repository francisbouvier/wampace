package serializer

type Serializer interface {
	Marshal([]interface{}) ([]byte, error)
	Unmarshal([]byte) ([]interface{}, error)
}
