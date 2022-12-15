package flux

type Line struct {
	Measurement string
	Tags        map[string]string
	Fields      map[string]interface{}
	Timestamp   int
}
