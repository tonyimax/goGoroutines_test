package main

import (
	"fmt"
	"sync"
	"time"
)

// 普通函数
func f(from string) {
	//输出三次传入的字符串
	for i := 0; i < 50; i++ {
		fmt.Println(from, ":", i)
	}
}

// 信道作为参数使用
func workWithChannel(done chan bool) {
	fmt.Println("信道开始工作...")
	for i := 0; i < 5; i++ {
		fmt.Println(">>>信道内容处理输出:", i)
	}
	time.Sleep(1 * time.Second)
	fmt.Println("信道工作结束")
	done <- true //向信道写消息
}

// 信道直通
// pings : 信道
// msg : 消息
func ping(pings chan<- string, msg string) {
	pings <- msg //向信道写入消息
}

// pings信道与pongs信道直接通信
// pings:  pings信道
// pongs:  pongs信道
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings //读取pings信道的消息存储到msg
	pongs <- msg   //向pongs信道写入msg
}

// 工作线程
// id : 线程号
// jobs : 任务通道 (chan)
// results: 完成结果通道 (chan)
func worker(id int, jobs <-chan int, results chan<- int) {
	//遍历任务
	for j := range jobs {
		fmt.Println("工作协程: ", id, "启动任务: ", j)
		fmt.Println(">>>>>>休眠2秒，模拟工作处理数据中...")
		time.Sleep(time.Second)
		fmt.Println("工作协程: ", id, "结束任务", j)
		results <- j * 2 //任务完成后写入结果到通道
	}
}

func main() {
	f("Hello World")
	//使用go协程调用函数
	go f("go routines ==> Hello World")
	//使用协和调用临时函数
	go func(str string) {
		for i := range 10 {
			fmt.Println(i, str)
		}
	}("调用中...") //传参给临时函数
	time.Sleep(time.Second) //不休眠，看不到协和输出
	fmt.Println("完成")

	//信道
	fmt.Println("===创建无缓冲区信道msg===")
	msg := make(chan string) //双向信道
	fmt.Println("msg:当前通道已使用大小:", len(msg), "通道容量", cap(msg))
	//向信道发送数据
	//msg <- "发送到信道中的消息" //不使用协和向信道写入消息会报deadlock
	go func() { msg <- "==发送到信道中的消息==" }() //使用协程向信道中写消息
	fmt.Println("向信道中发送消息：==发送到信道中的消息==")
	fmt.Println("msg:当前通道已使用大小:", len(msg), "通道容量", cap(msg))
	//读取信道中的消息
	recvMsg := <-msg
	fmt.Println("读取信道中消息: ==>", recvMsg)
	fmt.Println("msg:当前通道已使用大小:", len(msg), "通道容量", cap(msg))
	fmt.Println("成功从信道读取信息：===》 ", recvMsg)

	//指定信道缓冲区大小
	fmt.Println("===创建缓冲区大小为2的信道msg2===")
	msg2 := make(chan string, 2) //信道容量为2
	fmt.Println("msg2:当前通道已使用大小:", len(msg2), "通道容量", cap(msg2))
	//向信道写数据
	msg2 <- "write message by channel one"
	fmt.Println("向通道中写入数据：write message by channel one")
	fmt.Println("msg2:当前通道已使用大小:", len(msg2), "通道容量", cap(msg2))
	msg2 <- "write message by channel two"
	fmt.Println("向通道中写入数据：write message by channel two")
	fmt.Println("msg2:当前通道已使用大小:", len(msg2), "通道容量", cap(msg2))
	//读取信道数据
	fmt.Println("读取信道数据:", <-msg2)
	fmt.Println("msg2:当前通道已使用大小:", len(msg2), "通道容量", cap(msg2))
	fmt.Println("读取信道数据:", <-msg2)
	fmt.Println("msg2:当前通道已使用大小:", len(msg2), "通道容量", cap(msg2))

	//信道同步
	done := make(chan bool, 1) //初始化信道参数
	workWithChannel(done)      //调用函数并传入信道

	//信道直通测试
	pings := make(chan string, 1) //创建信道
	pongs := make(chan string, 1) //创建信道
	//调用信道通信方法
	ping(pings, "Hello from pings ")
	pong(pings, pongs)      //双信道进行通信
	msgfrompongs := <-pongs //读取pongs信道消息
	fmt.Println("pongs信道消息:", msgfrompongs)

	//通道选择
	c1 := make(chan string) //双向通道
	c2 := make(chan string) //双向通道
	//协程1向通道1写数据
	go func() { c1 <- "hello world from c1" }()
	//协程2向通道2写数据
	go func() { c2 <- "hello world from c2" }()
	//遍历通道并选择对应通道操作
	for i := 0; i < 2; i++ {
		select {
		case recv1 := <-c1: //选择通道1
			fmt.Println("接收通道1数据:", recv1)
		case recv2 := <-c2: //选择通道2
			fmt.Println("接收通道2数据:", recv2)
		}
	}

	//单通道超时操作
	go func() {
		time.Sleep(2 * time.Second) //休眠2秒，使通道超时有效
		c1 <- "===再次向通道c1写入数据==="
	}()
	select {
	case msg := <-c1:
		fmt.Println("接收通道1数据:", msg)
	case <-time.After(1 * time.Second): //通道操作超时
		fmt.Println("通道操作超时")
	}

	//非阻塞通道操作
	message := make(chan string)  //处理字符串的通道
	signalChan := make(chan bool) //处理布尔类型的通道
	//不阻塞处理message通道
	select {
	case msg := <-message: //接收
		fmt.Println("接收message通道消息:", msg)
	default:
		fmt.Println("message:没有要处理的消息")
	}

	//向通道发送消息
	msgSend := "这是向通道发送的消息"
	select {
	case message <- msgSend: //发送
		fmt.Println("成功向message通道发送消息", msgSend)
	default:
		fmt.Println("没有要向message通道发送的消息")
	}

	//不阻塞处理signalChan通道
	select {
	case msg := <-signalChan:
		fmt.Println("signalChan:处理消息:", msg)
	default:
		fmt.Println("signalChan:没有要处理的消息")
	}

	//向signalChan通道发送消息
	select {
	case signalChan <- true:
		fmt.Println("成功向通道signalChan发送消息", <-signalChan)
	default:
		fmt.Println("没有要向signalChan通道发送的消息")
	}

	go func() { message <- "msg for channel message" }()
	go func() { signalChan <- false }()

	fmt.Println(cap(message), len(message), cap(signalChan), len(signalChan))
	time.Sleep(1 * time.Second)
	//接收多通道消息,随机处理message/signalChan通道
	select {
	case mul_msg := <-message:
		fmt.Println("===>接收message通道消息:", mul_msg)
	case mul_sig := <-signalChan:
		fmt.Println("===>接收signalChan通道消息", mul_sig)
	default:
		fmt.Println("无通信活动")
	}

	//关闭通道测试
	jobs := make(chan int, 5) //通道缓存5条消息
	finish := make(chan bool) //任务完成状态

	//启动一个协程
	go func() {
		fmt.Println("协程已启动...")
		fmt.Println("jobs通道大小:", len(jobs), "容量:", cap(jobs))
		for {
			valueInChannel, getValueOK := <-jobs //取通道消息
			if getValueOK {
				fmt.Println(">>>从通道中读取的值:", valueInChannel)
			} else {
				fmt.Println("任务完成！") //调用close来关闭通道后，才会执行到这里
				finish <- true
				return
			}
		}
	}()

	//向通道插入5条消息
	for i := 0; i < 5; i++ {
		jobs <- (i + 1) * 5 //插入消息
	}
	fmt.Println("所有任务发送完成!!!")
	//关闭通道
	close(jobs)
	<-finish
	_, ok := <-jobs //再次读取通道数据
	if !ok {
		fmt.Println("所有任务已完成，通道无任何数据!", len(jobs), len(finish))
	}
	fmt.Println(ok)
	time.Sleep(time.Second * 1) //休眠1秒

	//范围通道使用
	queue_channel := make(chan string, 2)
	//向通道写消息
	queue_channel <- "one"
	queue_channel <- "two"
	//关闭通道
	close(queue_channel)
	//遍历通道
	for ch := range queue_channel {
		fmt.Println("通道数据：", ch)
	}

	//定时器使用
	t1 := time.NewTimer(2 * time.Second)
	<-t1.C
	fmt.Println("timer1 fired")
	t2 := time.NewTimer(5 * time.Second)
	go func() {
		fmt.Println("go协程处理中,等待5秒后输出...")
		<-t2.C
		fmt.Println("timer2 fired")
	}()
	fmt.Println("主线程等待10秒...")
	t2Stop := t2.Stop() //停止定时器
	if t2Stop {
		fmt.Println("t2定时器已停止...", t2Stop)
	}
	time.Sleep(time.Second * 10)
	fmt.Println("完成")

	tick := time.NewTicker(500 * time.Millisecond) //500毫秒
	timeEnd := make(chan bool)                     //结束通道
	//协程
	go func() {
		//循环
		for {
			//随机执行任务
			select {
			case <-timeEnd: //结束任务
				fmt.Println("===结束任务")
				return
			case t := <-tick.C: //超时任务
				fmt.Println("==500毫秒响应一次：", t)
			}
		}
	}()
	time.Sleep(5000 * time.Millisecond) //休眠5秒
	tick.Stop()                         //停止定时器
	timeEnd <- true                     //设置结束标志
	fmt.Println("===定时任务结束===")

	const numJobs = 5                  //通道容量
	jobarr := make(chan int, numJobs)  //任务通道
	results := make(chan int, numJobs) //任务执行结果通道
	//启动3个工程协程
	for w := 1; w <= 3; w++ {
		go worker(w, jobarr, results)
	}
	//向任务通道发送5个任务
	for j := 1; j <= numJobs; j++ {
		jobarr <- j //发送任务到任务通道
	}
	close(jobarr) //关闭任务通道
	//遍历执行结果
	for a := 1; a <= numJobs; a++ {
		<-results //读取通道数据不做处理
	}

	dt1 := time.Now()
	//使用WaitGroup等待多个协程执行完成
	var wg = sync.WaitGroup{} //初始化协程等待对象
	fmt.Println("===启动10个协程序===", dt1)
	//启动10个协程序
	for i := 0; i < 10; i++ {
		wg.Add(1) //添加等待计数
		go func(id int) {
			defer wg.Done() //延迟执行，当前函数执行完成后，执行该语句
			fmt.Println("工作协程序已启动", id)
			//time.Sleep(1 * time.Second)
			fmt.Println("===工作协程序已完成===", id)
		}(i + 1)
	}
	wg.Wait() //等待所有协程执行完成
	dt2 := time.Now()
	fmt.Println("<<所有协程已执行完成>>", dt2)
	dtRet := dt2.Sub(dt1)
	aUnit := ""
	fmt.Println(dtRet, "毫秒：", dtRet.Milliseconds(), "微秒:", dtRet.Microseconds(), "纳秒:", dtRet.Nanoseconds())
	if (dtRet.Nanoseconds() / 1000) < 1000 {
		aUnit = "微秒"
	} else {
		aUnit = "毫秒"
	}
	fmt.Println("10个协程执行完成花费时间:", dtRet, aUnit)

	fmt.Println("=============速率控制================")
	//速率控制
	requests := make(chan int, 5) //创建通道
	//向通道写入5条数据
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)                              //关闭通道
	limiter := time.Tick(200 * time.Millisecond) //控制时间单位,200毫秒读取一次
	//遍历通道
	for req := range requests {
		<-limiter                                      //限制读取时间为200毫秒
		fmt.Println("通道数据:", req, "读取时间:", time.Now()) //输出读取到的通道数据与时间
	}

	burstyLimiter := make(chan time.Time, 3) //存储时间的通道
	//向通道写入3个时间
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	//使用协程每200毫秒向通道中写入一个时间
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			fmt.Println("===> 向通道中写入数据 ：", t)
			burstyLimiter <- t
		}
	}()

	for ret := range burstyLimiter {
		fmt.Println(ret)
	}
}
