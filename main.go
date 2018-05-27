package main

import (
	"fmt"
	"math/rand"
	"time"

	myutil "github.com/lioncruise/ant-colony-algorithm-golang/util"
	"strconv"
)

const (
	TASK_NUM int = 1000 // 任务数量
	NODE_NUM int = 100 // 节点数量

	ITERATE_NUM int = 100 // 迭代次数
	ANT_NUM     int = 100 // 蚂蚁数量
)

const (
	INT_MAX = int(^uint(0) >> 1)
	INT_MIN = ^INT_MAX
)

var (
	tasks [TASK_NUM]int // 任务集合
	nodes [NODE_NUM]int // 节点集合

	timeMatrix [TASK_NUM][NODE_NUM]float64 // 任务处理时间矩阵
	pheromoneMatrix [TASK_NUM][NODE_NUM]float64 // 信息素矩阵
	maxPheromoneMatrix [TASK_NUM]int
	criticalPointMatrix [TASK_NUM]int // 一次迭代中，随机分配的蚂蚁临界编号
	resultData = make([][]float64, ITERATE_NUM)

	p float64 = 0.5 // 每次迭代信息素衰减的比例
	q float64 = 2.0 // 每次经过，信息素增加的比例
)

// 初始化任务和节点集合
func initTaskAndNode() {
	fmt.Println("initTaskAndNode >>>>>>")

	for i := 0; i < TASK_NUM; i++ {
		tasks[i] = rand.Intn(91) + 10 // 产生一个10到100之间的随机数
	}

	for i := 0; i < NODE_NUM; i++ {
		nodes[i] = rand.Intn(91) + 10 // 产生一个10到100之间的随机数
	}
}

// 初始化信息素矩阵
func initPheromoneMatrix() {
	for i := 0; i < TASK_NUM; i++ {
		for j := 0; j < NODE_NUM; j++ {
			pheromoneMatrix[i][j] = 1.0
		}
	}
}

// 初始化处理时间矩阵
func initTimeMatrix() {
	for i := 0; i < TASK_NUM; i++ {
		for j := 0; j < NODE_NUM; j++ {
			timeMatrix[i][j] = float64(tasks[i]) / float64(nodes[j])
		}
	}
}

func initMatrix(m int, n int, defaultNum int) [][]int {
	// 使用动态分配内存的slice，slice使用的引用类型，array使用的值类型
	matrix := make([][]int, m)
	for i := 0; i < m; i++ {
		matrix[i] = make([]int, n)
		for j := 0; j < n; j++ {
			matrix[i][j] = defaultNum
		}
 	}

 	return matrix
}

func assignOneTask(antCount int, taskCount int) int {
	// 如果当前蚂蚁编号在临界点之前，则采用最大信息素的方式分配
	if antCount <= criticalPointMatrix[taskCount] {
		return maxPheromoneMatrix[taskCount]
	}

	return rand.Intn(NODE_NUM)
}

func callTimeOneIt(pathMatrixAllAnt [][][]int) []float64 {
	timeAllAnt := make([]float64, ANT_NUM)
	for antIndex := 0; antIndex < len(pathMatrixAllAnt); antIndex++ {
		pathMatrix := pathMatrixAllAnt[antIndex]

		var maxTime float64 = -1
		for nodeIndex := 0; nodeIndex < NODE_NUM; nodeIndex++ {
			var time float64 = 0
			for taskIndex := 0; taskIndex < TASK_NUM; taskIndex++ {
				if pathMatrix[taskIndex][nodeIndex] == 1 {
					time += timeMatrix[taskIndex][nodeIndex]
				}
			}

			if time > maxTime {
				maxTime = time
			}
		}

		timeAllAnt[antIndex] = maxTime
	}

	return timeAllAnt
}

// 更新信息素矩阵
func updatePheromoneMatrix(pathMatrixAllAnt [][][]int, timeArrayOneIt []float64) {
	// 所有信息素衰减
	for i := 0; i < TASK_NUM; i++ {
		for j := 0; j < NODE_NUM; j++ {
			pheromoneMatrix[i][j] *= p
		}
	}

	// 找出处理任务时间最短的蚂蚁的编号
	var minTime float64 = float64(INT_MAX)
	var minIndex int = -1
	for antIndex := 0; antIndex < ANT_NUM; antIndex++ {
		if timeArrayOneIt[antIndex] < minTime {
			minTime = timeArrayOneIt[antIndex]
			minIndex = antIndex
		}
	}

	// 将本次迭代中最优路径的信息素增加
	for taskIndex := 0; taskIndex < TASK_NUM; taskIndex++ {
		for nodeIndex := 0; nodeIndex < NODE_NUM; nodeIndex++ {
			if pathMatrixAllAnt[minIndex][taskIndex][nodeIndex] == 1 {
				pheromoneMatrix[taskIndex][nodeIndex] *= q
			}
		}
	}

	for taskIndex := 0; taskIndex < TASK_NUM; taskIndex++ {
		var maxPheromone float64 = pheromoneMatrix[taskIndex][0]
		var maxIndex int = 0
		var sumPheromone float64 = pheromoneMatrix[taskIndex][0]
		var isAllSame bool = true

		for nodeIndex := 1; nodeIndex < NODE_NUM; nodeIndex++ {
			if pheromoneMatrix[taskIndex][nodeIndex] > maxPheromone {
				maxPheromone = pheromoneMatrix[taskIndex][nodeIndex]
				maxIndex = nodeIndex
			}

			if pheromoneMatrix[taskIndex][nodeIndex] != pheromoneMatrix[taskIndex][nodeIndex - 1] {
				isAllSame = false
			}

			sumPheromone += pheromoneMatrix[taskIndex][nodeIndex]
		}

		if isAllSame {
			maxIndex = rand.Intn(NODE_NUM)
			maxPheromone = pheromoneMatrix[taskIndex][maxIndex]
		}

		maxPheromoneMatrix[taskIndex] = maxIndex

		criticalPointMatrix[taskIndex] = int(myutil.Round(float64(ANT_NUM) * (maxPheromone / sumPheromone)))
	}
}

func acaSearch() {
	for itCount := 0; itCount < ITERATE_NUM; itCount++ {
		pathMatrixAllAnt := make([][][]int, ANT_NUM)

		for antCount := 0; antCount < ANT_NUM; antCount++ {
			var pathMatrixOneAnt [][]int = initMatrix(TASK_NUM, NODE_NUM, 0)

			for taskCount := 0; taskCount < TASK_NUM; taskCount++ {
				nodeCount := assignOneTask(antCount, taskCount)
				pathMatrixOneAnt[taskCount][nodeCount] = 1
			}

			pathMatrixAllAnt[antCount] = pathMatrixOneAnt
		}

		var timeArrayOneIt []float64 = callTimeOneIt(pathMatrixAllAnt)
		resultData[itCount] = timeArrayOneIt

		updatePheromoneMatrix(pathMatrixAllAnt, timeArrayOneIt)
	}
}

func aca() {
	rand.Seed(time.Now().UnixNano()) // 一定要设置随机数的种子，默认的种子是1，不设置的话每次随机的结果都是一样的

	initTaskAndNode()
	initTimeMatrix()
	initPheromoneMatrix()
	acaSearch()
}

func main() {
	t1 := time.Now()
	aca()
	elapsed := time.Since(t1)

	for i := 0; i < ITERATE_NUM; i++ {
		fmt.Println("第" + strconv.Itoa(i) + "轮：", resultData[i])
	}

	fmt.Println("耗时: ", elapsed)
}
