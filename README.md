接口说明
================
### /upload上传文件

方法：POST  

uri：/upload

##### 请求参数：
字段：fileName(String) fileType(String) 

//nginx-http：没有证书以http方式访问的配置文件，fileName规定http_proxy.conf，Type规定nginx-http  

//nginx-https：有证书以https方式访问的配置文件，fileName规定https_proxy.conf，Type规定nginx-https  

//cert-key：用于https访问是的证书key，fileName规定为cert.key，Type规定为cert-key  

//cert-crt：用于https访问是的证书crt，fileName规定为cert.crt，Type规定为cert-crty  

//rewrite-rule：该规则是规定当使用https方式访问时，需要跳转的https域名，fileName规定为rewrite-rule，Type规定为rewrite-rule  

//config-lua：主要配置一些防御规则开关，主要修改防御CC规则,fileName规定为config.lua，Type规定为config-lua  

//filebeat-yaml：日志filebeat配置文件，fileName规定为filebeat.yaml，Type规定为filebeat-yaml  


##### 返回参数：
status  状态 错误状态500  

type    标记为false  

msg     错误内容信息  


