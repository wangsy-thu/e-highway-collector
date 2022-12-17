package flux

type Line struct {
	Measurement string                 `json:"measurement"`
	Tags        map[string]string      `json:"tags"`
	Fields      map[string]interface{} `json:"Fields"`
	Timestamp   int                    `json:"timestamp"`
}
