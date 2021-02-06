package main

import (
	"fmt"
	"sync"
	"time"
)

const GB = 3
const max_thread = 4

func main() {
	var bytes [max_thread][GB * 1024 * 1024 * 1024 / 8 / max_thread]int64
	var waitGroup sync.WaitGroup
	waitGroup.Add(max_thread)
	fmt.Printf(" %vMB x %v Data Source Prepare ... ", GB*1024/max_thread, max_thread)
	t := time.Now()
	for thread := int64(0); thread < int64(max_thread); thread++ {
		go func(seed int64) {
			for i := int64(0); i < GB*1024*1024*1024/8/max_thread; i++ {
				bytes[seed][i] = seed * i
			}
			waitGroup.Done()
		}(thread)
	}
	waitGroup.Wait()
	fmt.Printf("Done in %v\n", time.Since(t))
	var new_bytes [max_thread][GB * 1024 * 1024 * 1024 / 8 / max_thread]int64
	waitGroup.Add(max_thread)
	fmt.Printf(" %vMB x %v Data Target Prepare ... ", GB*1024/max_thread, max_thread)
	t = time.Now()
	for thread := int64(0); thread < int64(max_thread); thread++ {
		go func(seed int64) {
			for i := int64(0); i < GB*1024*1024*1024/8/max_thread; i++ {
				new_bytes[seed][i] = 0
			}
			waitGroup.Done()
		}(thread)
	}
	waitGroup.Wait()
	fmt.Printf("Done in %v\n", time.Since(t))
	second, _ := time.ParseDuration("1s")
	waitGroup.Add(max_thread)
	fmt.Printf("\n %vMB x %v Memory Write ...", GB*1024/max_thread, max_thread)
	t = time.Now()
	total := t
	peak := 0
	for thread := int64(0); thread < int64(max_thread); thread++ {
		go func(seed int64) {
			for i := int64(0); i < GB*1024*1024*1024/8/max_thread; i++ {
				new_bytes[seed][i] = 1
			}
			fmt.Printf("\n  Thread %v Done in %v", seed, time.Since(t))
			total = total.Add(time.Since(t))
			if peak == 0 {
				peak ++
				fmt.Printf("  -> Peak %vGB/s", GB*1024*second/total.Sub(t))
			}
			waitGroup.Done()
		}(thread)
	}
	waitGroup.Wait()
	fmt.Printf("\n%v\n", 1024*second/total.Sub(t))
	fmt.Printf("\n %v GB Memory Write (%v Threads) Done in (AVG Time): %v (Bandwidth): %vGB/s\n", GB, max_thread, total.Sub(t)/max_thread, GB*1024*max_thread*second/total.Sub(t))
	waitGroup.Add(max_thread)
	fmt.Printf("\n %vMB x %v Memory Copy ...", GB*1024/max_thread, max_thread)
	t = time.Now()
	total = t
	peak = 0
	for thread := int64(0); thread < int64(max_thread); thread++ {
		go func(seed int64) {
			new_bytes[seed] = bytes[seed]
			fmt.Printf("\n  Thread %v Done in %v", seed, time.Since(t))
			total = total.Add(time.Since(t))
			if peak == 0 {
				peak ++
				fmt.Printf("  -> Peak %vGB/s", GB*1024*2*second/total.Sub(t))
			}
			waitGroup.Done()
		}(thread)
	}
	waitGroup.Wait()
	fmt.Printf("\n %v GB Memory Copy (%v Threads) Done in (AVG Time): %v (Bandwidth): %vGB/s\n", GB, max_thread, total.Sub(t)/max_thread, GB*1024*2*max_thread*second/total.Sub(t))
}
