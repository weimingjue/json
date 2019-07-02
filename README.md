# 由谷歌官方json直接复制而来：
1.增加结构体首字母默认小写功能：StructFirstLower=true表示首字母小写（tag的优先级最高，有tag则根据tag）

2.增加map、list、slide转结构体的功能：Map2Struct（map[string]interface{}类型转结构体）、List2Struct（[]map[string]interface{}转结构体）、Obj2Struct（上面两个的拓展，支持map[string]任意类型）

## 导入方式：go get github.com/weimingjue/json
