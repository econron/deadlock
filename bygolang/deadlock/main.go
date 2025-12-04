package main

import (
	"fmt"
	"sync"
	"time"
)

type Resource struct {
	mu sync.Mutex
	value int
}

type SharedResource struct {
	mu sync.RWMutex
	value int
}

func main() {
	// deadlockWithTwoResources()
	// resolvedDeadlockWithTwoResources()
	// deadlockWithOneResource()
	resolvedDeadlockWithOneResource()
}

// 2つのリソースに2つのTXが交互に排他ロックを取るのでデッドロックになる
func deadlockWithTwoResources() {
	var wg sync.WaitGroup
	wg.Add(2)
	var resourceA Resource
	var resourceB Resource
	// TX1ではリソースAを先に排他ロック
	go func() {
		defer wg.Done()
		fmt.Println("TX1: Start")
		resourceA.mu.Lock()
		// リソースBを排他ロックする隙をTX2に与える
		time.Sleep(1 * time.Second)
		resourceB.mu.Lock()
		resourceA.value++
		resourceB.value++
		resourceA.mu.Unlock()
		resourceB.mu.Unlock()
		fmt.Println("TX1: End")
	}()
	// TX2ではリソースBを先に排他ロック
	go func() {
		defer wg.Done()
		fmt.Println("TX2: Start")
		resourceB.mu.Lock()
		// リソースAを排他ロックする隙をTX1に与える
		time.Sleep(1 * time.Second)
		resourceA.mu.Lock()
		resourceB.value++
		resourceA.value++
		resourceB.mu.Unlock()
		resourceA.mu.Unlock()
		fmt.Println("TX2: End")
	}()
	wg.Wait()
	fmt.Println("All transactions completed")	
}

// 同じ順序でリソースにアクセスしているのでデッドロックが起きてない
func resolvedDeadlockWithTwoResources() {
	var wg sync.WaitGroup
	wg.Add(2)
	var resourceA Resource
	var resourceB Resource
	go func() {
		defer wg.Done()
		fmt.Println("TX1: Start")
		resourceA.mu.Lock()
		time.Sleep(1 * time.Second)
		resourceB.mu.Lock()
		resourceA.value++
		resourceB.value++
		resourceA.mu.Unlock()
		resourceB.mu.Unlock()
		fmt.Println("TX1: End")
	}()
	go func() {
		defer wg.Done()
		fmt.Println("TX2: Start")
		resourceA.mu.Lock()
		time.Sleep(1 * time.Second)
		resourceB.mu.Lock()
		resourceA.value++
		resourceB.value++
		resourceB.mu.Unlock()
		resourceA.mu.Unlock()
		fmt.Println("TX2: End")
	}()
	wg.Wait()
	fmt.Printf("resourceA.value: %d, resourceB.value: %d\n", resourceA.value, resourceB.value)
	fmt.Println("All transactions completed")
}

// 1つのリソースに対して2つのTXが共有ロックを交互に取り、そのあと書き込みを実行しようと排他ロックを取ろうとしてデッドロックになる
func deadlockWithOneResource() {
	var wg sync.WaitGroup
	var resource SharedResource
	wg.Add(2)
	go func(){
		defer wg.Done()
		fmt.Println("TX1: Start")
		resource.mu.RLock()
		time.Sleep(1 * time.Second)
		resource.mu.Lock()
		resource.value++
		resource.mu.Unlock()
		resource.mu.RUnlock()
		fmt.Println("TX1: End")
	}()
	go func(){
		defer wg.Done()
		fmt.Println("TX2: Start")
		resource.mu.RLock()
		time.Sleep(1 * time.Second)
		resource.mu.Lock()
		resource.value++
		resource.mu.Unlock()
		resource.mu.RUnlock()
		fmt.Println("TX2: End")
	}()
	wg.Wait()
	fmt.Printf("resource.value: %d\n", resource.value)
	fmt.Println("All transactions completed")
}

// 最初から排他ロックを使う。InnoDBで言うと同じレコードに対して排他ロックを最初から取って順番にTXを処理するイメージ。
func resolvedDeadlockWithOneResource() {
	var wg sync.WaitGroup
	var resource Resource
	wg.Add(2)
	go func(){
		defer wg.Done()
		fmt.Println("TX1: Start")
		resource.mu.Lock()
		time.Sleep(1 * time.Second)
		resource.value++
		resource.mu.Unlock()
		fmt.Println("TX1: End")
	}()
	go func(){
		defer wg.Done()
		fmt.Println("TX2: Start")
		resource.mu.Lock()
		time.Sleep(1 * time.Second)
		resource.value++
		resource.mu.Unlock()
		fmt.Println("TX2: End")
	}()
	wg.Wait()
	fmt.Printf("resource.value: %d\n", resource.value)
	fmt.Println("All transactions completed")
}