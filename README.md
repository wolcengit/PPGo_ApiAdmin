# XOP Admin
针对公司XOP平台的管理工具

对XOP产品进行API接口和文档进行管理

项目使用[PPGo_ApiAdmin](https://github.com/george518/PPGo_ApiAdmin)简化后作为基础框架，加上特定功能构建。



## 安装   

创建一个MySQL数据库后参考 docker-compose.yml.sample 进行容器化部署

## 使用   

访问：http://your_host:8081

用户名：admin 密码：admin    


## POST对接

主要用于导入数据

增加或修改时候使用form-data构建数据，请求地址 http://your_host:8081/xxx/ajaxsave?appkey=xxxxxx&product=xxx

删除记录请求地址 http://your_host:8081/xxx/ajaxsave?appkey=xxxxxx&product=xxx&id=xxx


### XOP产品
- URL：/xopproduct
- INSERT: code,name,detail => id
- UPDATE: id,code,name,detail => id

### XOP模块
- URL：/xopmodule
- INSERT: xop_name,code,name,detail => id
- UPDATE: id,xop_name,code,name,detail => id

### XOP类别
- URL：/xopcategory
- INSERT: mod_id,code,name,detail => id
- UPDATE: id,mod_id,code,name,detail => id

### XOP分组
- URL：/xopgroup
- INSERT: cat_id,mod_id,code,name,detail => id
- UPDATE: id,cat_id,mod_id,code,name,detail => id

### XOP函数
- URL：/xopfunction
- INSERT: grp_id,cat_id,mod_id,xop_name,code,name,detail,detail2 => id
- UPDATE: id,grp_id,cat_id,mod_id,xop_name,code,name,detail,detail2 => id

### XOP文档
- URL：/xopdocument
- INSERT: pid,sort,,name,detail => id
- UPDATE: id,pid,sort,name,detail => id

### 书籍类别
- URL：/bookcategory
- INSERT: name,detail => id
- UPDATE: id,name,detail => id

### XOP书籍
- URL：/booklibrary
- INSERT: cat_id,name,detail,opened,nodes_data => id
- UPDATE: id,cat_id,name,detail,opened,nodes_data => id


## Ref-API

通过唯一RefId进行数据导入

使用JSON构建数据，请求地址 http://your_host:8081/api/xopnode?appkey=xxxxxx&product=xxx

```json
{
"level": "1-module 2-category 3-group 4-function",
"mod":"ref-id",
"cat":"ref-id",
"grp":"ref-id",
"func":"ref-id",
"xop": "xop object name",
"code":"",
"name":"",
"detail":"",
"detail2":""
}

```