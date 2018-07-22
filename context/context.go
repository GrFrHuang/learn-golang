package main

import (
	"time"
	"fmt"
	"context"
	"net/http"
	"sync"
)

// Context 通常被译作 上下文 ，一般理解为程序单元的一个运行状态、现场、快照，而翻译中 上下 又很好地诠释了其本质，
// 上下上下则是存在上下层的传递， 上 会把内容传递给 下 。
// 在Go语言中，程序单元也就指的是Goroutine。每个Goroutine在执行之前，都要先知道程序当前的执行状态，
// 通常将这些执行状态封装在一个Context变量中，传递给要执行的Goroutine。
// 上下文则几乎已经成为传递与请求同生存周期变量的标准方法。

// context是一个在go中时常用到的程序包，google官方开发。
// 特别常见的一个应用场景是由一个请求衍生出的各个goroutine之间需要满足一定的约束关系，
// 以实现一些诸如有效期，中止routine树，传递请求全局变量之类的功能。
// 使用context实现上下文功能约定需要在方法的传入参数的第一个传入一个context.Context类型的变量。
// 给一个函数方法传递Context的时候，不要传递nil，如果不知道传递什么，就使用context.TODO，
// context是go程安全的，可以放心的在多个goroutine中传递。
// 我们将通过源代码的阅读和一些示例代码来说明context的用法。

// 模拟一个最小执行时间的阻塞函数
func inc(a int) int {
	res := a + 1                // 虽然我只做了一次简单的 +1 的运算,
	time.Sleep(8 * time.Second) //强行慢操作
	return res
}

// 如果计算被中断, 则返回 -1
func Add(ctx context.Context, a, b int) int {
	res := 0
	for i := 0; i < a; i++ {
		res = inc(res)
		select {
		// Done 方法在Context被取消或超时时返回一个close的channel,close的channel可以作为广播通知，
		// 告诉给context相关的goroutine函数要停止当前工作然后返回
		case <-ctx.Done():
			// Err方法返回context为什么被取消
			err := ctx.Err()
			fmt.Println("why cancel: ", err)
			// Deadline返回context何时会超时
			deadLineTime, ok := ctx.Deadline()
			fmt.Println("when deadline: ", deadLineTime, ok)
			fmt.Println("now time is : ", time.Now())
			return -1
		default:
			return res
			// 没有结束 ... 执行 ...
		}
	}
	for i := 0; i < b; i++ {
		res = inc(res)
		select {
		case <-ctx.Done():
			return -1
		default:
			fmt.Println("res b = ", res)
			return res
			// 没有结束 ... 执行 ...
		}
	}
	return res
}

// 测试并发协程不安全情况
// 读写共享元素,在并发情况下,则必然产生竞争,破坏共享元素数据, 所以要保护,要么加锁, 要么用channel将访问排队串行化
func main() {
	//c := make(map[string]int)
	//for i := 0; i < 100; i++ {
	//	go func() {
	//		for j := 0; j < 200; j++ {
	//			c[fmt.Sprintf("%d", j)] = j+1
	//		}
	//	}()
	//}
	var l sync.Mutex
	var a = 0
	for i := 0; i < 1000; i++ {
		if a < 500 {
			l.Lock()
			go func() {
				a = a + 1
				l.Unlock()
			}()
		}
	}
	time.Sleep(time.Second * 2)
	fmt.Println(a)
}

//func main() {
//	{
//		// 使用开放的 API 计算 a+b
//		a := 1
//		b := 2
//		fmt.Println("now is : ", time.Now())
//		timeout := 5 * time.Second
//		fmt.Println("The background context is : ", context.Background())
//		// WithTimeout 等价于 WithDeadline(parent, time.Now().Add(timeout))
//		ctx, _ := context.WithTimeout(context.Background(), timeout)
//		res := Add(ctx, 1, 2)
//		fmt.Println("child notify parent : ", <-ctx.Done())
//		fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)
//	}
//	{
//		// The WithCancel, WithDeadline, and WithTimeout functions take a Context (the
//		// parent) and return a derived Context (the child) and a CancelFunc. Calling
//		// the CancelFunc cancels the child and its children, removes the parent's
//		// reference to the child, and stops any associated timers.
//		_ctx := context.Background()
//		ctx := context.WithValue(_ctx, "GrFrHuang", "good")
//		fmt.Println("Found the key ", ctx.Value("GrFrHuang"))
//		fmt.Println("Not found ", ctx.Value("Baibaobao"))
//
//		time.Sleep(time.Second * 5)
//		fmt.Println("I am run here")
//	}
//	//{
//	//	// 手动取消
//	//	a := 1
//	//	b := 2
//	//	ctx, cancel := context.WithCancel(context.Background())
//	//	go func() {
//	//		time.Sleep(2 * time.Second)
//	//		cancel() // 在调用处主动取消
//	//	}()
//	//	res := Add(ctx, 1, 2)
//	//	fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)
//	//}
//}

// 要么ctx被取消，要么request请求出错。
func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}

// AddContextSupport是一个中间件，用来绑定一个context到原来的handler中，所有的请求都必须先经过该中间件后才能进入各自的路由处理中
func AddContextSupport(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, "-", r.RequestURI)
		cookie, _ := r.Cookie("username")
		if cookie != nil {
			ctx := context.WithValue(r.Context(), "username", cookie.Value)
			// WithContext returns a shallow copy of r with its context changed
			// to ctx. The provided ctx must be non-nil.
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
