package conf

var (
	Conf     *Config
	AliMDUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation"
	AliWSUrl = "wss://dashscope.aliyuncs.com/api-ws/v1/realtime"
	//AliWSUrl = "ws://localhost:8765"
)

// 0: 面试官
var SysPrompt = []string{
	"你是一位专业的面试官，正在为职位面试应聘者。",
}

var FastMsg = []string{
	"上面是一位面试者的面试过程，请根据面试者的表现，给出面试评价与总结。",
}

var Role = []string{
	"system",
	"user",
	"assistant",
}
