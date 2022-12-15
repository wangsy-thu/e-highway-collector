package flux

import "e-highway-collector/flux"

type LineWriter interface {
	Write(bytes flux.Line)
}
