package system

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

func CpuStatus() string {
	cpu, err := cpu.Percent(time.Second*1, false)
	if err != nil || len(cpu) < 1 {
		return ""
	}

	return fmt.Sprintf("[CPU:%.2f%%]", cpu[0])
}

func MemoryStatus() string {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("[MEM:%.2f%%]", vm.UsedPercent)
}

func LoadAverage() string {
	avg, err := load.Avg()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("[LoadAvg:%.2f %.2f %.2f]", avg.Load1, avg.Load5, avg.Load15)
}

func Info() string {
	return LoadAverage() + CpuStatus() + MemoryStatus()
}
