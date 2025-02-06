package opentelemetry

// OpenTelemetry 是一个开放标准, 用于跟踪和监控分布式系统.
// API 和 SDK 集合。
//
// OpenTelemetry 演进历史
// 1. Google 2010年发布 Dapper 论文, 介绍了分布式追踪的概念, 是分布式链路追踪的开端。
// 2. 2012 年 Twitter 开源 Zipkin, 用于跟踪分布式系统中的请求。
// 3. 2015 年 Uber 开源 Jaeger, 用于跟踪分布式系统中的请求。
// 4. 2015 年 OpenTracing 项目被 CNCF 接受为它的第三个托管项目，致力于标准化跨组件的分布式链路追踪。
// 5. 2017 年 Google 将内部的 Census 项目开源，随后 OpenCensus 在社区中流行起来。
// 6. 2019 年 OpenTracing 和 OpenCensus 合并为 OpenTelemetry。
// 7. 2021 年 OpenTelemetry 1.0 发布。 为客户端的链路追踪部分提供了稳定性保证。
// 8. 2023 年 OpenTelemetry 三个基本的功能，链路追踪、指标和日志。
//
// 每个网络调用都会被捕获并表示为一个跨度 span
// 分布式链路追踪工具将唯一的链路追踪上下文 traceID 插入到每个请求的标头中。
//
// OpenTelemetry 也称 OTel
// 用于检测、生成、收集和导出遥测数据。链路追踪（traces）、指标（metrics）和日志（logs）。
// - trace: 在分布式应用程序中的完整请求链路信息。
// - 指标: 在运行时捕获的服务指标，应用程序和请求指标是可用性和性能的重要指标。
// - 日志：系统或应用程序在特定时间点发生的事件的文本记录。
// - baggage：在信号之间传递的上下文信息。
