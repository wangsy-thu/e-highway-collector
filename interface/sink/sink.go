package sink

import "e-highway-collector/flux"

type Sink interface {
	Send(msg flux.Line)
}
