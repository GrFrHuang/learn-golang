package main

import (
	"io/ioutil"
	"net/http"
	"github.com/GrFrHuang/gox/log"
	"os"
	"fmt"
	"strings"
)

// golang爬虫的demo

type Spider struct {
	url    string
	header map[string]string
}

var header = map[string]string{
	"Host":                      "baike.baidu.com",
	"Connection":                "keep-alive",
	"Cache-Control":             "max-age=0",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"Referer":                   "https://www.baidu.com/link?url=sXReeP8-aZia2yEE1VzYOKwl9xv9c-iSYPcqOP2gjerdc_joGNeBqkKjHN2k0bhj3USE9g6DKXErnL5gz5CvUjhYvFs8HhxDjs292zHisURNYsNbLEMfOQyFz6Hf2VsrLSrRxB7b7yVGnGh67mquoUfQI_xrgwX4EjV6lC8Bwkd62mEX3-ZWuQLhJxvpYu4iJ2s0aH2TQjgjTZDO76qqW_&wd=&eqid=f2dd797500026812000000055b557450",
	//"Referer":                   "https://www.baidu.com/",
}

func NewSpider(url string, header map[string]string) (*Spider) {
	return &Spider{
		url:    url,
		header: header,
	}
}

func (spider *Spider) Fetch() (string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", spider.url, nil)
	if err != nil {
		log.Error(err)
		return ""
	}
	for key, value := range spider.header {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(body)
}

func main() {
	//创建excel文件
	f, err := os.Create("./haha3.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Sprintf("合江县介绍\n这里好多吃的啊")
	f.WriteString(fmt.Sprintf("\t合江县介绍\n\n\n这里好多吃的啊"))
	bts, _ := ioutil.ReadFile("./haha3.txt")
	fmt.Println(fmt.Sprintf("%t", strings.Contains(string(bts),"\t")))
	//f.WriteString("电影名称" + "\t" + "评分" + "\t" + "评价人数" + "\t" + "\r\n")

	//写入标题
	//f.WriteString("电影名称" + "\t" + "评分" + "\t" + "评价人数" + "\t" + "\r\n")
	//spider := NewSpider("https://baike.baidu.com/item/%E6%88%90%E9%83%BD%E9%AB%98%E6%96%B0%E6%8A%80%E6%9C%AF%E4%BA%A7%E4%B8%9A%E5%BC%80%E5%8F%91%E5%8C%BA/4149621?fr=aladdin", header)
	//log.Info(spider.Fetch())
	//for i := 0; i < 10; i++ {
	//	fmt.Println("正在抓取第" + strconv.Itoa(i) + "页......")
	//	spider := NewSpider("https://movie.douban.com/top250?start="+strconv.Itoa(i*25)+"&filter=", header)
	//	html := spider.Fetch()
	//	//评价人数
	//	pattern2 := `<span>(.*?)评价</span>`
	//	rp2 := regexp.MustCompile(pattern2)
	//	rp2.FindAllStringSubmatch(html, -1)
	//
	//	//评分
	//	pattern3 := `property="v:average">(.*?)</span>`
	//	rp3 := regexp.MustCompile(pattern3)
	//	 rp3.FindAllStringSubmatch(html, -1)
	//
	//	//电影名称
	//	//pattern4 := `img alt="(.*?)" src=`
	//	pattern4 := `img width="100" alt="(.*?)"`
	//	rp4 := regexp.MustCompile(pattern4)
	//	find_txt4 := rp4.FindAllStringSubmatch(html, -1)
	//	log.Info(find_txt4[20][1])
	//
	//	// 写入UTF-8 BOM
	//	f.WriteString("\xEF\xBB\xBF")
	//	//  打印全部数据和写入excel文件
	//	//for i := 0; i < len(find_txt2); i++ {
	//	//	fmt.Printf("%s %s %s\n", find_txt4[i][1], find_txt3[i][1], find_txt2[i][1], )
	//	//	f.WriteString(find_txt4[i][1] + "\t" + find_txt3[i][1] + "\t" + find_txt2[i][1] + "\t" + "\r\n")
	//	//
	//	//}
	//	//log.Info(find_txt4)
	//}
}
