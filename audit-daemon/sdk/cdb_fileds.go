package sdk

// Fields of ChubaoDB Table: raft
const (
	Raft_InodeId    = "Inode"
	Raft_ParentId   = "ParentId"
	Raft_InodeName  = "Name"
	Raft_OpType     = "r_op"
	Raft_OpKey      = "r_key"
	Raft_InsertTime = "r_insert_time"
	Raft_VolumeName = "r_volume"
	Raft_Term       = "r_term"
	Raft_Index      = "r_index"
	Raft_Pid        = "r_partition_id"
	Raft_NodeIP     = "r_node_ip"
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
