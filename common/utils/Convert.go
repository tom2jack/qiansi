package utils

import "reflect"

func SuperConvert(A interface{}, B interface{}) {
	At := reflect.TypeOf(A)
	Av := reflect.ValueOf(A)
	Bt := reflect.TypeOf(B)
	Bv := reflect.ValueOf(B)
	for i := 0; i < Bt.NumField(); i++ {
		println(Bt.Field(i).Name)
		for m := 0; m < At.NumField(); m++ {
			println(At.Field(m).Name)
			if Bt.Field(i).Name == At.Field(m).Name {
				Bv.Field(m).Set(Av.Field(i))
				Bv.Field(m).Interface()
			}
		}
	}
}
