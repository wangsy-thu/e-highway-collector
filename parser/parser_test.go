package parser

import (
	"e-highway-collector/flux"
	"fmt"
	"sync"
	"testing"
)

func TestParseLine(t *testing.T) {
	var p flux.Line
	fmt.Print("hello test")
	ch := make(chan flux.Line)
	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		p = <-ch
		fmt.Printf("point(measurement=%s, tags=%v, fields=%v, timestamp=%d)",
			p.Measurement, p.Tags, p.Fields, p.Timestamp)
		group.Done()
	}()
	ParseLine([]byte("*FLUX$testM,gateId=1,tag1=flux speed=134,plate=A36435 1671021217\\n"), ch)
	group.Wait()
}
