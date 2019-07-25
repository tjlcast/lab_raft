# 记录

## Golang

+ timer.After()
    - 返回一个 channel，并启动一个匿名 goroutine，向这个 channel 定时有且只输入一个 msg.
    - 可以与 select 结合.
        ```golang
             /**
                下面的第一个case永远不会执行。因为在执行select的时候，没有从第一个 case channel中读取到数据，
                转而执行default case
             **/
             select {
                  case timeon := <= time.After(3 * time.Seconds):
            	    fmt.Println("it's time.")
                  case default:
            	    fmt.Println("it's defualt.")
             }
        ```