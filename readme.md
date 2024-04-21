

服务采用垂直架构拆分，目录结构
```shell
./
├── README.md
├── application 微服务目录
│   ├── applet  BFF服务(api)
│   ├── article  文章
│   ├── chat  聊天
│   ├── concerned  关注
│   ├── member  会员
│   ├── message  消息
│   ├── qa  问答
│   └── user  用户
├── db  sql文件
│   └── user.sql
├── go.mod
├── go.sum
└── pkg  微服务共同依赖的方法
    ├── encrypt
    ├── jwt
    └── util
    
```