#### Golang 数据类型

###### uint：至少 32 位大小，它是一个独立的类型，不是 uint32 的别名

###### uintptr：它是一个整数类型，没有指定大小，只是说 large enough，可以表示任何类型的指针；也就是说，在 32 位的 os 中，uintptr 就是32位，在 64 的 os 中，uintptr 就是 64位大小

###### byte：它是 uint8 的别称，一般是用来区分字节值和 8 位无符号整数值的

###### rune：它是 uint32 的别名，一般用来区分字符值和整数值的

