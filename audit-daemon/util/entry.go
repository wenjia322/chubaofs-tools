package util

import (
	"time"
)

type RequestGetContent struct {
	Dir     string
	Name    string
	Pattern string
	Inode   uint64
	Start   int64
}

type RequestListFile struct {
	Dir     string
	Pattern string
}

type RequestListDir struct {
	Dir       string
	Exclusion string
}

type RequestCommand struct {
	Dir     string
	Command string
	LimitMB int // the size of result
}

type FileInfo struct {
	Inode uint64
	Name  string
	Size  int64
	Time  time.Time
}

type MachineState struct {
	Ip     string
	Time   time.Time
	Cpu    int32
	Memory int32
}

type RequestForwardCmdReq struct {
	AddrList []string
	Command  string
	LimitMB  int // the size of result
}

type ResponseForwardCmdReq struct {
	Code    int32
	Msg     string
	Results []string
}

type RequestSearch struct {
	Query   string
	Size    int8
	DBAddr  string
	DBTable string
}

type ResponseSearch struct {
	Code  int32
	Msg   string
	Total int32
	Hits  []*HitItem
}

type ResponseCDB struct {
	Code  int32
	Total int32
	Hits  []*HitItem
	Info  *CDBInfo
}

type CDBInfo struct {
	Success int32
	Error   int32
	Message string
}

type HitItem struct {
	Score float64
	Doc   *HitInfo
}

type HitInfo struct {
	Id      string      `json:"_id"`
	SortKey string      `json:"_sort_key"`
	Version int         `json:"_version"`
	Source  interface{} `json:"_source"`
}

type RaftItem struct {
	Op          int8   `json:"_op"`
	Key         string `json:"_key"`
	PartitionId string `json:"_partition_id"`
	VolName     string `json:"_volume"`
	NodeIP      string `json:"_node_ip"`
	//Crc			int64	`json:"_crc"`
	//DataSize	int		`json:"_dataSize"`
	Index string `json:"_index"`
	//OpType		int		`json:"_opType"`
	//RecType		int		`json:"_recType"`
	Term       int   `json:"_term"`
	InsertTime int64 `json:"_insert_time"`
	Data       interface{}
}
