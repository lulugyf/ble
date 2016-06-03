package q

// 预设的属性名称
const (
	/**
	 * 消息类型
	 */
	TYPE = "type"
	/**
	 * 主题
	 */
	TOPIC = "topic"
	/**
	 * 消息组<br/>
	 * 用于限定同一组消息的消费顺序
	 */
	GROUP = "group"
	/**
	 * 消息组<br/>
	 * 用于限定同一组消息的消费顺序
	 */
	PRIORITYNAME = "priority-name"

	/**
	 * 优先级<br/>
	 * 用于限定不同组消息的消费顺序
	 */
	PRIORITY = "priority"
	/**
	 * 客户端标识<br/>
	 * 用于标识客户端来源
	 */
	CLIENT_ID = "client-id"
	/**
	 * 目标主题
	 */
	TARGET_TOPIC = "target-topic"
	/**
	 * Broker地址列表
	 */
	ADDRESS = "address"
	/**
	 * 消息过期时间，相对于1970-1-1 00:00:00 的ms值
	 */
	EXPIRE_TIME = "expire-time"
	/**
	 * 消息生效时间，相对于1970-1-1 00:00:00 的ms值
	 */
	EFFECTIVE_TIME = "effective-time"
	/**
	 * 是否压缩
	 */
	COMPRESS = "compress"
	/**
	 * 是否Rest方式压缩
	 */
	REST_COMPRESS = "rest-compress"
	/**
	 * 是否压缩
	 */
	ENCRYPT = "encrypt"
	/**
	 * 生产者重试次数
	 */
	PRODUCER_RETRY = "producer-retry"
	/**
	 * 中间件重试次数
	 */
	BROKER_RETRY = "broker-retry"
	/**
	 * 消费者重试次数
	 */
	CONSUMER_RETRY = "consumer-retry"
	/**
	 * 提交时间
	 */
	COMMIT_TIME = "commit-time"
	/**
	 * 消息的唯一标识
	 */
	MESSAGE_ID = "message-id"
	/**
	 * 消息的唯一标识
	 */
	BATCH_MESSAGE_ID = "batch-message-id"
	/**
	 * 状态码<br/>
	 * 用于承载发送服务端处理结果码
	 */
	RESULT_CODE = "result-code"
	/**
	 * PULL消息状态码<br/>
	 * 用于承载发送消费者处理结果码
	 */
	PULL_CODE = "pull-code"
	/**
	 * 状态码、补充描述<br/>
	 * 通常情况下，状态码对应的状态描述是固定的，可能会出现不同的原因导致了相同状态码，此字段用于补充描述详细信息
	 */
	CODE_DESCRIPTION = "code-description"
	/**
	 * 消费者消费消息的预期时间，超出预期时间认为消费超时, Unixs时间戳记， 绝对时间
	 */
	PROCESSING_TIME = "processing-time"
	/**
	 * 消费者失败时要求过一定时间（单位秒）后再重新处理，用于标识出不立即处理
	 */
	RETRY_AFTER = "retry-after"
	/**
	 * 自定义流水号
	 */
	CUSTOM_SERIAL = "serial"

	/**
	 * 批量拉取的消息最大值
	 */
	PAGE_SIZE = "page-size"
	/**
	 * 客户端连接地址
	 */
	REMOTE_ADDRESS = "remote-address"
	/**
	 * 计算出当前目标主题的属性名
	 */
	CURRENT_PROPERTY_KEY = "current-property-key"
	/**
	 * 计算出当前目标主题的属性值（配置值，非消息所带属性值）
	 */
	CURRENT_PROPERTY_VALUE = "current-property-value"
	/**
	 * 当前主题的消费结果需要发到的主题上
	 */
	REPLY_TO = "reply-to"
	/**
	 * 标识被哪个消费者消费了
	 */
	CONSUMED_BY = "consumed-by"
	/**
	 * 生产者的消息id
	 */
	PRODUCER_MESSAGE_ID = "producer-message-id"
	/**
	 * 上一次发送产生的消息id，用于消息重发时的历史回溯
	 */
	LAST_MESSAGE_ID = "last-message-id"
)

// MessageType
const (
	/* 心跳：仅消息头 */
	BREAKHEART = 0
	/**
	 * 查询，仅消息头。<br/>
	 * 消息头， (*)表示可选：
	 * <ul>
	 * <li>客户端标识：{@link PropertyOption#CLIENT_ID}</li>
	 * </ul>
	 */
	QUERY = 1
	/**
	 * 应答消息，仅消息头。<br/>
	 * 消息头， (*)表示可选：
	 * <ul>
	 * <li>应答返回码：{@link PropertyOption#RESULT_CODE}</li>
	 * <li>*应答返回描述：{@link PropertyOption#RESULT_DESCRIPTION}</li>
	 * <li>*代理地址：{@link PropertyOption#ADDRESS}，当应答返回码为 {@link ResultCode#OK} 且是针对
	 * {@link MessageType#QUERY} 应答时必选</li>
	 * <li>*消息标识：{@link PropertyOption#MESSAGE_ID}，当应答返回码为 {@link ResultCode#OK} 且是针对
	 * {@link MessageType#SEND} 应答时必选</li>
	 * <li>*客户端流水：{@link PropertyOption#CUSTOM_SERIAL}，当应答返回码为 {@link ResultCode#OK} 且是针对
	 * {@link MessageType#SEND} 应答时可选返回原 {@link MessageType#SEND} 消息所带该值</li>
	 * </ul>
	 */
	ANSWER = 2
	/**
	 * 发送。<br/>
	 * 消息头列表， (*)表示可选：
	 * <ul>
	 * <li>客户端标识：{@link PropertyOption#CLIENT_ID}</li>
	 * <li>主题名称：{@link PropertyOption#TOPIC}</li>
	 * <li>过期时间：{@link PropertyOption#EXPIRE_TIME}</li>
	 * <li>生效时间：{@link PropertyOption#EFFECTIVE_TIME}</li>
	 * <li>是否压缩：{@link PropertyOption#COMPRESS}</li>
	 * <li>*消息组：{@link PropertyOption#GROUP}</li>
	 * <li>*优先级：{@link PropertyOption#PRIORITY}</li>
	 * <li>生产者重发次数：{@link PropertyOption#PRODUCER_RETRY}</li>
	 * <li>*客户端流水：{@link PropertyOption#CUSTOM_SERIAL}，当生产者重发次数大于0时，该值必须和第一次发送时的值相同</li>
	 * </ul>
	 */
	SEND = 11
	/**
	 * 发送后提交。Broker做BLE计算后直接转发<br/>
	 * 消息头列表， (*)表示可选：
	 * <ul>
	 * <li>客户端标识：{@link PropertyOption#CLIENT_ID}</li>
	 * <li>消息标识：{@link PropertyOption#MESSAGE_ID}</li>
	 * <li>原 {@link #SEND} 消息所带属性，由Broker添加（缓存没有时从存储取）</li>
	 * <li>*计算出目标主题的key：{@link PropertyOption#CURRENT_PROPERTY_KEY}，由Broker计算添加</li>
	 * <li>*计算出目标主题的value：{@link PropertyOption#CURRENT_PROPERTY_VALUE}，由Broker从配置取值添加</li>
	 * <li>*目标主题：{@link PropertyOption#TARGET_TOPIC}，由Broker计算添加</li>
	 * </ul>
	 */
	SEND_COMMIT = 12
	/**
	 * 发送后回滚。<br/>
	 * 消息头列表， (*)表示可选：
	 * <ul>
	 * <li>消息标识：{@link PropertyOption#MESSAGE_ID}</li>
	 * </ul>
	 */
	SEND_ROLLBACK = 13
	/**
	 * 拉取。Broker做BLE计算后直接转发<br/>
	 * 消息头列表， (*)表示可选：
	 * <ul>
	 * <li>客户端标识：{@link PropertyOption#CLIENT_ID}</li>
	 * <li>主题名称：{@link PropertyOption#TOPIC}</li>
	 * <li>预期处理时间：{@link PropertyOption#PROCESSING_TIME}</li>
	 * <li>*应答返回码：{@link PropertyOption#RESULT_CODE}，该值表示指定消息标识的消息的处理结果： {@link ResultCode#OK}
	 * 标识commit，否则rollback</li>
	 * <li>*应答返回描述：{@link PropertyOption#RESULT_DESCRIPTION}</li>
	 * <li>*消息标识：{@link PropertyOption#MESSAGE_ID}，当存在应答返回码时必选</li>
	 * </ul>
	 */
	PULL = 21
	/**
	 * 拉取应答。<br/>
	 * 消息头列表， (*)表示可选：
	 * <ul>
	 * <li>应答返回码：{@link PropertyOption#RESULT_CODE}</li>
	 * <li>*应答返回描述：{@link PropertyOption#RESULT_DESCRIPTION}</li>
	 * <li>*原 {@link #SEND} 消息所带属性，当应答返回码为 {@link ResultCode#OK} 时必选</li>
	 * <li>*消费者重发次数：{@link PropertyOption#CONSUMER_RETRY}，当应答返回码为 {@link ResultCode#OK} 时必选</li>
	 * <li>*生产者提交时间：{@link PropertyOption#COMMIT_TIME}，当应答返回码为 {@link ResultCode#OK} 时必选，需Broker添加</li>
	 * </ul>
	 */
	PULL_ANSWER = 22
	/**
	 * 删除消息。Broker做BLE计算后直接转发<br/>
	 * 消息头列表， (*)表示可选：
	 * <ul>
	 * <li>客户端标识：{@link PropertyOption#CLIENT_ID}</li>
	 * <li>消息标识：{@link PropertyOption#MESSAGE_ID}</li>
	 * <li>*目标主题：{@link PropertyOption#TARGET_TOPIC}，由Broker计算添加</li>
	 * </ul>
	 */
	DELETE = 30
	/**
	 * 解锁消息。BLE之间交互消息<br/>
	 * 消息头列表， (*)表示可选：
	 * <ul>
	 * <li>客户端标识：{@link PropertyOption#CLIENT_ID}</li>
	 * <li>消息标识：{@link PropertyOption#MESSAGE_ID}</li>
	 * <li>目标主题：{@link PropertyOption#TARGET_TOPIC}</li>
	 * </ul>
	 */
	UNLOCK = 40
	/**
	 * 未定义消息，此类消息也用 {@link MessageType#ANSWER} 返回。<br/>
	 * 使用 {@link PropertyOption#RESULT_CODE} 和 {@link PropertyOption#RESULT_DESCRIPTION} 进行描述
	 */
	UNDEFINED = 99
)

// ResultCode
const (
	/**
	 * 成功
	 */
	OK = "OK"
	/**
	 * pull无更多消息
	 */
	NO_MORE_MESSAGE = "NO_MORE_MESSAGE"
	/**
	 * 没有可用服务端地址，需要客户端检查地址配置
	 */
	SERVER_NOT_AVAILABLE = "SERVER_NOT_AVAILABLE"
	/**
	 * 负载地址无法获取，需要客户端检查地址配置
	 */
	SERVICE_ADDRESS_NOT_FOUND = "SERVICE_ADDRESS_NOT_FOUND"
	/**
	 * 参数缺少，需要客户端检查请求参数
	 */
	REQUIRED_PARAMETER_MISSING = "REQUIRED_PARAMETER_MISSING"
	/**
	 * 错误的请求，消息类型错误，消息字段缺少时等等场景使用
	 */
	BAD_REQUEST = "BAD_REQUEST"
	/**
	 * 请求消息类型错误
	 */
	UNSUPPORTED_MESSAGE_TYPE = "UNSUPPORTED_MESSAGE_TYPE"
	/**
	 * 请求消息长度过长
	 */
	MESSAGE_CONTENT_LENGTH_IS_TOO_LANG = "MESSAGE_CONTENT_LENGTH_IS_TOO_LANG"
	/**
	 * 服务器遇到了一个未曾预料的状况，导致了它无法完成对请求的处理。<br/>
	 * 一般来说，这个问题都会在服务器的程序码出错时出现。<br/>
	 * 此时客户端之后的重试无意义，会有同样的错误
	 */
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	/**
	 * 服务器上数据库相关操作发生了异常
	 */
	INTERNAL_DATA_ACCESS_EXCEPTION = "INTERNAL_DATA_ACCESS_EXCEPTION"
	/**
	 * 由于临时的服务器维护或者过载，服务器当前无法处理请求。<br/>
	 * 这个状况是临时的，并且将在一段时间以后恢复。<br/>
	 * 和 {@link ResultCode#INTERNAL_SERVER_ERROR} 不同的是客户端之后的重试都是有意义的。<br/>
	 * 客户端API中物理链路上发生的超时、断链、不可写、不可读等等状态都可以表示
	 */
	INTERNAL_SERVICE_UNAVAILABLE = "INTERNAL_SERVICE_UNAVAILABLE"
)

// PullCode
const (
	/**
	 * 提交消息
	 */
	PULL_COMMIT = "COMMIT"
	/**
	 * 提交消息并且获取下一个消息
	 */
	PULL_COMMIT_AND_NEXT = "COMMIT_AND_NEXT"
	/**
	 * 回滚当前消息<br/>
	 * 错误描述通过 {@link PropertyOption#CODE_DESCRIPTION} 指定<br/>
	 */
	PULL_ROLLBACK = "ROLLBACK"
	/**
	 * 回滚当前消息，并跳过且不再进行当前消息的处理<br/>
	 * 错误描述通过 {@link PropertyOption#CODE_DESCRIPTION} 指定<br/>
	 */
	PULL_ROLLBACK_AND_NEXT = "ROLLBACK_AND_NEXT"
	/**
	 * 回滚当前消息，并在一定时间后重试<br/>
	 * 一定时间通过 {@link PropertyOption#RETRY_AFTER} 指定<br/>
	 * 错误描述通过 {@link PropertyOption#CODE_DESCRIPTION} 指定<br/>
	 * 注：服务端锁定同group数据<br/>
	 */
	PULL_ROLLBACK_BUT_RETRY = "ROLLBACK_BUT_RETRY"
)
