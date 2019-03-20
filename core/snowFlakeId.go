package core

import (
	"github.com/silenMing/etcdInleft/lib"
	"sync"
	"time"
)

const startpoch = 1514736000000 //开始时间戳

//机器位
type MachineInfo struct {
	workerId     uint8
	dataCenterId uint8
}

//snowId结构
type worker struct {
	lock sync.Mutex
	identify     uint
	timestamp    int64
	machineInfo  *MachineInfo
	sequence uint64
}

var (
	workerIdBit uint8
	dataCenterIdBit uint8
	sequenceBit uint8
	maxWorkerId uint64
	maxDataCenterId uint64
	lastTimestamp int64
	sequenceMask uint64
	sequence uint64
)

//初始化位数
func init() {
	workerIdBit = 5
	dataCenterIdBit = 5
	sequenceBit = 12
	maxWorkerId = -1 ^ (-1 << workerIdBit)
	maxDataCenterId = -1 ^ (-1 << dataCenterIdBit)
	sequenceMask = -1 ^ (-1 << sequenceBit) //序列掩码
	sequence = 0
}


var cfg *lib.EtcdConfig
var timestamp  int64

func (machineInfo *MachineInfo) createMachine() {

	if cfg.Get("versionId") == "" {
		machineInfo.workerId = 0
	} else {
		data := []byte(cfg.Get("versionId"))
		machineInfo.workerId = data[0]
	}
}

func (worker *worker) createSnowInfo()  {
	worker.lock.Lock()
	timestamp = getTimeNano()

	if timestamp < lastTimestamp{
		//时钟回退
		//Todo
	}

	//同一时间则排序
	if timestamp == lastTimestamp {
		sequence = (sequence +1) & sequenceMask
		if sequence == 0 {
			timestamp = nextMillis(lastTimestamp)
		}
	}else{
		sequence = 0
	}

	lastTimestamp = timestamp

	worker.timestamp = timestamp
	worker.sequence = sequence

	defer worker.lock.Unlock()
}

func nextMillis(time int64) int64{
	timestamp = getTimeNano()
	for  {
		if timestamp <= time {
			timestamp = getTimeNano()
		}
	}

	return timestamp
}

func getTimeNano() int64 {
	return time.Now().UnixNano() / 1e6
}







