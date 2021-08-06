package convert

import "strconv"

type StrTo string

func (this StrTo) String() string {
	return string(this)
}

func (this StrTo) Int() (int, error) {
	return strconv.Atoi(this.String())
}

func (this StrTo) MustInt() int {
	v, _ := this.Int()
	return v
}

func (this StrTo) UInt32() (uint32, error) {
	v, err := strconv.Atoi(this.String())
	return uint32(v), err
}

func (this StrTo) MustUInt32() uint32 {
	v, _ := this.UInt32()
	return v
}