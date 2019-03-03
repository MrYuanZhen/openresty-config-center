package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type temp struct {
	fileName string
	fileType string
	content  string
}

var temps = map[string]temp{
	"nginx-http":    {"http_proxy.conf", "nginx-http", "/usr/local/openresty/nginx/conf/web/"},
	"nginx-https":   {"https_proxy.conf", "nginx-https", "/usr/local/openresty/nginx/conf/web/"},
	"cert-key":      {"cert.key", "cert-key", "/usr/local/openresty/nginx/conf/cert.d/"},
	"cert-crt":      {"cert.crt", "cert-crt", "/usr/local/openresty/nginx/conf/cert.d/"},
	"rewrite-rule":  {"rewrite.rule", "rewrite-rule", "/usr/local/openresty/nginx/conf/rule-config/"},
	"config-lua":    {"config.lua", "config-lua", "/usr/local/openresty/lualib/resty/upstream/"},
	"filebeat-yaml": {"filebeat.yaml", "filebeat-yaml", "/home/rancher/confd"}}

//nginx-http：没有证书以http方式访问的配置文件，文件名规定http_proxy.conf，Type规定nginx-http，生产环境路径为/usr/local/openresty/nginx/conf/web/；
//nginx-https：有证书以https方式访问的配置文件，文件名规定https_proxy.conf，Type规定nginx-https，生产环境路径为/usr/local/openresty/nginx/conf/web/，必须配合crt和key使用；
//cert-key：用于https访问是的证书key，文件名规定为cert.key，Type规定为cert-key，生产环境路径为/usr/local/openresty/nginx/conf/cert.d/；
//cert-crt：用于https访问是的证书crt，文件名规定为cert.crt，Type规定为cert-crty，生产环境路径为/usr/local/openresty/nginx/conf/cert.d/；
//rewrite-rule：该规则是规定当使用https方式访问时，需要跳转的https域名，文件名规定为rewrite-rule，Type规定为rewrite-rule，生产环境路径为/usr/local/openresty/nginx/conf/rule-config/；
//config-lua：主要配置一些防御规则开关，主要修改防御CC规则,文件名规定为config.lua，Type规定为config-lua，生产环境路径为/usr/local/openresty/lualib/resty/upstream/；
//filebeat-yaml：日志filebeat配置文件，文件名规定为filebeat.yaml，Type规定为filebeat-yaml，生产环境路径为/opt/filebeat/;

func main() {
	//绑定路由 如果访问 /upload 调用 Handler 方法
	http.HandleFunc("/upload", Handler)
	//使用 tcp 协议监听8888
	http.ListenAndServe(":8888", nil)
}

func Handler(w http.ResponseWriter, req *http.Request) {
	//输出对应的 请求方式
	fmt.Println(req.Method)
	//判断对应的请求来源。如果为get 显示对应的页面
	if req.Method == "GET" {
		fmt.Fprintln(w, "不支持这种调用方式!")
	} else if req.Method == "POST" {
		fileType := req.FormValue("fileType")

		//解析 form 中的file 上传名字
		file, file_head, file_err := req.FormFile("fileName")

		if file_err != nil {
			js := make(map[string]interface{})
			js["status"] = 500
			js["type"] = false
			js["msg"] = file_err
			upl, _ := json.Marshal(js)
			fmt.Fprintln(w, string(upl))
			return
		}

		if _, ok := temps[fileType]; !ok {
			js := make(map[string]interface{})
			js["status"] = 500
			js["type"] = false
			js["msg"] = "fileTpye_err"
			upl, _ := json.Marshal(js)
			fmt.Fprintln(w, string(upl))
			return
		}

		if file_head.Filename != temps[fileType].fileName {
			js := make(map[string]interface{})
			js["status"] = 500
			js["type"] = false
			js["msg"] = "fileName_err"
			upl, _ := json.Marshal(js)
			fmt.Fprintln(w, string(upl))
			return
		}

		file_save := temps[fileType].content + file_head.Filename
		//打开 已只读,文件不存在创建 方式打开  要存放的路径资源
		f, f_err := os.OpenFile(file_save, os.O_WRONLY|os.O_CREATE, 0666)
		if f_err != nil {
			fmt.Fprintf(w, "file open fail:%s", f_err)
			js := make(map[string]interface{})
			js["status"] = 500
			js["type"] = false
			js["msg"] = f_err
			upl, _ := json.Marshal(js)
			fmt.Fprintln(w, string(upl))
		}
		//文件 copy
		_, copy_err := io.Copy(f, file)
		if copy_err != nil {
			js := make(map[string]interface{})
			js["status"] = 500
			js["type"] = false
			js["msg"] = copy_err
			upl, _ := json.Marshal(js)
			fmt.Fprintln(w, string(upl))
		}
		//关闭对应打开的文件
		defer f.Close()
		defer file.Close()

		//返回上传结果
		js := make(map[string]interface{})
		js["status"] = 200
		js["type"] = true
		js["msg"] = "上传成功"
		upl, _ := json.Marshal(js)
		fmt.Fprintln(w, string(upl))

		checkNginx()

	} else { //如果有其他方式进行页面调用。http Status Code 500
		w.WriteHeader(500)
		fmt.Fprintln(w, "不支持这种调用方式!")

	}
}

func checkNginx() {
	check := exec.Command("/usr/local/openresty/nginx/sbin/nginx", "-t")
	out, err := check.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}

	var status bool = strings.Contains(string(out), "successful")

	if status == true {
		js := make(map[string]interface{})
		js["status"] = 200
		js["type"] = status
		js["msg"] = "检测通过"
		res, _ := json.Marshal(js)
		fmt.Println(string(res))

	} else {
		js := make(map[string]interface{})
		js["status"] = 500
		js["type"] = status
		js["msg"] = "检测失败"
		res, _ := json.Marshal(js)
		fmt.Println(string(res))
	}

}
