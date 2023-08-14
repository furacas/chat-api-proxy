# Chat-API-Proxy

## 项目介绍
这是一个免费的gpt api的代理程序，将请求转发到一些公益或白嫖站点，实现api免费使用。

## 服务提供者
[fakeopen](https://ai.fakeopen.com/)  
[xyhelper](https://xyhelper.cn/)  
chatgpt  
[ava](https://ava-ai-ef611.web.app/)  
[chimeragpt](https://chimeragpt.adventblocks.cc/)  
[chatanywhere](https://github.com/chatanywhere/GPT_API_free)  
## 使用方法

```bash
docker run -d --name=chat-api-proxy -p 8080:8080 -e API_KEYS=chat-api-proxy-api-key --rm furacas/chat-api-proxy:latest
```

```bash
curl --location 'http://127.0.0.1:8080/v1/chat/completions' \
--header 'Authorization: Bearer chat-api-proxy-api-key' \
--header 'Content-Type: application/json' \
--data '{
    "messages":[
        {
            "role":"system",
            "content":"You are ChatGPT, a large language model trained by OpenAI.\nCarefully heed the user'\''s instructions. \nRespond using Markdown."
        },
        {
            "role":"user",
            "content":"Hello"
        }
    ],
    "model":"gpt-3.5-turbo",
    "stream": true
}'
```

## 高级配置

### provider
Request Header里面传入`X-Provider`可以指定provider。  
不指定Provider的情况下如果调用失败会尝试下一个Provider，直到成功或者没有Provider可用，指定Provider的情况下如果调用失败会直接返回失败。
### 环境变量


| 变量名            | 默认  | 必填 | 备注                                                                                  |
|------------------|-------|------|---------------------------------------------------------------------------------------|
| FAKEOPEN_ENABLED | true  | 否   | 启用[fakeopen](https://ai.fakeopen.com/)                                             |
| XYHELPER_ENABLED | true  | 否   | 启用[xyhelper](https://xyhelper.cn/)                                                  |
| CHATGPT_ENABLED  | true  | 否   | 启用chatgpt                                                                           |
| AVA_ENABLED      | true  | 否   | 启用[ava](https://ava-ai-ef611.web.app/)                                              |
| CHATANYWHERE_KEY |       | 否   | [chatanywhere](https://github.com/chatanywhere/GPT_API_free) sk                       |
| CHIMERA_KEY      |       | 否   | [chimeragpt](https://chimeragpt.adventblocks.cc/) sk                                  |









