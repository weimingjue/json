package json

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//格式化int时间时是秒、毫秒还是纳秒
var TimeFormatType = time.Millisecond

//map转bean
//map里只能是基本类型、map、list、nil
//structs必须是结构体的指针
func Map2Struct(maps map[string]interface{}, structs interface{}) error {
	sValue := reflect.ValueOf(structs)
	if sValue.Kind() != reflect.Ptr {
		return errors.New("structs必须是引用类型")
	}
	return Obj2Struct(maps, sValue)
}

//list转bean
//slices必须是切片的指针
func List2Struct(list []map[string]interface{}, slices interface{}) error {
	sValue := reflect.ValueOf(slices)
	if sValue.Kind() != reflect.Ptr {
		return errors.New("slices必须是引用类型")
	}
	sValue = sValue.Elem() //取实际地址
	if sValue.Kind() != reflect.Slice {
		return errors.New("slices必须是切片")
	}
	return Obj2Struct(list, sValue)
}

//有些map可能不是Interface类型的
//obj：map或Array或Slice
//data：与之对应的Struct或Slice（不支持Array）
func Obj2Struct(obj interface{}, data reflect.Value) error {
	if obj == nil {
		return nil
	}
	if data.Kind() == reflect.Ptr { //指针取实际地址
		if data.IsNil() {
			if !data.CanSet() { //未导出的字段
				return nil
			}
			data.Set(reflect.New(data.Type().Elem()))
		}
		data = data.Elem()
	}
	objValue := reflect.ValueOf(obj)

	//根据obj类型判断
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
		obj = objValue.Interface()
	}
	objKind := objValue.Kind()
	switch data.Kind() {
	case reflect.Array:
		////////////////////////////////////////////////////////////////////////////////////////////////////////
		//结构体里面是Array将不做解析
		////////////////////////////////////////////////////////////////////////////////////////////////////////
		return errors.New("结构体里使用了Array，由于Array不可变，故无法解析")
	case reflect.Slice:
		////////////////////////////////////////////////////////////////////////////////////////////////////////
		//解析切片
		////////////////////////////////////////////////////////////////////////////////////////////////////////
		switch objKind {
		case reflect.Array, reflect.Slice:
			if objKind == reflect.Slice && objValue.IsNil() {
				return nil
			}
			length := objValue.Len();
			for i := 0; i < length; i++ {
				if i == 0 {
					data.Set(reflect.MakeSlice(data.Type(), length, length))
				}
				//循环调用
				if e := Obj2Struct(objValue.Index(i).Interface(), data.Index(i)); e != nil {
					return e
				}
			}
		default:
			return errors.New("类型不匹配：结构体类型" + data.Kind().String() + "，数据类型" + objKind.String())
		}
	case reflect.Struct:
		////////////////////////////////////////////////////////////////////////////////////////////////////////
		//解析结构体
		////////////////////////////////////////////////////////////////////////////////////////////////////////
		//时间类型
		if _, ok := obj.(time.Time); ok {
			data.Set(objValue)
			return nil
		}

		if objKind != reflect.Map {
			return errors.New("类型不匹配：结构体类型" + data.Kind().String() + "，数据类型" + objKind.String())
		}

		fields := cachedTypeFields(data.Type())

		for iter := objValue.MapRange(); iter.Next(); {
			key := []byte(iter.Key().Interface().(string))
			//匹配到的缓存数据
			var f *field
			for i := range fields {
				ff := &fields[i]
				if bytes.Equal(ff.nameBytes, key) {
					f = ff
					break
				}
				if f == nil && ff.equalFold(ff.nameBytes, key) {
					f = ff
				}
			}
			if f != nil {
				var subv reflect.Value
				//获取对应的反射
				for _, i := range f.index {
					subv = data.Field(i)
				}
				err := Obj2Struct(iter.Value().Interface(), subv) //循环设置数据
				if err != nil {
					return err
				}
			}
		}
		////////////////////////////////////////////////////////////////////////////////////////////////////////
		//解析基本类型
		////////////////////////////////////////////////////////////////////////////////////////////////////////
	case reflect.Bool:
		var dataSet bool
		switch objKind {
		case reflect.Bool:
			dataSet = obj.(bool)
		case reflect.String:
			st := obj.(string)
			dataSet = strings.Contains(st, "true") || strings.Contains(st, "1")
		case reflect.Int:
			dataSet = obj.(int) == 1
		case reflect.Int8:
			dataSet = obj.(int8) == 1
		case reflect.Int16:
			dataSet = obj.(int16) == 1
		case reflect.Int32:
			dataSet = obj.(int32) == 1
		case reflect.Int64:
			dataSet = obj.(int64) == 1
		case reflect.Uint:
			dataSet = obj.(uint) == 1
		case reflect.Uint8:
			dataSet = obj.(uint8) == 1
		case reflect.Uint16:
			dataSet = obj.(uint16) == 1
		case reflect.Uint32:
			dataSet = obj.(uint32) == 1
		case reflect.Uint64:
			dataSet = obj.(uint64) == 1
		case reflect.Uintptr:
			dataSet = obj.(uintptr) == 1
		case reflect.Float32:
			dataSet = obj.(float32) == 1
		case reflect.Float64:
			dataSet = obj.(float64) == 1
		default:
			return errors.New("类型不匹配：结构体类型" + data.Kind().String() + "，数据类型" + objKind.String())
		}
		data.SetBool(dataSet)
	case reflect.String:
		var dataSet string
		switch objKind {
		case reflect.Bool:
			if obj.(bool) {
				dataSet = "true"
			} else {
				dataSet = "false"
			}
		case reflect.String:
			dataSet = obj.(string)
		case reflect.Int:
			dataSet = strconv.FormatInt(int64(obj.(int)), 10)
		case reflect.Int8:
			dataSet = strconv.FormatInt(int64(obj.(int8)), 10)
		case reflect.Int16:
			dataSet = strconv.FormatInt(int64(obj.(int16)), 10)
		case reflect.Int32:
			dataSet = strconv.FormatInt(int64(obj.(int32)), 10)
		case reflect.Int64:
			dataSet = strconv.FormatInt(obj.(int64), 10)
		case reflect.Uint:
			dataSet = strconv.FormatUint(uint64(obj.(uint)), 10)
		case reflect.Uint8:
			dataSet = strconv.FormatUint(uint64(obj.(uint8)), 10)
		case reflect.Uint16:
			dataSet = strconv.FormatUint(uint64(obj.(uint16)), 10)
		case reflect.Uint32:
			dataSet = strconv.FormatUint(uint64(obj.(uint32)), 10)
		case reflect.Uint64:
			dataSet = strconv.FormatUint(uint64(obj.(uint64)), 10)
		case reflect.Uintptr:
			dataSet = strconv.FormatUint(uint64(obj.(uintptr)), 10)
		case reflect.Float32:
			dataSet = strconv.FormatFloat(float64(obj.(float32)), 'f', -1, 32)
		case reflect.Float64:
			dataSet = strconv.FormatFloat(obj.(float64), 'f', -1, 64)
		default:
			//时间类型
			if times, ok := obj.(time.Time); ok {
				dataSet = times.In(time.Local).Format("2006-01-02 15:04:05")
			} else {
				return errors.New("类型不匹配：结构体类型" + data.Kind().String() + "，数据类型" + objKind.String())
			}
		}
		data.SetString(dataSet)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var dataSet int64
		switch objKind {
		case reflect.Bool:
			if obj.(bool) {
				dataSet = 1
			} else {
				dataSet = 0
			}
		case reflect.String:
			st := obj.(string)
			switch st {
			case "true":
				dataSet = 1
			case "false":
				dataSet = 0
			default:
				it, err := strconv.ParseInt(st, 10, 64)
				if err != nil {
					return errors.New("无法将此字符串转成int：" + st)
				}
				dataSet = it
			}
		case reflect.Int:
			dataSet = int64(obj.(int))
		case reflect.Int8:
			dataSet = int64(obj.(int8))
		case reflect.Int16:
			dataSet = int64(obj.(int16))
		case reflect.Int32:
			dataSet = int64(obj.(int32))
		case reflect.Int64:
			dataSet = obj.(int64)
		case reflect.Uint:
			dataSet = int64(obj.(uint))
		case reflect.Uint8:
			dataSet = int64(obj.(uint8))
		case reflect.Uint16:
			dataSet = int64(obj.(uint16))
		case reflect.Uint32:
			dataSet = int64(obj.(uint32))
		case reflect.Uint64:
			dataSet = int64(obj.(uint64))
		case reflect.Uintptr:
			dataSet = int64(obj.(uintptr))
		case reflect.Float32:
			dataSet = int64(obj.(float32))
		case reflect.Float64:
			dataSet = int64(obj.(float64))
		default:
			//时间类型
			if times, ok := obj.(time.Time); ok {
				dataSet = times.UnixNano() / int64(TimeFormatType)
			} else {
				return errors.New("类型不匹配：结构体类型" + data.Kind().String() + "，数据类型" + objKind.String())
			}
		}
		data.SetInt(dataSet)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var dataSet uint64
		switch objKind {
		case reflect.Bool:
			if obj.(bool) {
				dataSet = 1
			} else {
				dataSet = 0
			}
		case reflect.String:
			st := obj.(string)
			switch st {
			case "true":
				dataSet = 1
			case "false":
				dataSet = 0
			default:
				it, err := strconv.ParseUint(st, 10, 64)
				if err != nil {
					return errors.New("无法将此字符串转成uint：" + st)
				}
				dataSet = it
			}
		case reflect.Int:
			dataSet = uint64(obj.(int))
		case reflect.Int8:
			dataSet = uint64(obj.(int8))
		case reflect.Int16:
			dataSet = uint64(obj.(int16))
		case reflect.Int32:
			dataSet = uint64(obj.(int32))
		case reflect.Int64:
			dataSet = uint64(obj.(int64))
		case reflect.Uint:
			dataSet = uint64(obj.(uint))
		case reflect.Uint8:
			dataSet = uint64(obj.(uint8))
		case reflect.Uint16:
			dataSet = uint64(obj.(uint16))
		case reflect.Uint32:
			dataSet = uint64(obj.(uint32))
		case reflect.Uint64:
			dataSet = obj.(uint64)
		case reflect.Uintptr:
			dataSet = uint64(obj.(uintptr))
		case reflect.Float32:
			dataSet = uint64(obj.(float32))
		case reflect.Float64:
			dataSet = uint64(obj.(float64))
		default:
			//时间类型
			if times, ok := obj.(time.Time); ok {
				dataSet = uint64(times.UnixNano()) / uint64(TimeFormatType)
			} else {
				return errors.New("类型不匹配：结构体类型" + data.Kind().String() + "，数据类型" + objKind.String())
			}
		}
		data.SetUint(dataSet)
	case reflect.Float32, reflect.Float64:
		var dataSet float64
		switch objKind {
		case reflect.Bool:
			if obj.(bool) {
				dataSet = 1
			} else {
				dataSet = 0
			}
		case reflect.String:
			st := obj.(string)
			switch st {
			case "true":
				dataSet = 1
			case "false":
				dataSet = 0
			default:
				it, err := strconv.ParseFloat(st, 64)
				if err != nil {
					return errors.New("无法将此字符串转成float：" + st)
				}
				dataSet = it
			}
		case reflect.Int:
			dataSet = float64(obj.(int))
		case reflect.Int8:
			dataSet = float64(obj.(int8))
		case reflect.Int16:
			dataSet = float64(obj.(int16))
		case reflect.Int32:
			dataSet = float64(obj.(int32))
		case reflect.Int64:
			dataSet = float64(obj.(int64))
		case reflect.Uint:
			dataSet = float64(obj.(uint))
		case reflect.Uint8:
			dataSet = float64(obj.(uint8))
		case reflect.Uint16:
			dataSet = float64(obj.(uint16))
		case reflect.Uint32:
			dataSet = float64(obj.(uint32))
		case reflect.Uint64:
			dataSet = float64(obj.(uint64))
		case reflect.Uintptr:
			dataSet = float64(obj.(uintptr))
		case reflect.Float32:
			dataSet = float64(obj.(float32))
		case reflect.Float64:
			dataSet = obj.(float64)
		default:
			//时间类型
			if times, ok := obj.(time.Time); ok {
				dataSet = float64(times.UnixNano()) / float64(TimeFormatType)
			} else {
				return errors.New("类型不匹配：结构体类型" + data.Kind().String() + "，数据类型" + objKind.String())
			}
		}
		data.SetFloat(dataSet)
	case reflect.Interface, reflect.Map: //这两个默认当做完全匹配
		data.Set(objValue)
	}
	return nil
}
