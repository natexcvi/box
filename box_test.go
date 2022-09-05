package box

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	A string
	B []int
}

func (s *TestStruct) Method() {

}

func (s *TestStruct) Unbox(data json.RawMessage) error {
	return json.Unmarshal(data, s)
}

type TestInterface interface {
	Method()
	Unbox(json.RawMessage) error
}

func TestBox(t *testing.T) {
	testCases := []struct {
		name         string
		data         TestInterface
		concreteType TestInterface
	}{
		{
			name: "struct",
			data: &TestStruct{
				"test",
				[]int{1, 2, 3},
			},
			concreteType: &TestStruct{},
		},
	}

	for _, tc := range testCases {
		box, err := NewBox(tc.data)
		assert.NoError(t, err, "box creation")
		payload, err := json.Marshal(box)
		assert.NoError(t, err, "JSON marshal box")
		unmarshalledBox, err := NewBox(tc.data)
		assert.NoError(t, err, "box creation")
		unmarshalledBox.Data = nil
		err = json.Unmarshal(payload, &unmarshalledBox)
		assert.NoError(t, err, "JSON unmarshal box")
		err = unmarshalledBox.Unbox(tc.concreteType)
		assert.NoError(t, err, "unboxing")
		assert.Equal(t, box.Data, unmarshalledBox.Data)
	}
}
