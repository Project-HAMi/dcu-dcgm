/******************************************************************************
 * Copyright 2016-2019 by SW Group, Chengdu Haiguang IC Design Co., Ltd.
 * All right reserved. See COPYRIGHT for detailed Information.
 *
 * @file        dmi_virtual.h
 *
 * @brief       Header file for hygon dcu device virtual interface.
 *
 * @author      Wang Yan<wangwy@hygon.cn>
 * @date        2022/09/30
 * @history     1.
 *
 * @modify      sheer-rey<xiarui@hygon.cn>
 *****************************************************************************/

#ifndef __INC_DMI_VIRTUAL_H__
#define __INC_DMI_VIRTUAL_H__

#include "dmi_error.h"
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    char name[DMI_NAME_SIZE];  // Device name.
    int compute_unit_count;    // Number of compute units in the device.
    size_t global_mem_size;    // Size of global memory region (in bytes).
    size_t usage_mem_size;     // Size of usage global memory (in bytes).
    uint64_t
        container_id;  // The container that is bound, only for virtual device.
    int device_id;  // Number of the physical card corresponding to the current
                    // device.
} dmiDeviceInfo;

/**
 *  @brief Set the log level.
 *
 *  @param[in] log_level the log level to set.
 *
 *  @return #DMI_STATUS_SUCCESS
 */
//dmiStatus dmiSetLogLevel(std::string log_level);

/**
 *  @brief Return number of devices.
 *
 *  @param[out] count written number of devices.
 *
 *  @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *          #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 */
dmiStatus dmiGetDeviceCount(int *count);

/**
 * @brief Return device infomation.
 *
 * @param [in]  device_id which device to query for information.
 * @param [out] device_info written with device infomation.
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_INVALID_ARGUMENTS
 */
dmiStatus dmiGetDeviceInfo(int device_id, dmiDeviceInfo *device_info);

/**
 * @brief Returns the maximum number of vdev per physical device supported.
 *
 * @param[out] count written number of virtual devices.
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 */
dmiStatus dmiGetMaxVDeviceCount(int *count);

/**
 * @brief Return the number of virtual devices for the specified
 *  physical device.
 *
 * @param[out] count written number of virtual devices.
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 */
dmiStatus dmiGetVDeviceCount(int *count);

/**
 * @brief Return remaining cus and memory size on the specified device for
 *        virtual device's further use.
 *
 * @param [in]  device_id which device to query for information.
 * @param [out] cus remaining compute unit counts.
 * @param [out] memories remaining memory size (in bytes).
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_INVALID_ARGUMENTS
 */
dmiStatus dmiGetDeviceRemainingInfo(int device_id, size_t *cus,
                                    size_t *memories);

/**
 * @brief Return virtual device information about the specified device.
 *
 * @param [in] vdevice_id which virtual device to query for information.
 * @param [out] device_info written with virtual device infomation.
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_INVALID_ARGUMENTS, #DMI_STATUS_NOT_SUPPORTED
 *         #DMI_STATUS_ERROR, #DMI_STATUS_VDEV_NOT_EXIST
 */
dmiStatus dmiGetVDeviceInfo(int vdevice_id, dmiDeviceInfo *device_info);

/**
 * @brief Initialize the current device for compute.
 *
 * @param[out] device_info written with device infomation.
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_INIT_DEVICE_FAILED, #DMI_STATUS_INVALID_ARGUMENTS
 */
dmiStatus dmiInitDevice(dmiDeviceInfo *device_info);

/**
 * @brief Create a specified number of virtual devices.
 *
 * @param [in] device_id which physical device to be created virtual device
 * @param [in] vdev_count number of virtual devices to be created
 * @param [in] vdev_cus vdev compute units array
 * @param [in] vdev_mem_size vdev memory size array
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_INVALID_ARGUMENTS, #DMI_STATUS_NOT_SUPPORTED
 *         #DMI_STATUS_DEVICE_NOT_SUPPORT, #DMI_STATUS_OUT_OF_RESOURCES
 *         #DMI_STATUS_ERROR
 */
dmiStatus dmiCreateVDevices(int device_id, int vdev_count, int *vdev_cus,
                            int *vdev_mem_size);

/**
 * @brief 销毁指定物理设备上的所有虚拟设备
 *
 * @param [in] deviceId physical device id
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_NOT_SUPPORTED, #DMI_STATUS_VDEV_NOT_EXIST
 */
dmiStatus dmiDestroyVDevices(int deviceId);

/**
 * @brief Destroy single virtual device.
 *
 * @param [in] vDeviceId virtual device id to be destroyed
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED,
 *         #DMI_STATUS_DEVICE_BUSY, #DMI_STATUS_MKFD_ALREADY_OPENED,
 *         #DMI_STATUS_SYS_NODE_NOT_EXIST, #DMI_STATUS_NOT_SUPPORTED,
 *         #DMI_STATUS_VDEV_NOT_EXIST
 */
dmiStatus dmiDestroySingleVDevice(int vDeviceId);

/**
 * @brief Destroy single virtual device.
 *
 * @param [in] vdeviceId virtual device id to be destroyed
 * @param [in] vdev_cus vdev compute units, -1 means not change
 * @param [in] vdev_mem_size vdev memory size (in MiBytes), -1 means not change
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED,
 *         #DMI_STATUS_DEVICE_BUSY, #DMI_STATUS_MKFD_ALREADY_OPENED,
 *         #DMI_STATUS_SYS_NODE_NOT_EXIST, #DMI_STATUS_NOT_SUPPORTED,
 *         #DMI_STATUS_VDEV_NOT_EXIST
 */
dmiStatus dmiUpdateSingleVDevice(int vdeviceId, int vdev_cus, int vdev_mem_size);

/**
 * @brief S启动虚拟设备.
 *
 * @param [in] deviceId virtual device id
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_INVALID_ARGUMENTS, #DMI_STATUS_NOT_SUPPORTED
 *         #DMI_STATUS_ERROR, #DMI_STATUS_VDEV_NOT_EXIST
 *
 * @warning while a virtual device has been started, it cannot be destroyed or
 *          updated until @c dmiStopVDevice has been invoked and returned
 *          successfully with the id of same virtual device.
 *
 * @note The @c dmiStartVDevice and @c dmiStopVDevice must be invoked in pairs
 *       within a single user process. If not, started virtual devices will be
 *       stopped automatically as a limited safeguard measure when user process
 *       has been normally exited. In extreme cases, like process has been
 *       killed unexpectedly, the behavior is undefined.
 */
dmiStatus dmiStartVDevice(int deviceId);

/**
 * @brief Stop virtual devices.
 *
 * @param [in] deviceId virtual device id
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_INVALID_ARGUMENTS, #DMI_STATUS_NOT_SUPPORTED
 *         #DMI_STATUS_ERROR, #DMI_STATUS_VDEV_NOT_EXIST
 *
 * @note The @c dmiStartVDevice and @c dmiStopVDevice must be invoked in pairs
 *       within a single user process. If not, started virtual devices will be
 *       stopped automatically as a limited safeguard measure when user process
 *       has been normally exited. In extreme cases, like process has been
 *       killed unexpectedly, the behavior is undefined.
 */
dmiStatus dmiStopVDevice(int deviceId);

/**
 * @brief Returns the device busy percent.
 *
 * @param[in] device_id which device to query.
 * @param[out] busy_percent written dev busy percent(0~100).
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_INVALID_ARGUMENTS, #DMI_STATUS_ERROR
 */
dmiStatus dmiGetDevBusyPercent(int device_id, int *busy_percent);

/**
 * @brief Returns the virtual device busy percent.
 *
 * @param[in] vdevice_id which device to query.
 * @param[out] busy_percent written dev busy percent(0~100).
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_INVALID_ARGUMENTS, #DMI_STATUS_ERROR
 *         #DMI_STATUS_VDEV_NOT_EXIST
 */
dmiStatus dmiGetVDevBusyPercent(int vdevice_id, int *busy_percent);

/**
 * @brief Set Encryption VM status.
 *
 * @param[in] status if status is true, enable encryption VM, else disable.
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_ERROR
 */
dmiStatus dmiSetEncryptionVMStatus(bool status);

/**
 * @brief Returns encryption VM status.
 *
 * @param[out] status the encryption VM status.
 *
 * @return #DMI_STATUS_SUCCESS, #DMI_STATUS_OPEN_MKFD_FAILED
 *         #DMI_STATUS_MKFD_ALREADY_OPENED, #DMI_STATUS_SYS_NODE_NOT_EXIST
 *         #DMI_STATUS_ERROR
 */
dmiStatus dmiGetEncryptionVMStatus(bool *status);

#ifdef __cplusplus
}  // extern "C"
#endif

#endif  // __INC_DMI_VIRTUAL_H__
