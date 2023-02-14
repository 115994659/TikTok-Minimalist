# 极简版抖音
## 基于gin+gorm的极简版抖音项目

<ul>
<li>架构说明</li>
<li>各模块代码详细说明</li>
<ul>
<li>Controller</li>
<li>Service</li>
<li>Models</li>
</ul>
<li>遇到的问题及对应解决方案
<ul>
<li>返回json数据的完整性和前端要求的一致性</li>
<li>视频的保存和封面的切片</li>
</ul>
</li>
<li>可改进的地方</li>
<ul>
<li>项目的社交部分Part还未完成</li>
<li>部分代码可以继续抽取</li>
</ul>
<li>项目运行</li>
<ul>
<li>ffmepg.exe(已放入lib自带，用于对视频切片得到封面)</li>
<li>需要gcc环境(主要用于cgo，windows请将mingw-w64设置到环境变量)</li>
</ul>


### 数据库关系说明
![Image text](C:\Users\asus\Desktop\Central Topic.png)

所有的表都有自己的id主键为唯一的标识。

user_logins：存下用户的用户名和密码

user_infos：存下用户的基本信息

videos：存下视频的基本信息

comment：存下每个评论的基本信息

### 架构说明

![Image text0](C:\Users\asus\Desktop\Central Topic.png)

以用户登录为例共需要经过以下过程：

进入中间件SHA1内的函数逻辑，得到password明文加密后再设置password。具体需要调用gin.Context的Set方法设置password。随后调用next()方法继续下层路由。
进入UserLoginController函数逻辑，获取username，并调用gin.Context的Get方法得到中间件设置的password。再调用service层的QueryUserLogin函数。
进入QueryUserLogin函数逻辑，执行三个过程：checkNum，prepareData，packData。也就是检查参数、准备数据、打包数据，准备数据的过程中会调用model层的UserLoginDAO。
进入UserLoginDAO的逻辑，执行最终的数据库请求过程，返回给上层。


#### Controllers
对于Controllerss层级的所有函数实现有如下规范：

所有的逻辑由代理对象进行，完成以下两个逻辑

- 解析得到参数
- 开始调用下层逻辑

#### Service
对于service层级的函数实现由如下规范：

同样由一个代理对象进行，完成以下三个或两个逻辑

当上层需要返回数据信息，则进行三个逻辑：

- 检查参数。
- 准备数据。
- 打包数据。

当上层不需要返回数据信息，则进行两个逻辑：

- 检查参数。
- 执行上层指定的动作。

#### Models
对于models层的各个操作，没有像service和Controller层针对前端发来的请求就行对应的处理，models层是面向于数据库的增删改查，不需要考虑和上层的交互。

而service层根据上层的需要来调用models层的不同代码请求数据库内的内容。

### 遇到的问题及对应解决方案
#### 返回json数据的完整性和前端要求的一致性
```
func FillCommentListFields(comments *[]*model.Comment) error {
	size := len(*comments)
	if comments == nil || size == 0 {
		return errors.New("util.FillCommentListFields comments为空")
	}
	dao := models.NewUserInfoDAO()
	for _, v := range *comments {
		_ = dao.QueryUserInfoById(v.UserInfoId, &v.User) //填充这条评论的作者信息
		v.CreateDate = v.CreatedAt.Format("1-2")         //转为前端要求的日期格式
	}
	return nil
}
```

#### 视频的保存和封面的切片

在本地建立static文件夹存储视频和封面图片。

具体逻辑如下：

- 检查视频格式
- 根据userId和该作者发布的视频数量产生唯一的名称
- 截取第一帧画面作为封面
- 保存视频基本信息到数据库(包括视频链接和封面链接)

### 可改进的地方
#### 项目的社交部分Part还未完成
整个项目在我的时间规划中是比较紧凑的，但队员的合作意愿太差，导致后续只有我一个人进行开发，打破了原定的计划
后续完成
#### 部分代码可以继续抽取
例如 在Controller层返回数据时，直接调用了c.json返回对应数据，这里可以再进行一次方法封装，抽离代码

### 项目运行
- mysql 5.7及以上
- redis 无版本限制
- ffmepg,已放入lib自带，用于对视频切片得到封面
- 需要gcc环境,主要用于cgo，windows请将mingw-w64设置到环境变量
` 在配置文件中进行本地配置
