/*******************************************************************************
 * Copyright 2016-2019 by SW Group, Chengdu Haiguang IC Design Co., Ltd.
 * All right reserved. See COPYRIGHT for detailed Information.
 *
 * @file        dmi_mig.h
 *
 * @brief       Header file for hygon multiple instance gpu interface.
 *
 * @author      sheer-rey<xiarui@hygon.cn>
 * @date        2023/10/25
 * @version     1.0
 *****************************************************************************/

#ifndef __INC_DMI_MIG_H__
#define __INC_DMI_MIG_H__

#ifdef __cplusplus
extern "C" {
#endif

///=============================================================================
/// @defgroup ↓ Return values for NVML API calls ↓
///=============================================================================
typedef enum nvmlReturn_enum {
    // cppcheck-suppress *
    //!< The operation was successful
    NVML_SUCCESS                         = 0,
    //!< NVML was not first initialized with nvmlInit()
    NVML_ERROR_UNINITIALIZED             = 1,
    //!< A supplied argument is invalid
    NVML_ERROR_INVALID_ARGUMENT          = 2,
    //!< The requested operation is not available on target device
    NVML_ERROR_NOT_SUPPORTED             = 3,
    //!< The current user does not have permission for operation
    NVML_ERROR_NO_PERMISSION             = 4,
    //!< Deprecated: Multiple initializations are now allowed through ref
    //!< counting
    NVML_ERROR_ALREADY_INITIALIZED       = 5,
    //!< A query to find an object was unsuccessful
    NVML_ERROR_NOT_FOUND                 = 6,
    //!< An input argument is not large enough
    NVML_ERROR_INSUFFICIENT_SIZE         = 7,
    //!< A device's external power cables are not properly attached
    NVML_ERROR_INSUFFICIENT_POWER        = 8,
    //!< NVIDIA driver is not loaded
    NVML_ERROR_DRIVER_NOT_LOADED         = 9,
    //!< User provided timeout passed
    NVML_ERROR_TIMEOUT                   = 10,
    //!< NVIDIA Kernel detected an interrupt issue with a GPU
    NVML_ERROR_IRQ_ISSUE                 = 11,
    //!< NVML Shared Library couldn't be found or loaded
    NVML_ERROR_LIBRARY_NOT_FOUND         = 12,
    //!< Local version of NVML doesn't implement this function
    NVML_ERROR_FUNCTION_NOT_FOUND        = 13,
    //!< infoROM is corrupted
    NVML_ERROR_CORRUPTED_INFOROM         = 14,
    //!< The GPU has fallen off the bus or has otherwise become inaccessible
    NVML_ERROR_GPU_IS_LOST               = 15,
    //!< The GPU requires a reset before it can be used again
    NVML_ERROR_RESET_REQUIRED            = 16,
    //!< The GPU control device has been blocked by the operating system/cgroups
    NVML_ERROR_OPERATING_SYSTEM          = 17,
    //!< RM detects a driver/library version mismatch
    NVML_ERROR_LIB_RM_VERSION_MISMATCH   = 18,
    //!< An operation cannot be performed because the GPU is currently in use
    NVML_ERROR_IN_USE                    = 19,
    //!< Insufficient memory
    NVML_ERROR_MEMORY                    = 20,
    //!< No data
    NVML_ERROR_NO_DATA                   = 21,
    //!< The requested vgpu operation is not available on target device, becasue
    //!< ECC is enabled
    NVML_ERROR_VGPU_ECC_NOT_SUPPORTED    = 22,
    //!< Ran out of critical resources, other than memory
    NVML_ERROR_INSUFFICIENT_RESOURCES    = 23,
    //!< Ran out of critical resources, other than memory
    NVML_ERROR_FREQ_NOT_SUPPORTED        = 24,
    //!< The provided version is invalid/unsupported
    NVML_ERROR_ARGUMENT_VERSION_MISMATCH = 25,
    //!< The requested functionality has been deprecated
    NVML_ERROR_DEPRECATED                = 26,
    //!< The system is not ready for the request
    NVML_ERROR_NOT_READY                 = 27,
    //!< An internal driver error occurred
    NVML_ERROR_UNKNOWN                   = 999
} nvmlReturn_t;

///=============================================================================
/// @defgroup ↓ Device Macro Definition ↓
///=============================================================================

#define NVML_DEVICE_PCI_BUS_ID_BUFFER_SIZE 32
#define NVML_DEVICE_SERIAL_BUFFER_SIZE     32
#define NVML_DEVICE_NAME_BUFFER_SIZE       256
#define NVML_DEVICE_UUID_BUFFER_SIZE       256

// Disable Multi Instance GPU mode.
#define NVML_DEVICE_MIG_DISABLE 0x0
// Enable Multi Instance GPU mode.
#define NVML_DEVICE_MIG_ENABLE  0x1

// GPU instance profiles.
// These macros should be passed to \ref nvmlDeviceGetGpuInstanceProfileInfo to
// retrieve the detailed information about a GPU instance such as profile ID,
// engine counts.
#define NVML_GPU_INSTANCE_PROFILE_1_SLICE 0x0
#define NVML_GPU_INSTANCE_PROFILE_2_SLICE 0x1
#define NVML_GPU_INSTANCE_PROFILE_3_SLICE 0x2
#define NVML_GPU_INSTANCE_PROFILE_4_SLICE 0x3
#define NVML_GPU_INSTANCE_PROFILE_COUNT   0x4

// Compute instance profiles.
// These macros should be passed to \ref
// nvmlGpuInstanceGetComputeInstanceProfileInfo to retrieve the detailed
// information about a compute instance such as profile ID, engine counts
#define NVML_COMPUTE_INSTANCE_PROFILE_1_SLICE 0x0
#define NVML_COMPUTE_INSTANCE_PROFILE_2_SLICE 0x1
#define NVML_COMPUTE_INSTANCE_PROFILE_3_SLICE 0x2
#define NVML_COMPUTE_INSTANCE_PROFILE_4_SLICE 0x3
#define NVML_COMPUTE_INSTANCE_PROFILE_COUNT   0x4

// Compute instance engine profiles.
#define NVML_COMPUTE_INSTANCE_ENGINE_PROFILE_SHARED 0x0
#define NVML_COMPUTE_INSTANCE_ENGINE_PROFILE_COUNT  0x1

///=============================================================================
/// @defgroup ↓ Global Type Definition ↓
///=============================================================================

// DMI device handle
typedef struct dmiDevice_struct *nvmlDevice_t;
// DMI gpu instance handle
typedef struct dmiGpuInstance_struct *nvmlGpuInstance_t;
// DMI compute instance handle
typedef struct dmiComputeInstance_struct *nvmlComputeInstance_t;

///=============================================================================
/// @defgroup ↓ Device Structs Definition ↓
///=============================================================================

typedef struct nvmlPciInfo_st {
    unsigned int pci_domain; // The PCI domain on which the device's bus resides
                             // 0 to 0xffffffff
    unsigned int pci_bus;    // The bus on which the device resides, 0 to 0xff
    unsigned int pci_device; // The device's id on the bus, 0 to 31
    unsigned int pci_function; // The device's function id, 0 to 7
    // PCI identifier <domain:bus:device.function>
    char bus_id[NVML_DEVICE_PCI_BUS_ID_BUFFER_SIZE];
} nvmlPciInfo;

typedef struct nvmlDeviceAttributes_st {
    unsigned int index;                      // the index of gpu or mig device
    unsigned int cu_count;                   // compute Unit count
    unsigned long long memory_size_MB;       // memory size in MBytes
    char uuid[NVML_DEVICE_UUID_BUFFER_SIZE]; // gpu or mig device's uuid string
    char name[NVML_DEVICE_NAME_BUFFER_SIZE]; // gpu or mig device's name

    // gpu slice count in parent gpu instance, for mig device valid only
    // for gpu device, it always equal to 0
    unsigned int gpu_instance_slice_count;
    // gpu slice count in parent compute instance, for mig device valid only
    // for gpu device, it always equal to 0
    unsigned int compute_instance_slice_count;

    // unsigned int sharedCopyEngineCount;     //!< Shared Copy Engine count
    // unsigned int sharedDecoderCount;        //!< Shared Decoder Engine count
    // unsigned int sharedEncoderCount;        //!< Shared Encoder Engine count
    // unsigned int sharedJpegCount;           //!< Shared JPEG Engine count
    // unsigned int sharedOfaCount;            //!< Shared OFA Engine count
} nvmlDeviceAttributes_t;

///=============================================================================
/// @defgroup ↓ Mig Structs Definition ↓
///=============================================================================

typedef struct nvmlGpuInstanceProfileInfo_st {
    unsigned int id;                   // unique profile ID within the device
    unsigned int gi_count_max;         // max supported gpu instance count
    unsigned int cu_count;             // compute unit count
    unsigned int gpu_slice_count;      // GPU Slice count the profile contain
    unsigned long long memory_size_MB; // memory size in MBytes
    char name[NVML_DEVICE_NAME_BUFFER_SIZE];

    // unsigned int jpegCount;       // JPEG Engine count.
    // unsigned int isP2pSupported;  // Peer-to-Peer support.
    // unsigned int copyEngineCount; // Copy Engine count.
    // unsigned int ofaCount;        // OFA Engine count.
} nvmlGpuInstanceProfileInfo_t;

typedef struct nvmlComputeInstanceProfileInfo_st {
    unsigned int id;              // unique profile ID within the GPU instance
    unsigned int ci_count_max;    // max supported compute instance count
    unsigned int cu_count;        // compute unit count
    unsigned int gpu_slice_count; // GPU Slice count the profile contain
    char name[NVML_DEVICE_NAME_BUFFER_SIZE];

    // unsigned int sharedCopyEngineCount; // Shared Copy Engine count.
    // unsigned int sharedDecoderCount;    // Shared Decoder Engine count.
    // unsigned int sharedEncoderCount;    // Shared Encoder Engine count.
    // unsigned int sharedJpegCount;       // Shared JPEG Engine count.
    // unsigned int sharedOfaCount;        // Shared OFA Engine count.
} nvmlComputeInstanceProfileInfo_t;

typedef struct nvmlGpuInstancePlacement_st {
    unsigned int start; // Index of first occupied memory slice
    unsigned int size;  // Number of memory slices occupied
} nvmlGpuInstancePlacement_t;

typedef struct nvmlComputeInstancePlacement_st {
    unsigned int start; // Index of first occupied compute slice
    unsigned int size;  // Number of compute slices occupied
} nvmlComputeInstancePlacement_t;

typedef struct nvmlGpuInstanceInfo_st {
    nvmlDevice_t device;     // Parent device handle
    unsigned int id;         // Unique instance ID within the device
    unsigned int profile_id; // Unique profile ID within the device
    nvmlGpuInstancePlacement_t placement; // Placement for this instance
} nvmlGpuInstanceInfo_t;

typedef struct nvmlComputeInstanceInfo_st {
    nvmlDevice_t device;            // Parent device handle
    nvmlGpuInstance_t gpu_instance; // Parent GPU instance handle
    unsigned int id;         // Unique instance ID within the GPU instance
    unsigned int profile_id; // Unique profile ID within the GPU instance

    // Placement for this instance within the GPU instance's compute slice range
    // {0, slice_count}
    nvmlComputeInstancePlacement_t placement;
} nvmlComputeInstanceInfo_t;

///=============================================================================
/// @defgroup ↓ Device Functions Definition ↓
///=============================================================================

/*******************************************************************************
 * @name nvmlDeviceGetCount
 *
 * @brief Retrieves the number of compute devices in the system. A compute
 *        device is a single GPU.
 *
 * @param[out] device_count Reference in which to return the number of
 *             accessible devices
 *
 * @return NVML_SUCCESS
 *         -> If \p device_count has been set successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device_count is invalid \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetCount(unsigned int *device_count);

/*******************************************************************************
 * @name nvmlDeviceGetHandleByIndex
 *
 * @brief Acquire the handle for a particular device, based on its index.
 *
 * @note Valid indices are derived from the accessible devices count returned by
 *       \p nvmlDeviceGetCount. For example, if accessible devices is 2, then
 *       the valid indices are 0 and 1, corresponding to GPU 0 and GPU 1.
 * @note The order in which DMI enumerates devices has no guarantees of
 *       consistency between reboots. For that reason it is recommended that
 *       devices be looked up by their PCI ids or UUID.
 *       See \p nvmlDeviceGetHandleByUUID and \p nvmlDeviceGetHandleByPciBusId
 * @note The DMI index may not correlate with other APIs, such as the HIP
 *       device index.
 *
 * @param[in]  index  The index of the target GPU, >= 0 and < accessible devices
 * @param[out] device Reference in which to return the device handle
 *
 * @return NVML_SUCCESS
 *         -> If device handle has been set successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p index or \p device are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If this query is not supported by the device \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 *
 * @see nvmlDeviceGetIndex
 * @see nvmlDeviceGetCount
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetHandleByIndex(unsigned int index,
                                        nvmlDevice_t *device);

/*******************************************************************************
 * @name nvmlDeviceGetHandleBySerial
 *
 * @brief Acquire the handle for a particular device, based on its serial number
 *
 * @note This number corresponds to the value to the value returned by
 *       \p nvmlDeviceGetSerial
 *
 * @param[in]  serial The serial number of the target GPU
 * @param[out] device Reference in which to return the device handle
 *
 * @return NVML_SUCCESS
 *         -> If device handle has been set successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p serial or \p device are invalid \
 * @return NVML_ERROR_NOT_FOUND
 *         -> If \p serial does not match a valid device on the system \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If this query is not supported by the device \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 *
 * @see nvmlDeviceGetSerial
 * @see nvmlDeviceGetHandleByUUID
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetHandleBySerial(const char *serial,
                                         nvmlDevice_t *device);

/*******************************************************************************
 * @name nvmlDeviceGetHandleByUUID
 *
 * @brief Acquire the handle for a particular device, based on its globally
 *        unique immutable UUID associated with each device.
 *
 * @param[in]  uuid   The UUID of the target GPU or MIG instance
 * @param[out] device Reference in which to return the device handle or MIG
 *                    device handle
 *
 * @return NVML_SUCCESS
 *         -> If device handle has been set successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p uuid or \p device are invalid \
 * @return NVML_ERROR_NOT_FOUND
 *         -> If \p uuid does not match a valid device on the system \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If this query is not supported by the device \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 *
 * @see nvmlDeviceGetUUID
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetHandleByUUID(const char *uuid, nvmlDevice_t *device);

/*******************************************************************************
 * @name nvmlDeviceGetHandleByPciBusId
 *
 * @brief Acquire the handle for a particular device, based on its PCI bus id.
 *
 * @note This value corresponds to the nvmlPciInfo_st::bus_id returned by
 *       \p nvmlDeviceGetPciInfo
 *
 * @param[in]  pci_bus_id The PCI bus id of the target GPU
 * @param[out] device     Reference in which to return the device handle
 *
 * @return NVML_SUCCESS
 *         -> If device handle has been set successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p pci_bus_id or \p device are invalid \
 * @return NVML_ERROR_NOT_FOUND
 *         -> If \p pci_bus_id does not match a valid device on the system \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If this query is not supported by the device \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetHandleByPciBusId(const char *pci_bus_id,
                                           nvmlDevice_t *device);

/*******************************************************************************
 * @name nvmlDeviceGetName
 *
 * @brief Retrieves the name of this device.
 *
 * @note The name is an alphanumeric string that denotes a particular product,
 *       e.g. Z100L. It will not exceed \p NVML_DEVICE_NAME_BUFFER_SIZE
 *       characters in length (including the NULL terminator).
 * @note When used with MIG device handles the API returns MIG device names
 *       which can be used to identify devices based on their attributes.
 *
 * @param[in]  device The identifier of the target gpu or mig device
 * @param[out] name   Reference in which to return the product name
 * @param[in]  length The maximum allowed length of the string returned in
 *                    \p name
 *
 * @return NVML_SUCCESS
 *         -> If device name has been set successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p name or \p length are invalid \
 * @return NVML_ERROR_INSUFFICIENT_SIZE
 *         -> If \p lenth is too small \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If this query is not supported by the device \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetName(nvmlDevice_t device, char *name,
                               unsigned int length);

/*******************************************************************************
 * @name nvmlDeviceGetIndex
 *
 * @brief Retrieves the DMI index of gpu or mig device.
 *
 * @note Valid indices are derived from the accessible devices count returned by
 *       \p nvmlDeviceGetCount. For example, if accessible devices is 2, then
 *the valid indices are 0 and 1, corresponding to GPU 0 and GPU 1.
 * @note The order in which DMI enumerates devices has no guarantees of
 *       consistency between reboots. For that reason it is recommended that
 *       devices be looked up by their PCI ids or UUID. See
 *       \p nvmlDeviceGetHandleByPciBusId and \p nvmlDeviceGetHandleByUUID
 * @note When used with MIG device handles this API returns indices that can be
 *       passed to \p nvmlDeviceGetMigDeviceHandleByIndex to retrieve an
 *       identical handle. MIG device indices are unique within a physical gpu
 *       device.
 * @note The DMI index may not correlate with other APIs, such as the HIP device
 *       index.
 *
 * @param[in]  device The identifier of the target gpu or mig device
 * @param[out] index  Reference in which to return the DMI index of the device
 *
 * @return NVML_SUCCESS
 *         -> If device index has been set successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device or \p index are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If this query is not supported by the device \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 *
 * @see nvmlDeviceGetHandleByIndex
 * @see nvmlDeviceGetCount
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetIndex(nvmlDevice_t device, unsigned int *index);

/*******************************************************************************
 * @name nvmlDeviceGetSerial
 *
 * @brief Retrieves the globally unique board serial number associated with the
 *        device.
 *
 * @note The serial number is an alphanumeric string that will not exceed
 *       \p NVML_DEVICE_SERIAL_BUFFER_SIZE characters (including the NULL
 *       terminator).
 *
 * @param[in]  device The identifier of the target device
 * @param[out] serial Reference in which to return the serial number
 * @param[in]  length The maximum allowed length of the string returned in
 *                    \p serial
 *
 * @return NVML_SUCCESS
 *         -> If device serial has been set successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p serial or \p length are invalid \
 * @return NVML_ERROR_INSUFFICIENT_SIZE
 *         -> If \p lenth is too small \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If this query is not supported by the device \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetSerial(nvmlDevice_t device, char *serial,
                                 unsigned int length);

/*******************************************************************************
 * @name nvmlDeviceGetPciInfo
 *
 * @brief Retrieves the PCI attributes of this device.
 *
 * @note See \p nvmlPciInfo_st for details on the available PCI info.
 *
 * @param[in]  device The identifier of the target device
 * @param[out] pci    Reference in which to return the PCI info
 *
 * @return NVML_SUCCESS
 *         -> If device pci has been set successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device or \p pci are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If this query is not supported by the device \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetPciInfo(nvmlDevice_t device, nvmlPciInfo *pci);

/*******************************************************************************
 * @name nvmlDeviceGetAttributes
 *
 * @brief Get attributes (compute unit counts etc.) for the given DMI device
 *        handle.
 *
 * @note This API currently supports GPU and MIG device handles.
 *
 * @param[in]  device     DMI gpu or mig device handle
 * @param[out] attributes Device attributes
 *
 * @return NVML_SUCCESS
 *         -> If device attributes has been set successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device or \p attributes are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> if this query is not supported by the device \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetAttributes(nvmlDevice_t device,
                                     nvmlDeviceAttributes_t *attributes);

///=============================================================================
/// @defgroup ↓ Mig Functions Definition ↓
///=============================================================================

/*******************************************************************************
 * @name nvmlDeviceSetMigMode
 *
 * @brief Set MIG mode for the device. This mode determines whether a GPU
 *        instance can be created.
 *
 * @note This API may unbind or reset the device to activate the requested mode.
 *       Thus, the attributes associated with the device, such as minor number,
 *       might change. The caller of this API is expected to query such
 *       attributes again.
 * @note On certain platforms like pass-through virtualization, where reset
 *       functionality may not be exposed directly, VM reboot is required. Then
 *       \p activation_status would return \p NVML_ERROR_RESET_REQUIRED for such
 *       cases.
 * @note \p activation_status would return the appropriate error code upon
 *       unsuccessful activation. For example, if device unbind fails because
 *       the device isn't idle, \p NVML_ERROR_IN_USE would be returned. The
 *       caller of this API is expected to idle the device and retry setting the
 *       \p mode.
 *
 * @param[in]  device            The identifier of the target device
 * @param[in]  mode              The mode to be set,
 *                               \p NVML_DEVICE_MIG_ENABLE or
 *                               \p NVML_DEVICE_MIG_DISABLE
 * @param[out] activation_status The activationStatus status
 *
 * @return NVML_SUCCESS
 *         -> If MIG activation status has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p mode or \p activation_status are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If \p device doesn't support MIG mode \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceSetMigMode(nvmlDevice_t device, unsigned int mode,
                                  nvmlReturn_t *activation_status);

/*******************************************************************************
 * @name nvmlDeviceGetMigMode
 *
 * @brief Get MIG mode for the device.
 *
 * @note Change MIG modes may require device unbind or reset. The "pending" MIG
 *       mode refers to the target mode following the next activation trigger.
 *
 * @param[in]  device       The identifier of the target device
 * @param[out] current_mode Returns the current mode,
 *                          \p NVML_DEVICE_MIG_ENABLE or
 *                          \p NVML_DEVICE_MIG_DISABLE
 * @param[out] pending_mode Returns the pending mode,
 *                          \p NVML_DEVICE_MIG_ENABLE,
 *                          \p NVML_DEVICE_MIG_DISABLE
 *
 * @return NVML_SUCCESS
 *         -> If MIG mode has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p currentMode or \p pendingMode are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If \p device doesn't support MIG mode \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetMigMode(nvmlDevice_t device,
                                  unsigned int *current_mode,
                                  unsigned int *pending_mode);

/*******************************************************************************
 * @name nvmlDeviceGetGpuInstanceProfileInfo
 *
 * @brief Get GPU instance profile information.
 *
 * @note Information provided by this API is immutable throughout the lifetime
 *       of a MIG mode.
 *
 * @param[in]  device  The identifier of the target device
 * @param[in]  profile One of the NVML_GPU_INSTANCE_PROFILE_*
 * @param[out] info    Returns detailed profile information of \p profile
 *
 * @return NVML_SUCCESS
 *         -> If \p info has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p profile or \p info are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If \p device doesn't have MIG mode enabled or
 *            \p profile isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetGpuInstanceProfileInfo(
    nvmlDevice_t device, unsigned int profile,
    nvmlGpuInstanceProfileInfo_t *info);

/*******************************************************************************
 * @name nvmlDeviceGetGpuInstancePossiblePlacements
 *
 * @brief Get GPU instance placements.
 *
 * @note A placement represents the location of a GPU instance within a device.
 *       This API only returns all the possible placements for the given
 *       profile. A created GPU instance occupies memory slices described by its
 *       placement. Creation of new GPU instance will fail if there is overlap
 *       with the already occupied memory slices.
 *
 * @param[in]  device     The identifier of the target device
 * @param[in]  profile_id The GPU instance profile ID.
 *                        See \p nvmlDeviceGetGpuInstanceProfileInfo
 * @param[out] placements Returns placements allowed for the profile. Can be
 *                        \b NULL to discover number of allowed placements for
 *                        this profile. If \b non-NULL must be large enough to
 *                        accommodate the placements supported by the profile.
 * @param[out] count      Returns number of allowed placemenets for the profile.
 *
 * @return NVML_SUCCESS
 *         -> If \p placements and \p count has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p profile or \p count are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If \p device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetGpuInstancePossiblePlacements(
    nvmlDevice_t device, unsigned int profile_id,
    nvmlGpuInstancePlacement_t *placements, unsigned int *count);

/*******************************************************************************
 * @name nvmlDeviceGetGpuInstanceRemainingCapacity
 *
 * @brief Get GPU instance profile capacity.
 *
 * @param[in]  device     The identifier of the target device
 * @param[in]  profile_id    The GPU instance profile ID.
 *                        See \p nvmlDeviceGetGpuInstanceProfileInfo
 * @param[out] count      Returns number of allowed placemenets for the profile.
 *
 * @return NVML_SUCCESS
 *         -> If \p count has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p profile_id or \p count are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If \p device doesn't have MIG mode enabled or
 *            \p profile_id isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetGpuInstanceRemainingCapacity(nvmlDevice_t device,
                                                       unsigned int profile_id,
                                                       unsigned int *count);

/*******************************************************************************
 * @name nvmlDeviceCreateGpuInstance
 *
 * @brief Create GPU instance.
 *
 * @note If the parent device is unbound, reset or the GPU instance is destroyed
 *       explicitly, the GPU instance handle would become invalid. The GPU
 *       instance must be recreated to acquire a valid handle.
 *
 * @param[in]  device       The identifier of the target device
 * @param[in]  profile_id   The GPU instance profile ID.
 *                          See \p nvmlDeviceGetGpuInstanceProfileInfo
 * @param[out] gpu_instance Returns the GPU instance handle
 *
 * @return NVML_SUCCESS
 *         -> If gpu instance has been created successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p profile_id or \p gpu_instance are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If \p device doesn't have MIG mode enabled or
 *            \p profile_id isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_INSUFFICIENT_RESOURCES
 *        -> If the requested GPU instance could not be created \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceCreateGpuInstance(nvmlDevice_t device,
                                         unsigned int profile_id,
                                         nvmlGpuInstance_t *gpu_instance);

/*******************************************************************************
 * @name nvmlDeviceCreateGpuInstanceWithPlacement
 *
 * @brief Create GPU instance with the specified placement.
 *
 * @note If the parent device is unbound, reset or the GPU instance is destroyed
 *       explicitly, the GPU instance handle would become invalid. The GPU
 *       instance must be recreated to acquire a valid handle.
 *
 * @param[in]  device       The identifier of the target device
 * @param[in]  profile_id   The GPU instance profile ID.
 *                          See \p nvmlDeviceGetGpuInstanceProfileInfo
 * @param[in]  placement    The requested placement.
 *                          See \p nvmlDeviceGetGpuInstancePossiblePlacements
 * @param[out] gpu_instance Returns the GPU instance handle
 *
 * @return NVML_SUCCESS
 *         -> If gpu instance has been created successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p profile_id, \p placement
 *            or \p gpu_instance are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If \p device doesn't have MIG mode enabled or
 *            \p profile_id isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_INSUFFICIENT_RESOURCES
 *        -> If the requested GPU instance could not be created \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceCreateGpuInstanceWithPlacement(
    nvmlDevice_t device, unsigned int profile_id,
    const nvmlGpuInstancePlacement_t *placement,
    nvmlGpuInstance_t *gpu_instance);

/*******************************************************************************
 * @name nvmlGpuInstanceDestroy
 *
 * @brief Destroy GPU instance.
 *
 * @param[in] gpu_instance The GPU instance handle
 *
 * @return NVML_SUCCESS
 *         -> If gpu instance has been destroied successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p gpu_instance are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If parent gpu device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_IN_USE
 *         -> If the GPU instance is in use. This error would be returned if
 *            processes (e.g. HIP application) or compute instances are active
 *            on the GPU instance. \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlGpuInstanceDestroy(nvmlGpuInstance_t gpu_instance);

/*******************************************************************************
 * @name nvmlDeviceGetGpuInstances
 *
 * @brief Get GPU instances for given profile ID.
 *
 * @param[in]  device        The identifier of the target device
 * @param[in]  profile_id    The GPU instance profile ID.
 *                           See \p nvmlDeviceGetGpuInstanceProfileInfo
 * @param[out] gpu_instances Returns pre-exiting GPU instance handles, the
 *                           buffer must be large enough to accommodate the
 *                           instances supported by the profile. See
 *                           \p nvmlDeviceGetGpuInstanceProfileInfo
 * @param[out] count         The count of returned GPU instance handles
 *
 * @return NVML_SUCCESS
 *         -> If \p gpu_instances and \p count has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p profile_id, \p gpu_instances or \p count
 *            are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If \p device doesn't have MIG mode enabled
 *            or \p profile_id isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetGpuInstances(nvmlDevice_t device,
                                       unsigned int profile_id,
                                       nvmlGpuInstance_t *gpu_instances,
                                       unsigned int *count);

/*******************************************************************************
 * @name nvmlDeviceGetGpuInstanceById
 *
 * @brief Get GPU instances for given instance ID.
 *
 * @param[in]  device          The identifier of the target device
 * @param[in]  gpu_instance_id The GPU instance ID
 * @param[out] gpu_instance    Returns GPU instance handle
 *
 * @return NVML_SUCCESS
 *         -> If gpu instance handle has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p gpu_instance_id, or \p gpu_instance are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If \p device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_NOT_FOUND
 *         -> If the GPU instance is not found \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetGpuInstanceById(nvmlDevice_t device,
                                          unsigned int gpu_instance_id,
                                          nvmlGpuInstance_t *gpu_instance);

/*******************************************************************************
 * @name nvmlGpuInstanceGetInfo
 *
 * @brief Get GPU instance information.
 *
 * @param[in]  gpu_instance The GPU instance handle
 * @param[out] info         Return GPU instance information
 *
 * @return NVML_SUCCESS
 *         -> If \p info has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p gpu_instance or \p info are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlGpuInstanceGetInfo(nvmlGpuInstance_t gpu_instance,
                                    nvmlGpuInstanceInfo_t *info);

/*******************************************************************************
 * @name nvmlGpuInstanceGetComputeInstanceProfileInfo
 *
 * @brief Get compute instance profile information.
 *
 * @note Information provided by this API is immutable throughout the lifetime
 *       of a MIG mode.
 *
 * @param[in]  gpu_instance The handle of the target GPU instance
 * @param[in]  profile      One of the \p NVML_COMPUTE_INSTANCE_PROFILE_*
 * @param[in]  eng_profile  One of the \p NVML_COMPUTE_INSTANCE_ENGINE_PROFILE_*
 * @param[out] info         Returns detailed profile information
 *
 * @return NVML_SUCCESS
 *         -> If \p info has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p gpu_instance, \p profile, \p eng_profile
 *            or \p gpu_instance are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If \p device doesn't have MIG mode enabled
 *            or \p profile isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlGpuInstanceGetComputeInstanceProfileInfo(
    nvmlGpuInstance_t gpu_instance, unsigned int profile,
    unsigned int eng_profile, nvmlComputeInstanceProfileInfo_t *info);

/*******************************************************************************
 * @name nvmlGpuInstanceGetComputeInstancePossiblePlacements
 *
 * @brief Get compute instance placements.
 *
 * @note A placement represents the location of a compute instance within a GPU
 *       instance. This API only returns all the possible placements for the
 *       given profile. A created compute instance occupies compute slices
 *       described by its placement. Creation of new compute instance will fail
 *       if there is overlap with the already occupied compute slices.
 *
 * @param[in]  gpu_instance The handle of the target GPU instance
 * @param[in]  profile_id   The compute instance profile ID.
 *                          See \p nvmlGpuInstanceGetComputeInstanceProfileInfo
 * @param[out] placements   Returns placements allowed for the profile. Can be
 *                          \b NULL to discover number of allowed placements for
 *                          this profile. If \b non-NULL must be large enough to
 *                          accommodate the placements supported by the profile.
 * @param[out] count        Returns number of allowed placemenets for the
 *                          profile.
 *
 * @return NVML_SUCCESS
 *         -> If \p placements and \p count has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p gpu_instance, \p profile_id, \p placements
 *            or \p count are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled
 *            or \p profile_id isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlGpuInstanceGetComputeInstancePossiblePlacements(
    nvmlGpuInstance_t gpu_instance, unsigned int profile_id,
    nvmlComputeInstancePlacement_t *placements, unsigned int *count);

/*******************************************************************************
 * @name nvmlGpuInstanceGetComputeInstanceRemainingCapacity
 *
 * @brief Get compute instance profile capacity.
 *
 * @param[in]  gpu_instance The handle of the target GPU instance
 * @param[in]  profile_id      The compute instance profile ID.
 *                          See \p nvmlGpuInstanceGetComputeInstanceProfileInfo
 * @param[out] count        Returns remaining instance count for the profile ID
 *
 * @return NVML_SUCCESS
 *         -> If \p count has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p gpu_instance, \p profile_id or \p count are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled
 *            or \p profile_id isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlGpuInstanceGetComputeInstanceRemainingCapacity(
    nvmlGpuInstance_t gpu_instance, unsigned int profile_id,
    unsigned int *count);

/*******************************************************************************
 * @name nvmlGpuInstanceCreateComputeInstance
 *
 * @brief Create compute instance
 *
 * @note If the parent device is unbound, reset or the parent GPU instance is
 *       destroyed or the compute instance is destroyed explicitly, the compute
 *       instance handle would become invalid. The compute instance must be
 *       recreated to acquire a valid handle.
 *
 * @param[in]  gpu_instance     The identifier of the target GPU instance
 * @param[in]  profile_id       The compute instance profile ID. See
 *                              \p nvmlGpuInstanceGetComputeInstanceProfileInfo
 * @param[out] compute_instance Returns the compute instance handle
 *
 * @return NVML_SUCCESS
 *         -> If compute instance has been successfully created \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p gpu_instance, \p profile_id,
 *            or \p compute_instance are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled
 *            or \p profile_id isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_INSUFFICIENT_RESOURCES
 *         -> If the requested compute instance could not be created \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlGpuInstanceCreateComputeInstance(
    nvmlGpuInstance_t gpu_instance, unsigned int profile_id,
    nvmlComputeInstance_t *compute_instance);

/*******************************************************************************
 * @name nvmlGpuInstanceCreateComputeInstanceWithPlacement
 *
 * @brief Create compute instance with the specified placement
 *
 * @note If the parent device is unbound, reset or the parent GPU instance is
 *       destroyed or the compute instance is destroyed explicitly, the compute
 *       instance handle would become invalid. The compute instance must be
 *       recreated to acquire a valid handle.
 *
 * @param[in]  gpu_instance     The handle of the target GPU instance
 * @param[in]  profile_id       The compute instance profile ID. See
 *                         \p nvmlGpuInstanceGetComputeInstanceProfileInfo
 * @param[in]  placement        The requested placement
 *                         \p
 *nvmlGpuInstanceGetComputeInstancePossiblePlacements
 * @param[out] compute_instance Returns the compute instance handle
 *
 * @return NVML_SUCCESS
 *         -> If compute instance has been created successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p gpu_instance, \p profile_id,
 *            or \p compute_instance are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled
 *            or \p profile_id isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_INSUFFICIENT_RESOURCES
 *         -> If the requested compute instance could not be created \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlGpuInstanceCreateComputeInstanceWithPlacement(
    nvmlGpuInstance_t gpu_instance, unsigned int profile_id,
    const nvmlComputeInstancePlacement_t *placement,
    nvmlComputeInstance_t *compute_instance);

/*******************************************************************************
 * @name nvmlComputeInstanceDestroy
 *
 * @brief Destroy compute instance
 *
 * @param[in] compute_instance The compute instance handle
 *
 * @return NVML_SUCCESS
 *         -> If compute instance has been destroied successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p compute_instance are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_IN_USE
 *         -> If the compute instance is in use. This error would be returned if
 *            processes (e.g. HIP application) are active on compute instance \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlComputeInstanceDestroy(nvmlComputeInstance_t compute_instance);

/*******************************************************************************
 * @name nvmlGpuInstanceGetComputeInstances
 *
 * @brief Get compute instances for given profile ID
 *
 * @param[in]  gpu_instance      The handle of the target GPU instance
 * @param[in]  profile           The compute instance profile ID. See
 *                               \p nvmlGpuInstanceGetComputeInstanceProfileInfo
 * @param[out] compute_instances Returns pre-exiting compute instance handles,
 *                               the buffer must be large enough to accommodate
 *                               the instances supported by the profile. See
 *                               \p nvmlGpuInstanceGetComputeInstanceProfileInfo
 * @param[out] count             The count of returned compute instance handles
 *
 * @return NVML_SUCCESS
 *         -> If \p compute_instances and \p count was successfully retrieved \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p gpu_instance, \p profile, \p compute_instance reference
 *            or \p count reference are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled
 *            or \p profile isn't supported \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlGpuInstanceGetComputeInstances(
    nvmlGpuInstance_t gpu_instance, unsigned int profile_id,
    nvmlComputeInstance_t *compute_instances, unsigned int *count);

/*******************************************************************************
 * @name nvmlGpuInstanceGetComputeInstanceById
 *
 * @brief Get compute instance for given instance ID
 *
 * @param[in]  gpu_instance        The handle of the target GPU instance
 * @param[in]  compute_instance_id The compute instance ID
 * @param[out] compute_instance    Returns compute instance handle
 *
 * @return NVML_SUCCESS
 *         -> if \p compute_instance was successfully retrieved \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p gpu_instance, \p compute_instance_id
 *            or \p compute_instance reference are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_NOT_FOUND
 *         -> If the compute instance is not found \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlGpuInstanceGetComputeInstanceById(
    nvmlGpuInstance_t gpu_instance, unsigned int id,
    nvmlComputeInstance_t *compute_instance);

/*******************************************************************************
 * @name nvmlComputeInstanceGetInfo
 *
 * @brief Get compute instance information
 *
 * @param[in]  compute_instance The compute instance handle
 * @param[out] info             Return compute instance information
 *
 * @return NVML_SUCCESS
 *         -> if \p info was successfully retrieved \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p compute_instance or \p info reference are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlComputeInstanceGetInfo(nvmlComputeInstance_t compute_instance,
                                        nvmlComputeInstanceInfo_t *info);

/*******************************************************************************
 * @name nvmlDeviceIsMigDeviceHandle
 *
 * @brief Test if the given handle refers to a MIG device.
 *
 * @note A MIG device handle is an DMI abstraction which maps to a MIG compute
 *       instance. These overloaded references can be used (with some
 *       restrictions) interchangeably with a GPU device handle to execute
 *       queries at a per-compute instance granularity.
 *
 * @param[in]  device        DMI device handle to test
 * @param[out] is_mig_device True when handle refers to a MIG device
 *
 * @return NVML_SUCCESS
 *         -> if \p device status was successfully retrieved \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device or \p is_mig_device reference are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceIsMigDeviceHandle(nvmlDevice_t device,
                                         unsigned int *is_mig_device);

/*******************************************************************************
 * @name nvmlDeviceGetGpuInstanceId
 *
 * @brief Get GPU instance ID for the given MIG device handle.
 *
 * @note GPU instance IDs are unique per device and remain valid until the GPU
 *       instance is destroyed.
 *
 * @param[in]  device          Target MIG device handle
 * @param[out] gpu_instance_id GPU instance ID
 *
 * @return NVML_SUCCESS
 *         -> if \p gpu_instance_id was successfully retrieved \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device or \p gpu_instance_id reference are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetGpuInstanceId(nvmlDevice_t device,
                                        unsigned int *gpu_instance_id);

/*******************************************************************************
 * @name nvmlDeviceGetComputeInstanceId
 *
 * @brief Get compute instance ID for the given MIG device handle.
 *
 * @note Compute instance IDs are unique per GPU instance and remain valid until
 *       the compute instance is destroyed.
 *
 * @param[in]  device              Target MIG device handle
 * @param[out] compute_instance_id Compute instance ID
 *
 * @return NVML_SUCCESS
 *         -> if \p compute_instance_id was successfully retrieved \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device or \p compute_instance_id reference are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetComputeInstanceId(nvmlDevice_t device,
                                            unsigned int *compute_instance_id);

/*******************************************************************************
 * @name nvmlDeviceGetMaxMigDeviceCount
 *
 * @brief Get the maximum number of MIG devices that can exist under a given
 *        parent DMI device.
 *
 * @param[in]  device Target device handle
 * @param[out] count  Count of MIG devices, returns zero if MIG is not supported
 *                    or enabled.
 *
 * @return NVML_SUCCESS
 *         -> if \p count was successfully retrieved \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device or \p count reference are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetMaxMigDeviceCount(nvmlDevice_t device,
                                            unsigned int *count);

/*******************************************************************************
 * @name nvmlDeviceGetMigDeviceHandleByIndex
 *
 * @brief Get MIG device handle for the given index under its parent DMI device
 *
 * @note If the compute instance is destroyed either explicitly or by
 *       destroying, resetting or unbinding the parent GPU instance or the GPU
 *       device itself the MIG device handle would remain invalid and must be
 *       requested again using this API. Handles may be reused and their
 *       properties can change in the process.
 *
 * @param[in]  device     Reference to the parent GPU device handle
 * @param[in]  index      Index of the MIG device
 * @param[out] mig_device Reference to the MIG device handle
 *
 * @return NVML_SUCCESS
 *         -> if \p mig_device handle was successfully created \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p device, \p index or \p mig_device reference are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_NOT_FOUND
 *         -> if no valid MIG device was found at \p index \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetMigDeviceHandleByIndex(nvmlDevice_t device,
                                                 unsigned int index,
                                                 nvmlDevice_t *mig_device);

/*******************************************************************************
 * @name nvmlDeviceGetDeviceHandleFromMigDeviceHandle
 *
 * @brief Get parent device handle from a MIG device handle
 *
 * @param[in]  mig_device MIG device handle
 * @param[out] device     Device handle
 *
 * @return NVML_SUCCESS
 *         -> if \p device handle was successfully created \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p mig_device or \p device reference are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If device doesn't have MIG mode enabled \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlDeviceGetDeviceHandleFromMigDeviceHandle(
    nvmlDevice_t mig_device, nvmlDevice_t *device);

/*******************************************************************************
 * @name nvmlSetSystemMigMode
 *
 * @brief Set MIG mode for the system level (MigManager's status)
 *
 * @note Change MIG modes may require devices unbind or reset. The "pending" MIG
 *       mode refers to the target mode following the next activation trigger.
 *
 * @param[in]  mode              The mode to be set for system level,
 *                               \p NVML_DEVICE_MIG_ENABLE or
 *                               \p NVML_DEVICE_MIG_DISABLE
 * @param[out] activation_status The activationStatus status
 *
 * @return NVML_SUCCESS
 *         -> If MIG mode has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p activation_status is invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If all gpu devices in system doesn't support MIG mode \
 * @return NVML_ERROR_IN_USE
 *         -> If there's any gpu device is in mig mode \p NVML_DEVICE_MIG_ENABLE
 *\
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlSetSystemMigMode(unsigned int mode,
                                  nvmlReturn_t *activation_status);
/*******************************************************************************
 * @name nvmlGetSystemMigMode
 *
 * @brief Get MIG mode for the system level (MigManager's status)
 *
 * @note Change MIG modes may require devices unbind or reset. The "pending" MIG
 *       mode refers to the target mode following the next activation trigger.
 *
 * @param[out] current_mode Returns the current system level mode of MigManager,
 *                          \p NVML_DEVICE_MIG_ENABLE or
 *                          \p NVML_DEVICE_MIG_DISABLE
 * @param[out] pending_mode Returns the pending system level mode of MigManager,
 *                          \p NVML_DEVICE_MIG_ENABLE,
 *                          \p NVML_DEVICE_MIG_DISABLE
 *
 * @return NVML_SUCCESS
 *         -> If MIG mode has been retrieved successfully \
 * @return NVML_ERROR_UNINITIALIZED
 *         -> If dmi runtime has not been successfully initialized \
 * @return NVML_ERROR_INVALID_ARGUMENT
 *         -> If \p current_mode or \p pending_mode are invalid \
 * @return NVML_ERROR_NOT_SUPPORTED
 *         -> If all gpu devices in system doesn't support MIG mode \
 * @return NVML_ERROR_NO_PERMISSION
 *         -> If user doesn't have permission to perform the operation \
 * @return NVML_ERROR_UNKNOWN
 *         -> on any unexpected error
 ******************************************************************************/
nvmlReturn_t nvmlGetSystemMigMode(unsigned int *current_mode,
                                  unsigned int *pending_mode);

#ifdef __cplusplus
} // extern "C"
#endif

#endif // __INC_DMI_MIG_H__