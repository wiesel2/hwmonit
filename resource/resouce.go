package resource

type ResourceType int

const (
	RTCPU ResourceType = iota // CPU
	RTMEM
	RTDISK
	RTNET
)

type Collector interface {
	Run() ResourceResult
	Stop()
}

type Resource struct {
	Name       string
	ResChan    chan [20]ResourceResult
	LastResult ResourceResult
	RrcType    ResourceType
	c          Collector
}

func (r *Resource) Run() {

}

type CPU struct{}

func (c *CPU) Run() ResourceResult {

	// 获取cpu占用时间
	// 获取cpu分配的周期
	if IsOnDocker() {
		// 重docker文件读取
	} else {
	}

	return ResourceResult{}

}
func (c *CPU) Stop() {

}

type MEM struct{}

type DISK struct{}

type NET struct{}
