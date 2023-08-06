# Chat-API-Proxy

## 项目介绍
这是一个gpt-3.5-turbo api的代理程序，将请求转发到一些公益站点，实现api免费使用。

## 服务提供者
[fakeopen](https://ai.fakeopen.com/)  
[xyhelper](https://xyhelper.cn/)  
chatgpt
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
    "temperature":1,
    "presence_penalty":0,
    "top_p":1,
    "frequency_penalty":0,
    "stream": true
}'
```










