package handle

import "github.com/itchin/proxy/server/config"

var Capacity capacity

type capacity struct{
    mark chan int
}

func init() {
    Capacity.mark = make(chan int, config.MAX_ACTIVE)
    for i := 0; i < config.MAX_ACTIVE; i++ {
        Capacity.mark <- i
    }
}

func (q *capacity) Shift() int {
    i := <- q.mark
    return i
}

func (q *capacity) Push(i int) {
    q.mark <- i
}
