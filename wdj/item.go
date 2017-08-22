package wdj

import (
	"os"
	"fmt"
	"time"
	"sync"
	"errors"
	"strings"
	"net/url"
	"encoding/json"
	"text/template"
)

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-pg/pg"
)

var ErrParse = errors.New("parse error")

const (
	AppPagePrefix   = "http://www.wandoujia.com/apps/"
	searchURLPrefix = "http://www.wandoujia.com/search?key="
	releaseTimeFmt  = "2006年01月02日"
	normalTimeFmt   = "20060102"
	appTemplate     = `
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
┃ 豌豆荚: {{ .ID }} {{ .Name }}
┃ {{.URL}}
┣┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈
┃ Icon        ┆ {{.Icon        }}
┃ Link        ┆ {{.Link        }}
┃ Version     ┆ {{.Version     }}
┃ Vendor      ┆ {{.Vendor      }}
┃ Genre       ┆ {{.Genre       }}
┃ Tags        ┆ {{.Tags        }}
┃ Categories  ┆ {{.Categories  }}
┃ Price       ┆ {{.Price       }}
┃ System      ┆ {{.System      }}
┃ Platform    ┆ {{.Platform    }}
┃ Permissions ┆ {{.Permissions }}
┃ Size        ┆ {{.Size        }}
┃ Rating      ┆ {{.Rating      }}
┃ InstallCnt  ┆ {{.InstallCnt  }}
┃ CommentCnt  ┆ {{.CommentCnt  }}
┃ Appkey      ┆ {{.Appkey      }}
┃ AppID       ┆ {{.AppID       }}
┃ ApkCode     ┆ {{.ApkCode     }}
┃ Subtitle    ┆ {{.Subtitle    }}
┃ Commentary  ┆ {{.Commentary  }}
┃ Reviews     ┆ {{.Reviews     }}
┃ News        ┆ {{.News        }}
┃ Extra       ┆ {{.Extra       }}
┃ Screenshots ┆ {{.Screenshots }}
┃ RelatedApps ┆ {{.RelatedApps }}
┃ SiblingApps ┆ {{.SiblingApps }}
┣┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈
┃ Description
{{.Description }}
┣┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈
┃ ReleaseNote
{{.ReleaseNote }}
┣┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈
┃ ReleaseTime ┆ {{.ReleaseTime }}
┃ CrawledTime ┆ {{.CrawledTime }}
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`
)

var appTmpl, _ = template.New("wdj.app").Parse(appTemplate)

// AppPageURL 根据 PkgName生成豌豆荚页面URL
func AppPageURL(id string) string {
	return AppPagePrefix + id
}

// searchURL 会根据关键词生成查询URL，并转义相关关键词
func searchURL(keyword string) string {
	return searchURLPrefix + url.QueryEscape(keyword)
}

// 应用定义
type App struct {
	Source      string                  // 数据来源，固定为豌豆荚`wdj`
	ID          string   `sql:",pk"`    // 标识 id
	Name        string                  // 名称 name
	URL         string                  // 页面 url
	Icon        string                  // 图标 icon
	Link        string                  // 下载 link
	Version     string                  // 版本 version
	Vendor      string                  // 厂商 vendor
	Genre       string                  // 分类,豌豆荚的分类等于类目categories第一项 genre
	Tags        []string  `pg:",array"` // 标签 tags
	Categories  []string  `pg:",array"` // 类目 categories
	Price       int64                   // 价格 price 留空
	System      string                  // 系统要求
	Platform    []string  `pg:",array"` // 平台
	Permissions []string  `pg:",array"` // 所需权限 permissions
	Size        int64                   // 大小 size
	Rating      int64                   // 评分 rank
	InstallCnt  int64                   // 安装数 install_cnt
	CommentCnt  int64                   // 评论数 comment_cnt
	Appkey      string                  // 友盟分配的Appkey，留空
	AppID       int64                   // 平台分配的应用ID,豌豆荚无 app_id
	ApkCode     int64                   // 平台分配的Apk代码,豌豆荚无 apk_code
	Subtitle    string                  // 副标题 subtitle
	Commentary  string                  // 编辑评论 commentary
	Description string                  // 应用描述,带有换行符 description
	Reviews     string                  // 客户评论,JSON数组,每项为三元组`["YYYY-MM-DD",<user>,<content>` reviews
	News        string                  // 新闻技巧与攻略，JSON数组,每项为[<title>,<src>,<vendor>] news
	Extra       string                  // 额外信息，目前置空。
	Screenshots []string  `pg:",array"` // 截图列表 screenshots
	RelatedApps []string  `pg:",array"` // 推荐的相关应用 related_apps
	SiblingApps []string  `pg:",array"` // 同一开发者的其他应用，豌豆荚无 sibling_apps
	ReleaseNote string                  // 最近更新日志,带有换行符 release_note
	ReleaseTime time.Time               // 最近更新时间 release_time
	CrawledTime time.Time               // 最近爬取时间 crawl_time
	tableName   struct{} `sql:"wdj"`
}

func (app *App) Parse(doc *goquery.Document) error {
	// quick selectors
	info := doc.Find("dl.infos-list")
	nums := doc.Find("div.num-list")

	// app.ID : apk package_name
	app.ID = getAttr(doc.Find("body"), "data-pn")

	// app.URL
	app.URL = AppPageURL(app.ID)

	// app.Name
	app.Name = getText(doc.Find("p.app-name span.title"))

	if app.ID == "" || app.Name == "" {
		return ErrParse
	}

	// app.Icon
	app.Icon = getAttr(doc.Find("div.app-icon img"), "src")

	// app.Link
	app.Link = getAttr(doc.Find("a.install-btn"), "href")

	// app.Size
	if size := getAttr(info.Find("meta[itemprop=fileSize]"), "content"); size != "" {
		app.Size, _ = bytesToInt(size)
	}

	// app.Rating
	if rank := getText(nums.Find("span.love i")); rank != "" && rank != "暂无" {
		app.Rating, _ = parsePercentInt(rank)
	}

	// app.AppID
	// app.ApkCode
	// 豌豆荚无此数据

	// app.InstallCnt
	if install := getText(nums.Find("i[itemprop=interactionCount]")); install != "" {
		app.InstallCnt, _ = parseZhNumber(install)
	}

	// app.CommentCnt
	if comment := getText(nums.Find("a.comment-open i")); comment != "" {
		app.CommentCnt, _ = parseZhNumber(comment)
	}

	// app.Genre
	// 题材分类是Categories的第一项，在Categories解析后补充设置

	// app.Vendor
	app.Vendor = getText(info.Find("span.dev-sites"))

	// app.System
	if system, ok := getFistTextNode(info.Find("dd.perms")); ok {
		app.System = strings.TrimRight(strings.TrimLeft(system, "Android "), " 以上")
	}

	// app.Device
	// 对于豌豆荚始终置空

	// app.Version
	app.Version = getText(info.Find("dd:nth-last-of-type(3)"))

	// app.Subtitle
	app.Subtitle = getText(doc.Find("p.tagline"))

	// app.Commentary
	app.Commentary = getRichText(doc.Find("div.editorComment div.con"))

	// app.Description
	app.Description = getRichText(doc.Find("div.desc-info div.con"))

	// app.ReleaseNote
	app.ReleaseNote = getRichText(doc.Find("div.change-info div"))

	// app.News
	var news [][3]string
	doc.Find("ul.app-news-list > li").Map(func(ind int, s *goquery.Selection) string {
		var entry [3]string
		entry[0] = getText(s.Find("p a"))
		entry[1] = getAttr(s.Find("p a"), "href")
		entry[2] = strings.TrimLeft(getText(s.Find("span")), "来自：")
		if entry[0] == "" || entry[1] == "" {
			return ""
		}
		news = append(news, entry)
		return ""
	})
	if body, err := json.Marshal(news); err == nil {
		if sb := string(body); sb != "" && sb != "null" {
			app.News = sb
		}
	}

	// app.Reviews
	var reviews [][3]string
	doc.Find("ul.comments-list li.normal-li").Map(func(ind int, s *goquery.Selection) string {
		var review [3]string
		review[0] = getText(s.Find("p.first span.name"))
		if t, err := time.Parse(releaseTimeFmt,
			getText(s.Find("p.first span:last-of-type"))); err == nil {
			review[1] = t.Format(normalTimeFmt)
		}
		review[2] = getText(s.Find("p.cmt-content span"))
		if review[0] == "" || review[1] == "" || review[2] == "" {
			return ""
		}
		reviews = append(reviews, review)
		return ""
	})
	if body, err := json.Marshal(reviews); err == nil {
		if sb := string(body); sb != "" && sb != "null" {
			app.Reviews = sb
		}
	}

	// app.Permissions
	app.Permissions = getTextList(doc.Find("span.perms"))

	// app.Screenshots
	app.Screenshots = getAttrList(doc.Find("img.screenshot-img"), "src")

	// app.Categories
	app.Categories = getTextList(info.Find("dd.tag-box a"))
	if len(app.Categories) > 0 {
		app.Genre = app.Categories[0]
	}

	// app.Tags
	app.Tags = getTextList(info.Find("div.tag-box a"))

	// app.RelatedApps
	app.RelatedApps = getAttrList(doc.Find("ul.relative-download li a.d-btn"), "data-app-pname")

	// app.SiblingApps
	// 豌豆荚无此数据

	// app.ReleaseTime
	if ts := getAttr(info.Find("#baidu_time"), "datetime"); ts != "" {
		if t, err := time.Parse(releaseTimeFmt, ts); err == nil {
			app.ReleaseTime = t
		}
	}

	// app.CrawledTime
	app.CrawledTime = time.Now()

	// app.Source
	app.Source = "wdj"
	return nil
}

// Print 打印出人类可读版本的应用信息
func (app *App) Print() {
	if err := appTmpl.Execute(os.Stdout, app); err != nil {
		fmt.Println(err.Error())
	}
	return
}

func (app *App) Save(db *pg.DB) error {
	_, err := db.Model(app).
		OnConflict("(id) DO UPDATE").
		Set("source= ?source").
		Set("name= ?name").
		Set("url= ?url").
		Set("icon= ?icon").
		Set("link= ?link").
		Set("version= ?version").
		Set("vendor= ?vendor").
		Set("genre= ?genre").
		Set("tags= ?tags").
		Set("categories= ?categories").
		Set("price= ?price").
		Set("system= ?system").
		Set("platform= ?platform").
		Set("permissions= ?permissions").
		Set("size= ?size").
		Set("rating= ?rating").
		Set("install_cnt= ?install_cnt").
		Set("comment_cnt= ?comment_cnt").
		Set("appkey= ?appkey").
		Set("app_id= ?app_id").
		Set("apk_code= ?apk_code").
		Set("subtitle= ?subtitle").
		Set("commentary= ?commentary").
		Set("description= ?description").
		Set("reviews= ?reviews").
		Set("news= ?news").
		Set("extra= ?extra").
		Set("screenshots= ?screenshots").
		Set("related_apps= ?related_apps").
		Set("sibling_apps= ?sibling_apps").
		Set("release_note= ?release_note").
		Set("release_time= ?release_time").
		Set("crawled_time= ?crawled_time").
		Insert()
	return err
}

// Search 会使用豌豆荚搜索，并返回所有搜索出的PackageName
func Search(keyword string) (apks []string, error error) {
	// 搜索结果第一页
	doc, err := goquery.NewDocument(searchURL(keyword))
	if err != nil {
		return nil, err
	}

	// 获取页面上所有的应用PkgName与其他的列表页URL
	apkMap := sync.Map{}
	initApks, pages := parseSearchPage(doc)
	for _, apk := range initApks {
		apkMap.Store(apk, nil)
	}

	// 处理后续的页面
	wg := sync.WaitGroup{}
	for _, pageURL := range pages {
		wg.Add(1)
		go func(pageURL string) {
			defer wg.Done()
			if doc, err := goquery.NewDocument(pageURL); err == nil {
				for _, apk := range getAttrList(doc.Find("li.search-item > a"),
					"data-app-pname") {
					apkMap.Store(apk, nil)
				}
			}
		}(pageURL)
	}
	wg.Wait()

	apkMap.Range(func(key, value interface{}) bool {
		apks = append(apks, key.(string))
		return true
	})

	// dedupe
	return
}

func parseSearchPage(doc *goquery.Document) (apks, pages []string) {
	// 不是想要的页面
	apks = getAttrList(doc.Find("li.search-item > a"), "data-app-pname")
	pages = getAttrList(doc.Find(`a.page-item:not(a.current):not(a.prev-page):not(a.next-page)`), "href")
	return
}
