/*
Copyright 2022 The Katalyst Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"net"
	"sync"
	"time"

	info "github.com/google/cadvisor/info/v1"
	v1 "k8s.io/api/core/v1"
)

// GetNodeInfo()
type NodeInfo struct {
	Meta   MachineMeta
	Status *MachineStatus
}

type MachineMeta struct {
	ReporterStartTime time.Time
	HostName          string
	HostIPv4          net.IP
	HostIPv6          net.IP
	ArchName          string
	CPUVendor         string
	CPUNum            int
	SocketNum         int
	PageSize          int
	GPUNum            int
	GPUUUIDMap        map[string]int
	PhysicalCluster   string // os.Getenv(TCE_CLUSTER)
	IDC               string
	Vendor            string
	NodeLevel         string
	Specification     string
	OSVersion         string
	CgroupVersion     string
	IsVM              *bool
	DiskMode          string
	// Topology Describes cpu/memory layout and hierarchy.
	Topology []info.Node
}

type MachineStatus struct {
	CPUInfoItems     CPUInfoItemList
	CPUUPowerMode    string
	NumaMaxBandwidth map[int]float64
}

type CPUInfoItem struct {
	FreqMhz  float64
	SocketID int
}

type CPUInfoItemList []CPUInfoItem

// GetMalachiteNode()
type MalachiteNode struct {
	Load    NodeLoad
	CPU     NodeCPU
	Memory  NodeMemory
	IO      NodeIO
	Network NodeNetwork

	QosLevelCgroups map[QosLevel]*CgroupResource

	RawCompute *SystemCompute
	RawMemory  *SystemMemory
	RawIO      *SystemIO
	RawNetwork *SystemNetwork
	RawEvents  *SystemEventData
	RawSensors *SystemSensorData
}

type QosLevel string

type NodeLoad struct {
	LoadOne     float64
	LoadFive    float64
	LoadFifteen float64

	CpuPressure Pressure
	MemPressure Pressure
	IoPressure  Pressure

	ProcessRunning uint64
	ProcessBlocked uint64
}

type NodeCPU struct {
	CPUTotal    int
	SocketNum   int
	CPUCodeName string
	CPUs        []CPU
	CpuGlobal   CPU
	Numas       map[int]*Numa
}

type NodeMemory struct {
	MemTotal               uint64 // unit: kB
	MemFree                uint64
	MemUsed                uint64
	MemShm                 uint64
	MemAvailable           uint64
	MemBuffers             uint64
	MemPageCache           uint64
	MemSlabReclaimable     uint64
	MemDirtyPageCache      uint64
	MemWritebackPageCache  uint64
	MemSwapTotal           uint64
	MemSwapFree            uint64
	MemActiveAnon          uint64
	MemInactiveAnon        uint64
	MemActiveFile          uint64
	MemInactiveFile        uint64
	VMWatermarkScaleFactor uint64
	VMStatPgStealKswapd    uint64
	VMStatPgStealDirect    uint64
	VMStatPgScanKswapd     uint64
	VMStatPgScanDirect     uint64
	VMStatCompactStall     uint64
	MemSockTcp             uint64
	MemSockUdp             uint64
	MemSockTcpLimit        uint64
	MemSockUdpLimit        uint64
}

type NodeIO struct {
	DiskIO    []DiskIO
	DiskUsage map[string]*NodeDiskUsage
}

type NodeDiskUsage struct {
	DiskUsage
	UsedBytes  uint64
	UsedInodes uint64
	IORead     uint64
	IOWrite    uint64
}

type NodeNetwork struct {
	TCP    TCP
	NicMap map[string]*Networkcard
}

// GetPodSummaryMap
type PodSummary struct {
	PodSummaryMeta
	// Pod is always not nil
	Pod            *v1.Pod `json:"-"`
	PodResource    PodResource
	CgroupResource CgroupResource
	// uuid is from kubelte, while metrics' uuid is from baca, they may be different, check this later
	GPUUUIDs []string
	PIDs     []string

	// ContainerSummary's length and order is strictly equals to Pod's Container Spec
	ContainerSummary []*ContainerSummary
}

type ContainerSummary struct {
	ContainerSummaryMeta

	Pod                *v1.Pod `json:"-"`
	ContainersResource ContainerResource
	CgroupResource     *CgroupResource // may be nil, because mala-core not found
	GPUUUIDs           []string
}

type ContainerResource struct {
	KubeResource
	DiskUsage  uint64
	InodeUsage uint64
}

type ContainerSummaryMeta struct {
	PodSummaryMeta

	ContainerIndex     int
	ContainerID        string
	ContainerName      string
	ContainerImage     string
	IsPrimaryContainer bool // if container is either primary container or sidecar container
	HasGPU             bool
	HasNPU             bool
	HasXPU             bool
}

type PodResource struct {
	KubeResource // asert not nil
	JVM          FileSystemPodJVM
	DiskUsage    uint64
	InodeUsage   uint64
}

type FileSystemPodJVM struct {
	Capacity int64
	Used     int64
}

type KubeResource struct {
	Requests v1.ResourceList
	Limits   v1.ResourceList
}

type PodSummaryMeta struct {
	PSM           string
	NodeLevel     string
	PodType       PodType
	RuntimeClass  RuntimeClass
	KatalystQos   KatalystQos
	QuotaSaleMode QuotaSaleMode
	TceClusterID  string
	TceCluster    string
}
type PodType string
type RuntimeClass string
type KatalystQos string
type QuotaSaleMode string

type CgroupResource struct {
	Path              string
	Version           string
	UpdateTimeAvg     int64
	UpdateTimeCPU     int64
	UpdateTimeMemory  int64
	UpdateTimeIO      int64
	UpdateTimeNetwork int64
	UpdateTimeResctrl int64

	CgroupCPUResource
	CgroupMemoryResource
	CgroupNetworkResource
	CgroupIOResource

	Raw CgroupData
}

type CgroupData struct {
	MountPoint string
	UserPath   string
	GroupsV1   *SubSystemGroupsV1
	GroupsV2   *SubSystemGroupsV2
	CgroupType string
	URL        string
}

type CgroupCPUResource struct {
	CPUSetMeta             string
	CPUSetCPUsInner        []int
	CPUSetMemsInner        []int
	CPUUtilization         float64
	CPUUserUtilization     float64
	CPUSysUtilization      float64
	CPULoadOne             float64
	CPULoadFive            float64
	CPULoadFifteen         float64
	CPULoadRunnableOne     float64
	CPULoadRunnableFive    float64
	CPULoadRunnableFifteen float64
	TaskNrUninterruptible  uint64
	TaskNrRunning          uint64
	TaskNrSleeping         uint64
	TaskNrIOWait           uint64
	CPUNrPeriods           uint64
	CPUNrThrottled         uint64
	CPUThrottledTime       uint64
	CPUInstructions        uint64
	CPUCycles              uint64
	CPUOcrReadDrams        uint64
	CPUL3CacheMiss         uint64
	CPUStoreAllIns         uint64
	CPUStoreIns            uint64
	CPUImcWrites           uint64
	CPUCfsQuotaUs          int64
	CPUCfsPeriodUs         int64
	CPUShares              uint64
	CPUNumaUsage           map[int]uint64
	CPUSchedWait           []uint64
}

type CgroupMemoryResource struct {
	MemoryMeta                      string
	MemoryMax                       uint64
	MemoryUsed                      uint64
	MemoryKernelUsed                uint64
	MemoryTotalRss                  uint64
	MemoryTotalCache                uint64
	MemoryTotalDirty                uint64
	MemorySwapMax                   uint64
	MemoryTotalSwap                 uint64
	MemoryTotalInactiveAnon         uint64
	MemoryTotalActiveAnon           uint64
	MemoryTotalInactiveFile         uint64
	MemoryTotalActiveFile           uint64
	MemoryPgfault                   uint64
	MemoryTotalPgfault              uint64
	MemoryPgmajfault                uint64
	MemoryTotalPgmajfault           uint64
	MemoryShmem                     uint64
	MemoryWriteback                 uint64
	MemoryOOMKillCnt                uint64
	MemoryOOMTriggerCnt             uint64
	MemoryNumaStats                 map[string]NumaStats
	MemoryStatSock                  uint64
	MemoryBalanceDirty              uint64
	MemoryReclaimCnt                uint64
	MemoryCompactCnt                uint64
	MemoryProactiveReclaimTargetSum uint64
	MemoryLow                       int64
	MemoryMin                       int64
}

type NumaStats struct {
	NumaName    string
	Total       uint64
	File        uint64
	Anon        uint64
	UnEvictable uint64
}

type CgroupNetworkResource struct {
	NetworkTcpRX      uint64
	NetworkTcpTX      uint64
	NetworkTcpRXBytes uint64
	NetworkTcpTXBytes uint64

	NetworkUdpRX      uint64
	NetworkUdpTX      uint64
	NetworkUdpRXBytes uint64
	NetworkUdpTXBytes uint64
}

type CgroupIOResource struct {
	IOFsCreated      uint64
	IOFsOpen         uint64
	IOFsRead         uint64
	IOFsReadBytes    uint64
	IOFsWrite        uint64
	IOFsWriteBytes   uint64
	IOFsync          uint64
	IOWait           uint64
	IOWriteBackPages uint64

	IOV2 struct {
		IOStat map[string]*DeviceIOStat
	}
}

type DeviceIOStat struct {
	ReadBytes  uint64
	WriteBytes uint64
	ReadIOs    uint64
	WriteIOs   uint64
	UpdateTime int64
}

// GetServiceSummaryMap
type ServiceSummary struct {
	ServiceSummaryMeta

	ServiceSummaryType CgroupServiceType
	CgroupResource     CgroupResource
}

type CgroupServiceType string

type ServiceSummaryMeta struct {
	Name               string
	PSM                string
	NodeLevel          string
	HostIP             string
	HostIPv6           string
	OSVersion          string
	CPUModel           string
	Specification      string
	PodType            PodType
	KatalystQos        KatalystQos
	IsPrimaryContainer bool
}

// GetPodRateMap
type PodRateEntity struct {
	PodRate        *PodRate
	ContainersRate []*ContainerRate
}

type PodRate struct {
	PodName string

	RateData
}

type RateData struct {
	CPUCyclesPerInstruction   float64
	CPUInstructionsPerSecond  float64
	CPUCyclesPerSecond        float64
	CPUNRThrottledPerSecond   float64
	CPUNRPeriodsPerSecond     float64
	CPUThrottledTimePerSecond float64
	CPUL3CacheMissesPerSecond float64
	CPUNUMAUsage              map[int]float64
	CPUSchedWaitPerSecondMax  float64
	CPUSchedWaitPerSecondMin  float64
	CPUSchedWaitPerSecondAvg  float64
	CPUSchedWaitPerSecondSum  float64

	MemoryReadBytesPerSecond             float64
	MemoryWriteBytesPerSecond            float64
	MemoryPGFaultPerSecond               float64
	MemoryPGMajFaultPerSecond            float64
	MemoryOOMCntPerSecond                float64
	MemoryAllocStallPerSecond            float64
	MemoryKswapdPgStealPerSecond         float64
	MemoryKswapdPgScanPerSecond          float64
	MemoryDirectPgStealPerSecond         float64
	MemoryDirectPgScanPerSecond          float64
	MemoryPgStealPerSecond               float64
	MemoryPgScanPerSecond                float64
	MemoryWorkingsetRefaultAnonPerSecond float64
	MemoryWorkingsetRefaultFilePerSecond float64
	MemoryWorkingsetRefaultPerSecond     float64
	MemoryProactiveReclaimPerSecond      float64

	IOFsCreatePerSecond     float64
	IOFsOpenPerSecond       float64
	IOFsSyncPerSecond       float64
	IOFsReadPerSecond       float64
	IOFsWritePerSecond      float64
	IOFsReadBytesPerSecond  float64
	IOFsWriteBytesPerSecond float64
	IODiskIOWaitPerSecond   float64

	TcpRXPerSecond        float64
	TcpTXPerSecond        float64
	TcpRXBytesPerSecond   float64
	TcpTXBytesPerSecond   float64
	TcpRxErrorsPerSecond  float64
	TcpTxErrorsPerSecond  float64
	TcpRxDroppedPerSecond float64
	TcpTxDroppedPerSecond float64

	UdpRXPerSecond      float64
	UdpTXPerSecond      float64
	UdpRXBytesPerSecond float64
	UdpTXBytesPerSecond float64

	DeviceIOPS map[string]*DeviceIOStatRate

	MBRateMap map[int]*MBMStat
}

type DeviceIOStatRate struct {
	ReadIOsPerSecond    uint64
	WriteIOsPerSecond   uint64
	ReadBytesPerSecond  uint64
	WriteBytesPerSecond uint64
}

type MBMStat struct {
	MBMLocalBytes uint64
	MBMTotalBytes uint64
}

type ContainerRate struct {
	ContainerID string
	Valid       bool

	RateData
	OOMLevelTrigger *LevelTriggerData
}

type LevelTriggerData struct {
	m            sync.Mutex
	accumulation int
	events       []LevelTriggerEvent
}

type LevelTriggerEvent struct {
	Timestamp time.Time
	Data      int
}

// GetGroupRateMap
type GroupRate struct {
	QoSLevel QosLevel

	RateData
}

// GetMalachiteNodeRate
type MalachiteNodeRate struct {
	TcpCloseWaitPerSecond      uint64
	TcpDelayAckPerSecond       uint64
	TcpListenOverflowPerSecond uint64
	TcpListenDropPerSecond     uint64
	TcpQFullDropPerSecond      uint64
	TcpAbortOnMemoryPerSecond  uint64
	TcpRetransPerOutSegs       float64
	TcpTimeoutsPerOutSegs      float64
	TcpMemPressurePct          float64

	SystemEventRate MalachiteSystemEventRate

	NicRateMap map[string]*MalachiteNicRate

	ComputeBpfProgRateMap map[string]*MalachiteComputeBpfProgRate

	MBRateMap map[int]*MalachiteMBRate
}

type MalachiteSystemEventRate struct {
	SoftLockupPS      float64
	RCUStallPS        float64
	BadPagePS         float64
	KernelWarn        float64
	MceUc             float64
	HungTask          float64
	CoreDump          float64
	IOError           float64
	Ext4Error         float64
	Ext4Abrt          float64
	NetXmitTimeout    float64
	TcpBadCsum        float64
	DevLinkDown       float64
	NetExceedBufLimit float64
	NetRcvqueueFull   float64
	MemAllocFailure   float64
	TotalOom          float64
	GlobalOom         float64
}

type MalachiteNicRate struct {
	RxBytePS  float64
	RxPPS     float64
	RxDropPPS float64
	TxBytePS  float64
	TxPPS     float64
	TxDropPPS float64
	Duplex    string
	Speeds    *uint64
}

type MalachiteComputeBpfProgRate struct {
	BpfProgCPUUsage          float64
	BpfProgRunCountPerSecond float64
	BpfProgRunTimePerCount   float64
}

type MalachiteMBRate struct {
	MBMLocalBytesPS  uint64
	MBMTotalBytesPS  uint64
	MBMVictimBytesPS uint64
	MBMMaxBytesPS    uint64
	NUMAID           int
}
