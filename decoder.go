package gotiny

import (
	"reflect"
	"unsafe"
)

type Decoder struct {
	buf     []byte //buf
	index   int    //下一个要读取的字节
	boolean byte   //下一次要读取的bool在buf中的下标,即buf[boolPos]
	boolBit byte   //下一次要读取的bool的buf[boolPos]中的bit位

	decEngs []decEng //解码器集合
	length  int      //解码器数量
}

func Decodes(buf []byte, is ...interface{}) int {
	d := NewDecoderWithPtr(is...)
	d.buf = buf
	d.Decodes(is...)
	return d.index
}

func NewDecoderWithPtr(is ...interface{}) *Decoder {
	l := len(is)
	if l < 1 {
		panic("must have argument!")
	}
	des := make([]decEng, l)
	for i := 0; i < l; i++ {
		des[i] = getDecEngine(reflect.TypeOf(is[i]).Elem())
	}
	return &Decoder{
		length:  l,
		decEngs: des,
	}
}

func NewDecoder(is ...interface{}) *Decoder {
	l := len(is)
	if l < 1 {
		panic("must have argument!")
	}
	des := make([]decEng, l)
	for i := 0; i < l; i++ {
		des[i] = getDecEngine(reflect.TypeOf(is[i]))
	}
	return &Decoder{
		length:  l,
		decEngs: des,
	}
}

func NewDecoderWithTypes(ts ...reflect.Type) *Decoder {
	l := len(ts)
	if l < 1 {
		panic("must have argument!")
	}
	des := make([]decEng, l)
	for i := 0; i < l; i++ {
		des[i] = getDecEngine(ts[i])
	}
	return &Decoder{
		length:  l,
		decEngs: des,
	}
}

func (d *Decoder) Reset() {
	d.index = 0
	d.boolean = 0
	d.boolBit = 0
}

func (d *Decoder) ResetWith(b []byte) {
	d.buf = b
	d.Reset()
}
func (d *Decoder) Decodes(is ...interface{}) {
	l, engs := d.length, d.decEngs
	for i := 0; i < l; i++ {
		engs[i](d, unsafe.Pointer(reflect.ValueOf(is[i]).Pointer()))
	}
}

// is is pointer of value
func (d *Decoder) DecodeByUPtr(ps ...unsafe.Pointer) {
	l, engs := d.length, d.decEngs
	for i := 0; i < l; i++ {
		engs[i](d, ps[i])
	}
}

func (d *Decoder) DecodeValues(vs ...reflect.Value) {
	l, engs := d.length, d.decEngs
	for i := 0; i < l; i++ {
		engs[i](d, unsafe.Pointer(vs[i].UnsafeAddr()))
	}
}
