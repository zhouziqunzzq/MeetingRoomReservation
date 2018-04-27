# A backend service for *Teacher Assessment APP*
This is a backend service for an android app
developed for our *Software Engineering Class*.
The backend is written in *pure Golang* and
follows *RESTful API* standard.
## API(v1) definition
### API URL base
    http|https://yourdomain.com/api/v1
### Heartbeat test
#### GET: /test
* Test with GET method
* Request: `null`
* Response:

        {
            "code":   200,
            "result": true,
            "msg":    "OK"
        }
#### POST: /test
* Test with POST method
* Request: `any data`
* Response:

        {
            "code":   200,
            "result": true,
            "data":   your request data in json format
        }
### Authentication
#### POST: /auth/login
* Login with an username and a password
* Request:
    * username  string
    * password  string
* Response:
    * success

            {
                "code":   200,
                "result": true,
                "msg":    "登陆成功",
                "access_token":  user's JWT
            }
    * error

            {
                "code":   error code,
                "result": false,
                "msg":    "用户名错误" or "密码错误" or "登录失败，未知错误",
            }
### User basic info(login required)
#### GET: /user/info
* Get user basic info
* Request:
    * access_token string
* Response:
    * success

            {
                "code":   200,
                "result": true,
                "msg":    "success",
                "data":   {
                                "username": username,
                                "name":     张三,
                                "tel":      telephone,
                                "identity": 0(管理员) or 1(教师)
                          }
            }
    * error

            {
                "code":   error code,
                "result": false,
                "msg":    "未授权的访问" or "无效的Token",
            }