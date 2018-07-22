package test

import (
	"testing"
	"fmt"
	"sync"
	"io/ioutil"
	"os"
	"strings"
	//"runtime"
	"time"
	"reflect"
	"os/signal"
	"syscall"
	"strconv"
	"net/rpc"
	"net"
	"log"
	"net/http"
)

var c chan bool
var wg sync.WaitGroup

func TestGoroutine(t *testing.T) {
	slice := []int{}
	value := 1
	n := 100
	c = make(chan bool, 100)
	for i := 0; i < n; i++ {
		slice = append(slice, i)
	}
	//协程与主协程速度严重不匹
	for _, v := range slice {
		//wg.Add(1)
		//设置最大线程数
		//runtime.GOMAXPROCS(2)
		//主线程在这个地方出让cpu时间片可以让goroutine都跑完
		//runtime.Gosched()
		go func(v1 int) {
			//defer wg.Done()
			value++
			fmt.Println("=", v1)
			c <- true
		}(v)
	}
	fmt.Print("run here")
	var a int
	e := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-c:
			fmt.Println("he", <-c)
		default:
			a = 1
			//fmt.Println("default")
		}
		select {
		case <-e.C:
			fmt.Printf("超时")
			break
		}
	}
	//wg.Wait()
	fmt.Println("total is: ", value)
	fmt.Print(a)
}

var b chan int

//从这里可以看出，对于无缓冲的channel，放入操作和取出操作不能再同一个routine中，
//而且应该是先确保有某个routine对它执行取出操作，然后才能在另一个routine中执行放入操作
func TestChannel(t *testing.T) {
	b = make(chan int, 1)
	fmt.Println("hello")
	b <- 1
	fmt.Println("huang")
	fmt.Println("输出:", <-b)
	//<-c
}

func TestIoutil(t *testing.T) {
	bs, err := ioutil.ReadFile("/home/huang/shell_script/shell/test2.sh")
	fmt.Println(string(bs), err)
	var mode os.FileMode = 755
	f, err := os.OpenFile("/home/huang/shell_script/shell/test2.sh", 1, mode)
	fmt.Println(f.Name(), err)
	//读取目录下的文件列表
	fs, _ := ioutil.ReadDir("/home/huang/shell_script/shell/")
	for k, v := range fs {
		fmt.Println("文件列表", k+1, ": ", v.Name())
	}
	//var b []byte
	r := strings.NewReader("hello Friend Huang")
	b, _ := ioutil.ReadAll(r)
	fmt.Println("读字符串", string(b))
	l := make([]byte, 50)
	_, err = r.Read(l)
	fmt.Println(l, len(l), err)
}

//func TestQRcode(t *testing.T) {
//	//err := qrcode.WriteFile("http://blog.csdn.net/wangshubo1989", qrcode.Medium, 256, "qr.png")
//	err := qrcode.WriteFile("weixin：//wxpay/s/An4baqw", qrcode.Medium, 256, "../static/qr.png")
//	if err != nil {
//		fmt.Println("write error")
//	}
//}

type People struct {
	Name    string `json:"name"`
	Age     int    `orm:"column(age)"`
	Address string `json:"address"`
}

func TestReflect(t *testing.T) {
	//操作变量
	var i int = 5
	value := reflect.ValueOf(&i)
	value = reflect.Indirect(value)
	fmt.Println("before set: ", value.Interface())
	value.SetInt(2)
	fmt.Println(value.Interface())

	p := People{
		Name:    "huang",
		Age:     23,
		Address: "zhong",
	}
	v := reflect.ValueOf(p)
	fmt.Println(v.Kind())
	fmt.Println(v.FieldByName("Address"))
	//	获取value的tag
	typ := v.Type()
	addr, _ := typ.FieldByName("Address")
	fmt.Println(addr.Tag.Get("json"))
	//addr.Name = "Address2"
	fmt.Println("===== ", addr.Name)
	Age, _ := typ.FieldByName("Age")
	fmt.Println(Age.Tag.Get("orm"))

}

//测试slice的初始值
func TestAppend(t *testing.T) {
	array := make([]string, 2)
	array = append(array, "1", "2", "3")
	fmt.Println(array)
	fmt.Println(len(array))
	fmt.Println(cap(array))
}

//panic:重复关闭channel||向已经关闭的channel写入
//以<-chan的形式读没有值的channel,输出channel默认类型的零值(非空)
func TestChannelReadAndWrite(t *testing.T) {
	cha := make(chan int)
	close(cha)
	//cha <- "hello"
	//go func() {
	//	cha <- 2
	//	close(cha)
	//	fmt.Println("hello")
	//}()
	fmt.Println(<-cha)
	fmt.Println(<-cha)
	//for v := range cha {
	//	fmt.Println(v)
	//}
	if val, isClose := <-cha; !isClose {
		fmt.Println(val, isClose)
	}
}

//defer在panic||return以后声明,不会触发这个钩子
//注意defer的声明顺序
func TestPanic(t *testing.T) {
	defer func() {
		fmt.Println("GrFrHuang")
		if err := recover(); err != nil {
			fmt.Println("run here ", err)
		}
	}()
	var a = 2
	var b = 0
	c := a / b
	fmt.Println(c)
	panic("hello")
}

//defer这个函数的默认参数是return返回的值
//注意defer匿名函数调用与普通函数调用的区别
func TestDefer(t *testing.T) {
	var a = 0
	var b = 0
	fmt.Println(a)
	a++
	defer func() {
		fmt.Println(b)
	}()
	b++
	return
}

type Operater struct {
	Lock *sync.RWMutex
	Wait *sync.WaitGroup
}

//todo test RWMutex
func TestRWMutex(t *testing.T) {
	opter := Operater{
		Lock: &sync.RWMutex{},
		Wait: &sync.WaitGroup{},
	}
	for i := 0; i < 3; i++ {
		opter.Wait.Add(2)
		go read(opter, i)
		go write(opter, i)
	}
	opter.Wait.Wait()
}

func read(opter Operater, i int) {
	defer opter.Wait.Done()
	fmt.Println("I am start read", i)
	opter.Lock.RLock()
	fmt.Println("I am reading", i)
	time.Sleep(time.Second * 15)
	opter.Lock.RUnlock()
	fmt.Println("I am over read", i)
}

func write(opter Operater, i int) {
	defer opter.Wait.Done()
	fmt.Println("I am start write", i)
	opter.Lock.Lock()
	fmt.Println("I am writing", i)
	time.Sleep(time.Second * 15)
	opter.Lock.Unlock()
	fmt.Println("I am over write", i)
}

var l sync.RWMutex

//重复解锁会fatal,加锁会死锁
func TestRepeatLockOperate(t *testing.T) {
	fmt.Println("I am run here 1")
	l.Lock()
	fmt.Println("I am run here 2")
	//l.Unlock()
	l.RLock()
	fmt.Println("I am run here 3")
}

func TestSystemSignal(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan int, 100)
	chSend := make(chan int)
	chConsume := make(chan int)
	sc := make(chan os.Signal, 1)

	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func(ch, chSend chan int) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("send to ch panic.===", err)
			}
		}()

		i := 0
		for {
			//select和switch非常相似,不过select的case里的操作语句只能是(IO操作)
			//有一个case可以执行,当前select监听结束
			select {
			case ch <- i:
				fmt.Println("send", i)
				time.Sleep(time.Second)
				i++
			case <-chSend:
				fmt.Println("send quit.")
				return
			}
		}
	}(ch, chSend)

	go func(ch, chConsume chan int) {
		wg.Add(1)
		for {
			select {
			case i, ok := <-ch:
				if ok {
					fmt.Println("read", i)
					time.Sleep(time.Second * 2)
				} else {
					fmt.Println("close ch1.")
				}

			case <-chConsume:
				for {
					select {
					case i, ok := <-ch:
						if ok {
							fmt.Println("read2", i)
							time.Sleep(time.Second * 2)
						} else {
							fmt.Println("close ch2.")
							goto L
						}
					}
				}
			L:
				fmt.Println("consume quit.")
				wg.Done()
				return
			}
		}
	}(ch, chConsume)

	<-sc

	close(ch)
	fmt.Println("close ch ")
	close(chSend)
	close(chConsume)
	wg.Wait()
}

// 函数形参直接声明为结构体或接口
func Structs(s struct {
	Sex string
}) {
	fmt.Println(s.Sex)
}

func TestStructs(t *testing.T) {
	type Sexs struct {
		Sex string
	}
	s := Sexs{Sex: "男"}
	Structs(s)
}

//golang使用接口来实现多态性
//可以将子接口类型的变量赋值给父接口类型的变量
type Human interface {
	getInfo()
}

type Huang struct {
	Name string
}

type Bai struct {
	Age int
}

func (H Huang) getInfo() {
	fmt.Println(H.Name)
}

func (H Bai) getInfo() {
	fmt.Println(H.Age)
}

func TestRedis(t *testing.T) {
	var human Human
	human = Huang{Name: "Hello GrFrHuang"}
	human.getInfo()
	human = Bai{Age: 24}
	human.getInfo()
}

// golang组合继承
type As struct {
	Name string
}

func (a *As) GetName() string {
	return "hello" + a.Name
}

type Ha struct {
	As
	Age  int
	Name <-chan string
}

func (h *Ha) GetName() int {
	return h.Age
}
func TestCount(t *testing.T) {
	c := make(chan string, 1)
	c <- "2"
	h := Ha{Age: 24, Name: c}
	fmt.Printf("%+v", h)
}

// <-chan read only
// chan <- write only
// 使用这样的方法查看channel是否关闭，会将里面的元素取出一个，channel长度减一
// if _, ok := <- channel; !ok {
//		log.Warn("channel has been closed")
//	}
// 每做一次这样的操作<-channel都会让channel里的值少一个

// 多维map的问题
// make一个二维 map，结果碰到了nil map
// m:= make(map[string]map[string]int)
// m["a"]["b"] = 2
// 必须这样才行
// m := make(map[string]map[string]int)
// m2:= make(map[string]int)
// m2["b"] = 1
// m["a"] = m2

// 已经提示很明确了，all goroutines are asleep，所有协程都在睡觉
// 所以go就认为是死锁了。至少有一个协程要是在干活的
// 对于死锁的检测非常麻烦，或许go就采用了这种比较简单粗暴的方法。算是一个小bug吧，但是关系不大

// 目标: 测试主协程跑完，子协程是否继续异步执行
// 结果: 顶层协程一旦执行完,其他协程也会被回收
func TestAsyncGoroutine(t *testing.T) {
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		go func() {
			for i := 0; i < 500; i++ {
				rsp, _ := http.Get("http://127.0.0.1:18000/")
				//f, err := os.Create("test.txt")
				//bts, err := ioutil.ReadAll(rsp.Body)
				fmt.Println(rsp)
			}
			time.Sleep(time.Millisecond * 5000)
			ch2 <- "[底级]: 我跑完了"
		}()
		ch1 <- "[次级]: 你们跑,跟我无关"
	}()
	fmt.Println(<-ch1)
	fmt.Println(<-ch2)
	time.Sleep(time.Millisecond * 10000)
	fmt.Println("[顶级]: 你们跑完了,我才可以安全返回")
}

var m = make(map[int]int)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	//var count int
	//count++
	m[1] = 100
	fmt.Println(m[1])
	fmt.Fprintln(w, "hello world")
}

func TestHttp(t *testing.T) {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe("127.0.0.1:18000", nil)
}

//RPC test
type Args struct {
	A, B int
}

type Bean int

func (t *Bean) Multiply(args *Args, reply *([]string)) error {
	*reply = append(*reply, strconv.Itoa(args.B), "GrFrHuang")
	return nil
}

func TestRpc(t *testing.T) {
	newServer := rpc.NewServer()
	newServer.Register(new(Bean))

	lst, e := net.Listen("tcp", "127.0.0.1:1234") // any available address
	if e != nil {
		log.Fatalf("net.Listen tcp :0: %v", e)
	}

	//Listen tcp port for rpc request.
	go newServer.Accept(lst)
	//newServer.HandleHTTP("/foo", "/bar")

	time.Sleep(5 * time.Second)

	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	conn, _ := net.DialTCP("tcp", nil, address)
	defer conn.Close()

	client := rpc.NewClient(conn)
	defer client.Close()

	args := &Args{7, 8}
	result := make([]string, 10)
	err = client.Call("Bean.Multiply", args, &result)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Println(result)
}

func TestMy(t *testing.T) {
	ch := make(chan string)
	forever := make(chan string)
	go func() {
		go func() {
			for {
				fmt.Println("顶级协程不返回的话我就还在跑哦")
				time.Sleep(time.Second * 1)
			}
		}()
		time.Sleep(time.Second * 2)
		ch <- "跑完了"
	}()
	fmt.Println("我知道: ", <-ch)
	<-forever
}

// 通过对栈上的引用表量声明类型,不仅可以让该变量存储对应类型的堆地址,还可以让该变量去静态方法区获取到变量对应的类型的成员属性和成员函数.
// 栈内存是在编译以前用代码事先分配(用完则该语言的虚拟机立即回收),堆内存是在程序运行时动态分配的(等待该语言的gc回收).
// 程序计数器/PC寄存器--记录当前时间片内,当前线程(协程)的机器指令所执行到的行号.
// 栈(先进后出)的存取速度比堆(链表结构,物理地址不需要连续,但是逻辑地址一定要连续)要快.
// 基本数据类型的值也保存在栈上(int,float).
// 方法区是静态区,是线程(协程)共享的.
// 每一个方法栈都是线程(协程)栈的栈帧.
// 并发时的堆对象操作.
func TestBingfa(t *testing.T) {
	var as = &As{
		Name: "123",
	}
	// as2与as指向的是同一个堆对象
	var as2 = as
	// as3已经被重新分配了内存地址
	var as3 = &As{
		Name: "123",
	}

	// 副本拷贝(在方法栈上重新分配堆的引用as5,再在堆上重新分配同一个对象)
	var as4 = As{
		Name: "nihao",
	}
	var as5 = as4
	for i := 0; i < 100; i++ {
		go func() {
			as5.Name = "520"

			as3.Name = "mei"
			as2.Name = "hee"
			as.Name = "gg"
			fmt.Println(as, as2, as3, as4)
		}()
	}
	time.Sleep(time.Second * 2)
}
