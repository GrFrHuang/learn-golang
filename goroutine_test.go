package models

import (
	"testing"
	"fmt"
	"sync"
	"io/ioutil"
	"os"
	"strings"
	//"runtime"
	"time"
	"github.com/skip2/go-qrcode"
	"reflect"
	"os/signal"
	"syscall"
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

func TestQRcode(t *testing.T) {
	//err := qrcode.WriteFile("http://blog.csdn.net/wangshubo1989", qrcode.Medium, 256, "qr.png")
	err := qrcode.WriteFile("weixin：//wxpay/s/An4baqw", qrcode.Medium, 256, "../static/qr.png")
	if err != nil {
		fmt.Println("write error")
	}
}

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
