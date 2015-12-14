package codetype

import (
	"io"
	"github.com/funny/link"
)

type CustomCodecTypeInterface interface {
	ReadOne(msg *[]byte) error
}


type CustomCodecType struct {
}

func (codec CustomCodecType) NewEncoder(w io.Writer) link.Encoder {
	return &customCodecTypeEncoder{w}
}

func (codec CustomCodecType) NewDecoder(r io.Reader) link.Decoder {
	return &customCodecTypeDecoder{r}
}




type customCodecTypeEncoder struct {
	w io.Writer
}

func (encoder *customCodecTypeEncoder) Encode(msg interface{}) error {
	_, err := encoder.w.Write(msg.([]byte))
	return err
}





type customCodecTypeDecoder struct {
	r io.Reader

}

func (decoder *customCodecTypeDecoder) Decode(msg interface{}) error {
	err := decoder.r.(CustomCodecTypeInterface).ReadOne(msg.(*[]byte))
	if err != nil {
		return err
	}
	return nil
}