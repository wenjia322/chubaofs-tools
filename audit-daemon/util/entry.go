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
	Fields  string
	DBAddr  string
	DBTable string
}

type ResponseSearch struct {
	Code int32
	Msg  string
	Hits []interface{}
}

type ResponseCDB struct {
	Code  int32
	Total int32
	Hits  []interface{}
	Info  CDBInfo
}

type CDBInfo struct {
	Success int32
	Error   int32
	Message string
}
