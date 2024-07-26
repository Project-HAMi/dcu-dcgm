/******************************************************************************
 * Copyright 2016-2019 by SW Group, Chengdu Haiguang IC Design Co., Ltd.
 * All right reserved. See COPYRIGHT for detailed Information.
 *
 * @file        dmi_error.h
 *
 * @brief       Header file for error of hygon dcu device virtual interface.
 *
 * @author      Liu Jiangting<liujiangting@hygon.cn>
 * @date        2023/04/03
 * @history     1.
 *
 * @modify      sheer-rey<xiarui@hygon.cn>
 *****************************************************************************/

#ifndef __INC_DMI_ERROR_H__
#define __INC_DMI_ERROR_H__

#include <stddef.h>
#include <stdint.h>
//#include <map>

#ifdef __cplusplus
extern "C" {
#endif

#define DMI_NAME_SIZE (256)

typedef enum _DMI_STATUS {
    // Operation successful
    DMI_STATUS_SUCCESS                = 0,
    // General error return if not otherwise specified
    DMI_STATUS_ERROR                  = 1,
    // Allocate memory error
    DMI_STATUS_NO_MEMORY              = 2,
    // Open driver mkfd failed
    DMI_STATUS_OPEN_MKFD_FAILED       = 3,
    // The driver mkfd has been opened
    DMI_STATUS_MKFD_ALREADY_OPENED    = 4,
    // The node does not exist
    DMI_STATUS_SYS_NODE_NOT_EXIST     = 5,
    // In current environment or device, the function is not support
    DMI_STATUS_NOT_SUPPORTED          = 6,
    // An unopened handle was accessed
    DMI_STATUS_MKFD_NOT_OPENED        = 7,
    // Create virtual device failed
    DMI_STATUS_CREATE_VDEV_FAILED     = 8,
    // Destroy virtual device failed
    DMI_STATUS_DESTROY_VDEV_FAILED    = 9,
    // Invalid args
    DMI_STATUS_INVALID_ARGUMENTS      = 10,
    // The required resources exceed the hardware limit
    DMI_STATUS_OUT_OF_RESOURCES       = 11,
    // Query virtual device information failed
    DMI_STATUS_QUERY_VDEV_INFO_FAILED = 12,
    // The device management runtime is not init
    DMI_STATUS_ERROR_NOT_INITIALIZED  = 13,
    // The current device is not support
    DMI_STATUS_DEVICE_NOT_SUPPORT     = 14,
    // The virtual device is not exist
    DMI_STATUS_VDEV_NOT_EXIST         = 15,
    // Init virtual device failed
    DMI_STATUS_INIT_DEVICE_FAILED     = 16,
    // Device busy
    DMI_STATUS_DEVICE_BUSY            = 17,
    // File read or write error
    DMI_STATUS_FILE_ERROR             = 18,
    // Permission denied
    DMI_STATUS_PERMISSION             = 19,
    // An internal exception was caught
    DMI_STATUS_INTERNAL_EXCEPTION     = 20,
    // The provided input is out of allowable or safe range
    DMI_STATUS_INPUT_OUT_OF_BOUNDS    = 21,
    // An error occurred when smi initializing internal data structures
    DMI_STATUS_SMI_INIT_ERROR         = 22,
    // An item was searched for but not found
    DMI_STATUS_NOT_FOUND              = 23,
    // Not enough resources were available for the operation
    DMI_STATUS_INSUFFICIENT_SIZE      = 24,
    // An interrupt occurred during execution of function
    DMI_STATUS_INTERRUPT              = 25,
    // An unexpected amount of data was read
    DMI_STATUS_UNEXPECTED_SIZE        = 26,
    // No data was found for a given input
    DMI_STATUS_NO_DATA                = 27,
    // The data read or provided to function is not what was expected
    DMI_STATUS_UNEXPECTED_DATA        = 28,
    // A resource or mutex could not be acquired because it is already being
    // used
    DMI_STATUS_SMI_BUSY               = 29,
    // An internal reference counter exceeded INT32_MAX
    DMI_STATUS_REFCOUNT_OVERFLOW      = 30,
    // The requested function has not yet been implemented in the current system
    // for the current devices
    DMI_STATUS_NOT_YET_IMPLEMENTED    = 31,
    // An unknown error occurred
    DMI_STATUS_UNKNOWN_ERROR          = 32
} dmiStatus;

/**
 * @brief Query additional information about a status code.
 *
 * @param[in] status Status code.
 * @param[out] status_string A NUL-terminated string that describes the error
 * status.
 *
 * @return ::DMI_STATUS_SUCCESS The function has been executed successfully.
 *         ::DMI_STATUS_ERROR The status is invalid.
 */
dmiStatus dmiGetStatusString(dmiStatus status, const char** status_string);

#ifdef __cplusplus
}  // extern "C"
#endif

#endif  // __INC_DMI_ERROR_H__
