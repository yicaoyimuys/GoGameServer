package codecType

import (
	"io"
	"github.com/funny/link"
)

type DummyCodecTypeInterface interface {
	ReadMsg(msg *[]byte) error
	WriteMsg(msg []byte) error
}


type DummyCodecType struct {
}

func (codec DummyCodecType) NewEncoder(w io.Writer) link.Encoder {
	return &dummyCodecTypeEncoder{w}
}

func (codec DummyCodecType) NewDecoder(r io.Reader) link.Decoder {
	return &dummyCodecTypeDecoder{r}
}




type dummyCodecTypeEncoder struct {
	w io.Writer
}

func (encoder *dummyCodecTypeEncoder) Encode(msg interface{}) error {
	return encoder.w.(DummyCodecTypeInterface).WriteMsg(msg.([]byte))
}





type dummyCodecTypeDecoder struct {
	r io.Reader

}

func (decoder *dummyCodecTypeDecoder) Decode(msg interface{}) error {
	return decoder.r.(DummyCodecTypeInterface).ReadMsg(msg.(*[]byte))
}