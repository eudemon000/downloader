package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestUrl(t *testing.T) {
	//url := "https://blog.csdn.net/weixin_33968104/article/details/92169944"
	//url := "https://dl.pstmn.io/download/latest/win64"
	//url := "http://cdimage.deepin.com/releases/20.3/deepin-desktop-community-20.3-amd64.iso"
	//url := "https://dldir1.qq.com/qqfile/qq/PCQQ9.7.3/QQ9.7.3.28946.exe"
	//url := "https://d1.music.126.net/dmusic/cloudmusicsetup2.9.6.199543..exe"
	//url := "https://t.alipayobjects.com/L1/71/100/and/alipay_wap_main.apk"
	//url := "https://cn.wordpress.org/wordpress-4.9.4-zh_CN.tar.gz"
	//url := "https://c2rsetup.officeapps.live.com/c2r/downloadEdge.aspx?platform=Default&source=EdgeStablePage&Channel=Stable&language=zh-cn&brand=M100"
	//url := "https://down.sandai.net/thunder11/XunLeiWebSetup11.4.1.2030gw.exe"
	//url := "https://download.manjaro.org/gnome/22.0.4/manjaro-gnome-22.0.4-230222-linux61.iso"
	//url := "https://releases.ubuntu.com/22.04/ubuntu-22.04.2-desktop-amd64.iso"
	url := "https://www.win-rar.com/fileadmin/winrar-versions/winrar/winrar-x64-621sc.exe"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	requ := resp.Request
	u := requ.URL
	path := u.Path
	fmt.Printf("path:%s", path)
	header := resp.Header
	disposition := header.Get("Content-Disposition")
	conType := header.Get("Content-Type")
	fmt.Printf("Content-Disposition：%s\n", disposition)
	fmt.Printf("Content-Type：%s\n", conType)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
