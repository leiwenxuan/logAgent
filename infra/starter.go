package infra

import "github.com/tietang/props/kvs"

const Keyprops  = "_conf"
//资源启动器上下文
// 用来在服务资源初始化、安装、启动和停止的生命周期中变量和对象的传递
type StarterContext map[string]interface{}

func (s StarterContext)Props() kvs.ConfigSource  {
	p := s[Keyprops]
	if p == nil {
		panic("配置文件还未初始化")
	}
	return p.(kvs.ConfigSource)
}

func (s StarterContext)SetProps(conf kvs.ConfigSource)  {
	s[Keyprops] = conf
}

type Starter interface {
	// 资源初始化
	Init(StarterContext)

	//安装资源
	Setup(StarterContext)

	// 启动资源
	Start(StarterContext)

	// 是否需要阻塞
	StartBlocking() bool

	//资源停止：
	// 通常在启动时遇到异常时或者启用远程管理时，用于释放资源和终止资源的使用，
	// 通常要优雅的释放，等待正在进行的任务继续，但不再接受新的任务
	Stop(StarterContext)
}

// 服务启动注册器
type starterRegister struct {
	nonBlockingStarters []Starter
	blockingStarters    []Starter
}

func (r *starterRegister) AllStarters() []Starter {
	starters := make([]Starter, 0)
	starters = append(starters, r.nonBlockingStarters...)
	starters = append(starters, r.blockingStarters...)
	return starters
}

// 注册器启动
func (r *starterRegister) Register(starter Starter) {
	if starter.StartBlocking() {
		r.blockingStarters = append(r.blockingStarters, starter)
	}else {
		r.nonBlockingStarters = append(r.nonBlockingStarters, starter)
	}
}
var StarterRegister *starterRegister = &starterRegister{}

type Starters []Starter
// 注册
func Register(starter Starter)  {
	StarterRegister.Register(starter)
}

func GetStarters() []Starter  {
	return StarterRegister.AllStarters()
}

//默认的空实现,方便资源启动器的实现
type BaseStarter struct {
}

func (s *BaseStarter) Init(ctx StarterContext)      {}
func (s *BaseStarter) Setup(ctx StarterContext)     {}
func (s *BaseStarter) Start(ctx StarterContext)     {}
func (s *BaseStarter) Stop(ctx StarterContext)      {}
func (s *BaseStarter) StartBlocking() bool          { return false }


