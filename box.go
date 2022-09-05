package main

import (
	"encoding/json"
	"fmt"
)

// `Box` can be used to make arbitrary
// interfaces serialisable.
type Box[T Boxable] struct {
	Data    T
	rawJSON json.RawMessage
}

type Boxable interface {
	Unbox(json.RawMessage) error
}

func NewBox[T Boxable](data T) (Box[T], error) {
	marshalledData, err := json.Marshal(data)
	if err != nil {
		return Box[T]{}, nil
	}
	return Box[T]{
		Data:    data,
		rawJSON: marshalledData,
	}, nil
}

func (b *Box[T]) UnmarshalJSON(data []byte) error {
	var v struct {
		Data json.RawMessage
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	b.rawJSON = v.Data
	return nil
}

func (b *Box[T]) Unbox(concreteTypes ...T) (err error) {
	for _, typ := range concreteTypes {
		err = typ.Unbox(b.rawJSON)
		if err == nil {
			b.Data = typ
			return nil
		}
	}
	return fmt.Errorf("no concrete types matched value inside box: %w", err)
}
