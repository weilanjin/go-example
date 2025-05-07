// protobuf 类型使用注意
// 1. int32 等类型使用时，对于负值编码效率不高. 注意 int32 等类型在 proto3 中已经被废弃，使用 sint32 等代替.
// 2. protobuf 是以编号来定义字段标识的. 注意 服务端和客户端 字段编号 必须一致.
package proto
