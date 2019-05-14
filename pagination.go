package Paginator

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

// Paginator
type Paginator struct {
	total       int
	pageSize    int
	current     int
	linkedCount int
	pageKey     string
	request     *http.Request
	url         string
	params      map[string]string
}

// Config
type Config struct {
	PageSize    int
	Current     int
	LinkedCount int
	PageKey     string
	Request     *http.Request
}

type Param struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	OrderBy []string
	ShowSQL bool
}

type PaginatorRecords struct {
	TotalRecord int         `json:"total_record"`
	TotalPage   int         `json:"total_page"`
	Records     interface{} `json:"records"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
	PrevPage    int         `json:"prev_page"`
	NextPage    int         `json:"next_page"`
}

const (
	defaultPageSize    = 15
	defaultCurrent     = 1
	defaultLinkedCount = 5

	defaultPageKey = "page"

	tempBegin        = `<ul class="pagination">`
	tempLinkFirst    = `<li><a href="%s">LinkFirst</a></li>`
	tempLinkPrevious = `<li><a href="%s"><i class="material-icons">chevron_left</i></a></li>`
	tempLinkPage     = `<li><a href="%s">%d</a></li>`
	tempLinkCurrent  = `<li class="active"><a href="%s">%d</a></li>`
	tempLinkNext     = `<li><a href="%s"><i class="material-icons">chevron_right</i></a></li>`
	tempLinkLast     = `<li><a href="%s">LinkLast</a></li>`
	tempEnd          = `</ul>`
)

func Paging(p *Param, result interface{}) *PaginatorRecords {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	done := make(chan bool, 1)
	var paginator PaginatorRecords
	var count int
	var offset int

	go countRecords(db, result, done, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	db.Limit(p.Limit).Offset(offset).Find(result)
	<-done

	paginator.TotalRecord = count
	paginator.Records = result
	paginator.Page = p.Page

	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int) {
	db.Model(anyType).Count(count)
	done <- true
}

// Пользовательские параметры по умолчанию
func Custom(c *Config, total int) *Paginator {
	if c.PageSize <= 0 {
		c.PageSize = defaultPageSize
	}
	p := &Paginator{
		total:       total,
		pageSize:    c.PageSize,
		current:     c.Current,
		linkedCount: c.LinkedCount,

		pageKey: c.PageKey,
		request: nil,
		url:     "",
		params:  map[string]string{},
	}
	if p.current > p.TotalPages() {
		p.current = p.TotalPages()
	}
	if len(p.pageKey) == 0 {
		p.pageKey = defaultPageKey
	}
	return p
}

// New Paginator инициализировать с параметрами по умолчанию
func New(total int) *Paginator {
	c := &Config{
		PageSize:    defaultPageSize,
		Current:     defaultCurrent,
		LinkedCount: defaultLinkedCount,
		PageKey:     defaultPageKey,
	}
	return Custom(c, total)
}

// Request can initialize a network request parameter, such as a specified URL for the page,
// before calling the interface. If the interface is passed in the pageKey request parameter,
// the current page number will be updated according to the parameter content.
func (p *Paginator) Request(r *http.Request) *Paginator {
	if r == nil {
		return p
	}
	current := 0
	params := map[string]string{}
	rawQuerys := strings.Split(r.URL.RawQuery, "&")

	for _, query := range rawQuerys {
		param := strings.Split(query, "=")
		if strings.EqualFold(param[0], p.pageKey) {
			current, _ = strconv.Atoi(param[1])
		} else {
			params[param[0]] = param[1]
		}
	}

	p.url = r.URL.Path
	p.params = params
	if current != 0 {
		p.current = current
	}
	return p
}

// IsFirst If the current page is the first page to return true
func (p *Paginator) IsFirst() bool {
	return p.current == 1
}

// Get path
func (p *Paginator) path(pageNum int) string {
	params := fmt.Sprintf("%s?", p.url)
	for key, value := range p.params {
		params = fmt.Sprintf("%s%s=%s&", params, key, value)
	}
	return fmt.Sprintf("%s%s=%d", params, p.pageKey, pageNum)
}

// FristURL
func (p *Paginator) FristURL() int {
	return defaultCurrent
}

// HasPrevious
func (p *Paginator) HasPrevious() bool {
	return p.current > 1
}

// Previous
func (p *Paginator) Previous() int {
	if !p.HasPrevious() {
		return p.current
	}
	return p.current - 1
}

// PreviousURL
func (p *Paginator) PreviousURL() int {
	return p.Previous()
}

// HasNext
func (p *Paginator) HasNext() bool {
	return p.total > p.current*p.pageSize
}

// Next
func (p *Paginator) Next() int {
	if !p.HasNext() {
		return p.current
	}
	return p.current + 1
}

// NextURL
func (p *Paginator) NextURL() int {
	return p.Next()
}

// IsLast
func (p *Paginator) IsLast() bool {
	if p.total == 0 {
		return true
	}
	return p.total > (p.current-1)*p.pageSize && !p.HasNext()
}

// Total
func (p *Paginator) Total() int {
	return p.total
}

// TotalPages
func (p *Paginator) TotalPages() int {
	if p.total == 0 {
		return 1
	}
	if p.total%p.pageSize == 0 {
		return p.total / p.pageSize
	}
	return p.total/p.pageSize + 1
}

// LastURL
func (p *Paginator) LastURL() int {
	return p.TotalPages()
}

// Current
func (p *Paginator) Current() int {
	return p.current
}

// CurrentURL
func (p *Paginator) CurrentURL() string {
	return p.path(p.Current())
}

// PageSize
func (p *Paginator) PageSize() int {
	return p.pageSize
}

// Page
type Page struct {
	num       int
	isCurrent bool
}

// Num
func (p *Page) Num() int {
	return p.num
}

// IsCurrent
func (p *Page) IsCurrent() bool {
	return p.isCurrent
}

func getMiddleIdx(linkedCount int) int {
	if linkedCount%2 == 0 {
		return linkedCount / 2
	}
	return linkedCount/2 + 1
}

// Pages
func (p *Paginator) Pages() []*Page {
	if p.linkedCount <= 0 {
		return []*Page{}
	}

	// Return only the current page
	if p.linkedCount == 1 || p.TotalPages() == 1 {
		return []*Page{{p.current, true}}
	}

	// The total number of bars is less than the number of bars to be returned, and all page number information is returned.
	if p.TotalPages() <= p.linkedCount {
		pages := make([]*Page, p.TotalPages())
		for i := range pages {
			pages[i] = &Page{i + 1, i+1 == p.current}
		}
		return pages
	}

	linkedRadius := p.linkedCount / 2

	// If linkedCount is odd, the number of pages before and after current is equal
	previousCount, nextCount := linkedRadius, linkedRadius
	// If linkedCount is even, the current page number is 1 more than the back
	if p.linkedCount%2 == 0 {
		nextCount--
	}

	// If current<=previousCount then the required page is from 1 to linkedCount
	// If current>previousCount and current>=TotalPages-nextCount then the required page is from TotalPages-linkedCount+1 to TotalPages
	// The rest is from current-previousCount to current+nextCount

	pages := make([]*Page, p.linkedCount)
	offsetIdx, maxIdx := 1, 1
	if p.current <= previousCount {
		offsetIdx = 1
		maxIdx = p.linkedCount
	} else if p.current > previousCount && p.current >= p.TotalPages()-nextCount {
		offsetIdx = p.TotalPages() - p.linkedCount + 1
		maxIdx = p.TotalPages()
	} else {
		offsetIdx = p.current - previousCount
		maxIdx = p.current + nextCount
	}

	for i := 0; i < maxIdx-offsetIdx+1; i++ {
		pages[i] = &Page{offsetIdx + i, offsetIdx+i == p.current}
	}
	return pages
}

// PageURLs returns the page number and URL information of the current page.
func (p *Paginator) PageURLs() []*PageURL {
	pages := p.Pages()
	pageURLs := make([]*PageURL, len(pages))
	for i := 0; i < len(pages); i++ {
		pageURLs[i] = &PageURL{
			Page:    pages[i],
			pageKey: p.pageKey,
			request: p.request,
			url:     p.url,
			params:  p.params,
		}
	}
	return pageURLs
}

// PageURL The current page's data content
type PageURL struct {
	*Page

	pageKey string            // Query request page number keyword for http request
	request *http.Request     // Network request
	url     string            // Http request url
	params  map[string]string // The set of query parameters for the http request, excluding the pageKey
}

// Num page number
func (p *PageURL) Num() int {
	return p.num
}

// IsCurrent is the current page
func (p *PageURL) IsCurrent() bool {
	return p.isCurrent
}

// Path Current web path address
func (p *PageURL) Path() string {
	params := fmt.Sprintf("%s?", p.url)
	for key, value := range p.params {
		params = fmt.Sprintf("%s%s=%s&", params, key, value)
	}
	return fmt.Sprintf("%s%s=%d", params, p.pageKey, p.num)
}

// PageTemp Gets a page template for pagination results that can be loaded directly in html
func (p *Paginator) PageTemp() template.HTML {
	paths := p.PageURLs()
	if len(paths) == 0 {
		return ""
	}
	if len(paths) == 1 {
		middle := fmt.Sprintf(tempLinkCurrent, paths[0].Path(), paths[0].Num())
		return template.HTML(fmt.Sprintf("%s\n%s\n%s", tempBegin, middle, tempEnd))
	}

	middle := ""
	for _, path := range paths {
		if path.IsCurrent() {
			middle += fmt.Sprintf(tempLinkCurrent, path.Path(), path.Num())
		} else {
			middle += fmt.Sprintf(tempLinkPage, path.Path(), path.Num())
		}
	}

	return template.HTML(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		tempBegin,
		fmt.Sprintf(tempLinkFirst, strconv.Itoa(p.FristURL())),
		fmt.Sprintf(tempLinkPrevious, strconv.Itoa(p.PreviousURL())),
		middle,
		fmt.Sprintf(tempLinkNext, strconv.Itoa(p.NextURL())),
		fmt.Sprintf(tempLinkLast, strconv.Itoa(p.LastURL())),
		tempEnd))
}
