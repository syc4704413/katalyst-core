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

type MalachiteSystemDiskIoResponse struct {
	Status int              `json:"status"`
	Data   SystemDiskIoData `json:"data"`
}

type SystemDiskIoData struct {
	DiskIo     []DiskIO `json:"disk_io"`
	UpdateTime int64    `json:"update_time"`
}

type MalachiteSystemNetworkResponse struct {
	Status int               `json:"status"`
	Data   SystemNetworkData `json:"data"`
}

type SystemNetworkData struct {
	SystemNetwork
	UpdateTime int64 `json:"update_time"`
}

type SystemNetwork struct {
	Networkcard []Networkcard `json:"networkcard"`
	TCP         TCP           `json:"tcp"`
}

type Networkcard struct {
	Name               string  `json:"name"`
	Duplex             string  `json:"duplex"`
	UpdateTimeSec      int64   `json:"update_time"`
	ReceiveBytes       uint64  `json:"receive_bytes"`
	ReceivePackets     uint64  `json:"receive_packets"`
	ReceiveErrs        uint64  `json:"receive_errs"`
	ReceiveDrop        uint64  `json:"receive_drop"`
	ReceiveFifo        uint64  `json:"receive_fifo"`
	ReceiveFrame       uint64  `json:"receive_frame"`
	ReceiveCompressed  uint64  `json:"receive_compressed"`
	ReceiveMulticast   uint64  `json:"receive_multicast"`
	TransmitBytes      uint64  `json:"transmit_bytes"`
	TransmitPackets    uint64  `json:"transmit_packets"`
	TransmitErrs       uint64  `json:"transmit_errs"`
	TransmitDrop       uint64  `json:"transmit_drop"`
	TransmitFifo       uint64  `json:"transmit_fifo"`
	TransmitColls      uint64  `json:"transmit_colls"`
	TransmitCarrier    uint64  `json:"transmit_carrier"`
	TransmitCompressed uint64  `json:"transmit_compressed"`
	Speeds             *uint64 `json:"speeds"`
}

type TCP struct {
	UpdateTimeSec         int64   `json:"update_time"`
	TCPDelayAcks          uint64  `json:"tcp_delay_acks"`
	TCPListenOverflows    uint64  `json:"tcp_listen_overflows"`
	TCPListenDrops        uint64  `json:"tcp_listen_drops"`
	TCPAbortOnMemory      uint64  `json:"tcp_abort_on_memory"`
	TCPReqQFullDrop       uint64  `json:"tcp_req_q_full_drop"`
	TCPRetran             float64 `json:"tcp_retran"`
	TCPRetransSegs        uint64  `json:"tcp_retrans_segs"`
	TCPOldRetransSegs     uint64  `json:"tcp_old_retrans_segs"`
	TCPOutSegs            uint64  `json:"tcp_out_segs"`
	TCPOldOutSegs         uint64  `json:"tcp_old_out_segs"`
	TCPCloseWait          uint64  `json:"tcp_close_wait"`
	TCPTimeouts           uint64  `json:"tcp_timeouts"`
	TCPMemPressures       uint64  `json:"tcp_mem_press"`
	TCPMemPressuresTimeMS uint64  `json:"tcp_mem_press_chrono"`
	TCPInPressure         uint64  `json:"tcp_in_pressure"`
	TCPV6InPressure       uint64  `json:"tcpv6_in_pressure"`
}

type MalachiteSystemInfoResponse struct {
	Status int            `json:"status"`
	Data   SystemInfoData `json:"data"`
}

type SystemInfoData struct {
	IsVM       bool  `json:"is_vm"`
	UpdateTime int64 `json:"update_time"`
}

type MalachiteSystemComputeResponse struct {
	Status int               `json:"status"`
	Data   SystemComputeData `json:"data"`
}

type SystemComputeData struct {
	SystemCompute
	UpdateTime int64 `json:"update_time"`
}

type SystemCompute struct {
	Load          Load           `json:"load"`
	CPU           []CPU          `json:"cpu"`
	CpuGlobal     CPU            `json:"global_cpu"`
	CPUCodename   string         `json:"cpu_codename"`
	ProcessStats  ProcessStats   `json:"process_stats"`
	CpuPressure   *Pressure      `json:"pressure"`
	BpfProgsStats *BpfProgsStats `json:"bpf_prog_stats"`
	Resctrl       *Resctrl       `json:"l3_mon"`
}

type SystemIO struct {
	DiskIo         []DiskIO        `json:"disk_io"`
	DiskUsage      []NodeDiskUsage `json:"disk_usage"`
	IoPressure     *Pressure       `json:"pressure"`
	ZramStat       []ZramStat      `json:"zram_stat"`
	Jbd2Info       []Jbd2Info      `json:"jbd_info"`
	WriteBackPages uint64          `json:"writeback_pages"`
}

type DiskUsage struct {
	MountPoint  string  `json:"mount_point"`
	FsType      string  `json:"filesystem_type"`
	DeviceName  string  `json:"device_name"`
	TotalBytes  uint64  `json:"total"`
	FreeBytes   uint64  `json:"free"`
	SpaceUsage  float64 `json:"usage"`
	TotalInodes uint64  `json:"total_inodes"`
	FreeInodes  uint64  `json:"free_inodes"`
}

type Jbd2Info struct {
	Partition              string `json:"partition"`
	TransacLockedTimeMS    uint64 `json:"trans_locked_ms"`
	TransacCommitAvgTimeUS uint64 `json:"avg_trans_commit_us"`
}

type ZramStat struct {
	Name           string `json:"name"`
	OrigDataSize   uint64 `json:"orig_data_size"`
	ComprDataSize  uint64 `json:"compr_data_size"`
	MemUsedTotal   uint64 `json:"mem_used_total"`
	MemLimit       uint64 `json:"mem_limit"`
	MemUsedMax     uint64 `json:"mem_used_max"`
	SamePages      uint64 `json:"same_pages"`
	PagesCompacted uint64 `json:"pages_compacted"`
	HugePages      uint64 `json:"huge_pages"`
}

type DiskIO struct {
	PrimaryDeviceID   int    `json:"primary_device_id"`
	SecondaryDeviceID int    `json:"secondary_device_id"`
	DeviceName        string `json:"device_name"`
	IoRead            uint64 `json:"io_read"`
	IoWrite           uint64 `json:"io_write"`
	IoBusy            uint64 `json:"io_busy"`
	DeviceType        string `json:"disk_type"`
	IoReadLat95       uint64 `json:"io_r_lat_95"`
	IoWriteLat95      uint64 `json:"io_w_lat_95"`
	IoReadLat90       uint64 `json:"io_r_lat_90"`
	IoWriteLat90      uint64 `json:"io_w_lat_90"`
	IoReadLat80       uint64 `json:"io_r_lat_80"`
	IoWriteLat80      uint64 `json:"io_w_lat_80"`
	WBTValue          int64  `json:"wbt_lat_usec"`
}

type BpfProgsStats struct {
	Stats      []BpfProgStats `json:"stats"`
	UpdateTime int64          `json:"update_time"`
}

type BpfProgStats struct {
	ID        uint32 `json:"id"`
	Name      []byte `json:"name"`
	RunTimeNS uint64 `json:"run_time_ns"`
	RunCount  uint64 `json:"run_cnt"`
	LoadTime  uint64 `json:"load_time"`
}

type ProcessStats struct {
	ProcessRunning uint64 `json:"procs_running"`
	ProcessBlocked uint64 `json:"procs_blocked"`
}

type Resctrl struct {
	L3Mon      []L3Mon `json:"l3mon"`
	UpdateTime int64   `json:"update_time"`
}

type L3Mon struct {
	ID                   int    `json:"id"`
	MbmTotalBytes        uint64 `json:"mbm_total_bytes"`
	MbmLocalBytes        uint64 `json:"mbm_local_bytes"`
	MbmVictimBytesPerSec uint64 `json:"mbm_victim_bytes_psec"`
	Llcoccupancy         uint64 `json:"llc_occupancy"`
}

type L3CacheBytesPS struct {
	NumaID           int
	MbmTotalBytesPS  uint64
	MbmLocalBytesPS  uint64
	MbmVictimBytesPS uint64
	MBMMaxBytesPS    uint64
}

type Load struct {
	One     float64 `json:"one"`
	Five    float64 `json:"five"`
	Fifteen float64 `json:"fifteen"`
}

type CPU struct {
	Name           string   `json:"name"`
	CPUUsage       float64  `json:"cpu_usage"`
	CPUSysUsage    float64  `json:"cpu_sys_usage"`
	CPUIowaitRatio float64  `json:"cpu_iowait_ratio"`
	CPUSchedWait   float64  `json:"cpu_sched_wait"`
	CpiData        *CpiData `json:"cpi_data"`
	CPUStealRatio  float64  `json:"cpu_steal_ratio"`
	CPUUsrRatio    float64  `json:"cpu_usr_ratio"`
	CPUIrqRatio    float64  `json:"cpu_irq_ratio"`
}

type CpiData struct {
	Cpi          float64 `json:"cpi"`
	Instructions float64 `json:"instructions"`
	Cycles       float64 `json:"cycles"`
	L3Misses     float64 `json:"l3_misses"`
	Utilization  float64 `json:"utilization"`
}

type MalachiteSystemMemoryResponse struct {
	Status int              `json:"status"`
	Data   SystemMemoryData `json:"data"`
}

type SystemMemoryData struct {
	SystemMemory
	UpdateTime int64 `json:"update_time"`
}

type SystemMemory struct {
	System      System    `json:"system"`
	Numa        []Numa    `json:"numa"`
	MemPressure *Pressure `json:"pressure"`
	ExtFrag     []ExtFrag `json:"extfrag"`
}

type System struct {
	MemTotal               uint64  `json:"mem_total"`
	MemFree                uint64  `json:"mem_free"`
	MemUsed                uint64  `json:"mem_used"`
	MemShm                 uint64  `json:"mem_shm"`
	MemAvailable           uint64  `json:"mem_available"`
	MemBuffers             uint64  `json:"mem_buffers"`
	MemPageCache           uint64  `json:"mem_page_cache"`
	MemSlabReclaimable     uint64  `json:"mem_slab_reclaimable"`
	MemDirtyPageCache      uint64  `json:"mem_dirty_page_cache"`
	MemWriteBackPageCache  uint64  `json:"mem_write_back_page_cache"`
	MemSwapTotal           uint64  `json:"mem_swap_total"`
	MemSwapFree            uint64  `json:"mem_swap_free"`
	MemUtil                float64 `json:"mem_util"`
	MemActiveAnon          uint64  `json:"mem_active_anon"`
	MemInactiveAnon        uint64  `json:"mem_inactive_anon"`
	MemActiveFile          uint64  `json:"mem_active_file"`
	MemInactiveFile        uint64  `json:"mem_inactive_file"`
	VMWatermarkScaleFactor uint64  `json:"vm_watermark_scale_factor"`
	VmstatPgstealKswapd    uint64  `json:"vmstat_pgsteal_kswapd"`
	MemSockTcp             uint64  `json:"mem_sock_tcp"`
	MemSockUdp             uint64  `json:"mem_sock_udp"`
	MemSockTcpLimit        uint64  `json:"mem_sock_tcp_limit"`
	MemSockUdpLimit        uint64  `json:"mem_sock_udp_limit"`
}

type CPUList struct {
	Meta  string `json:"meta"`
	Inner []int  `json:"inner"`
}

type Numa struct {
	ID                      int     `json:"id"`
	Path                    string  `json:"path"`
	CPUList                 CPUList `json:"cpu_list"`
	MemFree                 uint64  `json:"mem_free"`
	MemUsed                 uint64  `json:"mem_used"`
	MemTotal                uint64  `json:"mem_total"`
	MemShmem                uint64  `json:"mem_shmem"`
	MemAvailable            uint64  `json:"mem_available"`
	MemFilePages            uint64  `json:"mem_file_pages"`
	MemInactiveFile         uint64  `json:"mem_inactive_file"`
	MemMaxBandwidthMB       float64 `json:"mem_mx_bandwidth_mb"`
	MemReadBandwidthMB      float64 `json:"mem_read_bandwidth_mb"`
	MemReadLatency          float64 `json:"mem_read_latency"`
	MemTheoryMaxBandwidthMB float64 `json:"mem_theory_mx_bandwidth_mb"`
	MemWriteBandwidthMB     float64 `json:"mem_write_bandwidth_mb"`
	MemWriteLatency         float64 `json:"mem_write_latency"`
	AMDL3MissLatencyMax     float64 `json:"amd_l3_miss_latency_max"`
}

type ExtFrag struct {
	ID           int    `json:"id"`
	MemFragScore uint64 `json:"mem_frag_score"`
}

type Some struct {
	Avg10  float64 `json:"avg10"`
	Avg60  float64 `json:"avg60"`
	Avg300 float64 `json:"avg300"`
	Total  float64 `json:"total"`
}

type Full struct {
	Avg10  float64 `json:"avg10"`
	Avg60  float64 `json:"avg60"`
	Avg300 float64 `json:"avg300"`
	Total  float64 `json:"total"`
}

type Pressure struct {
	Some Some `json:"some"`
	Full Full `json:"full"`
}

type SystemEventData struct {
	Events     SystemEvent `json:"event_data"`
	UpdateTime int64       `json:"update_time"`
}

type SystemEvent struct {
	GeneralEvent SystemGeneralEvent `json:"gen"`
	IOEvent      SystemIOEvent      `json:"io"`
	FsEvent      SystemFsEvent      `json:"fs"`
	NetEvent     SystemNetEvent     `json:"net"`
	MemEvent     SystemMemEvent     `json:"mem"`
	SchedEvent   SystemSchedEvent   `json:"sched"`
}

type SystemGeneralEvent struct {
	SoftLockup uint64 `json:"soft_lockup"`
	RcuStall   uint64 `json:"rcu_stall"`
	BadPage    uint64 `json:"bad_page"`
	KernelWarn uint64 `json:"kernel_warn"`
	MceUC      uint64 `json:"mce_uc"`
	PanicTS    uint64 `json:"panic_ts"`
}

type SystemIOEvent struct {
	IOError uint64 `json:"io_error"`
}

type SystemFsEvent struct {
	Ext4Error uint64 `json:"ext4_error"`
	Ext4Abrt  uint64 `json:"ext4_abrt"`
}

type SystemNetEvent struct {
	XmitTimeout    uint64 `json:"xmit_timeout"`
	TcpBadCsum     uint64 `json:"tcp_bad_csum"`
	DevLinkDown    uint64 `json:"link_down"`
	ExceedBufLimit uint64 `json:"exceed_buf_limit"`
	RcvQueueFull   uint64 `json:"rcvqueue_full"`
}

type SystemMemEvent struct {
	AllocFailure uint64 `json:"alloc_failure"`
	TotalOOM     uint64 `json:"total_oom"`
	GlobalOOM    uint64 `json:"global_oom"`
}

type SystemSchedEvent struct {
	TaskHung uint64 `json:"hung_task"`
	CoreDump uint64 `json:"coredump"`
}

type SystemSensorData struct {
	Sensors SystemSensor `json:"sensors"`
}

type SystemSensor struct {
	TotalPower float64 `json:"total_power"`
	CPUPower   float64 `json:"cpu_power"`
	MemPower   float64 `json:"mem_power"`
	FanPower   float64 `json:"fan_power"`
	HDDPower   float64 `json:"hdd_power"`
	PSU0POut   float64 `json:"psu0_pout"`
	PSU1POut   float64 `json:"psu1_pout"`
	PSU0PIn    float64 `json:"psu0_pin"`
	PSU1PIn    float64 `json:"psu1_pin"`
	UpdateTime int64   `json:"update_time"`
}
