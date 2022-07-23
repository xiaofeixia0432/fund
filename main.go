package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"httpclient/utils"
)

// var fundreg = `^\[.*\]`

type Fund struct {
	ID                int     // 基金序号
	Code              string  // 基金代码
	Name              string  // 基金简称
	Date              string  // 日期
	UnitValue         float64 // 单位净值
	TotalValue        float64 // 累计净值
	DaySwellRate      string  // 日增长率
	WeekSwellRate     string  // 近一周增长率
	MothSwellRate     string  // 近一个月增长率
	TreeMothSwellRate string  // 近三个月增长率
	SixMothSwellRate  string  // 近六个月增长率
	YearSwellRate     string  // 近一年增长率
	TwoYearSwellRate  string  // 近两年增长率
	TreeYearSwellRate string  // 近三年增长率
	ThisYearSwellRate string  // 近年来增长率
	CreateSwellRate   string  // 成立以来
	CustomRate        string  // 自定义
	Fee               string  // 手续费
	IsBuy             bool    // 是否可以购买
}

func CreateExcel(file string) *os.File {
	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("create excel file is error, err: %v\n", err)
		return nil
	}
	return f
}

func WriteExcelHeader(s *os.File) {
	header := []string{"基金序号", "基金代码", "基金简称", "日期", "单位净值", "累计净值",
		"日增长率", "近一周增长率", "近一个月增长率", "近三个月增长率", "近六个月增长率",
		"近一年增长率", "近两年增长率", "近三年增长率", "近年来增长率", "成立以来",
		"自定义", "手续费", "是否可以购买"}
	_, err := s.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	if err != nil {
		fmt.Printf("writeExcelHeader is error, err: %v\n", err)
		return
	}
	w := csv.NewWriter(s)
	err = w.Write(header)
	if err != nil {
		fmt.Printf("writeExcelHeader is error, err: %v\n", err)
		return
	}
	w.Flush()
}

func CloseFile(s *os.File) {
	s.Close()
}

func InsertFund(s *os.File, fund *Fund) {
	strconv.FormatFloat(fund.TotalValue, 'E', -1, 64)
	data := []string{strconv.Itoa(fund.ID), fund.Code, fund.Name, fund.Date, strconv.FormatFloat(fund.UnitValue,
		'E', -1, 64),
		strconv.FormatFloat(fund.TotalValue, 'E', -1, 64), fund.DaySwellRate, fund.WeekSwellRate,
		fund.MothSwellRate, fund.TreeMothSwellRate, fund.SixMothSwellRate, fund.YearSwellRate, fund.TwoYearSwellRate,
		fund.TreeYearSwellRate, fund.ThisYearSwellRate, fund.CreateSwellRate, fund.CustomRate, fund.Fee, "Y"}
	w := csv.NewWriter(s)
	w.Write(data)
	w.Flush()
}

func StoreFund(f *Fund) error {
	sql := "insert into fund(code,name, date,unit_value ,total_value ,dayswell_rate,weekswell_rate ,monthswell_rate ,threemonthswell_rate ,sixmonthswell_rate,yearswell_rate,twoyearswell_rate ,threeyearswell_rate ,thisyearwell_rate ,createswell_rate ,custom_rate ,fee, isbuy) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := utils.Db.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		fmt.Printf("err: %v", err)
		return err
	}
	_, err = stmt.Exec(f.Code, f.Name, f.Date, f.UnitValue, f.TotalValue, f.DaySwellRate, f.WeekSwellRate,
		f.MothSwellRate, f.TreeMothSwellRate, f.SixMothSwellRate, f.YearSwellRate, f.TwoYearSwellRate,
		f.TreeYearSwellRate, f.ThisYearSwellRate, f.CreateSwellRate, f.CustomRate, f.Fee, f.IsBuy)
	if err != nil {
		return err
	}
	return nil
}

func CountPage(url string) (itemnum, pagenum int, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET",
		url,
		nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set(utils.HeaderCookieKey, utils.HeaderCookieValue)
	req.Header.Set(utils.HeaderContentTypeKey, utils.HeaderContentTypeValue)
	req.Header.Set(utils.HeaderHostKey, utils.HeaderHostValue)
	req.Header.Set(utils.HeaderRefererKey, utils.HeaderRefererValue)
	req.Header.Set(utils.HeaderUserAgentKey, utils.HeaderUserAgentValue)
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	// var i map[string]interface{}
	pagefirstindex := bytes.Index(body, []byte("pageNum:"))
	pagelastindex := bytes.LastIndex(body, []byte(",allPages"))
	firstindex := bytes.Index(body, []byte("allPages:"))
	lastindex := bytes.LastIndex(body, []byte(",allNum"))
	fmt.Println(string(body))
	page, err := strconv.Atoi(string(body[firstindex+9 : lastindex]))
	item, err := strconv.Atoi(string(body[pagefirstindex+8 : pagelastindex]))
	fmt.Printf("page:%v, totalpage:%v\n", item, page)
	return item, page, nil
}

func StoreFundData(url string, s *os.File, items, pages int) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET",
		url,
		nil)
	if err != nil {
		return err
	}
	req.Header.Set(utils.HeaderCookieKey, utils.HeaderCookieValue)
	req.Header.Set(utils.HeaderContentTypeKey, utils.HeaderContentTypeValue)
	req.Header.Set(utils.HeaderHostKey, utils.HeaderHostValue)
	req.Header.Set(utils.HeaderRefererKey, utils.HeaderRefererValue)
	req.Header.Set(utils.HeaderUserAgentKey, utils.HeaderUserAgentValue)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var i interface{}
	firstindex := bytes.Index(body, []byte{'['})
	lastindex := bytes.LastIndex(body, []byte{']'})

	data := body[firstindex : lastindex+1]
	// fmt.Printf("nn: %v\n", nn)
	err = json.Unmarshal(data, &i)
	if err != nil {
		fmt.Printf("json unmarshal is error, err:%v\n", err)
		return err
	}
	// fmt.Printf("i: %v\n", i)
	// b := re.FindAllStringSubmatch(string(body), -1)
	content := i.([]interface{})
	// fmt.Println(len(content))
	// fmt.Println(reflect.TypeOf(i))
	// fmt.Println(reflect.TypeOf(content[0]), reflect.TypeOf(content[1]))
	for index, value := range content {
		funddata, ok := value.(string)
		if ok {
			msg := strings.Split(funddata, ",")
			fund := &Fund{}
			fund.ID = items*(pages-1) + index + 1
			fund.Code = msg[0]
			fund.Name = msg[1]
			fund.Date = msg[3]
			fund.UnitValue, err = strconv.ParseFloat(msg[4], 64)
			if err != nil {
				fmt.Printf("unit value parsefloat error, err:%v\n", err)
				fund.UnitValue = 0.0
			}
			fund.TotalValue, err = strconv.ParseFloat(msg[5], 64)
			if err != nil {
				fmt.Printf("total value parsefloat error, err:%v\n", err)
				fund.TotalValue = 0.0
			}
			fund.DaySwellRate = msg[6]
			fund.WeekSwellRate = msg[7]
			fund.MothSwellRate = msg[8]
			fund.TreeMothSwellRate = msg[9]
			fund.SixMothSwellRate = msg[10]
			fund.YearSwellRate = msg[11]
			fund.TwoYearSwellRate = msg[12]
			fund.TreeYearSwellRate = msg[13]
			fund.ThisYearSwellRate = msg[14]
			fund.CreateSwellRate = msg[15]
			fund.CustomRate = msg[18]
			fund.Fee = msg[20]
			fund.IsBuy = true
			// 插入csv文件
			InsertFund(s, fund)
			// 插入数据库中
			StoreFund(fund)
		}
	}
	return nil
}

func main() {
	// 后续改成配置文件加载
	utils.DefaultDB = &utils.MysqlDBConf{
		Ip:           "10.10.2.66",
		Port:         3306,
		Database:     "fund",
		User:         "root",
		Password:     "xinze123",
		Charset:      "utf8",
		MaxConntions: 100,
		MaxIdles:     50,
	}
	err := utils.InitDB(utils.DefaultDB)
	if err != nil {
		fmt.Printf("init db is error, err: %v\n", err)
		return
	}
	// 创建excel,写入表头信息
	f := CreateExcel("fund.xls")
	WriteExcelHeader(f)

	items, pages, err := CountPage(utils.MastURL)
	if err != nil {
		fmt.Printf("countpage is error, err:%v\n", err)
		return
	}
	for i := 1; i <= pages; i++ {
		pageurl := fmt.Sprintf("http://fund.eastmoney.com/data/rankhandler.aspx?op=ph&dt=kf&ft=all&rs=&gs=0&sc=jc&st=asc&sd=2021-04-04&ed=2022-04-04&qdii=&tabSubtype=,,,,,&pi=%d&pn=50&dx=1&v=0.6160846806292868", i)
		StoreFundData(pageurl, f, items, i)
	}
	CloseFile(f)
}
