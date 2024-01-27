package mongo

import (
	"fmt"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

var tDecimal = reflect.TypeOf(decimal.Decimal{})

// Decimal primitive.Decimal128 -> decimal.Decimal
type Decimal decimal.Decimal

func (d Decimal) DecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	decimalType := reflect.TypeOf(decimal.Decimal{})
	if !val.IsValid() || !val.CanSet() || val.Type() != decimalType {
		return bsoncodec.ValueDecoderError{
			Name:     "decimalDecodeValue",
			Types:    []reflect.Type{decimalType},
			Received: val,
		}
	}

	decimal.DivisionPrecision = 128
	var value decimal.Decimal
	switch vr.Type() {
	case bson.TypeDecimal128:
		dec, err := vr.ReadDecimal128()
		if err != nil {
			return err
		}
		if value, err = decimal.NewFromString(dec.String()); err != nil {
			return err
		}
	default:
		return fmt.Errorf("received invalid BSON type to decode into decimal.Decimal: %s", vr.Type())
	}

	val.Set(reflect.ValueOf(value))
	return nil
}

func (d Decimal) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	decimalType := reflect.TypeOf(decimal.Decimal{})
	if !val.IsValid() || val.Type() != decimalType {
		return bsoncodec.ValueEncoderError{
			Name:     "decimalEncodeValue",
			Types:    []reflect.Type{decimalType},
			Received: val,
		}
	}

	decimal.DivisionPrecision = 128
	dec := val.Interface().(decimal.Decimal)
	dec128, err := primitive.ParseDecimal128(dec.String())
	if err != nil {
		return err
	}

	return vw.WriteDecimal128(dec128)
}
