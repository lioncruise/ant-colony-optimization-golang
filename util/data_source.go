package util

import (
	"os"
	"encoding/csv"
	"strconv"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetSimplifiedNodeArray() *[]Node {
	f, err := os.Open("E:/node.csv")
	checkError(err)
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	nodeAll := make([]Node, len(lines))
	for i, line := range lines {
		nodeAll[i].Cpu, err = strconv.Atoi(line[0])
		nodeAll[i].Mem, err = strconv.Atoi(line[1])
	}

	return &nodeAll
}

func GetSimplifiedPodArray() *[]Pod {
	f, err := os.Open("E:/pod.csv")
	checkError(err)
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	podAll := make([]Pod, len(lines))
	for i, line := range lines {
		podAll[i].Cpu, err = strconv.Atoi(line[0])
		podAll[i].Mem, err = strconv.Atoi(line[1])
	}

	return &podAll
}
