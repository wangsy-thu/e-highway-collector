package parser

import (
	"e-highway-collector/flux"
	"e-highway-collector/lib/logger"
)

func ParseLine(msg []byte, lineCh chan flux.Line) {
	// TODO: parse bytes to Line(Point)
	msg = msg[:len(msg)-1]
	logger.Info(string(msg))
	lineCh <- flux.Line(msg)
}
