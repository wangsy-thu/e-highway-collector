package reply

type OkReply struct {
}

func (o *OkReply) ToBytes() []byte {
	return []byte("+ok\n")
}

func MakeOkReply() *OkReply {
	return &OkReply{}
}
