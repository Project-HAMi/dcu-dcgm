package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L./lib -lrocm_smi64 -lhydmi -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
#include <dmi_virtual.h>
#include <dmi_error.h>
#include <dmi.h>
#include <dmi_mig.h>
*/
import "C"

type RSMIPcieBandwidth struct {
	TransferRate RSMIFrequencies
	lanes        [32]uint32
}

type RSMIFrequencies struct {
	NumSupported uint32
	Current      uint32
	Frequency    [32]uint64
}

type RSNIPowerProfilePresetMasks C.rsmi_power_profile_preset_masks_t

const (
	RSMI_PWR_PROF_PRST_CUSTOM_MASK       RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_CUSTOM_MASK       //!< Custom Power Profile
	RSMI_PWR_PROF_PRST_VIDEO_MASK        RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VIDEO_MASK        //!< Video Power Profile
	RSMI_PWR_PROF_PRST_POWER_SAVING_MASK RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_POWER_SAVING_MASK //!< Power Saving Profile
	RSMI_PWR_PROF_PRST_COMPUTE_MASK      RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_COMPUTE_MASK      //!< Compute Saving Profile
	RSMI_PWR_PROF_PRST_VR_MASK           RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VR_MASK           //!< VR Power Profile

	//!< 3D Full Screen Power Profile
	RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK
	RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT   RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT //!< Default Boot Up Profile
	RSMI_PWR_PROF_PRST_LAST             RSNIPowerProfilePresetMasks = RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT

	//!< Invalid power profile
	RSMI_PWR_PROF_PRST_INVALID RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_INVALID
)

type RSMIRetiredPageRecord struct {
	PageAddress uint64               //!< Start address of page
	PageSize    uint64               //!< Page size
	Status      RSMIMemoryPageStatus //!< Page "reserved" status
}

type RSMIMemoryPageStatus C.rsmi_memory_page_status_t

const (
	RSMI_MEM_PAGE_STATUS_RESERVED     RSMIMemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_RESERVED
	RSMI_MEM_PAGE_STATUS_PENDING      RSMIMemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_PENDING
	RSMI_MEM_PAGE_STATUS_UNRESERVABLE RSMIMemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_UNRESERVABLE
)

type RSMIFreqVoltRegion struct {
	FreqRange RSMIRange
	VoltRange RSMIRange
}

type RSMITemperatureMetric C.rsmi_temperature_metric_t

const (
	RSMI_TEMP_CURRENT        RSMITemperatureMetric = C.RSMI_TEMP_CURRENT
	RSMI_TEMP_FIRST          RSMITemperatureMetric = C.RSMI_TEMP_FIRST
	RSMI_TEMP_MAX            RSMITemperatureMetric = C.RSMI_TEMP_MAX
	RSMI_TEMP_MIN            RSMITemperatureMetric = C.RSMI_TEMP_MIN
	RSMI_TEMP_MAX_HYST       RSMITemperatureMetric = C.RSMI_TEMP_MAX_HYST
	RSMI_TEMP_MIN_HYST       RSMITemperatureMetric = C.RSMI_TEMP_MIN_HYST
	RSMI_TEMP_CRITICAL       RSMITemperatureMetric = C.RSMI_TEMP_CRITICAL
	RSMI_TEMP_CRITICAL_HYST  RSMITemperatureMetric = C.RSMI_TEMP_CRITICAL_HYST
	RSMI_TEMP_EMERGENCY      RSMITemperatureMetric = C.RSMI_TEMP_EMERGENCY
	RSMI_TEMP_EMERGENCY_HYST RSMITemperatureMetric = C.RSMI_TEMP_EMERGENCY_HYST
	RSMI_TEMP_CRIT_MIN       RSMITemperatureMetric = C.RSMI_TEMP_CRIT_MIN
	RSMI_TEMP_CRIT_MIN_HYST  RSMITemperatureMetric = C.RSMI_TEMP_CRIT_MIN_HYST
	RSMI_TEMP_OFFSET         RSMITemperatureMetric = C.RSMI_TEMP_OFFSET
	RSMI_TEMP_LOWEST         RSMITemperatureMetric = C.RSMI_TEMP_LOWEST
	RSMI_TEMP_HIGHEST        RSMITemperatureMetric = C.RSMI_TEMP_HIGHEST
	RSMI_TEMP_LAST           RSMITemperatureMetric = C.RSMI_TEMP_LAST
)

type RSMIVoltageType C.rsmi_voltage_type_t

const (
	RSMI_VOLT_TYPE_FIRST   RSMIVoltageType = C.RSMI_VOLT_TYPE_FIRST
	RSMI_VOLT_TYPE_VDDGFX  RSMIVoltageType = C.RSMI_VOLT_TYPE_VDDGFX
	RSMI_VOLT_TYPE_LAST    RSMIVoltageType = C.RSMI_VOLT_TYPE_LAST
	RSMI_VOLT_TYPE_INVALID RSMIVoltageType = C.RSMI_VOLT_TYPE_INVALID
)

type RSMIVoltageMetric C.rsmi_voltage_metric_t

const (
	RSMI_VOLT_CURRENT  RSMIVoltageMetric = C.RSMI_VOLT_CURRENT //!< Voltage current value.
	RSMI_VOLT_FIRST    RSMIVoltageMetric = C.RSMI_VOLT_FIRST
	RSMI_VOLT_MAX      RSMIVoltageMetric = C.RSMI_VOLT_MAX      //!< Voltage max value.
	RSMI_VOLT_MIN_CRIT RSMIVoltageMetric = C.RSMI_VOLT_MIN_CRIT //!< Voltage critical min value.
	RSMI_VOLT_MIN      RSMIVoltageMetric = C.RSMI_VOLT_MIN      //!< Voltage min value.
	RSMI_VOLT_MAX_CRIT RSMIVoltageMetric = C.RSMI_VOLT_MAX_CRIT //!< Voltage critical max value.
	RSMI_VOLT_AVERAGE  RSMIVoltageMetric = C.RSMI_VOLT_AVERAGE  //!< Average voltage.
	RSMI_VOLT_LOWEST   RSMIVoltageMetric = C.RSMI_VOLT_LOWEST   //!< Historical minimum voltage.
	RSMI_VOLT_HIGHEST  RSMIVoltageMetric = C.RSMI_VOLT_HIGHEST  //!< Historical maximum voltage.
	RSMI_VOLT_LAST                       = C.RSMI_VOLT_LAST
)

type RSMIUtilizationCounterType C.RSMI_UTILIZATION_COUNTER_TYPE

const (
	RSMI_UTILIZATION_COUNTER_FIRST RSMIUtilizationCounterType = C.RSMI_UTILIZATION_COUNTER_FIRST
	RSMI_COARSE_GRAIN_GFX_ACTIVITY RSMIUtilizationCounterType = C.RSMI_COARSE_GRAIN_GFX_ACTIVITY
	RSMI_COARSE_GRAIN_MEM_ACTIVITY RSMIUtilizationCounterType = C.RSMI_COARSE_GRAIN_MEM_ACTIVITY
	RSMI_UTILIZATION_COUNTER_LAST  RSMIUtilizationCounterType = C.RSMI_UTILIZATION_COUNTER_LAST
)

type RSMIUtilizationCounter struct {
	Type  RSMIUtilizationCounterType
	Value uint64
}

type RSMIClkType C.rsmi_clk_type_t

const (
	RSMI_CLK_TYPE_SYS  RSMIClkType = C.RSMI_CLK_TYPE_SYS
	RSMI_CLK_TYPE_DF   RSMIClkType = C.RSMI_CLK_TYPE_DF
	RSMI_CLK_TYPE_DCEF RSMIClkType = C.RSMI_CLK_TYPE_DCEF
	RSMI_CLK_TYPE_SOC  RSMIClkType = C.RSMI_CLK_TYPE_SOC
	RSMI_CLK_TYPE_MEM  RSMIClkType = C.RSMI_CLK_TYPE_MEM
	RSMI_CLK_TYPE_PCIE RSMIClkType = C.RSMI_CLK_TYPE_PCIE
	RSMI_CLK_INVALID   RSMIClkType = C.RSMI_CLK_INVALID
)

type RSMIOdVoltFreqData struct {
	CurrSclkRange  RSMIRange
	CurrMclkRange  RSMIRange
	SclkFreqLimits RSMIRange
	MclkFreqLimits RSMIRange
	Curve          RSMIOdVoltCurve
	NumRegions     uint32
}

type RSMIRange struct {
	LowerBound uint64
	UpperBound uint64
}

type RSMIOdVoltCurve struct {
	VcPoints [3]RSMIOdVddcPoint
}

type RSMIOdVddcPoint struct {
	Frequency uint64
	Voltage   uint64
}

type MetricsTableHeader struct {
	StructureSize   uint16
	FormatRevision  uint8
	ContentRevision uint8
}

type RSMIGPUMetrics struct {
	CommonHeader           MetricsTableHeader
	TemperatureEdge        uint16
	TemperatureHotspot     uint16
	TemperatureMem         uint16
	TemperatureVRGfx       uint16
	TemperatureVRSoc       uint16
	TemperatureVRMem       uint16
	AverageGfxActivity     uint16
	AverageUmcActivity     uint16
	AverageMmActivity      uint16
	AverageSocketPower     uint16
	EnergyAccumulator      uint64
	SystemClockCounter     uint64
	AverageGfxclkFrequency uint16
	AverageSocclkFrequency uint16
	AverageUclkFrequency   uint16
	AverageVclk0Frequency  uint16
	AverageDclk0Frequency  uint16
	AverageVclk1Frequency  uint16
	AverageDclk1Frequency  uint16
	CurrentGfxclk          uint16
	CurrentSocclk          uint16
	CurrentUclk            uint16
	CurrentVclk0           uint16
	CurrentDclk0           uint16
	CurrentVclk1           uint16
	CurrentDclk1           uint16
	ThrottleStatus         uint32
	CurrentFanSpeed        uint16
	PcieLinkWidth          uint16
	PcieLinkSpeed          uint16
	Padding                uint16
	GfxActivityAcc         uint32
	MemActivityAcc         uint32
	TemperatureHBM         [4]uint16
}

type RSMIDevPerfLevel C.rsmi_dev_perf_level_t

const (
	RSMI_DEV_PERF_LEVEL_AUTO            RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_AUTO
	RSMI_DEV_PERF_LEVEL_FIRST           RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_FIRST
	RSMI_DEV_PERF_LEVEL_LOW             RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_LOW
	RSMI_DEV_PERF_LEVEL_HIGH            RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_HIGH
	RSMI_DEV_PERF_LEVEL_MANUAL          RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_MANUAL
	RSMI_DEV_PERF_LEVEL_STABLE_STD      RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_STD
	RSMI_DEV_PERF_LEVEL_STABLE_PEAK     RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_PEAK
	RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK
	RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK
	RSMI_DEV_PERF_LEVEL_DETERMINISM     RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_DETERMINISM
	RSMI_DEV_PERF_LEVEL_LAST            RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_LAST
	RSMI_DEV_PERF_LEVEL_UNKNOWN         RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_UNKNOWN
)

type RSMIBitField C.rsmi_bit_field_t

type RSMIPowerProfilePresetMasks C.rsmi_power_profile_preset_masks_t

// 定义 power profile preset masks 的枚举类型
const (
	RSMIPowerProfPrstCustomMask      RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_CUSTOM_MASK       // Custom Power Profile
	RSMIPowerProfPrstVideoMask       RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VIDEO_MASK        // Video Power Profile
	RSMIPowerProfPrstPowerSavingMask RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_POWER_SAVING_MASK // Power Saving Profile
	RSMIPowerProfPrstComputeMask     RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_COMPUTE_MASK      // Compute Saving Profile
	RSMIPowerProfPrstVRMask          RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VR_MASK           // VR Power Profile
	RSMIPowerProfPrst3DFullScrMask   RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK  // 3D Full Screen Power Profile
	RSMIPowerProfPrstBootupDefault   RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT    // Default Boot Up Profile
	RSMIPowerProfPrstLast            RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_LAST              // Last Profile (same as Bootup Default)
	RSMIPowerProfPrstInvalid         RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_INVALID           // Invalid power profile
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
	RSMISwCompFirst  RSMISwComponent = C.RSMI_SW_COMP_FIRST
	RSMISwCompDriver RSMISwComponent = C.RSMI_SW_COMP_DRIVER
	RSMISwCompLast   RSMISwComponent = C.RSMI_SW_COMP_LAST
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

type EventHandle C.rsmi_event_handle_t

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

type RSMIStatus C.rsmi_status_t

const (
	RSMI_STATUS_SUCCESS             RSMIStatus = C.RSMI_STATUS_SUCCESS             //!< Operation was successful
	RSMI_STATUS_INVALID_ARGS        RSMIStatus = C.RSMI_STATUS_INVALID_ARGS        //!< Passed in arguments are not valid
	RSMI_STATUS_NOT_SUPPORTED       RSMIStatus = C.RSMI_STATUS_NOT_SUPPORTED       //!< The requested information or
	RSMI_STATUS_FILE_ERROR          RSMIStatus = C.RSMI_STATUS_FILE_ERROR          //!< Problem accessing a file. This
	RSMI_STATUS_PERMISSION          RSMIStatus = C.RSMI_STATUS_PERMISSION          //!< Permission denied/EACCESS file
	RSMI_STATUS_OUT_OF_RESOURCES    RSMIStatus = C.RSMI_STATUS_OUT_OF_RESOURCES    //!< Unable to acquire memory or other
	RSMI_STATUS_INTERNAL_EXCEPTION  RSMIStatus = C.RSMI_STATUS_INTERNAL_EXCEPTION  //!< An internal exception was caught
	RSMI_STATUS_INPUT_OUT_OF_BOUNDS RSMIStatus = C.RSMI_STATUS_INPUT_OUT_OF_BOUNDS //!< The provided input is out of
	RSMI_STATUS_INIT_ERROR          RSMIStatus = C.RSMI_STATUS_INIT_ERROR          //!< An error occurred when rsmi
	RSMI_INITIALIZATION_ERROR       RSMIStatus = C.RSMI_INITIALIZATION_ERROR
	RSMI_STATUS_NOT_YET_IMPLEMENTED RSMIStatus = C.RSMI_STATUS_NOT_YET_IMPLEMENTED //!< The requested function has not
	RSMI_STATUS_NOT_FOUND           RSMIStatus = C.RSMI_STATUS_NOT_FOUND           //!< An item was searched for but not
	RSMI_STATUS_INSUFFICIENT_SIZE   RSMIStatus = C.RSMI_STATUS_INSUFFICIENT_SIZE   //!< Not enough resources were
	RSMI_STATUS_INTERRUPT           RSMIStatus = C.RSMI_STATUS_INTERRUPT           //!< An interrupt occurred during
	RSMI_STATUS_UNEXPECTED_SIZE     RSMIStatus = C.RSMI_STATUS_UNEXPECTED_SIZE     //!< An unexpected amount of data
	RSMI_STATUS_NO_DATA             RSMIStatus = C.RSMI_STATUS_NO_DATA             //!< No data was found for a given
	RSMI_STATUS_UNEXPECTED_DATA     RSMIStatus = C.RSMI_STATUS_UNEXPECTED_DATA     //!< The data read or provided to
	RSMI_STATUS_BUSY                RSMIStatus = C.RSMI_STATUS_BUSY
	RSMI_STATUS_REFCOUNT_OVERFLOW   RSMIStatus = C.RSMI_STATUS_REFCOUNT_OVERFLOW   //!< An internal reference counter
	RSMI_STATUS_SETTING_UNAVAILABLE RSMIStatus = C.RSMI_STATUS_SETTING_UNAVAILABLE //!< Requested setting is unavailable
	RSMI_STATUS_AMDGPU_RESTART_ERR  RSMIStatus = C.RSMI_STATUS_AMDGPU_RESTART_ERR  //!< Could not successfully restart
	RSMI_STATUS_UNKNOWN_ERROR       RSMIStatus = C.RSMI_STATUS_UNKNOWN_ERROR
)

type MonitorInfo struct {
	MinorNumber     int
	PicBusNumber    string
	DeviceId        string
	SubSystemName   string
	Temperature     float64
	PowerUsage      float64
	PowerCap        float64
	MemoryCap       float64
	MemoryUsed      float64
	UtilizationRate float64
	PcieBwMb        float64
	Clk             float64
}

type DeviceInfo struct {
	DvInd        int
	DeviceId     string
	DevType      string
	DevTypeName  string
	PicBusNumber string
	MemoryTotal  float64
	MemoryUsed   float64
	ComputeUnit  float64
}

var type2name = map[string]string{
	"66a1": "WK100",
	"51b7": "Z100L",
	"52b7": "Z100L",
	"53b7": "Z100L",
	"54b7": "Z100L",
	"55b7": "Z100L",
	"56b7": "Z100L",
	"57b7": "Z100L",
	"61b7": "K100",
	"62b7": "K100",
	"63b7": "K100",
	"64b7": "K100",
	"65b7": "K100",
	"66b7": "K100",
	"67b7": "K100",
	"6210": "K100 AI",
	"6211": "K100 AI Liquid",
	"6212": "K100 AI Liquid",
}

var computeUnitType = map[string]float64{
	"K100 AI": 120,
	"K100":    120,
	"Z100":    60,
	"Z100L":   60,
}

var memoryTypeL = []string{"VRAM", "VIS_VRAM", "GTT"}

const DMI_NAME_SIZE = 256

type DMIDeviceInfo struct {
	Name                      string
	ComputeUnitCount          int
	ComputeUnitRemainingCount uintptr
	MemoryRemaining           uintptr
	GlobalMemSize             uintptr
	UsageMemSize              uintptr
	DeviceID                  int
	Percent                   int
	MaxVDeviceCount           int
}

type DMIVDeviceInfo struct {
	Name             string
	ComputeUnitCount int
	GlobalMemSize    uintptr
	UsageMemSize     uintptr
	ContainerID      uint64
	DeviceID         int
	Percent          int
}

type DMIStatus C.dmiStatus

const (
	DMI_STATUS_SUCCESS                DMIStatus = C.DMI_STATUS_SUCCESS
	DMI_STATUS_ERROR                  DMIStatus = C.DMI_STATUS_ERROR
	DMI_STATUS_NO_MEMORY              DMIStatus = C.DMI_STATUS_NO_MEMORY
	DMI_STATUS_OPEN_MKFD_FAILED       DMIStatus = C.DMI_STATUS_OPEN_MKFD_FAILED
	DMI_STATUS_MKFD_ALREADY_OPENED    DMIStatus = C.DMI_STATUS_MKFD_ALREADY_OPENED
	DMI_STATUS_SYS_NODE_NOT_EXIST     DMIStatus = C.DMI_STATUS_SYS_NODE_NOT_EXIST
	DMI_STATUS_NOT_SUPPORTED          DMIStatus = C.DMI_STATUS_NOT_SUPPORTED
	DMI_STATUS_MKFD_NOT_OPENED        DMIStatus = C.DMI_STATUS_MKFD_NOT_OPENED
	DMI_STATUS_CREATE_VDEV_FAILED     DMIStatus = C.DMI_STATUS_CREATE_VDEV_FAILED
	DMI_STATUS_DESTROY_VDEV_FAILED    DMIStatus = C.DMI_STATUS_DESTROY_VDEV_FAILED
	DMI_STATUS_INVALID_ARGUMENTS      DMIStatus = C.DMI_STATUS_INVALID_ARGUMENTS
	DMI_STATUS_OUT_OF_RESOURCES       DMIStatus = C.DMI_STATUS_OUT_OF_RESOURCES
	DMI_STATUS_QUERY_VDEV_INFO_FAILED DMIStatus = C.DMI_STATUS_QUERY_VDEV_INFO_FAILED
	DMI_STATUS_ERROR_NOT_INITIALIZED  DMIStatus = C.DMI_STATUS_ERROR_NOT_INITIALIZED
	DMI_STATUS_DEVICE_NOT_SUPPORT     DMIStatus = C.DMI_STATUS_DEVICE_NOT_SUPPORT
	DMI_STATUS_VDEV_NOT_EXIST         DMIStatus = C.DMI_STATUS_VDEV_NOT_EXIST
	DMI_STATUS_INIT_DEVICE_FAILED     DMIStatus = C.DMI_STATUS_INIT_DEVICE_FAILED
	DMI_STATUS_DEVICE_BUSY            DMIStatus = C.DMI_STATUS_DEVICE_BUSY
	DMI_STATUS_FILE_ERROR             DMIStatus = C.DMI_STATUS_FILE_ERROR
	DMI_STATUS_PERMISSION             DMIStatus = C.DMI_STATUS_PERMISSION
	DMI_STATUS_INTERNAL_EXCEPTION     DMIStatus = C.DMI_STATUS_INTERNAL_EXCEPTION
	DMI_STATUS_INPUT_OUT_OF_BOUNDS    DMIStatus = C.DMI_STATUS_INPUT_OUT_OF_BOUNDS
	DMI_STATUS_SMI_INIT_ERROR         DMIStatus = C.DMI_STATUS_SMI_INIT_ERROR
	DMI_STATUS_NOT_FOUND              DMIStatus = C.DMI_STATUS_NOT_FOUND
	DMI_STATUS_INSUFFICIENT_SIZE      DMIStatus = C.DMI_STATUS_INSUFFICIENT_SIZE
	DMI_STATUS_INTERRUPT              DMIStatus = C.DMI_STATUS_INTERRUPT
	DMI_STATUS_UNEXPECTED_SIZE        DMIStatus = C.DMI_STATUS_UNEXPECTED_SIZE
	DMI_STATUS_NO_DATA                DMIStatus = C.DMI_STATUS_NO_DATA
	DMI_STATUS_UNEXPECTED_DATA        DMIStatus = C.DMI_STATUS_UNEXPECTED_DATA
	DMI_STATUS_SMI_BUSY               DMIStatus = C.DMI_STATUS_SMI_BUSY
	DMI_STATUS_REFCOUNT_OVERFLOW      DMIStatus = C.DMI_STATUS_REFCOUNT_OVERFLOW
	DMI_STATUS_NOT_YET_IMPLEMENTED    DMIStatus = C.DMI_STATUS_NOT_YET_IMPLEMENTED
	DMI_STATUS_UNKNOWN_ERROR          DMIStatus = C.DMI_STATUS_UNKNOWN_ERROR
)

// 定义 PhysicalDeviceInfo 结构体
type PhysicalDeviceInfo struct {
	DeviceInfo     DMIDeviceInfo
	VirtualDevices []DMIVDeviceInfo
}
