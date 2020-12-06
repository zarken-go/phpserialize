package phpserialize

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
)

type EncodeSuite struct {
	suite.Suite
}

func (Suite *EncodeSuite) TestMarshalString() {
	Suite.assertMarshal(`Hello`, `s:5:"Hello";`)
	Suite.assertMarshalContained(`World`, `s:5:"World";`)
}

func (Suite *EncodeSuite) TestMarshalSignedInts() {
	Suite.assertMarshal(12345, `i:12345;`)
	Suite.assertMarshal(int8(-100), `i:-100;`)
	Suite.assertMarshal(int16(2134), `i:2134;`)
	Suite.assertMarshal(int32(-912745), `i:-912745;`)
	Suite.assertMarshal(int64(73912745), `i:73912745;`)

	type container struct {
		v int8 `php:"v"`
	}

	Suite.assertMarshal(container{v: 123}, `a:1:{s:1:"v";i:123;}`)
	Suite.assertMarshalContained(54321, `i:54321;`)
	Suite.assertMarshalContained(int8(120), `i:120;`)
}

func (Suite *EncodeSuite) TestMarshalUnsignedInts() {
	Suite.assertMarshal(uint8(255), `i:255;`)
	Suite.assertMarshal(uint16(4002), `i:4002;`)
	Suite.assertMarshal(uint32(98743), `i:98743;`)
	Suite.assertMarshal(uint64(9702398740), `i:9702398740;`)
	Suite.assertMarshal(uint(23702235398740), `i:23702235398740;`)
}

func (Suite *EncodeSuite) TestMarshalNils() {
	var str *string
	type container2 struct {
		str *string `php:"s"`
	}
	container := struct {
		str *string `php:"s"`
	}{}

	Suite.assertMarshal(str, `N;`)
	Suite.assertMarshal(nil, `N;`)
	Suite.assertMarshal(container, `a:1:{s:1:"s";N;}`)
	Suite.assertMarshal(&container2{}, `a:1:{s:1:"s";N;}`)
	Suite.assertMarshalContained(nil, `N;`)
}

func (Suite *EncodeSuite) TestMarshalFloats() {
	Suite.assertMarshal(15.35, `d:15.35;`)
	Suite.assertMarshal(-19275.1872, `d:-19275.1872;`)
	Suite.assertMarshal(math.Inf(-1), `d:-INF;`)
	Suite.assertMarshal(math.Inf(1), `d:INF;`)
	Suite.assertMarshal(math.NaN(), `d:NAN;`)

	Suite.assertMarshal(float32(math.Inf(-1)), `d:-INF;`)
	Suite.assertMarshal(float32(math.Inf(1)), `d:INF;`)
	Suite.assertMarshal(float32(math.NaN()), `d:NAN;`)

	Suite.assertMarshalContained(12.456, `d:12.456;`)
	Suite.assertMarshalContained(float32(math.NaN()), `d:NAN;`)
}

func (Suite *EncodeSuite) TestMarshalBool() {
	Suite.assertMarshal(true, `b:1;`)
	Suite.assertMarshal(false, `b:0;`)

	Suite.assertMarshalContained(true, `b:1;`)
	Suite.assertMarshalContained(false, `b:0;`)
}

func (Suite *EncodeSuite) TestMarshalSlices() {
	var nilSlice []string
	Suite.assertMarshal([]int{1, 3, 5}, `a:3:{i:0;i:1;i:1;i:3;i:2;i:5;}`)
	Suite.assertMarshal([]string{`one`, `three`, `five`}, `a:3:{i:0;s:3:"one";i:1;s:5:"three";i:2;s:4:"five";}`)
	Suite.assertMarshal(nilSlice, `N;`)
	Suite.assertMarshal([]string{}, `a:0:{}`)
}

func (Suite *EncodeSuite) TestMarshalIntKeyMap() {
	m := make(map[int]string)
	m[45] = `Hello`
	m[17] = `World`

	b, err := Marshal(m)
	Suite.Nil(err)

	expectedA := `a:2:{i:45;s:5:"Hello";i:17;s:5:"World";}`
	expectedB := `a:2:{i:17;s:5:"World";i:45;s:5:"Hello";}`
	if !assert.ObjectsAreEqual(expectedA, string(b)) &&
		!assert.ObjectsAreEqual(expectedB, string(b)) {
		Suite.Fail(fmt.Sprintf("Not equal: \n"+
			"expected: %s OR %s\n"+
			"actual  : %s", expectedA, expectedB, string(b)))
	}
}

func (Suite *EncodeSuite) TestMarshalStringKeyMap() {
	m := make(map[string]string)
	m[`a`] = `Hello`
	m[`b`] = `World`

	b, err := Marshal(m)
	Suite.Nil(err)

	expectedA := `a:2:{s:1:"a";s:5:"Hello";s:1:"b";s:5:"World";}`
	expectedB := `a:2:{s:1:"b";s:5:"World";s:1:"a";s:5:"Hello";}`
	if !assert.ObjectsAreEqual(expectedA, string(b)) &&
		!assert.ObjectsAreEqual(expectedB, string(b)) {
		Suite.Fail(fmt.Sprintf("Not equal: \n"+
			"expected: %s OR %s\n"+
			"actual  : %s", expectedA, expectedB, string(b)))
	}
}

func (Suite *EncodeSuite) assertMarshal(v interface{}, expected string) {
	b, err := Marshal(v)
	assert.Nil(Suite.T(), err)
	assert.Equal(Suite.T(), expected, string(b))
}

func (Suite *EncodeSuite) assertMarshalContained(v interface{}, expected string) {
	type container struct {
		Value interface{} `php:"value"`
	}

	b, err := Marshal(container{Value: v})
	assert.Nil(Suite.T(), err)
	assert.Equal(Suite.T(), `a:1:{s:5:"value";`+expected+`}`, string(b))
}

func TestEncodeSuite(t *testing.T) {
	suite.Run(t, new(EncodeSuite))
}
