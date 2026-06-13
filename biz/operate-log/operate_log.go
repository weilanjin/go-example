package operate_log

// 操作类型常量
const (
	OpAdd    = 1 // 创建
	OpUpdate = 2 // 更新
	OpDelete = 3 // 删除
)

type OperateLog struct {
	Id        int64  `json:"id"`         // 日志主键
	Operator  string `json:"operator"`   // 操作者名称
	OperateAt int64  `json:"operate_at"` // 操作时间，Unix时间戳
	IpAddress string `json:"ip_address"` // 操作者IP地址

	OpType     int    `json:"op_type"`     // 操作类型，1-创建，2-更新，3-删除
	ModuleName string `json:"module_name"` // 模块名称
	TargetId   string `json:"target_id"`   // 目标ID，例如被操作的对象的ID
	Content    string `json:"content"`     // 操作内容 json字符串，记录操作前后的数据等信息
	Remark     string `json:"remark"`      // 备注
}
