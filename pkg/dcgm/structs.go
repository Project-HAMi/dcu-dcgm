package dcgm

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
