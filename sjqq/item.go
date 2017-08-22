package sjqq

import (
	"os"
	"fmt"
	"time"
	"regexp"
	"errors"
	"strconv"
	"strings"
	"text/template"
)

import (
	"github.com/go-pg/pg"
	"github.com/PuerkitoBio/goquery"
)

const (
	AppPagePrefix = "http://sj.qq.com/myapp/detail.htm?apkName="
	appTemplate   = `
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
┃ 应用宝: {{ .ID }} {{ .Name }}
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

var ErrParse = errors.New("parse error")
var appTmpl, _ = template.New("wdj.app").Parse(appTemplate)
// 用于解析底层Script的正则表达式
var (
	pApkCode  = regexp.MustCompile(`apkCode\s*:\s*"(\d+)",`)
	pAppId    = regexp.MustCompile(`appId\s*:\s*"(\d+)",`)
	pDownTime = regexp.MustCompile(`downTimes\s*:\s*"(\d+)",`)
	pIconUrl  = regexp.MustCompile(`iconUrl\s*:\s*"(\S+)",`)
	pDownUrl  = regexp.MustCompile(`downUrl\s*:\s*"(\S+)",`)
)

func AppPageURL(id string) string {
	return AppPagePrefix + id
}

// App 包括手机QQ应用页中能获取的信息
type App struct {
	Source      string                  // 数据来源，固定为应用宝`sjqq`
	ID          string   `sql:",pk"`    // 标识，实质上是PkgName id
	Name        string                  // 名称 name
	URL         string                  // 页面 url
	Icon        string                  // 图标 icon
	Link        string                  // 下载 link
	Version     string                  // 版本 version
	Vendor      string                  // 厂商 vendor
	Genre       string                  // 分类 genre
	Tags        []string  `pg:",array"` // 标签，应用宝无 tags
	Categories  []string  `pg:",array"` // 类目，应用宝无 categories
	Price       int64                   // 价格，应用宝无 price 留空
	System      string                  // 系统要求，应用宝无
	Platform    []string  `pg:",array"` // 平台，应用宝无
	Permissions []string  `pg:",array"` // 所需权限，应用宝较为详细 permissions
	Size        int64                   // 大小 size
	Rating      int64                   // 评分 rating
	InstallCnt  int64                   // 安装数 install_cnt
	CommentCnt  int64                   // 评论数 comment_cnt
	Appkey      string                  // 友盟分配的Appkey，留空
	AppID       int64                   // 应用宝分配的应用ID app_id
	ApkCode     int64                   // 应用宝平台分配的Apk代码 apk_code
	Subtitle    string                  // 副标题，应用宝无 subtitle
	Commentary  string                  // 编辑评论，应用宝无 commentary
	Description string                  // 应用描述，带有换行符 description
	Reviews     string                  // 客户评论，应用宝暂无
	News        string                  // 新闻技巧与攻略，应用宝暂无
	Extra       string                  // 额外信息，目前置空。
	Screenshots []string  `pg:",array"` // 截图列表 screenshots
	RelatedApps []string  `pg:",array"` // 推荐的相关应用 related_apps
	SiblingApps []string  `pg:",array"` // 同一开发者的其他应用 sibling_apps
	ReleaseNote string                  // 最近更新日志,带有换行符 release_note
	ReleaseTime time.Time               // 最近更新时间 release_time
	CrawledTime time.Time               // 最近爬取时间 crawl_time
	tableName   struct{} `sql:"sjqq"`
}

// App_Valid
func (app *App) Valid() bool {
	return app != nil && app.ID != "" && app.Name != ""
}

// App_Print 打印出人类可读版本的应用信息
func (app *App) Print() {
	if err := appTmpl.Execute(os.Stdout, app); err != nil {
		fmt.Println(err.Error())
	}
	return
}

// App_Parse 核心解析逻辑
func (app *App) Parse(doc *goquery.Document) error {
	// app.ID 包名，必需存在
	app.ID = getAttr(doc.Find("a.det-down-btn"), "apk")

	// app.URL
	app.URL = AppPageURL(app.ID)

	// app.Name
	app.Name = getText(doc.Find("div.det-name-int"))

	if app.ID == "" || app.Name == "" {
		return ErrParse
	}

	// quick selector
	oi := doc.Find("div.det-othinfo-container")

	// app.Icon
	app.Icon = getAttr(doc.Find("div.app-icon img"), "src")

	// app.Link
	app.Link = getAttr(doc.Find("a.det-ins-btn"), "ex_url")

	// app.Size
	if size := getText(doc.Find("div.det-size")); size != "" {
		app.Size, _ = bytesToInt(size)
	}

	// app.Rank
	if rank := getText(doc.Find("div.com-blue-star-num")); rank != "" {
		if num, err := strconv.ParseFloat(strings.TrimRight(rank, "分"), 32); err == nil {
			app.Rating = int64(num * 20)
		}
	}

	// app.AppID from script json
	// app.ApkCode from script json

	// app.InstallCnt
	if install := getText(doc.Find("div.det-ins-num")); install != "" {
		app.InstallCnt, _ = parseZhNumber(strings.TrimRight(install, "下载"))
	}

	// app.CommentCnt TBD

	// app.Genre
	app.Genre = getText(doc.Find("#J_DetCate"))

	// app.Vendor
	app.Vendor = getText(oi.Find("div:nth-of-type(6)"))

	// app.Version
	if version := getText(oi.Find("div.det-othinfo-data:nth-of-type(2)")); version != "" {
		app.Version = strings.TrimLeft(version, "V")
	}

	// app.Description
	app.Description = getRichText(doc.Find("div.det-app-data-info:first-of-type"))

	// app.Permissions
	app.Permissions = getTextList(doc.Find("ul.det-othinfo-plist div.r"))

	// app.Screenshots
	app.Screenshots = getAttrList(doc.Find("div.pic-img-box img"), "data-src")

	// app.RelatedApps
	app.RelatedApps = getAttrList(doc.Find("li.det-about-app-box a.com-install-btn"), "apk")

	// app.SiblingApps
	app.SiblingApps = getAttrList(doc.Find("li.det-samedeve-app-box a.com-install-btn"), "apk")

	// app.ReleaseTime
	if rt := getAttr(doc.Find("#J_ApkPublishTime"), "data-apkpublishtime"); rt != "" {
		if num, err := strconv.ParseInt(rt, 10, 64); err == nil {
			app.ReleaseTime = time.Unix(num, 0)
		}
	}

	// 解析最后的script元素: AppId ApkCode IconURL DownURL InstallCnt
	if script := getText(doc.Find("script:last-of-type")); script != "" {
		if start, stop :=
			strings.IndexByte(script, '{'),
			strings.LastIndexByte(script, '}');
			start != -1 && stop != -1 {
			script = script[start:stop+1]

			// APK Code
			if res := pApkCode.FindStringSubmatch(script); res != nil && len(res) > 1 {
				if i, err := strconv.Atoi(res[1]); err == nil {
					app.ApkCode = int64(i)
				}
			}
			if res := pAppId.FindStringSubmatch(script); res != nil && len(res) > 1 {
				if i, err := strconv.Atoi(res[1]); err == nil {
					app.AppID = int64(i)
				}
			}
			if res := pIconUrl.FindStringSubmatch(script); res != nil && len(res) > 1 {
				if iconUrl := res[1]; iconUrl != "" {
					app.Icon = iconUrl
				}
			}
			if res := pDownUrl.FindStringSubmatch(script); res != nil && len(res) > 1 {
				if downUrl := res[1]; downUrl != "" {
					app.Link = downUrl
				}
			}
			if res := pDownTime.FindStringSubmatch(script); res != nil && len(res) > 1 {
				if i, err := strconv.Atoi(res[1]); err == nil {
					app.InstallCnt = int64(i)
				}
			}
		}
	}

	// app.CrawledTime
	app.CrawledTime = time.Now()

	// app.Source
	app.Source = "sjqq"

	return nil
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
