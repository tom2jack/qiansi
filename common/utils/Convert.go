package utils

import "reflect"

// 使用反射，转换结构体
func SuperConvert(sourceStruct interface{}, targetStruct interface{}) {
	source := structToMap(sourceStruct)
	targetV := reflect.ValueOf(targetStruct).Elem()
	targetT := reflect.TypeOf(targetStruct).Elem()
	for i := 0; i < targetV.NumField(); i++ {
		fieldName := targetT.Field(i).Name
		sourceVal := source[fieldName]
		if !sourceVal.IsValid() {
			continue
		}
		targetVal := targetV.Field(i)
		targetVal.Set(sourceVal)
	}
}

func structToMap(structName interface{}) map[string]reflect.Value {
	t := reflect.TypeOf(structName).Elem()
	v := reflect.ValueOf(structName).Elem()
	fieldNum := t.NumField()
	resMap := make(map[string]reflect.Value, fieldNum)
	for i := 0; i < fieldNum; i++ {
		resMap[t.Field(i).Name] = v.Field(i)
	}
	return resMap
}
