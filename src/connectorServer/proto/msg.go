package proto

type Msg interface{
	Encode() []byte
	Decode(msg []byte)
}
