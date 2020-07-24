package gather

import (
	"encoding/json"
	"github.com/google/uuid"
	"path"
	"strconv"
	"time"

	"github.com/chubaofs/chubaofs-tools/audit-daemon/sdk"
	. "github.com/chubaofs/chubaofs-tools/audit-daemon/util"
	"github.com/chubaofs/chubaofs-tools/audit-daemon/util/raft"
)

const rootInode = "1"

type dentryInfo struct {
	ParentInode string `json:"parent_inode"`
	Inode       string `json:"inode"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	VolName     string `json:"vol_name"`
	PartitionId string `json:"partition_id"`
	InsertTime  int64  `json:"insert_time"`
}

func InsertDentryInfo(parentID, inode, name, partitionID, vol string, dbConfig *sdk.DBConfig) {
	var (
		path string
		body []byte
		err  error
	)
	queryMap := make(map[string]interface{})
	queryMap[sdk.Raft_VolumeName] = vol
	queryMap[sdk.Raft_ParentId] = parentID
	queryMap[sdk.Raft_InodeName] = name
	if objs, _ := dbConfig.QueryAnd(queryMap, 10); len(objs) > 0 {
		LOG.Debugf("dentry exists in dentry table: vol[%v], parentID[%v], name[%v]", vol, parentID, name)
		return
	}
	if path, err = FindDentryPath(parentID, name, vol, dbConfig); err != nil {
		LOG.Errorf("find dentry path err: vol[%v], parentID[%v], name[%v]", vol, parentID, name)
		return
	}
	dInfo := &dentryInfo{
		ParentInode: parentID,
		Inode:       inode,
		Name:        name,
		Path:        path,
		VolName:     vol,
		PartitionId: partitionID,
		InsertTime:  time.Now().UnixNano() / 1e6,
	}
	if body, err = json.Marshal(dInfo); err != nil {
		LOG.Errorf("json marshal err: dentry[%v], err[%v]", dInfo, err)
		return
	}
	index := getUUid()
	if err = dbConfig.Insert(dbConfig.DentryTable, index, body); err != nil {
		LOG.Errorf("insert chubaodb err: table[%v], index[%v], body[%v]", dbConfig.DentryTable, index, dInfo)
		return
	}
	return
}

func getUUid() (index string) {
	var id uuid.UUID
	var err error
	if id, err = uuid.NewRandom(); err != nil {
		index = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		return
	}
	index = id.String()
	return
}

func FindDentryPath(parentID, name, vol string, dbc *sdk.DBConfig) (dentryPath string, err error) {
	dentryPath = name

	if parentID == rootInode {
		return
	}

	for {
		if parentID == "1" {
			break
		}

		// 1. search in dentry table
		queryDentryMap := make(map[string]interface{})
		queryDentryMap[sdk.D_Vol] = vol
		queryDentryMap[sdk.D_Inode] = parentID

		var obj interface{}
		if obj, err = dbc.QuerySortTop(dbc.DentryTable, queryDentryMap, sdk.D_InsertTime, sdk.DESC); err != nil {
			LOG.Errorf("query chubaodb err: table[%v] queryMap[%v]", dbc.DentryTable, queryDentryMap)
			return "", err
		}
		if obj != nil {
			dentryPath = path.Join(obj.(dentryInfo).Path, dentryPath)
			return
		}

		// 2. if not found, search in raft table
		queryRaftMap := make(map[string]interface{})
		queryRaftMap[sdk.Raft_OpType] = opFSMCreateDentry
		queryRaftMap[sdk.Raft_VolumeName] = vol
		queryRaftMap[sdk.Raft_InodeId] = parentID
		if obj, err = dbc.QuerySortTop(dbc.RaftTable, queryRaftMap, sdk.Raft_InsertTime, sdk.DESC); err != nil {
			LOG.Errorf("query chubaodb err: table[%v], queryMap[%v]", dbc.RaftTable, queryRaftMap)
			return "", err
		}
		if obj != nil {
			dentry := obj.(RaftItem).Data.(raft.Dentry)
			dentryPath = path.Join(dentry.Name, dentryPath)
			parentID = strconv.FormatUint(dentry.ParentId, 10)
			continue
		}

		// todo 3. if not found, search by metanode api
		// Note: Dentry that has been deleted does not exist in the map. So some file modification records may be lost.
		//var dentryMap map[string]*raft.Dentry
		//if dentryMap, err = sdk.GetAllDentry(metaAddr, partitionID); err != nil {	//todo partitionid 不一直是一个，使用metawrapper的接口获取
		//	LOG.Errorf("get all dentry from metanode err: addr[%v], pid[%v]", metaAddr, partitionID)
		//	return
		//}
	}
	return dentryPath, nil
}

func newMetaWrapper(masters []string, vol string) {
	//var metaConfig = &meta.MetaConfig{
	//	Volume:        vol,
	//	Masters:       masters,
	//	Authenticate:  false,
	//	ValidateOwner: false,
	//}
}
