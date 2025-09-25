package models

// PlexSessionResponse represents the Plex API response for active sessions
type PlexSessionResponse struct {
	MediaContainer struct {
		Size     int                `json:"size"`
		Metadata []PlexSessionMedia `json:"Metadata"`
	} `json:"MediaContainer"`
}

// PlexSessionMedia represents metadata for a Plex session
type PlexSessionMedia struct {
	TranscodeSession map[string]interface{} `json:"TranscodeSession"`
}

// CPUInfo contains processor information
type CPUInfo struct {
	Usage    string  `json:"usage"`
	Temp     string  `json:"temp"`
	TempDeg  float64 `json:"tempDeg"`
}

// NetCounters represents raw network counters
type NetCounters struct {
	RxBytes float64
	TxBytes float64
}

// CPUTimes represents CPU times for usage calculation
type CPUTimes struct {
	Idle  float64
	Total float64
}

// RAMInfo contains RAM memory information
type RAMInfo struct {
	Used       string  `json:"used"`
	Total      string  `json:"total"`
	Percent    string  `json:"percent"`
	PercentNum float64 `json:"percentNum"`
}

// DiskInfo contains storage information
type DiskInfo struct {
	Total      string  `json:"total"`
	Used       string  `json:"used"`
	Free       string  `json:"free"`
	Percent    string  `json:"percent"`
	PercentNum float64 `json:"percentNum"`
	MountPoint string  `json:"mountPoint"`
}

// NetTraffic represents formatted network traffic
type NetTraffic struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

// ZPoolVdev represents a ZFS vdev
type ZPoolVdev struct {
	Name    string   `json:"name"`
	Status  string   `json:"status"`
	Devices []string `json:"devices"`
}

// ZFSConfig contains ZFS configuration
type ZFSConfig struct {
	PoolName   string      `json:"poolName"`
	PoolStatus string      `json:"poolStatus"`
	DataVdevs  []ZPoolVdev `json:"dataVdevs"`
	CacheVdev  *ZPoolVdev  `json:"cacheVdev,omitempty"`
}

// ARCCache contains ZFS ARC cache information
type ARCCache struct {
	ARCSize       string  `json:"arcSize"`
	ARCMaxSize    string  `json:"arcMaxSize"`
	ARCHitRate    string  `json:"arcHitRate"`
	ARCHitRateNum float64 `json:"arcHitRateNum"`
	L2ARCSize     string  `json:"l2arcSize"`
	L2ARCHitRate  string  `json:"l2arcHitRate"`
}

// SystemInfo contains general system information
type SystemInfo struct {
	OS     string `json:"os"`
	Kernel string `json:"kernel"`
	CPU    string `json:"cpu"`
	Uptime string `json:"uptime"`
}

// DockerInfo contains Docker information
type DockerInfo struct {
	Containers int `json:"containers"`
	Images     int `json:"images"`
	Volumes    int `json:"volumes"`
}

// StreamingInfo contains streaming services information
type StreamingInfo struct {
	Films       int `json:"films"`
	Series      int `json:"series"`
	Animes      int `json:"animes"`
	Playing     int `json:"playing"`
	Transcoding int `json:"transcoding"`
}

// GlobalMetrics groups all metrics sent to the frontend
type GlobalMetrics struct {
	CPU       CPUInfo       `json:"cpu"`
	RAM       RAMInfo       `json:"ram"`
	Disk      DiskInfo      `json:"disk"`
	Net       NetTraffic    `json:"net"`
	ZFSConfig ZFSConfig     `json:"zfsConfig"`
	ARCCache  ARCCache      `json:"arcCache"`
	System    SystemInfo    `json:"system"`
	Docker    DockerInfo    `json:"docker"`
	Streaming StreamingInfo `json:"streaming"`
}