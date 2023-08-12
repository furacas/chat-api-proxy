# Chat-API-Proxy

## 项目介绍
这是一个免费的gpt api的代理程序，将请求转发到一些公益或白嫖站点，实现api免费使用。

## 服务提供者
[fakeopen](https://ai.fakeopen.com/)  
[xyhelper](https://xyhelper.cn/)  
chatgpt  
[ava](https://ava-ai-ef611.web.app/)
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
在Reqiest Header里面传入`X-Provider`可以指定provider
### 不启用部分provider

通过设置环境变量来控制

| 变量名            | provider名   |
|-----------------|-------------|
| FAKEOPEN_ENABLED | fakeopen  |
| XYHELPER_ENABLED | xyhelper  |
| CHATGPT_ENABLED | chatgpt  |
| AVA_ENABLED | ava  |










