package sdk

// Fields of ChubaoDB Table: raft
const (
	Raft_InodeId    = "Inode"
	Raft_ParentId   = "ParentId"
	Raft_InodeName  = "Name"
	Raft_OpType     = "_op"
	Raft_OpKey      = "_key"
	Raft_InsertTime = "_insert_time"
	Raft_VolumeName = "_volume"
	Raft_Term       = "_term"
	Raft_Index      = "_index"
	Raft_Pid        = "_partition_id"
	Raft_NodeIP     = "_node_ip"
)

// Fields of ChubaoDB Table: dentry
const (
	D_Inode       = "inode"
	D_Name        = "name"
	D_Path        = "path"
	D_Vol         = "vol_name"
	D_Pid         = "partition_id"
	D_ParentInode = "parent_inode"
	D_InsertTime  = "insert_time"
)
