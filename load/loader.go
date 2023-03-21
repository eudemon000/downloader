package load

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// 默认协程下载数量
const QUEUENUM = 10

// 下载文件的结构体
type Read struct {
	io.Reader
	total   int64
	current int64
}

type FileStruct struct {
	url    string
	length int
	name   string
	start  int
	end    int
}

type DownloadContent struct {
	reader  io.Reader
	total   int64
	current int64
}

type FileName interface {
	CreateName()
}

func (read *Read) Read(p []byte) (n int, err error) {
	n, err = read.Reader.Read(p)
	read.current += int64(n)
	num := read.current * 100 / read.total
	//fmt.Printf("\r当前下载%d%%", num)
	fmt.Printf("\r%d%%", num)
	return n, err
}

func Run(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	s, isExist := createName(*resp)
	if !isExist {
		log.Fatal("不存在")
	}
	file, err := os.Create(s)
	r := &Read{
		Reader: resp.Body,
		total:  resp.ContentLength,
	}
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	io.Copy(file, r)
	return nil
}

/*
Content-Type:application/octet-stream表示是一个字节流，大多搭配
content-disposition（文件名）出现。可以下载的文件包括文件，视频，
压缩包等
application/octet-stream
application/zip
video/***
application/vnd.android.package-archive
*/
func createName(r http.Response) (string, bool) {
	header := r.Header
	contentType := header.Get("content-type")
	if contentType != "application/octet-stream" && contentType != "application/zip" &&
		!strings.Contains(contentType, "video/") &&
		contentType != "application/vnd.android.package-archive" &&
		contentType != "application/x-iso9660-image" {
		return "", false
	}
	//"attachment; filename=Postman-win64-Setup.exe"
	//"inline; filename=\"cloudmusicsetup2.9.6.199543..exe\"
	//"attachment; filename=wordpress-4.9.4-zh_CN.tar.gz"
	//"attachment; filename=MicrosoftEdgeSetup.exe"
	//"application/x-iso9660-image"
	contentDisposition := header.Get("content-disposition")
	if contentDisposition != "" {
		name := contentDispositionName(contentDisposition)
		return name, true
	}
	fmt.Println(contentDisposition)
	name := noContentDispositionName(*r.Request)
	return name, true
}

func contentDispositionName(str string) string {
	strs := strings.Split(str, "=")
	if len(strs) == 0 {
		return ""
	}
	temp := strs[len(strs)-1]
	temp = strings.ReplaceAll(temp, "\"", "")
	fmt.Println(temp)
	return temp
}

func noContentDispositionName(r http.Request) string {
	path := r.URL.Path
	strs := strings.Split(path, "/")
	name := strs[len(strs)-1]
	return name
}

func Test() {
	for i := 0; i < 10000; i++ {
		fmt.Printf("\r%s,", "进度")
	}
}

// 先获取header，判断是不是可以下载的URL
func Header(url string) {
	respone, err := http.Head(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	header := respone.Header
	contentType := header.Get("content-type")
	if contentType != "application/octet-stream" && contentType != "application/zip" &&
		!strings.Contains(contentType, "video/") &&
		contentType != "application/vnd.android.package-archive" &&
		contentType != "application/x-iso9660-image" &&
		contentType != "binary/octet-stream" {
		fmt.Fprintln(os.Stderr, "不支持的下载类型")
		return
	}
	name := getName(*respone)
	length := header.Get("content-length")
	contentLength, _ := strconv.Atoi(length)
	subSize := contentLength / QUEUENUM
	start := 0
	end := 0
	var w sync.WaitGroup
	for i := 0; i < QUEUENUM; i++ {
		start = end
		if i == QUEUENUM-1 {
			end = contentLength
		} else {
			end += subSize
		}
		f := &FileStruct{
			url:    url,
			length: contentLength,
			//name:   name + strconv.Itoa(i),
			name:  fmt.Sprintf("%s.part%d", name, i),
			start: start,
			end:   end,
		}
		fmt.Fprintf(os.Stdout, "start:%d, end:%d\n", f.start, f.end)
		end++
		go download(f, &w)
		w.Add(1)
	}
	w.Wait()
	fmt.Println("是否打印")
	mergeFile(name)
}

func getName(response http.Response) string {
	contentDisposition := response.Header.Get("content-disposition")
	if contentDisposition != "" {
		name := contentDispositionName(contentDisposition)
		return name
	}
	fmt.Println(contentDisposition)
	name := noContentDispositionName(*response.Request)
	fmt.Println(name)
	return name
}

// 下载
func downloadForWriteFile(f *FileStruct) {
	fmt.Fprintln(os.Stdout, "url:", f.url, "	name:", f.name, "	开始下载")
	client := &http.Client{}
	request, err := http.NewRequest("GET", f.url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	request.Header.Add("range", fmt.Sprintf("bytes=%d-%d", f.start, f.end))
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	err = os.WriteFile(f.name, b, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

// 下载
func download(f *FileStruct, w *sync.WaitGroup) {
	defer w.Done()
	fmt.Fprintln(os.Stdout, "url:", f.url, "	name:", f.name, "	开始下载")
	client := &http.Client{}
	request, err := http.NewRequest("GET", f.url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	request.Header.Add("range", fmt.Sprintf("bytes=%d-%d", f.start, f.end))
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	b := bufio.NewReader(response.Body)
	defer response.Body.Close()

	file, err := os.Create("temp\\" + f.name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()
	d := &DownloadContent{
		reader: b,
		total:  int64(f.length),
	}
	io.Copy(file, d)

}

func (d *DownloadContent) Read(p []byte) (n int, err error) {
	n, err = d.reader.Read(p)
	//fmt.Fprint(os.Stdout, n, "\t")
	d.current += int64(n)

	num := d.current * 100 / d.total
	//fmt.Printf("\r当前下载%d%%", num)
	//fmt.Printf("\r%d%%", num)
	fmt.Printf("%d%%\n", num)

	return
}

func mergeFile(name string) {
	fs, err := os.ReadDir("temp")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	f, err := os.Create("./" + name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer f.Close()
	for _, v := range fs {
		if strings.Contains(v.Name(), name) {
			t, err := os.ReadFile("temp/" + v.Name())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			f.Write(t)
			os.Remove("temp/" + v.Name())
		}
	}
}
