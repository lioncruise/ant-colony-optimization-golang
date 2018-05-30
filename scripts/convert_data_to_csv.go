package scripts

import "github.com/lioncruise/ant-colony-algorithm-golang/mongo"
import (
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func convertPodDataToCSV() {
	podAll := mongo.GetAllPods()

	f, err := os.Create("E:/pod.csv")
	check(err)
	defer f.Close()

	for _, pod := range *podAll {
		// 只需要容器的CPU和内存即可
		f.WriteString(strconv.Itoa(pod.Cpu) + "," + strconv.Itoa(pod.Mem) + "\n")
	}
}

func convertNodeDataToCSV() {
	nodeAll := mongo.GetAllNodes()

	f, err := os.Create("E:/node.csv")
	check(err)
	defer f.Close()

	for _, node := range *nodeAll {
		// 只需要机器的CPU和内存即可
		f.WriteString(strconv.Itoa(node.Cpu) + "," + strconv.Itoa(node.Mem) + "\n")
	}
}

func ConvertDataToCSV() {
	convertNodeDataToCSV()
	convertPodDataToCSV()
}