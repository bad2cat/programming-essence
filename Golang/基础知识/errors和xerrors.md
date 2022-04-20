### errors 和 xerrors

#### errors.Unwrap(err error) error

```
func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}
```

这个方法就是查看当前的这个`error`是否实现了`unWrap() error`方法，如果实现了，则返回该 `err.Unwrap()`，否则返回`nil`，下面是一个代码示例：

```
type WrapErr struct{
    msg string
    err error
}

func (w *WrapErr) Unwrap() error{
    return fmt.Errorf("err:%s",w.msg)
}

func main(){
    err := errors.New("This is a error")
    we := &WrapErr{
    	msg:"This is err",
    	err:err}
    if newErr := we.Unwrap();newErr != nil{
        fmt.Println(newErr.String())
    }
}

//output:
err: This is a error
```

#### errors.Is(err,target error) bool

是否在 `err`链中存在与 `target` 匹配的`error`，存在则返回`true`，否则返回`false`

这个`err`链包含`err`本身和后续连续调用`Unwrap`方法获取得到的`error`组成

`err`和`target`匹配意味着`err`和`target`相等，或者是它实现了一个方法`Is(error)`，并且有`Is(target)`为`true`

#### errors.As(err,target interface{}) bool

找到在这个`err`链中第一个与`target`匹配的`error`，并且将`error`的值赋给`target`，然后返回`true`

这个`err`链包含`err`本身和后续连续调用`Unwrap`方法获取得到的`error`组成

`err`和`target`匹配是值，`err`的值是指针`target`类型的值，或者是有一个方法`As(interfalce)bool`使得`As(target)`返回为 `true`

#### errors 和 xerrors

我感觉 xerrors 就是在 errors 的基础上添加了栈追踪的功能