package conf

var (
	Conf     *Config
	AliMDUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation"
)

// 0: 面试官
var SysPrompt = []string{
	"你是一位专业的面试官，正在为职位面试应聘者。",
}
