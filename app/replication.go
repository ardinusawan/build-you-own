package main

type Replication interface {
	GetStatus() Status
}

type Status struct {
	Role                       string
	ConnectedSlaves            int32
	MasterReplid               string
	MasterReplOffset           int32
	SecondReplOffset           int32
	ReplBacklogActive          int32
	ReplBacklogSize            int64
	ReplBacklogFirstByteOffset int32
	ReplBacklogHistlen         int32
}

type ReplicationStorage struct {
	status Status
}

func (r *ReplicationStorage) GetStatus() Status {
	return r.status
}

func NewReplicationStorage() *ReplicationStorage {
	r := ReplicationStorage{
		status: Status{
			Role: "master",
		},
	}
	return &r
}
