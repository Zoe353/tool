package main

import (
	"io/ioutil"
	"net/http"
	"flag"
	"fmt"
	"container/list"
	"time"
	"container/heap"
)


type IntHeap []int64

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int64))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Remove(h *IntHeap, i int) interface{} {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !down(h, i, n) {
			up(h, i)
		}
	}
	return h.Pop()
}
func up(h *IntHeap, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func down(h *IntHeap, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}


func main() {
	// add command line args
	url := flag.String("url", "https://zsunapp.zoey353.workers.dev/links", "url")
	profile := flag.Int("profile", 3, "the number of requests to the website")
	
	flag.Parse()

	err_codes := list.New()
	smallest_resp := -1
	largest_resp := -1
	var times int64 = 0
	h := &IntHeap{}
	heap.Init(h)

	// make requests to the website
	for i:=0; i < *profile; i++{

		// use http library
		start := time.Now()
		resp, err := http.Get(*url)
		end := time.Now()
		if err != nil{
		  err_codes.PushBack(err.Error())
		  continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		  err_codes.PushBack(err.Error())
		  continue
		}
		 
		sb := string(body)
		duration := end.Sub(start).Milliseconds()
		times += duration
		heap.Push(h,duration)

		if largest_resp == -1 || largest_resp < len(sb){
			largest_resp = len(sb)
		}
		if smallest_resp == -1 || smallest_resp > len(sb){
			smallest_resp = len(sb)
		}

		fmt.Println("********response body--begin********")
		fmt.Println(sb)
		fmt.Println("********response body--end********")

	}

	mean := 0
	if h.Len() > 0{
		mean = int(times) / h.Len()
	}
	var median int64 = 0
	var min int64 = 0
	var max int64 = 0
	l := h.Len()
	for i:= 0; i < l; i++ {
		if i == 0 {
			min = heap.Pop(h).(int64)
		}
		if i == l/2{
			median = heap.Pop(h).(int64)
		}
		if i == l-1 {
			max = heap.Pop(h).(int64)
		}
	}

	

	// print results
    fmt.Println("***************************Evaluation logic***************************")
	fmt.Println("The number of requests: ",*profile)
	fmt.Println("The fastest time: ", min,"ms")
	fmt.Println("The slowest time: ", max,"ms")
	fmt.Println("The mean & median times: ",mean,"ms, ", median, "ms")
	fmt.Println("The percentage of requests that succeeded: ", (*profile-err_codes.Len()) / *profile * 100,"%")
	fmt.Print("Any error codes returned that weren't a success: ")
	for i := err_codes.Front(); i != nil; i = i.Next() {
	    fmt.Println(i.Value)
	}
	fmt.Println(" ")
	fmt.Println("The size in bytes of the smallest response: ", smallest_resp)
	fmt.Println("The size in bytes of the largest response: ",largest_resp)
	fmt.Println("**********************************************************************")


}