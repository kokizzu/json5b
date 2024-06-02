// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json5b

import (
	"math"
	"testing"
)

func TestNumberIsValid(t *testing.T) {
	validTests := []string{
		"0",
		"-0",
		"1",
		"+1",
		"-1",
		"0.1",
		"-0.1",
		"+0.1",
		"1234",
		"-1234",
		"+1234",
		"12.34",
		"-12.34",
		"+12.34",
		"12E0",
		"12E1",
		"12e34",
		"12E-0",
		"12e+1",
		"12e-34",
		"-12E0",
		"-12E1",
		"-12e34",
		"-12E-0",
		"-12e+1",
		"-12e-34",
		"+12E0",
		"+12E1",
		"+12e34",
		"+12E-0",
		"+12e+1",
		"+12e-34",
		"1.2E0",
		"1.2E1",
		"1.2e34",
		"1.2E-0",
		"1.2e+1",
		"1.2e-34",
		"-1.2E0",
		"-1.2E1",
		"-1.2e34",
		"-1.2E-0",
		"-1.2e+1",
		"-1.2e-34",
		"+1.2E0",
		"+1.2E1",
		"+1.2e34",
		"+1.2E-0",
		"+1.2e+1",
		"+1.2e-34",
		"0E0",
		"0E1",
		"0e34",
		"0E-0",
		"0e+1",
		"0e-34",
		"-0E0",
		"-0E1",
		"-0e34",
		"-0E-0",
		"-0e+1",
		"-0e-34",
		"+0E0",
		"+0E1",
		"+0e34",
		"+0E-0",
		"+0e+1",
		"+0e-34",
		"1.",
		"+1.",
		"-1.",
		"1.e1",
		"+1.e1",
		"-1.e1",
		".5",
		"-.5",
		"+.5",
		"0xa",
		"0xA",
		"0XA",
		"-0XA",
		"+0XA",
		"0x0",
		"0x0",
		"0X0",
		"-0X0",
		"+0X0",
		"0x2",
		"0x2",
		"0X2",
		"-0X2",
		"+0X2",
		"0xDEADBeef",
		"-0xDEADBeef",
		"+0xDEADBeef",
		"0XDEADBeef",
		"-0XDEADBeef",
		"+0XDEADBeef",
		"0xDEAD3eef",
		"-0xDEAD3eef",
		"+0xDEAD3eef",
		"0XDEAD3eef",
		"-0XDEAD3eef",
		"+0XDEAD3eef",
		"NaN",
		"+Infinity",
		"-Infinity",
		"Infinity",
	}

	for _, test := range validTests {
		if !isValidNumber(test) {
			t.Errorf("%s should be valid", test)
		}

		var f float64
		if err := Unmarshal([]byte(test), &f); err != nil {
			t.Errorf("%s should be valid but Unmarshal failed: %v", test, err)
		}
	}

	invalidTests := []string{
		"",
		"invalid",
		"1.0.1",
		"1..1",
		"-1-2",
		"012a42",
		"01.2",
		"012",
		"12E12.12",
		"1e2e3",
		"1e+-2",
		"1e--23",
		"1e",
		"e1",
		"1e+",
		"1ea",
		"1a",
		"1.a",
		"01",
		"0xDsADBeef",
		".0xDEADBeef",
		"0XDsADBeef",
		".0XDEADBeef",
		"+NaN",
		"-NaN",
		".NaN",
		".Infinity",
		"0xs",
	}

	for _, test := range invalidTests {
		if isValidNumber(test) {
			t.Errorf("%s should be invalid", test)
		}

		var f float64
		if err := Unmarshal([]byte(test), &f); err == nil {
			t.Errorf("%s should be invalid but unmarshal wrote %v", test, f)
		}
	}
}

func BenchmarkNumberIsValid(b *testing.B) {
	s := "-61657.61667E+61673"
	for i := 0; i < b.N; i++ {
		isValidNumber(s)
	}
}

func TestNumberFloat64(t *testing.T) {
	tests := map[string]float64{
		"0xDeADb":   0xdeadb,
		"+0xDeADb":  0xdeadb,
		"-0xDeADb":  -0xdeadb,
		"-0XDeADb":  -0xdeadb,
		"-0x0":      math.Copysign(0, -1),
		".5":        0.5,
		"-.5":       -0.5,
		"+1.e1":     1.e1,
		"-1.e1":     -1.e1,
		"-0":        math.Copysign(0, -1),
		"-Infinity": math.Inf(-1),
		"Infinity":  math.Inf(0),
		"+Infinity": math.Inf(1),
		"NaN":       math.NaN(),
	}

	for s, f := range tests {
		res, err := Number(s).Float64()
		if err != nil {
			t.Errorf("failed to parse %s: %s", s, err)
		}
		if s == "NaN" {
			if !math.IsNaN(res) {
				t.Errorf("expected NaN")
			}
		} else {
			if res != f {
				t.Errorf("wanted %v, got %v", f, res)
			}
		}
	}
}

func TestNumberInt64(t *testing.T) {
	tests := map[string]int64{
		"0xDeADb":  0xdeadb,
		"+0xDeADb": 0xdeadb,
		"-0xDeADb": -0xdeadb,
		"-0XDeADb": -0xdeadb,
		"0x0":      0,
	}

	for s, i := range tests {
		res, err := Number(s).Int64()
		if err != nil {
			t.Errorf("failed to parse %s: %s", s, err)
		}
		if res != i {
			t.Errorf("wanted %v, got %v", i, res)
		}
	}
}

func TestDecodeStringFloat(t *testing.T) {
	tc := `{
fr: 1.2,
fs: '2.3',
fd: "4.5",
}`
	t.Run(`toMap`, func(t *testing.T) {
		m := make(map[string]any)
		err := Unmarshal([]byte(tc), &m)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if m["fr"].(float64) != 1.2 {
			t.Errorf("fr: %v", m["fr"])
		}
		if m["fs"].(string) != `2.3` {
			t.Errorf("fs: %v", m["fs"])
		}
		if m["fd"].(string) != `4.5` {
			t.Errorf("fd: %v", m["fd"])
		}
	})

	t.Run(`toFloat32Struct`, func(t *testing.T) {
		type Floats struct {
			Fr float32 `json5:"fr"`
			Fs float32 `json5:"fs"`
			Fd float32 `json5:"fd"`
		}
		f := Floats{}
		err := Unmarshal([]byte(tc), &f)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if f.Fr != 1.2 {
			t.Errorf("fr: %v", f.Fr)
		}
		if f.Fs != 2.3 {
			t.Errorf("fs: %v", f.Fs)
		}
		if f.Fd != 4.5 {
			t.Errorf("fd: %v", f.Fd)
		}
	})

	t.Run(`toFloat64Struct`, func(t *testing.T) {
		type Floats struct {
			Fr float64 `json5:"fr"`
			Fs float64 `json5:"fs"`
			Fd float64 `json5:"fd"`
		}
		f := Floats{}
		err := Unmarshal([]byte(tc), &f)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if f.Fr != 1.2 {
			t.Errorf("fr: %v", f.Fr)
		}
		if f.Fs != 2.3 {
			t.Errorf("fs: %v", f.Fs)
		}
		if f.Fd != 4.5 {
			t.Errorf("fd: %v", f.Fd)
		}
	})
}

func TestDecodeStringFloatJsonTag(t *testing.T) {
	tc := `{
fr: 1.2,
fs: '2.3',
fd: "4.5",
}`

	t.Run(`toFloat32Struct`, func(t *testing.T) {
		type Floats struct {
			Fr float32 `json:"fr"`
			Fs float32 `json:"fs"`
			Fd float32 `json:"fd"`
		}
		f := Floats{}
		err := Unmarshal([]byte(tc), &f)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if f.Fr != 1.2 {
			t.Errorf("fr: %v", f.Fr)
		}
		if f.Fs != 2.3 {
			t.Errorf("fs: %v", f.Fs)
		}
		if f.Fd != 4.5 {
			t.Errorf("fd: %v", f.Fd)
		}
	})

	t.Run(`toFloat64Struct`, func(t *testing.T) {
		type Floats struct {
			Fr float64 `json:"fr"`
			Fs float64 `json:"fs"`
			Fd float64 `json:"fd"`
		}
		f := Floats{}
		err := Unmarshal([]byte(tc), &f)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if f.Fr != 1.2 {
			t.Errorf("fr: %v", f.Fr)
		}
		if f.Fs != 2.3 {
			t.Errorf("fs: %v", f.Fs)
		}
		if f.Fd != 4.5 {
			t.Errorf("fd: %v", f.Fd)
		}
	})
}

func TestDecodeStringInt(t *testing.T) {
	tc := `{
ir: 1,
is: '2',
id: "3",
}`
	t.Run(`toMap`, func(t *testing.T) {
		m := make(map[string]any)
		err := Unmarshal([]byte(tc), &m)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if m["ir"].(float64) != 1 {
			t.Errorf("ir: %v", m["ir"])
		}
		if m["is"].(string) != `2` {
			t.Errorf("is: %v", m["is"])
		}
		if m["id"].(string) != `3` {
			t.Errorf("id: %v", m["id"])
		}
	})

	t.Run(`toInt8Struct`, func(t *testing.T) {
		type Ints struct {
			Ir int8 `json5:"ir"`
			Is int8 `json5:"is"`
			Id int8 `json5:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toInt16Struct`, func(t *testing.T) {
		type Ints struct {
			Ir int16 `json5:"ir"`
			Is int16 `json5:"is"`
			Id int16 `json5:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toInt32Struct`, func(t *testing.T) {
		type Ints struct {
			Ir int32 `json5:"ir"`
			Is int32 `json5:"is"`
			Id int32 `json5:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toInt64Struct`, func(t *testing.T) {
		type Ints struct {
			Ir int64 `json5:"ir"`
			Is int64 `json5:"is"`
			Id int64 `json5:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toUint8Struct`, func(t *testing.T) {
		type Ints struct {
			Ir uint8 `json5:"ir"`
			Is uint8 `json5:"is"`
			Id uint8 `json5:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toUint16Struct`, func(t *testing.T) {
		type Ints struct {
			Ir uint16 `json5:"ir"`
			Is uint16 `json5:"is"`
			Id uint16 `json5:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toUint32Struct`, func(t *testing.T) {
		type Ints struct {
			Ir uint32 `json5:"ir"`
			Is uint32 `json5:"is"`
			Id uint32 `json5:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toUint64Struct`, func(t *testing.T) {
		type Ints struct {
			Ir uint64 `json5:"ir"`
			Is uint64 `json5:"is"`
			Id uint64 `json5:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toUintStruct`, func(t *testing.T) {
		type Ints struct {
			Ir uint `json5:"ir"`
			Is uint `json5:"is"`
			Id uint `json5:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toIntStruct`, func(t *testing.T) {
		type Ints struct {
			Ir int `json5:"ir"`
			Is int `json5:"is"`
			Id int `json5:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})
}

func TestDecodeStringIntJsonTag(t *testing.T) {
	tc := `{
ir: 1,
is: '2',
id: "3",
}`
	t.Run(`toMap`, func(t *testing.T) {
		m := make(map[string]any)
		err := Unmarshal([]byte(tc), &m)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if m["ir"].(float64) != 1 {
			t.Errorf("ir: %v", m["ir"])
		}
		if m["is"].(string) != `2` {
			t.Errorf("is: %v", m["is"])
		}
		if m["id"].(string) != `3` {
			t.Errorf("id: %v", m["id"])
		}
	})

	t.Run(`toInt8Struct`, func(t *testing.T) {
		type Ints struct {
			Ir int8 `json:"ir"`
			Is int8 `json:"is"`
			Id int8 `json:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toInt16Struct`, func(t *testing.T) {
		type Ints struct {
			Ir int16 `json:"ir"`
			Is int16 `json:"is"`
			Id int16 `json:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toInt32Struct`, func(t *testing.T) {
		type Ints struct {
			Ir int32 `json:"ir"`
			Is int32 `json:"is"`
			Id int32 `json:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toInt64Struct`, func(t *testing.T) {
		type Ints struct {
			Ir int64 `json:"ir"`
			Is int64 `json:"is"`
			Id int64 `json:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toUint8Struct`, func(t *testing.T) {
		type Ints struct {
			Ir uint8 `json:"ir"`
			Is uint8 `json:"is"`
			Id uint8 `json:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toUint16Struct`, func(t *testing.T) {
		type Ints struct {
			Ir uint16 `json:"ir"`
			Is uint16 `json:"is"`
			Id uint16 `json:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toUint32Struct`, func(t *testing.T) {
		type Ints struct {
			Ir uint32 `json:"ir"`
			Is uint32 `json:"is"`
			Id uint32 `json:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toUint64Struct`, func(t *testing.T) {
		type Ints struct {
			Ir uint64 `json:"ir"`
			Is uint64 `json:"is"`
			Id uint64 `json:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toUintStruct`, func(t *testing.T) {
		type Ints struct {
			Ir uint `json:"ir"`
			Is uint `json:"is"`
			Id uint `json:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})

	t.Run(`toIntStruct`, func(t *testing.T) {
		type Ints struct {
			Ir int `json:"ir"`
			Is int `json:"is"`
			Id int `json:"id"`
		}
		i := Ints{}
		err := Unmarshal([]byte(tc), &i)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if i.Ir != 1 {
			t.Errorf("ir: %v", i.Ir)
		}
		if i.Is != 2 {
			t.Errorf("is: %v", i.Is)
		}
		if i.Id != 3 {
			t.Errorf("id: %v", i.Id)
		}
	})
}
