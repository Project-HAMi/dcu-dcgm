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

type RSMIBitField uint64

type RSMIPowerProfilePresetMasks uint64

// 定义 power profile preset masks 的枚举类型
const (
	RSMIPowerProfPrstCustomMask      RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_CUSTOM_MASK   // Custom Power Profile
	RSMIPowerProfPrstVideoMask       RSMIPowerProfilePresetMasks = C.RSMIPowerProfPrstVideoMask       // Video Power Profile
	RSMIPowerProfPrstPowerSavingMask RSMIPowerProfilePresetMasks = C.RSMIPowerProfPrstPowerSavingMask // Power Saving Profile
	RSMIPowerProfPrstComputeMask     RSMIPowerProfilePresetMasks = C.RSMIPowerProfPrstComputeMask     // Compute Saving Profile
	RSMIPowerProfPrstVRMask          RSMIPowerProfilePresetMasks = C.RSMIPowerProfPrstVRMask          // VR Power Profile
	RSMIPowerProfPrst3DFullScrMask   RSMIPowerProfilePresetMasks = C.RSMIPowerProfPrst3DFullScrMask   // 3D Full Screen Power Profile
	RSMIPowerProfPrstBootupDefault   RSMIPowerProfilePresetMasks = C.RSMIPowerProfPrstBootupDefault   // Default Boot Up Profile
	RSMIPowerProfPrstLast            RSMIPowerProfilePresetMasks = C.RSMIPowerProfPrstLast            // Last Profile (same as Bootup Default)
	RSMIPowerProfPrstInvalid         RSMIPowerProfilePresetMasks = C.RSMIPowerProfPrstInvalid         // Invalid power profile
)

// 定义 power profile status 结构体
type RSMPowerProfileStatus struct {
	AvailableProfiles RSMIBitField                // 哪些配置文件被系统支持
	Current           RSMIPowerProfilePresetMasks // 当前激活的电源配置文件
	NumProfiles       uint32                      // 可用的电源配置文件数量
}

type RSMIVersion struct {
	Major uint32
	Minor uint32
	Patch uint32
	Build string
}

type RSMISwComponent C.rsmi_sw_component_t

const (
	RSMISwCompFirst  RSMISwComponent = C.RSMISwCompFirst
	RSMISwCompDriver RSMISwComponent = C.RSMISwCompDriver
	RSMISwCompLast   RSMISwComponent = C.RSMISwCompLast
)

// 用于识别各种固
type RSMIFwBlock C.rsmi_fw_block_t

const (
	RSMIFwBlockFirst    RSMIFwBlock = C.RSMI_FW_BLOCK_FIRST
	RSMIFwBlockASD      RSMIFwBlock = C.RSMI_FW_BLOCK_ASD
	RSMIFwBlockCE       RSMIFwBlock = C.RSMI_FW_BLOCK_CE
	RSMIFwBlockDMCU     RSMIFwBlock = C.RSMI_FW_BLOCK_DMCU
	RSMIFwBlockMC       RSMIFwBlock = C.RSMI_FW_BLOCK_MC
	RSMIFwBlockME       RSMIFwBlock = C.RSMI_FW_BLOCK_ME
	RSMIFwBlockMEC      RSMIFwBlock = C.RSMI_FW_BLOCK_MEC
	RSMIFwBlockMEC2     RSMIFwBlock = C.RSMI_FW_BLOCK_MEC2
	RSMIFwBlockPFP      RSMIFwBlock = C.RSMI_FW_BLOCK_PFP
	RSMIFwBlockRLC      RSMIFwBlock = C.RSMI_FW_BLOCK_RLC
	RSMIFwBlockRLC_SRLC RSMIFwBlock = C.RSMI_FW_BLOCK_RLC_SRLC
	RSMIFwBlockRLC_SRLG RSMIFwBlock = C.RSMI_FW_BLOCK_RLC_SRLG
	RSMIFwBlockRLC_SRLS RSMIFwBlock = C.RSMI_FW_BLOCK_RLC_SRLS
	RSMIFwBlockSDMA     RSMIFwBlock = C.RSMI_FW_BLOCK_SDMA
	RSMIFwBlockSDMA2    RSMIFwBlock = C.RSMI_FW_BLOCK_SDMA2
	RSMIFwBlockSMC      RSMIFwBlock = C.RSMI_FW_BLOCK_SMC
	RSMIFwBlockSOS      RSMIFwBlock = C.RSMI_FW_BLOCK_SOS
	RSMIFwBlockTA_RAS   RSMIFwBlock = C.RSMI_FW_BLOCK_TA_RAS
	RSMIFwBlockTA_XGMI  RSMIFwBlock = C.RSMI_FW_BLOCK_TA_XGMI
	RSMIFwBlockUVD      RSMIFwBlock = C.RSMI_FW_BLOCK_UVD
	RSMIFwBlockVCE      RSMIFwBlock = C.RSMI_FW_BLOCK_VCE
	RSMIFwBlockVCN      RSMIFwBlock = C.RSMI_FW_BLOCK_VCN
	RSMIFwBlockLast     RSMIFwBlock = C.RSMI_FW_BLOCK_LAST
)

// 保存错误计
type RSMIErrorCount struct {
	CorrectableErr   uint64
	UncorrectableErr uint64
}

// 用于标识不同的GPU
type RSMIGpuBlock C.rsmi_gpu_block_t

const (
	RSMIGpuBlockInvalid  RSMIGpuBlock = C.RSMI_GPU_BLOCK_INVALID
	RSMIGpuBlockFirst    RSMIGpuBlock = C.RSMI_GPU_BLOCK_FIRST
	RSMIGpuBlockUMC      RSMIGpuBlock = C.RSMI_GPU_BLOCK_UMC
	RSMIGpuBlockSDMA     RSMIGpuBlock = C.RSMI_GPU_BLOCK_SDMA
	RSMIGpuBlockGFX      RSMIGpuBlock = C.RSMI_GPU_BLOCK_GFX
	RSMIGpuBlockMMHUB    RSMIGpuBlock = C.RSMI_GPU_BLOCK_MMHUB
	RSMIGpuBlockATHUB    RSMIGpuBlock = C.RSMI_GPU_BLOCK_ATHUB
	RSMIGpuBlockPCIEBIF  RSMIGpuBlock = C.RSMI_GPU_BLOCK_PCIE_BIF
	RSMIGpuBlockHDP      RSMIGpuBlock = C.RSMI_GPU_BLOCK_HDP
	RSMIGpuBlockXGMIWAFL RSMIGpuBlock = C.RSMI_GPU_BLOCK_XGMI_WAFL
	RSMIGpuBlockDF       RSMIGpuBlock = C.RSMI_GPU_BLOCK_DF
	RSMIGpuBlockSMN      RSMIGpuBlock = C.RSMI_GPU_BLOCK_SMN
	RSMIGpuBlockSEM      RSMIGpuBlock = C.RSMI_GPU_BLOCK_SEM
	RSMIGpuBlockMP0      RSMIGpuBlock = C.RSMI_GPU_BLOCK_MP0
	RSMIGpuBlockMP1      RSMIGpuBlock = C.RSMI_GPU_BLOCK_MP1
	RSMIGpuBlockFuse     RSMIGpuBlock = C.RSMI_GPU_BLOCK_FUSE
	RSMIGpuBlockMCA      RSMIGpuBlock = C.RSMI_GPU_BLOCK_MCA
	RSMIGpuBlockLast     RSMIGpuBlock = C.RSMI_GPU_BLOCK_LAST
	RSMIGpuBlockReserved RSMIGpuBlock = C.RSMI_GPU_BLOCK_RESERVED
)

// 当前ECC状态
type RSMIRasErrState C.rsmi_ras_err_state_t

const (
	RSMIRasErrStateNone     RSMIRasErrState = C.RSMI_RAS_ERR_STATE_NONE
	RSMIRasErrStateDisabled RSMIRasErrState = C.RSMI_RAS_ERR_STATE_DISABLED
	RSMIRasErrStateParity   RSMIRasErrState = C.RSMI_RAS_ERR_STATE_PARITY
	RSMIRasErrStateSingC    RSMIRasErrState = C.RSMI_RAS_ERR_STATE_SING_C
	RSMIRasErrStateMultUC   RSMIRasErrState = C.RSMI_RAS_ERR_STATE_MULT_UC
	RSMIRasErrStatePoison   RSMIRasErrState = C.RSMI_RAS_ERR_STATE_POISON
	RSMIRasErrStateEnabled  RSMIRasErrState = C.RSMI_RAS_ERR_STATE_ENABLED
	RSMIRasErrStateLast     RSMIRasErrState = C.RSMI_RAS_ERR_STATE_LAST
	RSMIRasErrStateInvalid  RSMIRasErrState = C.RSMI_RAS_ERR_STATE_INVALID
)

// 事件组枚举值
type RSMIEventGroup C.rsmi_event_group_t

const (
	RSMI_EVNT_GRP_XGMI          RSMIEventGroup = C.RSMI_EVNT_GRP_XGMI
	RSMI_EVNT_GRP_XGMI_DATA_OUT RSMIEventGroup = C.RSMI_EVNT_GRP_XGMI_DATA_OUT
	RSMI_EVNT_GRP_INVALID       RSMIEventGroup = C.RSMI_EVNT_GRP_INVALID
)

type RSMIEventType C.rsmi_event_type_t

const (
	RSMIEventFirst RSMIEventType = C.RSMI_EVNT_FIRST

	RSMIEventXGmiFirst       RSMIEventType = C.RSMI_EVNT_XGMI_FIRST
	RSMIEventXGmi0NopTx      RSMIEventType = C.RSMI_EVNT_XGMI_0_NOP_TX
	RSMIEventXGmi0RequestTx  RSMIEventType = C.RSMI_EVNT_XGMI_0_REQUEST_TX
	RSMIEventXGmi0ResponseTx RSMIEventType = C.RSMI_EVNT_XGMI_0_RESPONSE_TX
	RSMIEventXGmi0BeatsTx    RSMIEventType = C.RSMI_EVNT_XGMI_0_BEATS_TX
	RSMIEventXGmi1NopTx      RSMIEventType = C.RSMI_EVNT_XGMI_1_NOP_TX
	RSMIEventXGmi1RequestTx  RSMIEventType = C.RSMI_EVNT_XGMI_1_REQUEST_TX
	RSMIEventXGmi1ResponseTx RSMIEventType = C.RSMI_EVNT_XGMI_1_RESPONSE_TX
	RSMIEventXGmi1BeatsTx    RSMIEventType = C.RSMI_EVNT_XGMI_1_BEATS_TX

	RSMIEventXGmiLast RSMIEventType = C.RSMI_EVNT_XGMI_LAST

	RSMIEventXGmiDataOutFirst RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_FIRST

	RSMIEventXGmiDataOut0    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_0
	RSMIEventXGmiDataOut1    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_1
	RSMIEventXGmiDataOut2    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_2
	RSMIEventXGmiDataOut3    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_3
	RSMIEventXGmiDataOut4    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_4
	RSMIEventXGmiDataOut5    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_5
	RSMIEventXGmiDataOutLast RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_LAST

	RSMIEventLast RSMIEventType = C.RSMI_EVNT_LAST
)

type EventHandle uintptr

type RSMICounterCommand C.rsmi_counter_command_t

const (
	RSMI_CNTR_CMD_START RSMICounterCommand = C.RSMI_CNTR_CMD_START
	RSMI_CNTR_CMD_STOP  RSMICounterCommand = C.RSMI_CNTR_CMD_STOP
)

// 计数器值
type RSMICounterValue struct {
	Value       uint64
	TimeEnabled uint64
	TimeRunning uint64
}

// 进程的信息
type RSMIProcessInfo struct {
	ProcessID   uint32
	Pasid       uint32
	VramUsage   uint64
	SdmaUsage   uint64
	CuOccupancy uint32
}

// XGMI状态
type RSMIXGMIStatus C.rsmi_xgmi_status_t

const (
	RSMIXGMIStatusNoErrors       RSMIXGMIStatus = C.RSMI_XGMI_STATUS_NO_ERRORS
	RSMIXGMIStatusError          RSMIXGMIStatus = C.RSMI_XGMI_STATUS_ERROR
	RSMIXGMIStatusMultipleErrors RSMIXGMIStatus = C.RSMI_XGMI_STATUS_MULTIPLE_ERRORS
)

// IO链路类型
type RSMIIOLinkType C.RSMI_IO_LINK_TYPE

const (
	RSMIIOLinkTypeUndefined      RSMIIOLinkType = C.RSMI_IOLINK_TYPE_UNDEFINED
	RSMIIOLinkTypePCIExpress     RSMIIOLinkType = C.RSMI_IOLINK_TYPE_PCIEXPRESS
	RSMIIOLinkTypeXGMI           RSMIIOLinkType = C.RSMI_IOLINK_TYPE_XGMI
	RSMIIOLinkTypeNumIOLinkTypes RSMIIOLinkType = C.RSMI_IOLINK_TYPE_NUMIOLINKTYPES
	RSMIIOLinkTypeSize           RSMIIOLinkType = C.RSMI_IOLINK_TYPE_SIZE
)

type RSMIFuncIDIterHandle C.rsmi_func_id_iter_handle_t

type RSMIMemoryType C.rsmi_memory_type_t

const (
	RSMI_MEM_TYPE_FIRST    RSMIMemoryType = C.RSMI_MEM_TYPE_FIRST
	RSMI_MEM_TYPE_VRAM     RSMIMemoryType = C.RSMI_MEM_TYPE_VRAM
	RSMI_MEM_TYPE_VIS_VRAM RSMIMemoryType = C.RSMI_MEM_TYPE_VIS_VRAM
	RSMI_MEM_TYPE_GTT      RSMIMemoryType = C.RSMI_MEM_TYPE_GTT
	RSMI_MEM_TYPE_LAST     RSMIMemoryType = C.RSMI_MEM_TYPE_LAST
)

type RSMIFuncIDValue struct {
	ID         uint64
	Name       string
	MemoryType RSMIMemoryType
	TempMetric RSMITemperatureMetric
	EventType  RSMIEventType
	EventGroup RSMIEventGroup
	ClkType    RSMIClkType
	FwBlock    RSMIFwBlock
	GpuBlock   RSMIGpuBlock
}

type RSMIEvtNotificationType C.rsmi_evt_notification_type_t

const (
	RSMI_EVT_NOTIF_VMFAULT          RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_VMFAULT
	RSMI_EVT_NOTIF_FIRST            RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_FIRST
	RSMI_EVT_NOTIF_THERMAL_THROTTLE RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_THERMAL_THROTTLE
	RSMI_EVT_NOTIF_GPU_PRE_RESET    RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_GPU_PRE_RESET
	RSMI_EVT_NOTIF_GPU_POST_RESET   RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_GPU_POST_RESET
	RSMI_EVT_NOTIF_LAST             RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_LAST
)

type RSMIEEvtNotificationData struct {
	DvInd   uint32
	Event   RSMIEvtNotificationType
	Message [64]byte
}
