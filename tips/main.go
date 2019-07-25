package main

import (
	"time"
	"fmt"
)

func main() {
	//closeChannel()
	c := make(chan int)
	timeout := time.After(time.Second * 2) //
	t1 := time.NewTimer(time.Second * 3) // 效果相同 只执行一次
	var i int
	go func() {
		for {
			select {
			case <-c:
				fmt.Println("channel sign")
				return
			case <-t1.C: // 代码段2
				fmt.Println("3s定时任务")
			case <-timeout: // 代码段1
				i++
				fmt.Println(i, "2s定时输出")
			case <-time.After(time.Second * 4): // 代码段3
				fmt.Println("4s timeout。。。。")
			default:    // 代码段4
				fmt.Println("default")
				time.Sleep(time.Second * 1)
			}
		}
	}()
	time.Sleep(time.Second * 6)
	close(c)
	time.Sleep(time.Second * 2)
	fmt.Println("main退出")
}


//发送者
func sender(c chan int) {
	for i := 0; i < 100; i++ {
		c <- i
		if i >= 5 {
			time.Sleep(time.Second * 7)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func test2() {
	c := make(chan int)
	go sender(c)
	timeout := time.After(time.Second * 3)
	for {
		select {
		case d := <-c:
			fmt.Println(d)
		case <-timeout:
			fmt.Println("这是定时操作任务 >>>>>")
		case dd := <-time.After(time.Second * 3):
			fmt.Println(dd, "这是超时*****")
		}

		fmt.Println("for end")
	}
}
