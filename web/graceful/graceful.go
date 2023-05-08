package graceful

// TODO 这里的逻辑需要优化
// 1. 如果提供到公共部分，那么这就是一个启动和关闭的模板
// 2. 需要一个channel在stop时用来通知所有启动的服务开始关闭回收资源,保存数据等等
// 3. 需要以goroutine的方式来启动所有服务，服务中就只管自己的逻辑，不再重新单独自己去起协程
// 4. 需要考虑服务之间的依赖关系，如果出现依赖，怎么办，怎么处理顺序问题(这需要重新考虑第三点)

var (
	engines = make(map[string]Engine, 0)
)

type Engine interface {
	Init()
	Run()
	Stop()
}

// Register 注册服务引擎
func Register(name string, engine Engine) {
	engines[name] = engine
}

func InitAll() {
	for _, engine := range engines {
		engine.Init()
	}
}

// StartAll 启动所有服务引擎
func StartAll() {
	for _, engine := range engines {
		engine.Run()
	}
}

// ShutdownAll 关闭所有服务引擎
func ShutdownAll() {
	for _, engine := range engines {
		engine.Stop()
	}
}
