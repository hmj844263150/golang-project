package test

import (
	"time"
	"utilProf"
)

func preReportRusage(interval time.Duration) func(usage1 *syscall.Rusage, utilProf *SystemUtilProf) int {
	var lastUtime int64
	var lastStime int64
	var waitInterval time.Duration = interval

	a := func(usage1 *syscall.Rusage, utilProf *SystemUtilProf) int {
		if err := syscall.Getrusage(syscall.RUSAGE_SELF, usage1); err == nil {
			utime := usage1.Utime.Sec*1000000000 + usage1.Utime.Usec
			stime := usage1.Stime.Sec*1000000000 + usage1.Stime.Usec

			utilProf.userCPUUtil = float64(utime-lastUtime) * 100 / float64(waitInterval) //用户cpu时间
			utilProf.sysCPUUtil = float64(stime-lastStime) * 100 / float64(waitInterval)  //系统cpu时间
			utilProf.memUtil = uint64(usage1.Maxrss * 1024)                               //系统内存使用情况
			utilProf.memStats = runtime.MemStats{}

			runtime.ReadMemStats(&utilProf.memStats)

			lastUtime = utime
			lastStime = stime

			return kSuccess
		} else {
			return kFail
		}
	}

	return a
}

func funcC() {
　　　　rusageReport := preReportRusage(time.Minute)
    //print rusage status...
    usage1 := &syscall.Rusage{}
    utilProf := &SystemUtilProf{}
 
    if result := rusageReport(usage1, utilProf); result == kSuccess {
        content = fmt.Sprintf("%3.2f,%3.2f,%s,%s,%s,%s\n",
        utilProf.userCPUUtil, utilProf.sysCPUUtil, toH(utilProf.memUtil),
        toH(utilProf.memStats.HeapSys), toH(utilProf.memStats.HeapAlloc), toH(utilProf.memStats.HeapIdle))
    }
}
