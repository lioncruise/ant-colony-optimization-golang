package mongo

type Node struct {
	// Id string `bson:"_id"`
	UID string `bson:"uid"`
	Cpu int `bson:"cpu"`
	Mem int `bson:"memory"`
	DiskSize int `bson:"disk_size"`
	MaxKeyInstance int `bson:"max_key_instances"`
	MatchTag []string `bson:"match_tag"`
	ExcludeTag []string `bson:"exclude_tag"`
}

type Pod struct {
	// Id string `bson:"_id"`
	UID string `bson:"uid"`
	Cpu int `bson:"cpu"`
	Mem int `bson:"memory"`
	DiskSize int `bson:"disk_size"`
	MatchTag []string `bson:"match_tag"`
	ExcludeTag []string `bson:"exclude_tag"`
	ExclusiveTag []string `bson:"exclusive_tag"`
	IsKeyInstance bool `bson:"is_key_instance"`
	AppGroup string `bson:"app_group"`
	HostId string `bson:"host_id"`
}