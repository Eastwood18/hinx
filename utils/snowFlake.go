package utils

import (
	"errors"
	"hinx/utils/commandLine/args"
	"sync"
	"time"
)

const (
	workerIDBits = uint64(10) // 10bit 工作机器ID中的 5bit workerID
	//dataCenterIDBits = uint64(5) // 10 bit 工作机器ID中的 5bit dataCenterID
	sequenceBits = uint64(12)

	maxWorkerID = int64(-1) ^ (int64(-1) << workerIDBits) //节点ID的最大值 用于防止溢出
	//maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
	maxSequence = int64(-1) ^ (int64(-1) << sequenceBits)

	timeLeft = uint8(22) // timeLeft = workerIDBits + sequenceBits // 时间戳向左偏移量
	//dataLeft = uint8(17) // dataLeft = dataCenterIDBits + sequenceBits
	workLeft = uint8(12) // workLeft = sequenceBits // 节点IDx向左偏移量

	twepoch = int64(1589923200000) // 常量时间戳(毫秒) 2020-05-20 08:00:00 +0800 CST
)

type worker struct {
	mu        sync.Mutex
	LastStamp int64 // 记录上一次ID的时间戳
	WorkerID  int64 // 该节点的ID
	//DataCenterID int64 // 该节点的 数据中心ID
	Sequence int64 // 当前毫秒已经生成的ID序列号(从0 开始累加) 1毫秒内最多生成4096个ID
}

var w *worker = new(worker)

func newWorker(workerID uint) *worker {

	workerID = workerID & (1<<timeLeft - 1)
	return &worker{
		WorkerID:  int64(workerID),
		LastStamp: 0,
		Sequence:  0,
	}
}
func (w *worker) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (w *worker) nextID() (uint64, error) {
	timeStamp := w.getMilliSeconds()
	if timeStamp < w.LastStamp {
		return 0, errors.New("time is moving backwards,waiting until")
	}

	if w.LastStamp == timeStamp {

		w.Sequence = (w.Sequence + 1) & maxSequence

		if w.Sequence == 0 {
			for timeStamp <= w.LastStamp {
				timeStamp = w.getMilliSeconds()
			}
		}
	} else {
		w.Sequence = 0
	}

	w.LastStamp = timeStamp
	id := ((timeStamp - twepoch) << timeLeft) |
		(w.WorkerID << workLeft) |
		w.Sequence

	return uint64(id), nil
}

func init() {
	w = newWorker(args.Args.MachineCode)

}

func GetSnowFlakeID() (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.nextID()
}
