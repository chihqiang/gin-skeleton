package httpx

// controllers 全局变量，用于保存已注册的控制器实例
// 存储实现了 IController 接口的所有控制器
var (
	controllers = []IController{}
)

// GetControllers 返回当前已注册的所有控制器
// 通常在框架初始化（如 fx.Invoke）时使用，用于批量注册路由
// 返回值: 已注册的控制器列表
func GetControllers() []IController {
	return controllers
}

// RegisterController 注册一个控制器实例
// 通常在 init 函数或模块初始化中调用
// 所有注册的控制器会保存在 controllers 切片中
// 参数 controller: 要注册的控制器实例
func RegisterController(controller IController) {
	controllers = append(controllers, controller)
}
