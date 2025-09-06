1 接口功能
单向流式API为用户提供文本转语音的能力，支持多语种、多方言，同时支持websocket协议流式输出。

1.1 最佳实践
推荐使用链接复用，可降低耗时约70ms左右。
对比v1单向流式接口，不同的音色优化程度不同，以具体测试结果为准，理论上相对会有几十ms的提升。

2 接口说明

2.1 请求Request

请求路径
wss://openspeech.bytedance.com/api/v3/tts/unidirectional/stream

建连&鉴权

Request Headers
Key

说明

是否必须

Value示例

X-Api-App-Id

使用火山引擎控制台获取的APP ID，可参考 控制台使用FAQ-Q1

是

your-app-id

X-Api-Access-Key

使用火山引擎控制台获取的Access Token，可参考 控制台使用FAQ-Q1

是

your-access-key

X-Api-Resource-Id

表示调用服务的资源信息 ID

大模型语音合成：volc.service_type.10029、
volc.service_type.10048（并发版）

声音复刻2.0：volc.megatts.default（字符版）
volc.megatts.concurr（并发版）
是

大模型语音合成：volc.service_type.10029、
volc.service_type.10048（并发版）

声音复刻2.0：volc.megatts.default（字符版）
volc.megatts.concurr（并发版）
（不支持声音复刻1.0）

X-Api-Request-Id

标识客户端请求ID，uuid随机字符串

否

67ee89ba-7050-4c04-a3d7-ac61a63499b3


Response Headers
Key

说明

Value示例

X-Tt-Logid

服务端返回的 logid，建议用户获取和打印方便定位问题

2025041513355271DF5CF1A0AE0508E78C


WebSocket 二进制协议
WebSocket 使用二进制协议传输数据。
协议的组成由至少 4 个字节的可变 header、payload size 和 payload 三部分组成，其中

header 描述消息类型、序列化方式以及压缩格式等信息；
payload size 是 payload 的长度；
payload 是具体负载内容，依据消息类型不同 payload 内容不同；
需注意：协议中整数类型的字段都使用大端表示。

二进制帧
Byte

Left 4-bit

Right 4-bit

说明

0 - Left half

Protocol version

目前只有v1，始终填0b0001

0 - Right half

Header size (4x)

目前只有4字节，始终填0b0001

1 - Left half

Message type

固定为0b001

1 - Right half

Message type specific flags

在sendText时，为0
在finishConnection时，为0b100

2 - Left half

Serialization method

0b0000：Raw（无特殊序列化方式，主要针对二进制音频数据）0b0001：JSON（主要针对文本类型消息）

2 - Right half

Compression method

0b0000：无压缩0b0001：gzip

3

Reserved

留空（0b0000 0000）

[4 ~ 7]

[Optional field,like event number,...]

取决于Message type specific flags，可能有、也可能没有

...

Payload

可能是音频数据、文本数据、音频文本混合数据


payload请求参数

字段

描述

是否必须

类型

默认值

user

用户信息

user.uid

用户uid

event

请求的事件

namespace

请求方法

string

BidirectionalTTS

req_params.text

输入文本

string

req_params.model

模型版本，传seed-tts-1.1较默认版本音质有提升，并且延时更优，不传为默认效果。
注：若使用1.1模型效果，在复刻场景中会放大训练音频prompt特质，因此对prompt的要求更高，使用高质量的训练音频，可以获得更优的音质效果。

否

string

——

req_params.ssml

当文本格式是ssml时，需要将文本赋值为ssml，此时文本处理的优先级高于text。ssml和text字段，至少有一个不为空

string

req_params.speaker

发音人，具体见发音人列表

√

string

req_params.audio_params

音频参数，便于服务节省音频解码耗时

√

object

req_params.audio_params.format

音频编码格式，mp3/ogg_opus/pcm。接口传入wav并不会报错，在流式场景下传入wav会多次返回wav header，这种场景建议使用pcm。

string

mp3

req_params.audio_params.sample_rate

音频采样率，可选值 [8000,16000,22050,24000,32000,44100,48000]

number

24000

req_params.audio_params.bit_rate

音频比特率，可传16000、32000等。
bit_rate默认设置范围为64k～160k，传了disable_default_bit_rate为true后可以设置到64k以下
GoLang示例：additions = fmt.Sprintf("{"disable_default_bit_rate":true}")
注：bit_rate只针对MP3格式，wav计算比特率跟pcm一样是 比特率 (bps) = 采样率 × 位深度 × 声道数
目前大模型TTS只能改采样率，所以对于wav格式来说只能通过改采样率来变更音频的比特率

number

req_params.audio_params.emotion

设置音色的情感。示例："emotion": "angry"
注：当前仅部分音色支持设置情感，且不同音色支持的情感范围存在不同。
详见：大模型语音合成API-音色列表-多情感音色

string

req_params.audio_params.emotion_scale

调用emotion设置情感参数后可使用emotion_scale进一步设置情绪值，范围1~5，不设置时默认值为4。
注：理论上情绪值越大，情感越明显。但情绪值1~5实际为非线性增长，可能存在超过某个值后，情绪增加不明显，例如设置3和5时情绪值可能接近。

number

4

req_params.audio_params.speech_rate

语速，取值范围[-50,100]，100代表2.0倍速，-50代表0.5倍数

number

0

req_params.audio_params.loudness_rate

音量，取值范围[-50,100]，100代表2.0倍音量，-50代表0.5倍音量（mix音色暂不支持）

number

0

req_params.audio_params.enable_timestamp

设置 "enable_timestamp": true 返回字与音素时间戳（默认为 flase，参数传入 true 即表示启用）

bool

false

req_params.additions

用户自定义参数

jsonstring

req_params.additions.silence_duration

设置该参数可在句尾增加静音时长，范围0~30000ms。（注：增加的句尾静音主要针对传入文本最后的句尾，而非每句话的句尾）

number

0

req_params.additions.enable_language_detector

自动识别语种

bool

false

req_params.additions.disable_markdown_filter

是否开启markdown解析过滤，
为true时，解析并过滤markdown语法，例如，你好，会读为“你好”，
为false时，不解析不过滤，例如，你好，会读为“星星‘你好’星星”

bool

false

req_params.additions.disable_emoji_filter

开启emoji表情在文本中不过滤显示，默认为false，建议搭配时间戳参数一起使用。
GoLang示例：additions = fmt.Sprintf("{"disable_emoji_filter":true}")

bool

false

req_params.additions.mute_cut_remain_ms

该参数需配合mute_cut_threshold参数一起使用，其中：
"mute_cut_threshold": "400", // 静音判断的阈值（音量小于该值时判定为静音）
"mute_cut_remain_ms": "50", // 需要保留的静音长度
注：参数和value都为string格式
Golang示例：additions = fmt.Sprintf("{"mute_cut_threshold":"400", "mute_cut_remain_ms": "1"}")
特别提醒：

因MP3格式的特殊性，句首始终会存在100ms内的静音无法消除，WAV格式的音频句首静音可全部消除，建议依照自身业务需求综合判断选择
string

req_params.additions.enable_latex_tn

是否可以播报latex公式，需将disable_markdown_filter设为true

bool

false

req_params.additions.max_length_to_filter_parenthesis

是否过滤括号内的部分，0为不过滤，100为过滤

int

100

req_params.additions.explicit_language（明确语种）

仅读指定语种的文本
精品音色和 ICL 声音复刻场景：

不给定参数，正常中英混
crosslingual 启用多语种前端（包含zh/en/ja/es-ms/id/pt-br）
zh 中文为主，支持中英混
en 仅英文
ja 仅日文
es-mx 仅墨西
id 仅印尼
pt-br 仅巴葡
DIT 声音复刻场景：
当音色是使用model_type=2训练的，即采用dit标准版效果时，建议指定明确语种，目前支持：

不给定参数，启用多语种前端zh,en,ja,es-mx,id,pt-br,de,fr
zh,en,ja,es-mx,id,pt-br,de,fr 启用多语种前端
zh 中文为主，支持中英混
en 仅英文
ja 仅日文
es-mx 仅墨西
id 仅印尼
pt-br 仅巴葡
de 仅德语
fr 仅法语
当音色是使用model_type=3训练的，即采用dit还原版效果时，必须指定明确语种，目前支持：

zh 中文为主，支持中英混
en 仅英文
GoLang示例：additions = fmt.Sprintf("{"explicit_language": "zh"}")

string

req_params.additions.context_language（参考语种）

给模型提供参考的语种

不给定 西欧语种采用英语
id 西欧语种采用印尼
es 西欧语种采用墨西
pt 西欧语种采用巴葡
string

req_params.additions.unsupported_char_ratio_thresh

默认: 0.3，最大值: 1.0
检测出不支持语种超过设置的比例，则会返回错误码或者返回兜底音频。

float

0.3

req_params.additions.cache_config（缓存相关参数）

开启缓存，开启后合成相同文本时，服务会直接读取缓存返回上一次合成该文本的音频，可明显加快相同文本的合成速率，缓存数据保留时间1小时。
（通过缓存返回的数据不会附带时间戳）
Golang示例：additions = fmt.Sprintf("{"disable_default_bit_rate":true, "cache_config": {"text_type": 1,"use_cache": true}}")

object

req_params.additions.cache_config.text_type（缓存相关参数）

和use_cache参数一起使用，需要开启缓存时传1

int

1

req_params.additions.cache_config.use_cache（缓存相关参数）

和text_type参数一起使用，需要开启缓存时传true

bool

true

req_params.additions.post_process

后处理配置
Golang示例：additions = fmt.Sprintf("{"post_process":{"pitch":12}}")

object

req_params.additions.post_process.pitch

音调取值范围是[-12,12]

int

0

req_params.mix_speaker

混音参数结构

object

req_params.mix_speaker.speakers

混音音色名以及影响因子列表

最多支持3个音色混音
混音影响因子和必须=1
使用复刻音色时，需要使用查询接口获取的icl_的speakerid，而非S_开头的speakerid
音色风格差异较大的两个音色（如男女混），以0.5-0.5同等比例混合时，可能出现偶发跳变，建议尽量避免
注意：使用Mix能力时，req_params.speaker = custom_mix_bigtts

list

null

req_params.mix_speaker.speakers[i].source_speaker

混音源音色名（支持大小模型音色和复刻2.0音色）

string

""

req_params.mix_speaker.speakers[i].mix_factor

混音源音色名影响因子

float

0

单音色请求参数示例：

{
    "user": {
        "uid": "12345"
    },
    "req_params": {
        "text": "明朝开国皇帝朱元璋也称这本书为,万物之根",
        "speaker": "zh_female_shuangkuaisisi_moon_bigtts",
        "audio_params": {
            "format": "mp3",
            "sample_rate": 24000
        },
      }
    }
}
mix请求参数示例：

{
    "user": {
        "uid": "12345"
    },
    "req_params": {
        "text": "明朝开国皇帝朱元璋也称这本书为万物之根",
        "speaker": "custom_mix_bigtts",
        "audio_params": {
            "format": "mp3",
            "sample_rate": 24000
        },
        "mix_speaker": {
            "speakers": [{
                "source_speaker": "zh_male_bvlazysheep",
                "mix_factor": 0.3
            }, {
                "source_speaker": "BV120_streaming",
                "mix_factor": 0.3
            }, {
                "source_speaker": "zh_male_ahu_conversation_wvae_bigtts",
                "mix_factor": 0.4
            }]
        }
    }
}

2.2 响应Response

建连响应
主要关注建连阶段 HTTP Response 的状态码和 Body

建连成功：状态码为 200
建连失败：状态码不为 200，Body 中提供错误原因说明

WebSocket 传输响应

二进制帧 - 正常响应帧
Byte

Left 4-bit

Right 4-bit

说明

0 - Left half

Protocol version

目前只有v1，始终填0b0001

0 - Right half

Header size (4x)

目前只有4字节，始终填0b0001

1 - Left half

Message type

音频帧返回：0b1011
其他帧返回：0b1001

1 - Right half

Message type specific flags

固定为0b0100

2 - Left half

Serialization method

0b0000：Raw（无特殊序列化方式，主要针对二进制音频数据）0b0001：JSON（主要针对文本类型消息）

2 - Right half

Compression method

0b0000：无压缩0b0001：gzip

3

Reserved

留空（0b0000 0000）

[4 ~ 7]

[Optional field,like event number,...]

取决于Message type specific flags，可能有、也可能没有

...

Payload

可能是音频数据、文本数据、音频文本混合数据


payload响应参数
字段

描述

类型

data

返回的二进制数据包

[]byte

event

返回的事件类型

number

res_params.text

经文本分句后的句子

string


二进制帧 - 错误响应帧
Byte

Left 4-bit

Right 4-bit

说明

0 - Left half

Protocol version

目前只有v1，始终填0b0001

0 - Right half

Header size (4x)

目前只有4字节，始终填0b0001

1

Message type

Message type specific flags

0b11110000

2 - Left half

Serialization method

0b0000：Raw（无特殊序列化方式，主要针对二进制音频数据）0b0001：JSON（主要针对文本类型消息）

2 - Right half

Compression method

0b0000：无压缩0b0001：gzip

3

Reserved

留空（0b0000 0000）

[4 ~ 7]

Error code

错误码

...

Payload

错误消息对象


2.3 event定义
在发送文本转TTS阶段，不需要客户端发送上行的event帧。event类型如下：

Event code

含义

事件类型

应用阶段：上行/下行

152

SessionFinished，会话已结束（上行&下行）
标识语音一个完整的语音合成完成

Session 类

下行

350

TTSSentenceStart，TTS 返回句内容开始

数据类

下行

351

TTSSentenceEnd，TTS 返回句内容结束

数据类

下行

352

TTSResponse，TTS 返回句的音频内容

数据类

下行

在关闭连接阶段，需要客户端传递上行event帧去关闭连接。event类型如下：

Event code

含义

事件类型

应用阶段：上行/下行

2

FinishConnection，结束连接

Connect 类

上行

52

ConnectionFinished 结束连接成功

Connect 类

下行

交互示例：
Image

2.4 不同类型帧举例说明

SendText

请求Request
Byte

Left 4-bit

Right 4-bit

说明

0

0001

0001

v1

4-byte header

1

0001

0000

Full-client request

with no event number

2

0001

0000

JSON

no compression

3

0000

0000

4 ~ 7

uint32(...)

len(payload_json)

8 ~ ...

{...}

文本

payload

{
    "user": {
        "uid": "12345"
    },
    "req_params": {
        "text": "明朝开国皇帝朱元璋也称这本书为,万物之根",
        "speaker": "zh_female_shuangkuaisisi_moon_bigtts",
        "audio_params": {
            "format": "mp3",
            "sample_rate": 24000
        },
      }
    }
}

响应Response

TTSSentenceStart
Byte

Left 4-bit

Right 4-bit

说明

0

0001

0001

v1

4-byte header

1

1001

0100

Full-client request

with event number

2

0001

0000

JSON

no compression

3

0000

0000

4 ~ 7

TTSSentenceStart

event type

8 ~ 11

uint32(12)

len(<session_id>)

12 ~ 23

nxckjoejnkegf

session_id

24 ~ 27

uint32( ...)

len(text_binary)

28 ~ ...

{...}

text_binary


TTSResponse
Byte

Left 4-bit

Right 4-bit

说明

0

0001

0001

v1

4-byte header

1

1011

0100

Audio-only response

with event number

2

0001

0000

JSON

no compression

3

0000

0000

4 ~ 7

TTSResponse

event type

8 ~ 11

uint32(12)

len(<session_id>)

12 ~ 23

nxckjoejnkegf

session_id

24 ~ 27

uint32( ...)

len(audio_binary)

28 ~ ...

{...}

audio_binary


TTSSentenceEnd
Byte

Left 4-bit

Right 4-bit

说明

0

0001

0001

v1

4-byte header

1

1001

0100

Full-client request

with event number

2

0001

0000

JSON

no compression

3

0000

0000

4 ~ 7

TTSSentenceEnd

event type

8 ~ 11

uint32(12)

len(<session_id>)

12 ~ 23

nxckjoejnkegf

session_id

24 ~ 27

uint32( ...)

len(payload)

28 ~ ...

{...}

payload


SessionFinished
Byte

Left 4-bit

Right 4-bit

说明

0

0001

0001

v1

4-byte header

1

1001

0100

Full-client request

with event number

2

0001

0000

JSON

no compression

3

0000

0000

4 ~ 7

SessionFinished

event type

8 ~ 11

uint32(12)

len(<session_id>)

12 ~ 23

nxckjoejnkegf

session_id

24 ~ 27

uint32( ...)

len(response_meta_json)

28 ~ ...

{
"status_code": 20000000,
"message": "ok"
}

response_meta_json

仅含status_code和message字段

FinishConnection

请求request
Byte

Left 4-bit

Right 4-bit

说明

0

0001

0001

v1

4-byte header

1

0001

0100

Full-client request

with event number

2

0001

0000

JSON

no compression

3

0000

0000

4-7

uint32(...)

len(payload_json)

8 ~ ...

{...}

payload_json
扩展保留，暂留空JSON


响应response
Byte

Left 4-bit

Right 4-bit

说明

0

0001

0001

v1

4-byte header

1

1001

0100

Full-client request

with event number

2

0001

0000

JSON

no compression

3

0000

0000

4 ~ 7

ConnectionFinished

event type

8 ~ 11

uint32(7)

len(<connection_id>)

12 ~ 15

uint32(58)

len(<response_meta_json>)

28 ~ ...

{
"status_code": 20000000,
"message": "ok"
}

response_meta_json

仅含status_code和message字段

3 错误码
Code

Message

说明

20000000

ok

音频合成结束的成功状态码

45000000

speaker permission denied: get resource id: access denied

音色鉴权失败，一般是speaker指定音色未授权或者错误导致

quota exceeded for types: concurrency

并发限流，一般是请求并发数超过限制

55000000

服务端一些error

服务端通用错误


4 调用示例
Python调用示例
Java调用示例
Go调用示例
C#调用示例
TypeScript调用示例

前提条件
调用之前，您需要获取以下信息：
<appid>：使用控制台获取的APP ID，可参考 控制台使用FAQ-Q1。
<access_token>：使用控制台获取的Access Token，可参考 控制台使用FAQ-Q1。
<voice_type>：您预期使用的音色ID，可参考 大模型音色列表。

Python环境
Python：3.9版本及以上。
Pip：25.1.1版本及以上。您可以使用下面命令安装。
python3 -m pip install --upgrade pip

下载代码示例

volcengine_unidirectional_stream_demo.tar.gz
未知大小


解压缩代码包，安装依赖
mkdir -p volcengine_unidirectional_stream_demo
tar xvzf volcengine_unidirectional_stream_demo.tar.gz -C ./volcengine_unidirectional_stream_demo
cd volcengine_unidirectional_stream_demo
python3 -m venv .venv
source .venv/bin/activate
python3 -m pip install --upgrade pip
pip3 install -e .

发起调用
<appid>替换为您的APP ID。
<access_token>替换为您的Access Token。
<voice_type>替换为您预期使用的音色ID，例如zh_female_cancan_mars_bigtts。

python3 examples/volcengine/unidirectional_stream.py --appid <appid> --access_token <access_token> --voice_type <voice_type> --text "你好，我是火山引擎的语音合成服务。这是一个美好的旅程。"