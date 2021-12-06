# pdd_fanli
Third party extension interface for sillyGirl
按照sillyGirl的源码，自己拼凑了一个扩展插件，方便在pdd上获取推广佣金。

https://github.com/hdbjlizhe/pdd_fanli

具体步骤如下：

1.登入账号 https://jinbao.pinduoduo.com/

2.推广管理->推广者登记

新增”媒体登记“、”推广位“等。可以得到XXXXXXXX-XXXXXXXXX格式的推广位pid

3.注册或进入 多多进宝 API开放平台https://open.pinduoduo.com/

创建应用会得到client_id和client_key

还要做一下绑定备案，这个绑定备案我是用pdd开放平台的API测试工具完成的。扩展插件里包含了绑定的功能，但是此功能未测试，不能保证其可用性。

4.https://gituhub.com/hdbjlizhe/pdd_fanli，文件夹下载到sillyGirl/develop/下面

将sillyGirl/dev.go增加一行pdd_fanli的代码"github.com/cdle/sillyGirl/develop/pdd_fanli"

用命令go build编译一下，然后启动。

启动后用管理账号输入命令

set pinduoduo client_id XXXX

set pinduoduo client_key XXXX

set pinduoduo pid XXX

最后再重启一下sillyGirl

5.效果图：
![144696501-3c914fc0-1152-48a6-99a7-d1e95c4e755b](https://user-images.githubusercontent.com/22290807/144797607-dca114a5-6be3-4385-af0f-eeab99c60237.jpg)

6.请我喝杯咖啡（buy me a cup of coffee）
![微信图片_20211206140154](https://user-images.githubusercontent.com/22290807/144797154-6516c74e-fec8-4342-a628-20997eef5826.png)
