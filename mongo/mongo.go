package mongo

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

const (
	MONGODB_URL = "47.104.16.133:27017"
)

func GetAllPods() *[]Pod {
	session, err := mgo.Dial(MONGODB_URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	db := session.DB("cluster") // 数据库名称
	c := db.C("newinstanceinfo") // 集合名称（集合相当于关系型数据库中的二维表）

	var podAll []Pod
	err = c.Find(bson.M{}).All(&podAll)

	// fmt.Println(podAll[0]) // 格式：{cybertron_0RE9NZW 4 8192 61440 [sigma_public] [] [cybertronhost] false cybertronhost }
	// fmt.Println("一共有", len(podAll), "个容器")  // 第一版数据一共有 34596 个容器

	return &podAll
}

func GetAllNodes() *[]Node {
	//创建连接
	session, err := mgo.Dial(MONGODB_URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	db := session.DB("cluster") // 数据库名称
	c := db.C("newnodeinfo") // 集合名称（集合相当于关系型数据库中的二维表）

	var nodeAll []Node
	err = c.Find(bson.M{}).All(&nodeAll)

	for i := 0; i < len(nodeAll); i++ {
		fmt.Println(nodeAll[i])
	}

	// fmt.Println(nodeAll[0]) // 格式：{9b849f22-cac0-4af3-91b6-eef22cb78d34 88 354544 863408 9 [sigma_public alios7u2] []}
	// fmt.Println("一共有", len(nodeAll), "个节点") // 第一版数据一共有4338个节点

	return &nodeAll
}