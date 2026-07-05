package debugger

var bufferedDebugInfo []DebugInfo = []DebugInfo{}

func AddDebugInfo(info DebugInfo) {
	bufferedDebugInfo = append(bufferedDebugInfo, info)
}

func GetDebugInfo() []DebugInfo {
	return bufferedDebugInfo
}

type DebugInfo struct {
	Message string
}
