package container

import (
	"github.com/lioncruise/ant-colony-algorithm-golang/util"
	"fmt"
	"math/rand"
	"time"
)

var (
	podArray  = *util.GetSimplifiedPodArray()  // 容器数组
	nodeArray = *util.GetSimplifiedNodeArray() // 机器数组

	// availableNodeArray = nodeArray // 这样是拷贝的指针的值，两个指针指向同一个对象
	availableNodeArray = make([]util.Node, nodeNum)

	podNum = len(podArray) // 34596
	nodeNum = len(nodeArray)  // 4338

	pheromoneMatrix = make([][]float64, podNum) // 信息素矩阵(记录每个pod)，在10000个pod 1000个node的情况下数组大小是80MB
	maxPheromoneMap = make([]int, podNum)
	criticalPointMatrix = make([]int, podNum)

	unscheduledPods = make(map[int]int)

	iteratorNum = 10
	antNum = 100

	p = 0.5 // 每次迭代信息素衰减比例
	q = 2.0 // 每次经过信息素增加的比例
)

func CheckResources() {
	podCpuSum := 0
	podMemSum := 0
	for i := 0; i < podNum; i++ {
		podCpuSum += podArray[i].Cpu
		podMemSum += podArray[i].Mem
	}

	nodeCpuSum := 0
	nodeMemSum := 0
	for i := 0; i < nodeNum; i++ {
		nodeCpuSum += nodeArray[i].Cpu
		nodeMemSum += nodeArray[i].Mem
	}

	fmt.Println("pod cpu sum:", podCpuSum)
	fmt.Println("node cpu sum:", nodeCpuSum)
	fmt.Println("pod memory sum:", podMemSum)
	fmt.Println("node memory sum:", nodeMemSum)
	/*
	目前第一版数据容器申请的CPU、内存大于所有机器可用的CPU、内存
		pod cpu sum: 312058
		node cpu sum: 239560
		pod memory sum 978734258
		node memory sum 913737606
	 */

	if nodeCpuSum < podCpuSum {
		fmt.Println("cpu not enough!")
	}
	if nodeMemSum < podMemSum {
		fmt.Println("memory not enough!")
	}
}

func podResourceSum() (int, int) {
	podCpuSum := 0
	podMemSum := 0
	for i := 0; i < podNum; i++ {
		podCpuSum += podArray[i].Cpu
		podMemSum += podArray[i].Mem
	}
	return podCpuSum, podMemSum
}

// 第一次版数据有问题，将容器的CPU、内存除2
func NormolizePodArray() {
	for i := 0; i < podNum; i++ {
		if podArray[i].Cpu == 1 {

		} else {
			podArray[i].Cpu = podArray[i].Cpu / 2
		}

		podArray[i].Mem = podArray[i].Mem / 2
	}
}

func resetAvailableNodeArray() {
	copy(availableNodeArray, nodeArray)
}

func initPheromoneMatrix() {
	for i := 0; i < podNum; i++ {
		pheromoneMatrix[i] = make([]float64, nodeNum)
		for j := 0; j < nodeNum; j++ {
			pheromoneMatrix[i][j] = 1.0
		}
	}
}

func assignOnePod(antCount int, podCount int) int {
	pod := podArray[podCount]

	if antCount <= criticalPointMatrix[podCount] {
		nodeIndex := maxPheromoneMap[podCount]
		node := &availableNodeArray[nodeIndex]
		if node.Cpu >= pod.Cpu && node.Mem >= pod.Mem {
			return nodeIndex
		}
	}

	nodeIndex := rand.Intn(nodeNum)
	node := &availableNodeArray[nodeIndex]
	retryCount := 0
	// 随机重试3次
	for retryCount < 3 && (pod.Cpu > node.Cpu || pod.Mem > node.Mem) {
		nodeIndex := rand.Intn(nodeNum)
		node = &availableNodeArray[nodeIndex]

		if pod.Cpu <= node.Cpu && pod.Mem <= node.Mem {
			return nodeIndex
		}

		retryCount++
	}

	for i := 0; i < nodeNum; i++ {
		node = &availableNodeArray[i]

		if pod.Cpu <= node.Cpu && pod.Mem <= node.Mem {
			return i
		}
	}

	// 返回0表示没有调度pod到一个node上
	return -1
}

func updatePheromoneMatrix(minPathOneAnt []int) {
	// 将所有信息素衰减
	for i := 0; i < podNum; i++ {
		for j := 0; j < nodeNum; j++ {
			pheromoneMatrix[i][j] *= p
		}
	}

	// 将本次迭代中最优路径的信息素增加
	for podIndex := 0; podIndex < podNum; podIndex++ {
		nodeIndex := minPathOneAnt[podIndex]
		pheromoneMatrix[podIndex][nodeIndex] *= q
	}

	for podIndex := 0; podIndex < podNum; podIndex++ {
		maxPheromone := pheromoneMatrix[podIndex][0]
		maxIndex := 0
		sumPheromone := pheromoneMatrix[podIndex][0]
		isAllSame := true

		for nodeIndex := 1; nodeIndex < nodeNum; nodeIndex++ {
			if pheromoneMatrix[podIndex][nodeIndex] > maxPheromone {
				maxPheromone = pheromoneMatrix[podIndex][nodeIndex]
				maxIndex = nodeIndex
			}

			if pheromoneMatrix[podIndex][nodeIndex] != pheromoneMatrix[podIndex][nodeIndex - 1] {
				isAllSame = false
			}

			sumPheromone += pheromoneMatrix[podIndex][nodeIndex]
		}

		if isAllSame == true {
			maxIndex = rand.Intn(nodeNum)
			maxPheromone = pheromoneMatrix[podIndex][maxIndex]
		}

		maxPheromoneMap[podIndex] = maxIndex

		criticalPointMatrix[podIndex] = int(util.Round(float64(antNum) * (maxPheromone / sumPheromone)))
	}

}

func acaSearch() {
	podCpuSum, podMemSum := podResourceSum()
	for itCount := 0; itCount < iteratorNum; itCount++ {

		minNodeNum := 10000
		var minPathOneAnt []int

		for antCount := 0; antCount < antNum; antCount++ {
			resetAvailableNodeArray() // 重置可用资源数组
			pathOneAnt := make([]int, podNum) // 重置当前蚂蚁的路径
			unscheduledPods = make(map[int]int) // 重置未调度的pod的数组

			hasPodUnscheduled := false
			for podCount := 0; podCount < podNum; podCount++ {
				nodeCount := assignOnePod(antCount, podCount)
				if nodeCount >= 0 {
					pathOneAnt[podCount] = nodeCount
					node := &availableNodeArray[nodeCount]
					node.Cpu = node.Cpu - podArray[podCount].Cpu
					node.Mem = node.Mem - podArray[podCount].Mem
				} else {
					unscheduledPods[podCount] = -1
					pathOneAnt[podCount] = -1 // -1表示没有调度该pod
					hasPodUnscheduled = true
				}

			}

			// 如果当前路径中有pod没有调度，不参与比较
			if hasPodUnscheduled == false {
				var nodeSet = make(map[int]int)// 使用map实现set
				for i := 0; i < podNum; i++ {
					nodeSet[pathOneAnt[i]] = 0
				}

				if len(nodeSet) < minNodeNum {
					minNodeNum = len(nodeSet)
					minPathOneAnt = pathOneAnt // pathOneAnt中如果含有-1的话，不能将其赋值给minPathOneAnt，否则在后面更新信息素可能会数组越界
				}
			}
		}

		fmt.Println("第", itCount + 1, "轮最小机器数:", minNodeNum)

		// 计算资源利用率
		costNodeCpuSum := 0
		costNodeMemSum := 0
		usedNodeSet := make(map[int]int)
		for i := 0; i < podNum; i++ {
			nodeIndex := minPathOneAnt[i]
			usedNodeSet[nodeIndex] = 1
		}

		for nodeIndex := range usedNodeSet {
			costNodeCpuSum += nodeArray[nodeIndex].Cpu
			costNodeMemSum += nodeArray[nodeIndex].Mem
		}
		fmt.Println("第", itCount + 1, "轮最小机器数结果下的CPU利用率:", float64(podCpuSum) / float64(costNodeCpuSum))
		fmt.Println("第", itCount + 1, "轮最小机器数结果下的内存利用率:", float64(podMemSum) / float64(costNodeMemSum))
		fmt.Println("======")

		updatePheromoneMatrix(minPathOneAnt)
	}
}

func Aca() {
	rand.Seed(time.Now().UnixNano()) // 设置随机种子

	t1 := time.Now()

	NormolizePodArray()
	resetAvailableNodeArray()
	initPheromoneMatrix()
	acaSearch()

	elapsed := time.Since(t1)

	fmt.Println("耗时: ", elapsed)
}