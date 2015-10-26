介绍
====

这是一个单元测试工具库，用来降低Go语言项目单元测试的代码复杂度。

工具接口
=======

以下先做一个简单对比，使用unitest库之前的单元测试如下：

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

使用unitest库重构后的单元测试如下：

```go
func VerifyBuffer(t *testing.T, buffer InBuffer) {
	unitest.Pass(t, buffer.ReadByte() == 99)
	unitest.Pass(t, buffer.ReadInt8() == -2)
	unitest.Pass(t, buffer.ReadUint8() == 1)
	unitest.Pass(t, buffer.ReadInt16() == 0x7FEE)
	unitest.Pass(t, buffer.ReadUint16() == 0xFFEE)
	unitest.Pass(t, buffer.ReadInt32() == 0x7FEEDDCC)
	unitest.Pass(t, buffer.ReadUint32() == 0xFFEEDDCC)
	unitest.Pass(t, buffer.ReadInt64() == 0x7FEEDDCCBBAA9988)
	unitest.Pass(t, buffer.ReadUint64() == 0xFFEEDDCCBBAA9988)
	unitest.Pass(t, buffer.ReadRune() == '好')
	unitest.Pass(t, buffer.ReadString(6) == "Hello1")
	unitest.Pass(t, bytes.Equal(buffer.ReadBytes(6), []byte("Hello2")))
	unitest.Pass(t, bytes.Equal(buffer.ReadSlice(6), []byte("Hello3")))
}
```

当以上单元测试判断失败时，unitest将自动从代码中提取对应行号的代码作为错误信息。

在不牺牲单元测试结果输出的清晰性的前提下，unitest可以减少很多不必要的判断语句和错误信息。

进程监控
=======

此外，在进行一些复杂的多线程单元测试的时候，可能出现死锁的情况，或者进行benchmark的时候，需要知道过程中内存的增长情况和GC情况。

unitest为这些情况提供了一个统一的监控功能，在单元测试运行目录下使用以下方法可以获取到单元测试过程中的信息：

```shell
echo 'lookup goroutine' > unitest.cmd
```

以上shell脚本将使unitest自动输出goroutine堆栈跟踪信息到`unitest.goroutine`文件。

unitest支持以下几种监控命令：

```
lookup goroutine  -  获取当前所有goroutine的堆栈跟踪信息，输出到unitest.goroutine文件，用于排查死锁等情况
lookup heap       -  获取当前内存状态信息，输出到unitest.heap文件，包含内存分配情况和GC暂停时间等
lookup threadcreate - 获取当前线程创建信息，输出到unitest.thread文件，通常用来排查CGO的线程使用情况
```

此外你还可以通过注册`unitest.CommandHandler`回调来添加自己的监控命令支持。
