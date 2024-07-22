package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L./lib -lrocm_smi64 -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
*/
import "C"
import (
	"fmt"

	"github.com/golang/glog"
)

// rsmiInit 初始化rocm_smi
func rsmiInit() (err error) {
	ret := C.rsmi_init(0)
	glog.Info("go_rsmi_init_ret:", ret)
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error go_rsmi_init: %s", err)
	}
	return nil
}

// rsmiShutdown 关闭rocm_smi
func rsmiShutdown() (err error) {
	ret := C.rsmi_shut_down()
	glog.Info("go_rsmi_shutdown_ret:", ret)
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_shutdown: %s", err)
	}
	return nil
}
