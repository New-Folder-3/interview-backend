package conf

var (
	Conf     *Config
	AliMDUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation"
	AliWSUrl = "wss://dashscope.aliyuncs.com/api-ws/v1/realtime"
	//AliWSUrl = "ws://localhost:8765"
)

// 0: 面试官
var PromptTemplate = []string{
	"你是一位专业的面试官，正在为职位面试应聘者。",
	`你是一位专业的面试官，上面是一位面试者的面试过程，请根据面试者的表现，直接给出硬技能，软技能，潜力，自信心，发展潜力，岗位适配度的评分，满分为十分，其中岗位为%s。按照顺序直接输出五个一位小数，中间用空格分开，不要给出其他内容。`,
	`你是一位专业的面试官，上面是一位面试者的面试过程，请根据面试者的表现，在%s这些关键词中，给出面试者的关键词，按照顺序直接输出关键词，中间用空格分开，不要给出其他内容。`,
}

var Keywords = []string{
	"团队合作",
}

const (
	InterviewerPrompt = 0 // 面试官
	DimensionPrompt   = 1 // 维度评分
	DimensionCount    = 6 // 维度评分长度
	KeywordsPrompt    = 2 // 关键词提取
)

const (
	DefaultJob = "任何职位" // 默认职位
	DefaultAge = 18     // 默认年龄
)

var Role = []string{
	"system",
	"user",
	"assistant",
}
