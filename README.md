介绍
====

这是一个单元测试工具库，用来降低Go语言项目单元测试的代码复杂度。

工具接口
=======

以下先做一个简单对比，使用库之前的单元测试如下：

```go
func VerifyBuffer(t *testing.T, buffer InBuffer) {
	if buffer.ReadUint8() != 1 {
		t.Fatal("buffer.ReadUint8() != 1")
	}

	if buffer.ReadByte() != 99 {
		t.Fatal("buffer.ReadByte() != 99")
	}

	if buffer.ReadInt8() != -2 {
		t.Fatal("buffer.ReadInt8() != -2")
	}

	if buffer.ReadUint16() != 0xFFEE {
		t.Fatal("buffer.ReadUint16() != 0xFFEE")
	}

	if buffer.ReadInt16() != 0x7FEE {
		t.Fatal("buffer.ReadInt16() != 0x7FEE")
	}

	if buffer.ReadUint32() != 0xFFEEDDCC {
		t.Fatal("buffer.ReadUint32() != 0xFFEEDDCC")
	}

	if buffer.ReadInt32() != 0x7FEEDDCC {
		t.Fatal("buffer.ReadInt32() != 0x7FEEDDCC")
	}

	if buffer.ReadUint64() != 0xFFEEDDCCBBAA9988 {
		t.Fatal("buffer.ReadUint64() != 0xFFEEDDCCBBAA9988")
	}

	if buffer.ReadInt64() != 0x7FEEDDCCBBAA9988 {
		t.Fatal("buffer.ReadInt64() != 0x7FEEDDCCBBAA9988")
	}

	if buffer.ReadRune() != '好' {
		t.Fatal(`buffer.ReadRune() != '好'`)
	}

	if buffer.ReadString(6) != "Hello1" {
		t.Fatal(`buffer.ReadString() != "Hello"`)
	}

	if bytes.Equal(buffer.ReadBytes(6), []byte("Hello2")) != true {
		t.Fatal(`bytes.Equal(buffer.ReadBytes(5), []byte("Hello")) != true`)
	}

	if bytes.Equal(buffer.ReadSlice(6), []byte("Hello3")) != true {
		t.Fatal(`bytes.Equal(buffer.ReadSlice(5), []byte("Hello")) != true`)
	}
}
```

使用库重构后的单元测试如下：

```go
func VerifyBuffer(t *testing.T, buffer InBuffer) {
	.Equal(t, buffer.ReadByte(), 99)
	.Equal(t, buffer.ReadInt8(), -2)
	.Equal(t, buffer.ReadUint8(), 1)
	.Equal(t, buffer.ReadInt16(), 0x7FEE)
	.Equal(t, buffer.ReadUint16(), 0xFFEE)
	.Equal(t, buffer.ReadInt32(), 0x7FEEDDCC)
	.Equal(t, buffer.ReadUint32(), 0xFFEEDDCC)
	.Equal(t, buffer.ReadInt64(), 0x7FEEDDCCBBAA9988)
	.Equal(t, buffer.ReadUint64(), 0xFFEEDDCCBBAA9988)
	.Equal(t, buffer.ReadRune(), '好')
	.Equal(t, buffer.ReadString(6), "Hello1")
	.Equal(t, buffer.ReadBytes(6), []byte("Hello2"))
	.Equal(t, buffer.ReadSlice(6), []byte("Hello3"))
}
```

在不牺牲单元测试结果输出的清晰性的前提下，可以减少很多不必要的判断语句和错误信息。

同时还会在测试失败是输出必要数据方便调试，例如：

```
$ go test -v
=== RUN   Test_All
    .go:35: check fail
        args[0] = 1
        args[1] = 2
        args[2] = 3
    _test.go:10: not nil
        v = "io: read/write on closed pipe"
    _test.go:11: not equal
        a = 1
        b = 2
    _test.go:12: not equal
        a = 1.233
        b = 3.333
    _test.go:13: not equal
        a = '你' = 20320
        b = '好' = 22909
    _test.go:14: not equal
        a = "sadkfjsl"
        b = "sdfs*\r\n"
    _test.go:15: not equal
        a = []byte{0x1, 0x2, 0x3, 0x3}
        b = []byte{0x3, 0x4, 0x5, 0x6}
    _test.go:17: not equal
        a = []int{1, 2, 3}
        b = []int{3, 4, 5}
    _test.go:18: not deep equal
        a = []int{1, 2, 3}
        b = []int{3, 4, 5}
--- FAIL: Test_All (0.00s)
FAIL
exit status 1
FAIL	vendor/github.com/funny/	0.006s
```

如果单纯靠testing包做单元测试，就需要加判断和打印才能做到。

支持以下几种测试检查：

```go
// 自定义条件的断言，断言失败时测试立即终止
// 支持变长参数的数据打印，测试失败时这些数据将被打印出来
.Assert(t, condition, args...)

// 同Assert，区别是Check失败时测试不会立即终止执行
.Check(t, condition, args...)

// 断言v必须为nil，否则测试失败
.IsNil(t, v)
// 同IsNil，失败时测试立即终止
.IsNilNow(t, v)

// 断言v不能为nil，否则测试失败
.NotNil(t, v)
.NotNilNow(t, v)

//
// 断言a和b必须相等，此函数只支持以下数据类型：
//   int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, rune, byte, string, []byte
//   []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []rune
// 或者实现了.Equals接口的自定义类型。
// 当为简单类型时，a和b必须是同类型，只有当b为int类型时，函数内部才会尝试做类型转换，这个设计是为了适应常量值的用法：
//   .Equal(t, GetUint64(), 123)
//
.Equal(t, a, b)
.EqualNow(t, a, b)

// 当Equal无法满足测试需要时可以使用次函数，但是请记得次函数开销较大
.DeepEqual(t, a, b)
.DeepEqualNow(t, a, b)
```

进程监控
=======

此外，在进行一些复杂的多线程单元测试的时候，可能出现死锁的情况，或者进行benchmark的时候，需要知道过程中内存的增长情况和GC情况。

为这些情况提供了一个统一的监控功能，在单元测试运行目录下使用以下方法可以获取到单元测试过程中的信息：

```shell
echo 'lookup goroutine' > .cmd
```

以上shell脚本将使自动输出goroutine堆栈跟踪信息到`.goroutine`文件。

支持以下几种监控命令：

```
lookup goroutine  -  获取当前所有goroutine的堆栈跟踪信息，输出到.goroutine文件，用于排查死锁等情况
lookup heap       -  获取当前内存状态信息，输出到.heap文件，包含内存分配情况和GC暂停时间等
lookup threadcreate - 获取当前线程创建信息，输出到.thread文件，通常用来排查CGO的线程使用情况
```

此外你还可以通过注册`.CommandHandler`回调来添加自己的监控命令支持。
