package wdj

import (
	"os"
	"io/ioutil"
	"strings"
	"unicode"
)

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

/**************************************************************\
* Auxiliary functions
***************************************************************/

// getText will extract text from selector and trim space
func getText(selection *goquery.Selection) (s string) {
	return strings.TrimSpace(selection.Text())
}

// getAttr will extract attr according attrName from selector and trim space
func getAttr(selection *goquery.Selection, attrName string) (s string) {
	s, _ = selection.Attr(attrName)
	return strings.TrimSpace(s)
}

// removeEmpty remove empty string from a string slice
func removeEmpty(input []string) (output []string) {
	for _, str := range input {
		if str != "" {
			output = append(output, str)
		}
	}
	return
}

// getRichText handles multiline text
func getRichText(selection *goquery.Selection) (s string) {
	if s, err := selection.Html(); s != "" && err == nil {
		s = strings.Replace(s, "<br>", "\n", -1)
		s = strings.Replace(s, "<br/>", "\n", -1)
		s = strings.TrimSpace(s)
		return s
	}
	return
}

// getTextList will fetch a list of text of selectors
func getTextList(selection *goquery.Selection) []string {
	res := selection.Map(func(ind int, s *goquery.Selection) string {
		return strings.TrimSpace(s.Text())
	})
	return removeEmpty(res)
}

// getAttrList will fetch a list of attr of selectors
func getAttrList(selection *goquery.Selection, attrName string) []string {
	res := selection.Map(func(ind int, s *goquery.Selection) string {
		attr, _ := s.Attr(attrName)
		return attr
	})
	return removeEmpty(res)
}

// getFistTextNode will fetch data of element node's first child text node
func getFistTextNode(selection *goquery.Selection) (s string, valid bool) {
	if len(selection.Nodes) > 0 {
		node := selection.Nodes[0]
		if node != nil && node.FirstChild != nil {
			if node = node.FirstChild; node.Type == 1 {
				s := strings.TrimSpace(node.Data)
				return s, s != ""
			}
		}
	}
	return
}

// getLastTextNode will fetch data of element node's last child text node
func getLastTextNode(selection *goquery.Selection) (s string, valid bool) {
	if len(selection.Nodes) > 0 {
		node := selection.Nodes[0]
		if node != nil && node.LastChild != nil {
			if node = node.LastChild; node.Type == 1 {
				s := strings.TrimSpace(node.Data)
				return s, s != ""
			}
		}
	}
	return
}

// bytesToInt turns "128k, 25 MB" to bytes count
func bytesToInt(s string) (res int64, ok bool) {
	var i, nFrac int
	var val int64
	var c byte
	var dot bool

	// parse numeric val (omit dot), and length of frac part
Loop:
	for i < len(s) {
		c = s[i]
		switch {
		case '0' <= c && c <= '9':
			val *= 10
			val += int64(c - '0')
			if dot {
				nFrac ++
			}
			i++
		case c == '.':
			dot = true
			i++
		default:
			break Loop
		}
	}
	unit := strings.ToUpper(strings.TrimSpace(s[i:]))

	switch unit {
	case "", "B":
	case "KB", "K":
		val <<= 10
	case "MB", "M":
		val <<= 20
	case "GB", "G":
		val <<= 30
	case "TB", "T":
		val <<= 40
	case "PB", "P":
		val <<= 50
	case "EB", "E":
		val <<= 60
	default:
		return 0, false
	}

	// handle frac
	for j := 0; j < nFrac; j++ {
		val /= 10
	}

	return val, true
}

// parseZhNumber transform "1.28亿" to corresponding integer
func parseZhNumber(s string) (res int64, ok bool) {
	r := []rune(s)
	n := len(r)

	var mutiplier float64;
	switch r[n-1] {
	case rune('万'):
		mutiplier = 10000
		r = r[0:n-1]
	case rune('亿'):
		mutiplier = 100000000
		r = r[0:n-1]
	default:
		mutiplier = 1
	}

	numStr := string(r)
	if dotInd := strings.Index(numStr, "."); dotInd == -1 {
		// not float dot
		if i, err := strconv.Atoi(numStr); err != nil {
			return 0, false
		} else {
			return int64(float64(i) * mutiplier), true
		}
	} else {
		// there's a dot, find it's position and shift value
		for i := 0; i < len(numStr)-dotInd-1; i++ {
			mutiplier /= 10
		}

		numStr = strings.Replace(numStr, ".", "", 1)
		if i, err := strconv.Atoi(numStr); err != nil {
			return 0, false
		} else {
			return int64(float64(i) * mutiplier), true
		}
	}
}

// parsePercentInt will parse "97.00%" to integer 97
func parsePercentInt(s string) (i int64, ok bool) {
	s = strings.TrimSpace(strings.TrimRight(s, "%"))
	if f, err := strconv.ParseFloat(s, 32); err != nil {
		return 0, false
	} else {
		return int64(f), true
	}
}

// squeezeTime will parse "2015年05月10日" to "20150510"
func squeezeTime(s string) string {
	var nb []rune
	for _, ch := range s {
		if unicode.IsDigit(ch) {
			nb = append(nb, ch)
		}
	}
	return string(nb)
}

// buildDocumentFromFile will load a goquery document from filepath
func buildDocumentFromFile(filename string) (doc *goquery.Document, err error) {
	if f, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		return goquery.NewDocumentFromReader(f)
	}
}

// buildDocumentFromURL will load a goquery document from url
func buildDocumentFromURL(url string) (doc *goquery.Document, err error) {
	return goquery.NewDocument(url)
}

func buildDocumentFromApk(id string) (doc *goquery.Document, err error) {
	return goquery.NewDocument(AppPageURL(id))
}

// ReadAllFilename will return a []string contains all file path in that dir
// returns nil when error occurs
func ReadAllFilename(dirname string) []string {
	if !strings.HasSuffix(dirname, "/") {
		dirname = dirname + "/"
	}

	var buf []string
	if files, err := ioutil.ReadDir(dirname); err != nil {
		return nil
	} else {
		for _, file := range files {
			buf = append(buf, dirname+file.Name())
		}
	}

	return buf
}

// Parse will return wandoujia app by PkgName
func Parse(id string) (app *App, err error) {
	doc, err := buildDocumentFromURL(AppPageURL(id))
	if err != nil {
		return nil, err
	}
	app = new(App)
	err = app.Parse(doc)
	return
}
