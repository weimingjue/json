# 由谷歌官方json直接复制而来
1.增加结构体首字母默认小写功能

    StructFirstLower=true表示首字母小写（tag的优先级最高，有tag则根据tag）
    示例：
    type Test struct {
	    Data int
	    Data2 int `json:"D2"`
    }
    赋值并解析成json：{"data":10,"D2":20}
2.增加map、list、slide转结构体的功能

    Map2Struct（map[string]interface{}类型转结构体）
    List2Struct（[]map[string]interface{}转结构体）
    Obj2Struct（上面两个的拓展，支持map[string]任意类型）
    TimeFormatType（Time类型格式化为整数时是按秒毫秒还是纳秒来计算，默认毫秒）
3.空数据默认不返回

    移除原始的tag"omitempty"，新增"keepEmpty"，想返回空数据加上此tag即可
    
    示例，第一个data是必填，第二个keepEmpty表示空也返回：
    type Test struct {
	    Data interface{} `json:"data,keepEmpty"` //空也返回
    }
    
    全局增加空数据返回：StructKeepType
    示例，想全局让boolean、string、整形这三个空都返回
    type TestKeepStruct struct {
	    Bo  bool
	    It  uint
	    St  string
	    Ar  []string
    }
	json.StructKeepType = json.KeepEmptyBool | json.KeepEmptyNumber | json.KeepEmptyString
	结构体数据均为空，调用以上代码后返回的json：{"bo":false,"it":0,"st":""}
	
4.将Time转json时默认格式为YY-MM-DD HH:mm:ss
    
    修改格式：BaseTimeFormat
    修改时区将time.Local改为对应的时区即可，如修改为0时区：time.Local = time.UTC

	
## 导入方式
导入或更新：go get -u github.com/weimingjue/json
