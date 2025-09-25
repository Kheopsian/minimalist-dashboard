package models

// PlexSessionResponse représente la réponse de l'API Plex pour les sessions actives
type PlexSessionResponse struct {
	MediaContainer struct {
		Size     int                `json:"size"`
		Metadata []PlexSessionMedia `json:"Metadata"`
	} `json:"MediaContainer"`
}

// PlexSessionMedia représente les métadonnées d'une session Plex
type PlexSessionMedia struct {
	TranscodeSession map[string]interface{} `json:"TranscodeSession"`
}

// CPUInfo contient les informations sur le processeur
type CPUInfo struct {
	Usage    string  `json:"usage"`
	Temp     string  `json:"temp"`
	TempDeg  float64 `json:"tempDeg"`
}

// NetCounters représente les compteurs réseau bruts
type NetCounters struct {
	RxBytes float64
	TxBytes float64
}

// CPUTimes représente les temps CPU pour le calcul d'utilisation
type CPUTimes struct {
	Idle  float64
	Total float64
}

// RAMInfo contient les informations sur la mémoire RAM
type RAMInfo struct {
	Used       string  `json:"used"`
	Total      string  `json:"total"`
	Percent    string  `json:"percent"`
	PercentNum float64 `json:"percentNum"`
}

// DiskInfo contient les informations sur le stockage
type DiskInfo struct {
	Total      string  `json:"total"`
	Used       string  `json:"used"`
	Free       string  `json:"free"`
	Percent    string  `json:"percent"`
	PercentNum float64 `json:"percentNum"`
	MountPoint string  `json:"mountPoint"`
}

// NetTraffic représente le trafic réseau formaté
type NetTraffic struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

// ZPoolVdev représente un vdev ZFS
type ZPoolVdev struct {
	Name    string   `json:"name"`
	Status  string   `json:"status"`
	Devices []string `json:"devices"`
}

// ZFSConfig contient la configuration ZFS
type ZFSConfig struct {
	PoolName   string      `json:"poolName"`
	PoolStatus string      `json:"poolStatus"`
	DataVdevs  []ZPoolVdev `json:"dataVdevs"`
	CacheVdev  *ZPoolVdev  `json:"cacheVdev,omitempty"`
}

// ARCCache contient les informations sur le cache ARC ZFS
type ARCCache struct {
	ARCSize       string  `json:"arcSize"`
	ARCMaxSize    string  `json:"arcMaxSize"`
	ARCHitRate    string  `json:"arcHitRate"`
	ARCHitRateNum float64 `json:"arcHitRateNum"`
	L2ARCSize     string  `json:"l2arcSize"`
	L2ARCHitRate  string  `json:"l2arcHitRate"`
}

// SystemInfo contient les informations système générales
type SystemInfo struct {
	OS     string `json:"os"`
	Kernel string `json:"kernel"`
	CPU    string `json:"cpu"`
	Uptime string `json:"uptime"`
}

// DockerInfo contient les informations sur Docker
type DockerInfo struct {
	Containers int `json:"containers"`
	Images     int `json:"images"`
	Volumes    int `json:"volumes"`
}

// StreamingInfo contient les informations sur les services de streaming
type StreamingInfo struct {
	Films       int `json:"films"`
	Series      int `json:"series"`
	Animes      int `json:"animes"`
	Playing     int `json:"playing"`
	Transcoding int `json:"transcoding"`
}

// GlobalMetrics regroupe toutes les métriques envoyées au frontend
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