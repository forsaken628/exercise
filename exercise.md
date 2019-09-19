# 网络传输分析、服务设计与实现

## 要求
- 所有解答源码请放入公开 github 仓库，并给出可访问的 url

1. 网络传输分析: wireshark 抓包文件 mysterious_networking_behavior.pcapng/mysterious_networking_behavior.txt 内容描述了一次客户端/服务端架构下的网络行为产生的网络传输数据 (.pcapng 与 .txt 文件内容等价), 请根据抓包文件回答以下问题，细节越多越好
    - 详细描述客户端发起的每一次 DNS 请求和结果
    - 说明客户端与服务端建立了多少个 TCP channel，分别是哪些 frame，分别完成了什么传输任务，为什么存在多个 TCP channel
    - 选择几个 frame 详细说明一次 TCP 握手流程，需要包含具体 frame 内容
    - 请说明服务端程序可以如何优化，以提升单个用户访问延迟，以及并发吞吐量

2. 服务设计与实现: 仅使用 TCP socket 库，实现一个 HTTP 服务程序，监听 localhost:8080 端口，使用任意网页浏览器 (Chrome/Firefox/Safari等) 打开 http://localhost:8080，显示用户名/密码和对应输入框，以及登录按钮，点击后跳转页面，显示刚刚输入的用户名及密码
   - 不限编程语言
   - 整个服务程序为单个文件，不能包含资源文件
   - 登录一次后，下次访问直接跳转至显示用户名密码的页面
