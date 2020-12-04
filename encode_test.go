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
}

func (Suite *EncodeSuite) TestMarshalBool() {
	Suite.assertMarshal(true, `b:1;`)
	Suite.assertMarshal(false, `b:0;`)
}

func (Suite *EncodeSuite) TestMarshalSlices() {
	Suite.assertMarshal([]int{1, 3, 5}, `a:3:{i:0;i:1;i:1;i:3;i:2;i:5;}`)
	Suite.assertMarshal([]string{`one`, `three`, `five`}, `a:3:{i:0;s:3:"one";i:1;s:5:"three";i:2;s:4:"five";}`)
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

func TestEncodeSuite(t *testing.T) {
	suite.Run(t, new(EncodeSuite))
}
