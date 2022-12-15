package reply

type ErrorReply struct {
	Reason string
}

func (o *ErrorReply) ToBytes() []byte {
	return []byte(o.Reason)
}

func MakeErrReply(reason string) *ErrorReply {
	return &ErrorReply{
		Reason: reason,
	}
}
