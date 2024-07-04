package dcgm

/*
#cgo CFLAGS: -Wall -I/opt/dtk-24.04/rocm_smi/include/rocm_smi
#cgo LDFLAGS: -L/opt/dtk-24.04/rocm_smi/lib -lrocm_smi64 -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
*/
import "C"
import "fmt"

type RSMIStatus C.rsmi_status_t

const (
	RSMI_STATUS_SUCCESS       RSMIStatus = C.RSMI_STATUS_SUCCESS       //!< Operation was successful
	RSMI_STATUS_INVALID_ARGS  RSMIStatus = C.RSMI_STATUS_INVALID_ARGS  //!< Passed in arguments are not valid
	RSMI_STATUS_NOT_SUPPORTED RSMIStatus = C.RSMI_STATUS_NOT_SUPPORTED //!< The requested information or
	//!< action is not available for the
	//!< given input, on the given system
	RSMI_STATUS_FILE_ERROR RSMIStatus = C.RSMI_STATUS_FILE_ERROR //!< Problem accessing a file. This
	//!< may because the operation is not
	//!< supported by the Linux kernel
	//!< version running on the executing
	//!< machine
	RSMI_STATUS_PERMISSION RSMIStatus = C.RSMI_STATUS_PERMISSION //!< Permission denied/EACCESS file
	//!< error. Many functions require
	//!< root access to run.
	RSMI_STATUS_OUT_OF_RESOURCES RSMIStatus = C.RSMI_STATUS_OUT_OF_RESOURCES //!< Unable to acquire memory or other
	//!< resource
	RSMI_STATUS_INTERNAL_EXCEPTION  RSMIStatus = C.RSMI_STATUS_INTERNAL_EXCEPTION  //!< An internal exception was caught
	RSMI_STATUS_INPUT_OUT_OF_BOUNDS RSMIStatus = C.RSMI_STATUS_INPUT_OUT_OF_BOUNDS //!< The provided input is out of
	//!< allowable or safe range
	RSMI_STATUS_INIT_ERROR RSMIStatus = C.RSMI_STATUS_INIT_ERROR //!< An error occurred when rsmi
	//!< initializing internal data
	//!< structures
	RSMI_INITIALIZATION_ERROR       RSMIStatus = C.RSMI_INITIALIZATION_ERROR
	RSMI_STATUS_NOT_YET_IMPLEMENTED RSMIStatus = C.RSMI_STATUS_NOT_YET_IMPLEMENTED //!< The requested function has not
	//!< yet been implemented in the
	//!< current system for the current
	//!< devices
	RSMI_STATUS_NOT_FOUND RSMIStatus = C.RSMI_STATUS_NOT_FOUND //!< An item was searched for but not
	//!< found
	RSMI_STATUS_INSUFFICIENT_SIZE RSMIStatus = C.RSMI_STATUS_INSUFFICIENT_SIZE //!< Not enough resources were
	//!< available for the operation
	RSMI_STATUS_INTERRUPT RSMIStatus = C.RSMI_STATUS_INTERRUPT //!< An interrupt occurred during
	//!< execution of function
	RSMI_STATUS_UNEXPECTED_SIZE RSMIStatus = C.RSMI_STATUS_UNEXPECTED_SIZE //!< An unexpected amount of data
	//!< was read
	RSMI_STATUS_NO_DATA RSMIStatus = C.RSMI_STATUS_NO_DATA //!< No data was found for a given
	//!< input
	RSMI_STATUS_UNEXPECTED_DATA RSMIStatus = C.RSMI_STATUS_UNEXPECTED_DATA //!< The data read or provided to
	//!< function is not what was expected
	RSMI_STATUS_BUSY RSMIStatus = C.RSMI_STATUS_BUSY
	//!< A resource or mutex could not be
	//!< acquired because it is already
	//!< being used
	RSMI_STATUS_REFCOUNT_OVERFLOW RSMIStatus = C.RSMI_STATUS_REFCOUNT_OVERFLOW //!< An internal reference counter
	//!< exceeded INT32_MAX
	RSMI_STATUS_SETTING_UNAVAILABLE RSMIStatus = C.RSMI_STATUS_SETTING_UNAVAILABLE //!< Requested setting is unavailable
	//!< for the current device
	RSMI_STATUS_AMDGPU_RESTART_ERR RSMIStatus = C.RSMI_STATUS_AMDGPU_RESTART_ERR //!< Could not successfully restart
	//!< the amdgpu driver
	RSMI_STATUS_UNKNOWN_ERROR RSMIStatus = C.RSMI_STATUS_UNKNOWN_ERROR
)

func errorString(result C.rsmi_status_t) error {
	if RSMIStatus(result) == RSMI_STATUS_SUCCESS {
		return nil
	}
	var cStatusString *C.char
	statusCode := C.rsmi_status_string(result, &cStatusString)
	if RSMIStatus(statusCode) != RSMI_STATUS_SUCCESS {
		return fmt.Errorf("error: %v", statusCode)
	}
	goStatusString := C.GoString(cStatusString)
	return fmt.Errorf("%v", goStatusString)
}
