# 由谷歌官方json直接复制而来
1.增加结构体首字母默认小写功能：

    StructFirstLower=true表示首字母小写（tag的优先级最高，有tag则根据tag）

2.增加map、list、slide转结构体的功能：

    Map2Struct（map[string]interface{}类型转结构体）
    List2Struct（[]map[string]interface{}转结构体）
    Obj2Struct（上面两个的拓展，支持map[string]任意类型）
    TimeFormatType（Time类型格式化为整数时是按秒毫秒还是纳秒来计算，默认毫秒）
3.空数据默认不返回：

    移除原始的tag"omitempty"，新增"keepEmpty"，想返回空数据加上此tag即可
    
    不返回的数据有（bool一定会返回）：
            Array、Map、Slice、String的长度是0
            Int、Uint、Float是0
            Interface、Ptr是nil

## 导入方式
导入：go get github.com/weimingjue/json

更新（已经导入过)：go get -u github.com/weimingjue/json
