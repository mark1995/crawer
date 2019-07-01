package parser

import (
	"bookcrawer/models"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/lexkong/log"
	"regexp"
	"strconv"
	"time"
)

const (
	MAXTRYCOUNT          = 5
	SINGLEBOOKGOROUTINES = 10
)

type SBook struct {
	Name     string
	Author   string
	Url      string
	Chapters []*SChapter
}

type SChapter struct {
	Title   string
	Url     string
	Order   int
	Pre     int
	Next    int
	Content string
}
type BookTextSpider struct {
}

func (self *BookTextSpider) ParserUrl(url string) error {
	book := SBook{}
	doc, err := getDoc(url)
	if err != nil {
		return nil
	}
	book.Name = doc.Find(".book-info").Find("h1").Find("em").Text()
	book.Author = doc.Find("a[class=writer]").Text()
	book.Url = url
	fmt.Println(book)
	b, err := models.GetBookByName(book.Name)
	fmt.Println(b, err)
	if err != nil || b == nil {
		b := models.Book{Name: book.Name, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		models.BookAdd(&b)
	}
	reg := "(\\d+)"
	re := regexp.MustCompile(reg)
	doc.Find("div[class=volume]").Each(func(index int, sel *goquery.Selection) {
		sel.Find("ul").Each(func(i int, childSel *goquery.Selection) {
			chapname := sel.Find("li>a").Text()
			texturl, _ := sel.Find("li>a").Attr("href")
			chapnums := re.FindAllString(chapname, -1)
			chapnum := 0
			if len(chapnums) > 0 {
				chapnum, _ = strconv.Atoi(chapnums[0])
			}
			fmt.Println(chapnum)
			schapter := &SChapter{Title: chapname, Url: texturl, Order: chapnum, Pre: chapnum - 1, Next: chapnum + 1}
			book.Chapters = append(book.Chapters, schapter)
		})
	})
	log.Infof("%v", book)
	// 让一个文章的goroutines数量控制在10个
	channel := make(chan struct{}, SINGLEBOOKGOROUTINES)
	for _, chapter := range book.Chapters {
		channel <- struct{}{}
		go SpiderChapter(b.Id, chapter, channel)
	}

	for i := 0; i < SINGLEBOOKGOROUTINES; i++ {
		channel <- struct{}{}
	}
	close(channel)
	return nil
}

func SpiderChapter(id int, sch *SChapter, c chan struct{}) {
	defer func() {
		<-c
	}()
	doc, err := getDoc(sch.Url)
	if err != nil {
		log.Errorf(err, "chapter url %s ", sch.Url)
	}
	var buf bytes.Buffer
	doc.Find("div[read-content j_readCount>p]").Each(func(i int, seq *goquery.Selection) {
		buf.WriteString("<p>")
		buf.WriteString(seq.Text())
		buf.WriteString("</p>")
	})
	sch.Content = buf.String()

	ch := models.Chapter{
		ChapterId: sch.Order,
		Title:     sch.Title,
		Content:   sch.Content,
		Sort:      sch.Order,
		Pre:       sch.Pre,
		Next:      sch.Next,
	}
	err = models.ChapterAdd(&ch)
	if err != nil {
		log.Errorf(err, "chapter add fail %v", ch)
	}

}
