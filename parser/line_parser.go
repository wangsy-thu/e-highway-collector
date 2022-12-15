package parser

import (
	reply2 "e-highway-collector/core/reply"
	"e-highway-collector/flux"
	"e-highway-collector/interface/reply"
	"e-highway-collector/lib/logger"
	"fmt"
	"strconv"
	"strings"
)

// ParseLine 协议解析器
// Rule: *FLUX$measurement,tag1=t1,tag2=t2 field1=value1,field2=value2 timestamp
// Example: *FLUX$testM,gateId=1,tag1=flux speed=134,plate=A36435 1671021217\n
func ParseLine(msg []byte, lineCh chan flux.Line) reply.Reply {
	l := string(msg[6 : len(msg)-2])
	logger.Info("receive line: " + l)
	// 1,按空格切分字符串
	strArr := strings.Split(l, " ")
	if len(strArr) > 3 || len(strArr) < 2 {
		return reply2.MakeErrReply("-token length error\n")
	}
	// 2,解析时间戳
	var ts int
	var err error
	if len(strArr) == 3 {
		ts, err = strconv.Atoi(strArr[2])
		fmt.Println(ts)
		if err != nil {
			logger.Error("protocol")
			return reply2.MakeErrReply("protocol error\n")
		}
	}
	//3,解析measurement和tag
	meaWithTag := strings.Split(strArr[0], ",")
	measurement := meaWithTag[0]
	var tags map[string]string
	if len(meaWithTag) > 1 {
		tags = parseTags(meaWithTag[1:])
	}

	//4,解析fields
	fieldsStr := strings.Split(strArr[1], ",")
	var fields map[string]interface{}
	fields = parseFields(fieldsStr)

	p := flux.Line{
		Measurement: measurement,
		Tags:        tags,
		Fields:      fields,
		Timestamp:   ts,
	}
	lineCh <- p

	return reply2.MakeOkReply()
}

func parseTags(tags []string) map[string]string {
	res := make(map[string]string, len(tags))
	for _, tagStr := range tags {
		kv := strings.Split(tagStr, "=")
		res[kv[0]] = kv[1]
	}
	return res
}

func parseFields(fields []string) map[string]interface{} {
	res := make(map[string]interface{})
	for _, tagStr := range fields {
		kv := strings.Split(tagStr, "=")
		var value interface{}
		iValue, err := strconv.Atoi(kv[1])
		if err != nil {
			value = kv[1]
		} else {
			value = iValue
		}
		res[kv[0]] = value
	}
	return res
}
