package main

import (
	"fmt"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

func main() {

	mem, _ := linuxproc.ReadMemInfo("/proc/meminfo")

	totalMem := fmt.Sprintf("%f", float64(mem.MemTotal)*0.0009765625)
	usedMeme := fmt.Sprintf("%f", float64(mem.MemTotal-mem.MemFree)*0.0009765625)

	fmt.Println("Memeory usage: " + usedMeme + "/" + totalMem)

	// var prevIdle, prevTot uint64
	// for i := 0; i < 10; i++ {
	// 	stat, _ := linuxproc.ReadStat("/proc/stat")
	// 	cpu := stat.CPUStatAll
	// 	tot := cpu.User + cpu.Nice + cpu.System + cpu.Idle + cpu.IOWait + cpu.IRQ + cpu.SoftIRQ + cpu.Steal + cpu.Guest + cpu.GuestNice
	// 	if i > 0 {
	// 		deltaIdle := cpu.Idle - prevIdle
	// 		deltaTot := tot - prevTot
	// 		cpuUsage := (1.0 - float64(deltaIdle)/float64(deltaTot)) * 100.0
	// 		fmt.Printf("%d : %6.3f\n", i, cpuUsage)
	// 	}
	// 	prevIdle = cpu.Idle
	// 	prevTot = tot
	// 	time.Sleep(time.Second)

	// }

}
