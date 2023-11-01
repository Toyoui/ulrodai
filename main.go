package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

type Response struct {
	Code int `json:"code"`
	Data struct {
		Items []Item `json:"items"`
	} `json:"data"`
}

type Item struct {
	ID int `json:"id"`
	//ColumnID     int    `json:"column_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Published_at string `json:"published_At"`
}

var uniitems []Item
var apiURL = "https://xizhi.qqoq.net/XZ0683393781ffb434c949f89b1acf4acf.channel"

func clearScreen() {
	switch goos := runtime.GOOS; goos {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func isItemIDExist(itemID int) bool {
	file, err := os.Open("itemid.txt")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == strconv.Itoa(itemID) {
			return true
		}
	}

	return false
}

func main() {
	for {
		resp, err := http.Get("https://www.odaily.news/api/pp/api/info-flow/newsflash_columns/newsflashes?b_id=&per_page=3")
		if err != nil {
			fmt.Println("请求出错：", err)
		} else {
			defer resp.Body.Close()

			// 将响应体的数据转换为UTF-8编码
			utf8Body, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
			if err != nil {
				fmt.Println("转换编码出错：", err)
			} else {
				result, err := ioutil.ReadAll(utf8Body)
				if err != nil {
					fmt.Println("读取转换后的数据出错：", err)
				} else {
					// 解析JSON数据
					var response Response
					err := json.Unmarshal(result, &response)
					if err != nil {
						fmt.Println("解析JSON数据出错：", err)
					} else {
						fmt.Println("解析到的数据：")

						for _, item := range response.Data.Items {
							if isItemIDExist(item.ID) {
								clearScreen()
								fmt.Println("存在相同文章")
								continue
							}

							// 将item信息添加到结构体数组中
							uniitems = append(uniitems, Item{
								ID:           item.ID,
								Published_at: item.Published_at,
								Title:        item.Title,
								Description:  item.Description,
							})
						}

						// 遍历输出结构体数组中的item信息
						for _, item := range uniitems {
							fmt.Println("ID:", item.ID)
							fmt.Println("Published_at:", item.Published_at)
							fmt.Println("Title:", item.Title)
							fmt.Println("Description:", item.Description)
							fmt.Println("==============================")

							// 准备要发送的数据
							title := item.Title
							content := "发布时间: " + "\n" + item.Published_at + "\n" + "\n" + "概述: " + "\n" + "\n" + item.Description
							// 使用url.Values来保存POST请求的参数
							payload := url.Values{}
							payload.Set("title", title)
							payload.Set("content", content)
							// 发送POST请求
							resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))
							if err != nil {
								fmt.Println("发送POST请求失败：", err)
								continue
							}
							defer resp.Body.Close()
							// 检查响应状态码
							if resp.StatusCode != http.StatusOK {
								fmt.Println("POST请求返回非200状态码：", resp.StatusCode)
								continue
							}
							// 输出响应结果
							fmt.Println("POST请求成功")

							// 将item.ID写入itemid.txt文件
							file, err := os.OpenFile("itemid.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
							if err != nil {
								fmt.Println("写入itemid.txt文件出错：", err)
								continue
							}
							defer file.Close()

							_, err = file.WriteString(strconv.Itoa(item.ID) + "\n")
							if err != nil {
								fmt.Println("写入itemid.txt文件出错：", err)
							}
						}

						// 清空items切片
						uniitems = nil

					}
				}
			}
		}
		time.Sleep(20 * time.Second)
	}
}
