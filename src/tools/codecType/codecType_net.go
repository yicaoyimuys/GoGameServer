package codecType

import (
	"io"
	"github.com/funny/link"
	"io/ioutil"
)


type NetCodecType struct {
}

func (codec NetCodecType) NewEncoder(w io.Writer) link.Encoder {
	return &netCodecTypeEncoder{w}
}

func (codec NetCodecType) NewDecoder(r io.Reader) link.Decoder {
	return &netCodecTypeDecoder{r}
}





type netCodecTypeEncoder struct {
	w io.Writer
}

func (encoder *netCodecTypeEncoder) Encode(msg interface{}) error {
	_, err := encoder.w.Write(msg.([]byte))
	return err
}





type netCodecTypeDecoder struct {
	r io.Reader
}

func (decoder *netCodecTypeDecoder) Decode(msg interface{}) error {
	d, err := ioutil.ReadAll(decoder.r)
	if err != nil {
		return err
	}
	*(msg.(*[]byte)) = d
	return nil
}