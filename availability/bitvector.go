package availability

import (
	"bytes"
	"encoding/json"
	"time"
)

type BitVector struct {
	Resolution         TimeResolution `json:"resolution"`
	InternalResolution TimeResolution `json:"internal_resolution"`
	From               time.Time      `json:"from"`
	To                 time.Time      `json:"to"`
	Data               []byte         `json:"available"`
}

func NewBitVector(res, intRes TimeResolution, data []byte, from time.Time) *BitVector {
	to := from.Add(time.Duration(len(data)*int(res)) * time.Second)
	return &BitVector{
		Resolution:         res,
		InternalResolution: res,
		Data:               data,
		From:               from,
		To:                 to,
	}
}

func (b *BitVector) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("BitVector {")
	buffer.WriteString("res: ")
	buffer.WriteString(b.Resolution.String())
	buffer.WriteString(", ")
	buffer.WriteString("from: ")
	buffer.WriteString(b.From.Format(time.RFC3339))
	buffer.WriteString(",")
	buffer.WriteRune('\n')
	buffer.WriteString("data:[")
	for i := 0; i < len(b.Data); i++ {
		if b.Data[i] == 1 {
			buffer.WriteString("1, ")
		} else {
			buffer.WriteString("0, ")
		}

	}
	buffer.WriteString("],")
	buffer.WriteRune('\n')
	buffer.WriteRune('}')

	return buffer.String()
}

func (bv *BitVector) MarshalJSON() ([]byte, error) {

	intdata := make([]int, len(bv.Data))
	for k, v := range bv.Data {
		intdata[k] = int(v)
	}

	return json.Marshal(struct {
		Resolution         string    `json:"resolution"`
		InternalResolution string    `json:"internal_resolution"`
		From               time.Time `json:"from"`
		To                 time.Time `json:"to"`
		Data               []int     `json:"available"`
	}{
		Resolution:         bv.Resolution.String(),
		InternalResolution: bv.InternalResolution.String(),
		From:               bv.From,
		To:                 bv.To,
		Data:               intdata,
	})
}

func (bitVector *BitVector) All() bool {
	for _, b := range bitVector.Data {
		if b == 0 {
			return false
		}
	}
	return true
}

func (bitVector *BitVector) Any() bool {
	for _, b := range bitVector.Data {
		if b == 1 {
			return true
		}
	}
	return false
}

func (bitVector *BitVector) Count() int {
	count := 0
	for _, b := range bitVector.Data {
		if b == 1 {
			count++
		}
	}
	return count
}
