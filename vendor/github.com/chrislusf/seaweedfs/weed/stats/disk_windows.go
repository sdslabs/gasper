package stats

import (
	"github.com/chrislusf/seaweedfs/weed/pb/volume_server_pb"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	kernel32           = windows.NewLazySystemDLL("Kernel32.dll")
	getDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")
)

func fillInDiskStatus(disk *volume_server_pb.DiskStatus) {

	ptr, err := syscall.UTF16PtrFromString(disk.Dir)

	if err != nil {
		return
	}
	var _temp uint64
	/* #nosec */
	r, _, e := syscall.Syscall6(
		getDiskFreeSpaceEx.Addr(),
		4,
		uintptr(unsafe.Pointer(ptr)),
		uintptr(unsafe.Pointer(&disk.Free)),
		uintptr(unsafe.Pointer(&disk.All)),
		uintptr(unsafe.Pointer(&_temp)),
		0,
		0,
	)

	if r == 0 {
		if e != 0 {
			return
		}

		return
	}
	disk.Used = disk.All - disk.Free
	disk.PercentFree = float32((float64(disk.Free) / float64(disk.All)) * 100)
	disk.PercentUsed = float32((float64(disk.Used) / float64(disk.All)) * 100)

	return
}
