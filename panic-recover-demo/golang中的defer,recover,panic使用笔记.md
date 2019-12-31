### golang中的defer,recover,panic使用笔记

#### 1,在没有任何recover的情况下，panic是啥情况

我们先写一个简单的demo来演示一下当panic出现的时候的样子

```go
package main

import (
	"fmt"
)

func main() {
	defer fmt.Println("in main")
	demo0()
	fmt.Println("process main")
}

func demo0() {
	panic("demo0 unknown err")
    fmt.Println("demo0 process")
}

```

运行以后:

in main

panic: demo0 unknown err

goroutine 1 [running]:

main.demo0()

​	F:/Golang/src/development/panic-recover-demo/panic_reconver.go:16 +0x40

main.main()

​	F:/Golang/src/development/panic-recover-demo/panic_reconver.go:10 +0x8f
exit status 2

[Finished in 1.3s with exit code 1]

可以看到如果在panic发生的时候，没有recover，那么process main的输出因为被中断了，同时panic发生以后，执行所有的defer，所以就是上面这个打印输出的信息了。



####2，panic，遇到defer recover的时候,这种是大多数人的写法

照样我们先写一个简单的demo程序来演示一下:

```go
package main

import (
	"fmt"
)

func main() {
	defer fmt.Println("in main")

    demo1()
    fmt.Println("process main")
}

func demo1() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

    panic("demo1 unknown err")
    
    fmt.Println("demo1 process")
}
```

运行以后：

demo1 unknown err

process main

in main

[Finished in 1.4s]

从运行情况看，panic终止了demo1()中的demo1 process输出，但同时当demo1中的panic被demo1()中的defer recover()以后，到了main的流程，他是走的正常逻辑，所以输出了prcocess main。

所以这种大多数人的写法，针对一些可以能被recover的panic，他是可以被recover的。这种写法也是很多教材上推荐的写法。

####3，懒人的做法，将recover函数封装起来进行处理

有时候，为了简化代码，少敲几行代码，我们将一些常用代码进行封装，然后再进行调用。

同样我们也写一个简单的demo2的代码来运行一次，看下结果:

```go
package main

import (
	"fmt"
)

func main() {
	defer fmt.Println("in main")

    demo2()
    fmt.Println("process main")
}
//将recover进行封装起来
func handpanic() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}

func demo2() {
	defer handpanic()//注意这里需要去defer处理

	panic("demo2 unknown err")

	fmt.Println("demo2 process")
}
```

运行结果:

demo2 unknown err

process main

in main

[Finished in 1.7s]

可以看到这种情况和上面demo1()例子中那种不进行封装的处理结果是一样的，他也是可以被recover的。



#### 4，犯错的recover封装做法

我们在demo2()例子的基础上，对代码稍加修改，形成如下的代码，你看这种用法是否正确?

```go
package main

import (
	"fmt"
)

func main() {
	defer fmt.Println("in main")

    demo3()
    fmt.Println("process main")
}

func handpanic() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}

func handle() {
	//直接调用,无法被recover
	handpanic()
	//在通过defer 调用，也无法recover
	// defer handpanic()
}

func demo3() {
	defer handle()

	panic("demo3 unknown err")

	fmt.Println("demo3 process")
}
```

运行结果:

in main

panic: demo3 unknown err

goroutine 1 [running]:

main.demo3()

​	F:/Golang/src/development/panic-recover-demo/panic_reconver.go:57 +0x5c

main.main()

​	F:/Golang/src/development/panic-recover-demo/panic_reconver.go:10 +0x8f

exit status 2

[Finished in 1.3s]

从运行的结果看，这种代码是错误的，因为他不能将能被recover的错误或者panic进行recover。

原因是什么？

**因为recover 只有在 defer 函数中才有用，在 defer 的函数调用的函数中 recover 不起作用**

如果你的用法也犯了类似的这种错误，请立即修改你的代码。

这里给大家推荐一个文章，https://mp.weixin.qq.com/s/aMKhU9rG_Al-sA5DAFji_g ，在之前刚刚入手golang的时候，我也犯了此种错误，后面因为看到他的这个文章，让我受益颇多，所以写下了这个笔记，避免自己以后再犯错。