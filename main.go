package main

import (
	"net/http"
	"path/filepath"
	"os"
	"fmt"
	"container/list"
	"log"
	"io"
	"net"
	"time"
	"crypto/md5"
	"strconv"
)

func getFilelist(path string) *list.List {
	pathList := list.New();
	pathList.PushBack("<h1>目录</h1>")
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if ( f == nil ) {return err}
		if f.IsDir() {return nil}
		println(path + f.ModTime().String())
		pathList.PushBack("<p>" + "<span class=\"text-muted\">" +
			f.ModTime().String() + "</span> - <a href=\"dist/" +
			f.Name() + "\">" + f.Name() + "</a>" +
			"<button type=\"button\" class=\"btn btn-primary\" data-toggle=\"modal\" data-target=\"#myModal\">删除</button></p>")
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return pathList
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		fmt.Sprintf("%x", h.Sum(nil))
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./static/dist/" + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)  // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func hrefList(w http.ResponseWriter, r *http.Request){
	l := getFilelist("./static/dist")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Fprint(w, e.Value)
	}
}

func removeFile(w http.ResponseWriter, r *http.Request){
	name := r.FormValue("name")
	os.Remove("./static"+name)
}

func getCurrentIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""		
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func printVersion() {
	println("-------------------------------------------")
	println("*********** Bypass File Server ************")
	println("author: leefanquan@gmail.com")
	println("verison: 20171029.12.24")
	println("Listing: ", getCurrentIP() + ":9090")
	println("-------------------------------------------")
}

func main() {
	printVersion()
	http.HandleFunc("/list", hrefList)
	http.HandleFunc("/uploader", upload)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	err := http.ListenAndServe(":9090", nil)
	//h := http.FileServer(http.Dir("/home/aey/project/static"))
	//err := http.ListenAndServe(":9090", h)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}