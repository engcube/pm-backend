### PM System设计

前后端分离，前端由html+js搭建；后端提供RESTFUL API接口；前后端通过http协议通信。

#### 后端设计

- 语言和框架
  - golang实现，引用go-json-rest框架（提供RESTful JSON APIs）

- 数据库存储
  - sqlite3，轻量，方便移植
  - 封装go-sqlite3库
  - 可以考虑内存数据库（如redis、mongodb中），也可以是传统的mysql

- API version
  - 只实现了一个接口的API多版本支持（info接口），go-semver支持

- 支持静态文件

- 可选实现（时间限制，没有实现）
  - 单元测试
    - 只做了session的简单单元测试
  - 后端无状态，可以通过load balance服务（云商的lb、nginx、haproxy等）做负载均衡
  - 消息未做分页处理
  - 本地logging
  - 集中式日志（graylog等）
  - debug
  - prometheus等监控服务接入
  - https接入
  - auth
  - docker化部署

- RESTFUL API
  - 基本信息
    - GET /api/status；获取服务器状态
    - GET /api/#version/info；获取版本信息
  - 会话信息
    - GET /api/#version/session；获取会话信息
      - header中指定SessionID
      - 返回Session结构体
    - POST /api/#version/session；创建新的会话（登入） 
      - body中指定Session结构体
      - 返回Session结构体
    - PUT /api/#version/session；更新会话信息
      - body中指定Session结构体
      - 返回Session结构体
    - DELETE /api/#version/session；销毁当前会话（登出）  
      - body中指定Session结构体
      - 返回Session结构体
  - 用户信息
    - GET /api/#version/user； 获取自身用户的信息
    - POST /api/#version/user；创建新的用户（注册） 
    - PUT /api/#version/user；更新用户的信息 
    - DELETE /api/#version/user；删除用户（注销）
    - PUT /api/#version/user/:id/password；更新id用户密码（未实现）
  - 联系人信息
    - GET /api/#version/friend/:id；获取联系人信息
    - GET /api/#version/friend；获取所有联系人信息
    - POST /api/#version/friend；创建新联系人
    - DELETE /api/#version/friend；删除指定联系人
    - PUT /api/#version/friend；更新指定联系人nickname信息（未实现）
    - GET /api/#version/friend/message
  - 私信信息
    - GET /api/#version/message/amount；获取私信数目
    - GET /api/#version/message/amount/:id；获取z指定用户的私信数目
    - GET /api/#version/message；获取所有私信信息
    - GET /api/#version/message/:id；获取指定用户的私信
    - POST /api/#version/message；发送私信
    - DELETE /api/#version/message；删除指定私信
    - PUT /api/#version/message；阅读发送给自己的指定私信

- 数据库设计
  - t_session 会话信息表（可以考虑存储在redis等介质上）
    - create table t_session(session_id text primary key, user_id integer not null, insert_time integer, is_deleted integer default 0, update_time integer)
    - session_id text
    - user_id integer
    - insert_time integer
    - is_deleted integer
    - update_time integer
  - t_user 用户信息表
    - create table t_user(user_id integer primary key AUTOINCREMENT, email text unique not null, username text not null, password text not null, insert_time integer, is_deleted integer default 0, update_time integer)
    - user_id integer AUTO_INCREMENT
    - email text
    - username text
    - password text
    - insert_time integer
    - is_deleted integer
    - update_time integer
  - t_friend 联系人信息表
    - create table t_friend(friend_id integer primary key autoincrement, user_id integer not null, friend_user_id integer not null, nickname text,  insert_time integer, is_deleted integer default 0, update_time integer)
    - friend_id integer AUTO_INCREMENT
    - user_id integer
    - friend_user_id integer
    - insert_time integer
    - is_deleted integer
    - update_time integer
  - t_message 消息表
    - create table t_message(message_id integer primary key autoincrement, user_id integer not null, to_user_id integer not null, context text, is_viewed integer default 0, insert_time integer, is_deleted integer default 0, update_time integer)
    - message_id integer AUTO_INCREMENT
    - user_id integer
    - to_user_id integer
    - content text
    - is_viewed integer
    - insert_time integer
    - is_deleted integer
    - update_time integer

- API范例

  - 注册用户

    - POST localhost:8080/api/1.0.0/user

      - request

        - ```
          {
          "Email": "abc@ucloud.cn",
          "Password": "123456",
          "Username": "abc"
          }
          ```

      - response

        - ```
          {
            "UserID": 4,
            "Email": "abc@ucloud.cn",
            "Username": "abc",
            "Password": "",
            "InsertTime": 1495885976,
            "UpdateTime": 1495885976,
            "IsDeleted": false,
            "SessionID": ""
          }
          ```

  - 用户登录

    -  POST localhost:8080/api/1.0.0/session

      - request

        - ```
          {	
          	"Email": "abc@ucloud.cn",
          	"Password": "123456"
          }
          ```

      - response

        - ```
          {
            "UserID": 4,
            "Email": "abc@ucloud.cn",
            "Username": "abc",
            "Password": "",
            "InsertTime": 1495885976,
            "UpdateTime": 1495885976,
            "IsDeleted": false,
            "SessionID": "d1cf1e83-2bc9-4457-a54b-cf149f171281"
          }
          ```

  - 用户信息获取

    - GET localhost:8080/api/1.0.0/session

      - request

        - header中添加 SessionID:d1cf1e83-2bc9-4457-a54b-cf149f171281

      - response

        - ```
          {
            "SessionID": "d1cf1e83-2bc9-4457-a54b-cf149f171281",
            "UserID": 4,
            "UpdateTime": 1495886590
          }
          ```

  - 更新会话信息

    - PUT localhost:8080/api/1.0.0/session

      - request

        - header中添加 SessionID:d1cf1e83-2bc9-4457-a54b-cf149f171281

      - response

        - ​

          ```
          {
            "SessionID": "d1cf1e83-2bc9-4457-a54b-cf149f171281",
            "UserID": 4,
            "UpdateTime": 1495886590
          }
          ```

  - 注销

    - DELETE localhost:8080/api/1.0.0/session
      - request
        - header中添加 SessionID:d1cf1e83-2bc9-4457-a54b-cf149f171281

      - response

        - ```
          {
            "SessionID": "d1cf1e83-2bc9-4457-a54b-cf149f171281",
            "UserID": 4,
            "UpdateTime": 1495886590
          }
          ```

  - 获取用户信息

    - GET localhost:8080/api/1.0.0/user
      - request
        - header中添加 SessionID:e333b02c-334c-4407-a67d-a4ecaa3e7fe6

      - response

        - ```
          {
            "UserID": 4,
            "Email": "abc@ucloud.cn",
            "Username": "abc",
            "Password": "",
            "InsertTime": 1495885976,
            "UpdateTime": 1495885976,
            "IsDeleted": false,
            "SessionID": ""
          }
          ```

  - 修改用户名

    - PUT localhost:8080/api/1.0.0/user

      - request

        - header中添加 SessionID:e333b02c-334c-4407-a67d-a4ecaa3e7fe6

        - ```
          {
          	"UserID": 4,
          	"Username": "test4"
          }
          ```

      - response

        - ```
          {
            "UserID": 4,
            "Email": "",
            "Username": "test4",
            "Password": "",
            "InsertTime": 0,
            "UpdateTime": 0,
            "IsDeleted": false,
            "SessionID": "e333b02c-334c-4407-a67d-a4ecaa3e7fe6"
          }
          ```

  - 删除用户

    - DELETE localhost:8080/api/1.0.0/user

      - request

        - header中添加 SessionID:e333b02c-334c-4407-a67d-a4ecaa3e7fe6

        - ```
          {
          	"UserID": 4
          }
          ```

      - response

        - ```
          {
            "UserID": 4,
            "Email": "",
            "Username": "",
            "Password": "",
            "InsertTime": 0,
            "UpdateTime": 0,
            "IsDeleted": true,
            "SessionID": "e333b02c-334c-4407-a67d-a4ecaa3e7fe6"
          }
          ```

  - 添加好友

    - POST localhost:8080/api/1.0.0/friend

      - request

        - header中添加 SessionID:e333b02c-334c-4407-a67d-a4ecaa3e7fe6

        - ```
          {
          	"Email": "ab@ucloud.cn"
          }
          ```

      - response

        - ```
          {
            "FriendID": 3,
            "FriendUserID": 6,
            "Email": "abcdef@ucloud.cn",
            "Nickname": "abcdef",
            "InsertTime": 1495889188,
            "UpdateTime": 0,
            "IsDeleted": false,
            "SentUnreadMsgs": null,
            "SentReadMsgs": null,
            "RecieveUnreadMsgs": null,
            "RecieveReadMsgs": null,
            "UnreadCount": 0
          }
          ```

  - 删除好友

    - DELETE localhost:8080/api/1.0.0/friend

      - request 

        - header中添加 SessionID:e333b02c-334c-4407-a67d-a4ecaa3e7fe6

        - ```
          {
          	"FriendID": 2
          }
          ```

      - response

        - ```
          {
            "FriendID": 2,
            "FriendUserID": 0,
            "Email": "",
            "Nickname": "",
            "InsertTime": 0,
            "UpdateTime": 1495889734,
            "IsDeleted": true,
            "SentUnreadMsgs": null,
            "SentReadMsgs": null,
            "RecieveUnreadMsgs": null,
            "RecieveReadMsgs": null,
            "UnreadCount": 0
          }
          ```

  - 发送私信

    - POST localhost:8080/api/1.0.0/message

      - request

        - header中添加 SessionID:e333b02c-334c-4407-a67d-a4ecaa3e7fe6

        - ```
          {
          	"RecieverEmail": "ab@ucloud.cn",
          	"Content": "133443"
          }
          ```

      - response

        - ```
          {
            "MessageID": 2,
            "RecieverEmail": "ab@ucloud.cn",
            "Sender": 5,
            "Reciever": 2,
            "Content": "133443",
            "IsViewed": false,
            "InsertTime": 1495899191,
            "UpdateTime": 0,
            "IsDeleted": false
          }
          ```

  - 阅读私信

    - PUT localhost:8080/api/1.0.0/message

  - 删除私信

    - DELETE localhost:8080/api/1.0.0/message

      - request

        - header中添加 SessionID:e333b02c-334c-4407-a67d-a4ecaa3e7fe6

        - ```
          {
          	"MessageId": 2
          }
          ```

      - response

        - ```
          {
            "MessageID": 2,
            "RecieverEmail": "",
            "Sender": 5,
            "Reciever": 2,
            "Content": "133443",
            "IsViewed": false,
            "InsertTime": 1495899191,
            "UpdateTime": 1495899332,
            "IsDeleted": true
          }
          ```

  - 获取私信数目