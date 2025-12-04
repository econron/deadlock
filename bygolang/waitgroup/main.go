package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	// 2つのgoroutineが終わるまで待てという
	wg.Add(2)
	// goroutine 1
	go func() {
		// 仕事の完了を報告するもの
		defer wg.Done()
		fmt.Println("Start TX1")
		time.Sleep(1 * time.Second)
		fmt.Println("End TX1")
	}()
	// goroutine 2
	go func() {
		// 仕事の完了を報告するもの
		defer wg.Done()
		fmt.Println("Start TX2")
		time.Sleep(1 * time.Second)
		fmt.Println("End TX2")
	}()
	// ここでgoroutineが終わるまで待つよう指示する
	wg.Wait()
	fmt.Println("All transactions completed")
}