package proto

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/chubaofs/chubaofs/util/auth"
	"github.com/chubaofs/chubaofs/util/config"
)

// For client
const (
	// Mandatory
	MountPoint int = iota
	VolName
	Owner
	Master
	// Optional
	LogDir
	WarnLogDir
	LogLevel
	ProfPort
	IcacheTimeout
	LookupValid
	AttrValid
	ReadRate
	WriteRate
	EnSyncWrite
	AutoInvalData
	Rdonly
	WriteCache
	KeepCache
	FollowerRead
	Authenticate
	ClientKey
	TicketHost
	EnableHTTPS
	CertFile
	TokenKey
	AccessKey
	SecretKey
	DisableDcache
	SubDir
	FsyncOnClose
	MaxCPUs
	EnableXattr

	MaxMountOption
)

// For server
const (
	MasterAddr       = "masterAddr"
	ListenPort       = "listen"
	ObjectNodeDomain = "objectNodeDomain"
)

type MountOption struct {
	keyword      string
	description  string
	cmdlineValue string
	value        interface{}
}

func (opt MountOption) String() string {
	return fmt.Sprintf("[%v] %T: %v", opt.keyword, opt.value, opt.value)
}

func NewMountOptions() []MountOption {
	opts := make([]MountOption, MaxMountOption)
	return opts
}

func InitMountOptions(opts []MountOption) {
	opts[MountPoint] = MountOption{"mountPoint", "Mount Point", "", ""}
	opts[VolName] = MountOption{"volName", "Volume Name", "", ""}
	opts[Owner] = MountOption{"owner", "Owner", "", ""}
	opts[Master] = MountOption{MasterAddr, "Master Address", "", ""}
	opts[LogDir] = MountOption{"logDir", "Log Path", "", ""}
	opts[WarnLogDir] = MountOption{"warnLogDir", "Warn Log Path", "", ""}
	opts[LogLevel] = MountOption{"logLevel", "Log Level", "", ""}
	opts[ProfPort] = MountOption{"profPort", "PProf Port", "", ""}
	opts[IcacheTimeout] = MountOption{"icacheTimeout", "Inode Cache Expiration Time", "", int64(-1)}
	opts[LookupValid] = MountOption{"lookupValid", "Lookup Valid Duration", "", int64(-1)}
	opts[AttrValid] = MountOption{"attrValid", "Attr Valid Duration", "", int64(-1)}
	opts[ReadRate] = MountOption{"readRate", "Read Rate Limit", "", int64(-1)}
	opts[WriteRate] = MountOption{"writeRate", "Write Rate Limit", "", int64(-1)}
	opts[EnSyncWrite] = MountOption{"enSyncWrite", "Enable Sync Write", "", int64(-1)}
	opts[AutoInvalData] = MountOption{"autoInvalData", "Auto Invalidate Data", "", int64(-1)}
	opts[Rdonly] = MountOption{"rdonly", "Mount as readonly", "", false}
	opts[WriteCache] = MountOption{"writecache", "Enable FUSE writecache feature", "", false}
	opts[KeepCache] = MountOption{"keepcache", "Enable FUSE keepcache feature", "", false}
	opts[FollowerRead] = MountOption{"followerRead", "Enable read from follower", "", false}

	opts[Authenticate] = MountOption{"authenticate", "Enable Authenticate", "", false}
	opts[ClientKey] = MountOption{"clientKey", "Client Key", "", ""}
	opts[TicketHost] = MountOption{"ticketHost", "Ticket Host", "", ""}
	opts[EnableHTTPS] = MountOption{"enableHTTPS", "Enable HTTPS", "", false}
	opts[CertFile] = MountOption{"certFile", "Cert File", "", ""}

	opts[TokenKey] = MountOption{"token", "Token Key", "", ""}
	opts[AccessKey] = MountOption{"accessKey", "Access Key", "", ""}
	opts[SecretKey] = MountOption{"secretKey", "Secret Key", "", ""}

	opts[DisableDcache] = MountOption{"disableDcache", "Disable Dentry Cache", "", false}
	opts[SubDir] = MountOption{"subdir", "Mount sub directory", "", ""}
	opts[FsyncOnClose] = MountOption{"fsyncOnClose", "Perform fsync upon file close", "", true}
	opts[MaxCPUs] = MountOption{"maxcpus", "The maximum number of CPUs that can be executing", "", int64(-1)}
	opts[EnableXattr] = MountOption{"enableXattr", "Enable xattr support", "", false}

	for i := 0; i < MaxMountOption; i++ {
		flag.StringVar(&opts[i].cmdlineValue, opts[i].keyword, "", opts[i].description)
	}
}

func ParseMountOptions(opts []MountOption, cfg *config.Config) {
	for i := 0; i < MaxMountOption; i++ {
		switch v := opts[i].value.(type) {
		case string:
			if opts[i].cmdlineValue != "" {
				opts[i].value = opts[i].cmdlineValue
			} else {
				if value, present := cfg.CheckAndGetString(opts[i].keyword); present {
					opts[i].value = value
				} else {
					opts[i].value = v
				}
			}
			fmt.Println(fmt.Sprintf("keyword[%v] value[%v] type[%T]", opts[i].keyword, opts[i].value, v))

		case int64:
			if opts[i].cmdlineValue != "" {
				opts[i].value = parseInt64(opts[i].cmdlineValue)
			} else {
				if value, present := cfg.CheckAndGetString(opts[i].keyword); present {
					opts[i].value = parseInt64(value)
				} else {
					opts[i].value = v
				}
			}
			fmt.Println(fmt.Sprintf("keyword[%v] value[%v] type[%T]", opts[i].keyword, opts[i].value, v))

		case bool:
			if opts[i].cmdlineValue != "" {
				opts[i].value = parseBool(opts[i].cmdlineValue)
			} else {
				if value, present := cfg.CheckAndGetBool(opts[i].keyword); present {
					opts[i].value = value
				} else {
					opts[i].value = v
				}
			}
			fmt.Println(fmt.Sprintf("keyword[%v] value[%v] type[%T]", opts[i].keyword, opts[i].value, v))

		default:
			fmt.Println(fmt.Sprintf("keyword[%v] unknown type[%T]", opts[i].keyword, v))
		}
	}
}

func parseInt64(s string) int64 {
	var ret int64 = -1

	if s != "" {
		val, err := strconv.Atoi(s)
		if err == nil {
			ret = int64(val)
		}
	}
	return ret
}

func parseBool(s string) bool {
	var ret bool = false

	if s == "true" {
		ret = true
	}
	return ret
}

func (opt *MountOption) GetString() string {
	val, ok := opt.value.(string)
	if !ok {
		return ""
	}
	return val
}

func (opt *MountOption) GetBool() bool {
	val, ok := opt.value.(bool)
	if !ok {
		return false
	}
	return val
}

func (opt *MountOption) GetInt64() int64 {
	val, ok := opt.value.(int64)
	if !ok {
		return int64(-1)
	}
	return val
}

type MountOptions struct {
	Config        *config.Config
	MountPoint    string
	Volname       string
	Owner         string
	Master        string
	Logpath       string
	Loglvl        string
	Profport      string
	IcacheTimeout int64
	LookupValid   int64
	AttrValid     int64
	ReadRate      int64
	WriteRate     int64
	EnSyncWrite   int64
	AutoInvalData int64
	UmpDatadir    string
	Rdonly        bool
	WriteCache    bool
	KeepCache     bool
	FollowerRead  bool
	Authenticate  bool
	TicketMess    auth.TicketMess
	TokenKey      string
	AccessKey     string
	SecretKey     string
	DisableDcache bool
	SubDir        string
	FsyncOnClose  bool
	MaxCPUs       int64
	EnableXattr   bool
}
