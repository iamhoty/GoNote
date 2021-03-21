# Golang基础五

## 1. goroutine

### 1.1 进程和线程

1）进程就是程序在操作系统中的一次执行过程，是**系统进行资源分配和调度的基本单位**

2）线程是进程的一个执行实例，是**程序执行的最小单元**，它是比进程更小的能独立运行的基本单位，是**cpu调度的最小单位**

3）一个进程可以创建核销毁多个线程，同一个进程中的多个线程可以**并发执行**

4）一个程序至少有一个进程，一个进程至少有一个线程

### 1.2 并发和并行

1）多线程程序在**单核**上运行，就是并发

2）多线程程序在**多核**上运行，利用计算机的多核，就是并行

![go1.14.7.darwin-amd64.tar](./asset_5/go1.14.7.darwin-amd64.tar.png)

并发：

​	因为是在一个cpu上，比如有10个线程，每个线程执行10毫秒（进行轮询操作），从人的角度看，好像这10个线程都在运行，但是从微观上看，**在某一个时间点看，其实只有一个线程在执行**，这就是并发

并行：

​	因为是在多个cpu上（比如有10个cpu），比如有10个线程，每个线程执行10毫秒（各自在不同cpu上执行），从人的角度看，这10个线程都在运行，但是从微观上看，**在某一个时间点看，也同时有10个线程在执行**，这就是并行

### 1.3 Go协程和Go主线程

Go主线程（有程序员直接称为线程/也可以理解成进程）：

一个Go线程上，可以**起多个协程**，你可以这样理解，协程是轻量级的线程【编译器做优化】。

#### 1.3.1 Go协程的特点

1）有独立的栈空间

2）共享程序堆空间

3）调度由用户控制

4）协程是轻量级的线程

![Snipaste_2021-03-18_23-27-14](./asset_5/Snipaste_2021-03-18_23-27-14.png)

#### 1.3.2 案例

在主线程(可以理解成进程)中，开启一个goroutine, 该协程每隔1秒输出 "hello,world"
在主线程中也每隔一秒输出"hello,golang", 输出10次后，退出程序
要求主线程和goroutine同时执行

```go
func test() {
	for i := 1; i <= 10; i++ {
		fmt.Println("hello world" + strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}

func main() {
	go test() // 开启协程
	// 主函数和test一起执行 
	for i := 1; i <= 10; i++ {
		fmt.Println("main hello world" + strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}
// hello world8
// main hello world8
// hello world9
// main hello world9
// main hello world10
// hello world10
```

主线程和协程执行图

![Snipaste_2021-03-18_23-43-31](./asset_5/Snipaste_2021-03-18_23-43-31.png)

#### 1.3.3 案例小结

1）主线程是一个物理线程，直接作用在cpu上的。是重量级的，非常耗费cpu资源。

2）协程从主线程开启的，是轻量级的线程，是逻辑态，**在应用层上执行**。对资源消耗相对小。

3） Golang的协程机制是重要的特点，可以轻松的**开启上万个协程**。其它编程语言的并发机制是一般基于线程的，开启过多的线程，资源耗费大，这里就突显Golang在并发上的优势了

### 1.4 goroutine的调度模式

#### 1.4.1 MPG模式

（1）M：Machine 操作系统的主线程

（2）P：Processor 上下文环境（运行时所需要的资源、内存、状态等）

（3）G：Goroutine 协程

M主线程可以执行在一个cpu上，也可以执行在多个cpu上；当有一个G阻塞时，会来回切换其他的G协程去执行，充分利用cpu的资源

### 1.5 设置Golang运行的cpu数

runtime

![Snipaste_2021-03-19_00-14-41](./asset_5/Snipaste_2021-03-19_00-14-41.png)

```go
func main()  {
	cpuNum := runtime.NumCPU()
	fmt.Println("cpu num >>>", cpuNum) // 12
	// 可以自己设置使用的cpu核数
	runtime.GOMAXPROCS(cpuNum - 2)
	fmt.Println("ok")
}
```

## 2. channel

### 2.1 案例

需求：现在要计算 1-200 的各个数的阶乘，并且把各个数的阶乘放入到map中。最后显示出来。要求使用goroutine完成

```go
// 1. 编写一个函数，来计算各个数的阶乘，并放入到 map中.
// 2. 我们启动的协程多个，统计的将结果放入到 map中
// 3. map 应该做出一个全局的.
var (
	myMap = make(map[int]int, 10)
)

// 计算n的阶乘 并将结果放入到map中
func test1(n int) {
	res := 1
	for i := 1; i <= n; i++ {
		res *= i
	}
	// 将res放入myMap中
	myMap[n] = res
}

func main() {
	// 开启多个协程
	for i := 1; i <= 200; i++ {
		go test1(i)
	}
  time.Sleep(time.Second * 10)
	// 输出结果
	for i, v := range myMap {
		fmt.Printf("map[%v]=%v\n", i, v)
	}
}
```

**问题：**执行时会报错 fatal error: concurrent map writes 原因是多个test1协程，同一时刻操作map空间进行写入操作

![Snipaste_2021-03-19_01-28-12](./asset_5/Snipaste_2021-03-19_01-28-12.png)

### 2.2 sync

#### package sync

```
import "sync"
```

sync包提供了基本的同步基元，如互斥锁。除了Once和WaitGroup类型，**大部分都是适用于低水平程序线程，高水平的同步使用channel通信更好一些。**

![Snipaste_2021-03-19_01-40-41](./asset_5/Snipaste_2021-03-19_01-40-41.png)

**使用全局变量加锁同步优化**

```go
var (
	myMap = make(map[int]int, 10)
	// 声明一个全局的互斥锁
	// Mutex : 是互斥
	lock sync.Mutex
)

// 计算n的阶乘 并将结果放入到map中
func test1(n int) {
	res := 1
	for i := 1; i <= n; i++ {
		res *= i
	}
	// 将res放入myMap中
	lock.Lock() //  加锁
	myMap[n] = res
	lock.Unlock() // 解锁
}
```

### 2.3 不同goroutine协程通信方式

（1）全局变量的互斥锁 sync.Mutex.lock/unlock

（2）使用管道channel

### 2.4 为什么需要使用channel

1）前面使用全局变量加锁同步来解决goroutine的通讯，但不完美

2）主线程在等待所有goroutine全部完成的时间很难确定，我们这里设置10秒，仅仅是估算。

3）如果主线程休眠时间长了，会加长等待时间，如果等待时间短了，可能还有goroutine处于工作状态，这时也会随主线程的退出而销毁

4）通过全局变量加锁同步来实现通讯，也并不利用多个协程对全局变量的读写操作。

5）上面种种分析都在呼唤一个新的通讯机制-channel

### 2.5 管道channel基本介绍

1） channle本质就是一个数据结构-队列

![Snipaste_2021-03-19_23-33-59](./asset_5/Snipaste_2021-03-19_23-33-59.png)

2）管道的数据是**先进先出**【FIFO： first in first out】**栈：先进后出**

3）线程安全，多goroutine访问时，**无需加锁**，底层是用锁的机制维护的，就是说channel本身就是线程安全的

4） **channel是有类型的**，一个string的channel只能存放string类型数据。

### 2.6 定义/声明channel

var 管道名 chan 数据类型

```go
var intChan chan int
var mapChan chan map[int]string
var perChan Person
```

说明：

**channel是引用类型，需要make进行初始化**

```go
func main() {
	// 演示管道的使用
	var intChan chan int
	intChan = make(chan int, 3)
	fmt.Println("intChan>>>", intChan)
	// intChan>>> 0xc000102000
	fmt.Printf("intChan本身的地址%p\n", &intChan)
	// intChan本身的地址0xc0000ae018

	// 向管道写入数据
	intChan <- 10
	num := 20
	intChan <- num
	// 管道的长度和容量 管道不可以扩容 写入超过容量的数据库会报错
	fmt.Printf("len=%v cap=%v\n", len(intChan), cap(intChan))
	// len=2 cap=3

	// 从管道取出数据 如果管道没数据 取出时会报错deadlock
	n1 := <-intChan
	fmt.Println("n1 >>> ", n1) // 10
	fmt.Printf("len=%v cap=%v\n", len(intChan), cap(intChan))
	// len=1 cap=3
}
```

### 2.7 channel使用注意事项

1） channel中只能存放指定的数据类型

2） channle的数据放满后，就不能再放入了，不能动态扩容

3）如果从channel取出数据后，可以继续放入

4）**在没有使用协程的情况下**，如果channel数据取完了，再取，就会报dead lock

```go
func main()  {
	// 定义一个可以存放任意类型的管道
	//var allChan chan interface{}
	allChan := make(chan interface{}, 3)
	allChan <- 10
	allChan <- "tangyu"
	cat := Cat{"汤姆猫", 4}
	allChan <- cat
	<-allChan
	<-allChan
	// 获取对列第三个 要先将前2个推出
	newCat := <-allChan
	fmt.Printf("newCat=%T, newCat=%v\n", newCat, newCat)
  // newCat=main.Cat, newCat={汤姆猫 4}
	a := newCat.(Cat) // 类型断言 因为从interface的channel取出的值认为是空接口类型 需要类型断言进行转换
	fmt.Printf("newCat.Name=%v\n", a.Name)
	// newCat.Name=汤姆猫
}
```

### 2.8 channel的遍历和关闭

#### 2.8.1 channel的关闭

使用内置的close函数可以关闭channel，当channel**关闭后**，就**不能在向channel写数据了**，但是仍然**可以读出该channel的数据**

```go
func main() {
	intChan := make(chan int, 10)
	intChan <- 10
	intChan <- 100
	close(intChan)
	// close后不能再写入数据到channel
	intChan <- 200
	// 报错 panic: send on closed channel
}
```

#### 2.8.2 channel的遍历

1）在遍历时，如果channel没有关闭，则回出现**deadlock的错误**
2）在遍历时，如果channel已经关闭，则会正常遍历数据，遍历完后，就会退出遍历

```go
func main() {
	intChan2 := make(chan int, 100)
	for i := 0; i < 100; i++ {
		intChan2 <- i * 2
	}
	close(intChan2)
	// 遍历管道时不能用for循环 要用range
	for v := range intChan2 {
		fmt.Println("v=", v)
	}
}
```

## 3. goroutine和channel

### 3.1 应用案例一

开启一个writeData协程，向管道intChan中写入50个整数
开启一个readData协程，从管道intChan中读取writeData写入的数据
注意：writeData和readDate操作的是同一个管道
主线程需要等待writeData和readDate协程都完成工作才能退出

![Snipaste_2021-03-21_17-43-27](./asset_5/Snipaste_2021-03-21_17-43-27.png)

```go
func writeData(intChan chan int)  {
	for i:=1; i <= 50 ; i++{
		// 写入数据
		intChan <- i
	}
	close(intChan) // 关闭管道
}

func readData(intChan chan int, exitChan chan bool)  {
	for {
    // 读取不到数据库 会阻塞
		v, ok:=<-intChan
		// ok 为false时管道无数据
		if !ok {
			break
		}
		fmt.Println("读取到数据 >>> ", v)
	}
	// 读取完数据 任务完成
	exitChan <- true
	close(exitChan)
}

func main()  {
	// 创建两个管道
	intChan := make(chan int, 50)
	exitChan := make(chan bool, 1)
	go writeData(intChan)
	go readData(intChan, exitChan)
	for {
		_, ok := <-exitChan
		fmt.Println("ok>>>", ok)
		if !ok {
			break
		}
	}
}
```

### 3.2 阻塞

协程写的快读的慢会，写的时候会阻塞，不会报deadlock，但是如果编译器发现只有写没有读取数据的协程，编译器阻塞后会报错deadlock

读写协程的频率不一致，不会发生死锁

### 3.3 求素数

![Snipaste_2021-03-21_19-29-56](./asset_5/Snipaste_2021-03-21_19-29-56.png)

```go
func putNum(intChan chan int) {
	for i := 0; i < 1000; i++ {
		intChan <- i
	}
	close(intChan)
	fmt.Println("put over!!!")
}

func primeNum(intChan chan int, primeChan chan int, exitChan chan bool)  {
	var flag bool
	for {
		num, ok := <-intChan
		if !ok {
			break // 取不到数据退出
		}
		flag = true
		// 判断素数
		for i :=2 ; i < num; i++ {
			if num % i == 0 {
				// 不是素数
				flag = false
				break
			}
		}
		if flag {
			// 是素数 方式到primeChan
			primeChan <- num
		}
	}
	fmt.Println("有一个primeNum 协程因为取不到数据，退出")
	// 这里我们还不能关闭 primeChan
	// 向 exitChan 写入true
	exitChan<- true
}

func main() {
	intChan := make(chan int, 1000)
	primeChan := make(chan int, 2000) // 放入素数结果
	// 标识管道退出
	exitChan := make(chan bool, 4)

	// 放入数字
	go putNum(intChan)
	// 开启4个协程 判断是否为素数 如果是就放入到primeChan
	for i := 0; i < 4; i ++ {
		go primeNum(intChan, primeChan, exitChan)
	}
	go func() {
		for i := 0; i < 4; i ++ {
      // 不是range遍历 可以不用close管道
			<-exitChan
		}
		// 4个线程结束后 素数取完 关闭管道
		close(primeChan)
	}()
	for {
		res, ok := <- primeChan
		// primeChan管道关闭 ，素数判断完成
		if !ok{
			break
		}
		fmt.Printf("素数=%d\n", res)
	}
	fmt.Println("main线程退出")
}
```

注：使用range遍历管道时要使用close关闭数组，不然会包deadlock错误

### 3.4 channel使用细节

（1）channel可以只声明为只读或者只写

```go
func main() {
	// 默认情况下 管道是双向的 可读可写
	// 1.声明为只写
	var chan1 chan<- int
	chan1 = make(chan int, 3)
	chan1 <- 10

	// 2.声明为只读
	var chan2 <-chan int
}
```

![Snipaste_2021-03-21_23-51-03](./asset_5/Snipaste_2021-03-21_23-51-03.png)

（2）使用select可以解决管道读取数据的阻塞问题

```go
func main() {
	// 使用select解决管道取数据的阻塞问题
	intChan := make(chan int, 10)
	for i := 0; i < 10; i++ {
		intChan <- i
	}

	stringChan := make(chan string, 5)
	for i := 0; i < 10; i++ {
		stringChan <- "hello" + fmt.Sprintf("%d", i)
	}

	// 传统方式在遍历管道时 如果不关闭就会阻塞而导致deadlock
	// 在实际开发中 不好确实确定什么时候关闭管道
	// 使用select解决
	lable: //跳出for循环
	for {
		select {
		// 如果intChan一直没有关闭 不会一直阻塞而deadlock 会自动向下一个case
		case v := <-intChan:
			fmt.Println("从intChan读取的数据>>>", v)
		case v := <-stringChan:
			fmt.Println("从stringChan读取的数据>>>", v)
		default:
			fmt.Println("都取不到...")
			break lable // return
		}
	}
}
```

（3）goroutine中使用recove，解决协程中出现的panic，这样不会影响到主线程和其他协程执行

```go
func sayHello() {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		fmt.Println("i>>>", i)
	}
}

func test() {
	// defer + recover
	defer func() {
		// 捕获panic
		if err := recover(); err != nil {
			fmt.Println("test发生错误", err)
		}
	}()
	var myMap map[int]string
	// 没有make 就赋值发生错误
	myMap[0] = "tangyu"
}

func main() {
	go sayHello()
	go test()
	for i := 0; i < 10; i++ {
		fmt.Println("main() ok=", i)
		time.Sleep(time.Second)
	}
}
```

## 4. 反射





























































