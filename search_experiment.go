package main

import "github.com/sclevine/agouti"
import "log"
import "time"
import "regexp"
import "strings"



// 既存のPageを埋め込む
type ExPage struct {
	*agouti.Page
}

// ページ情報
type PageInfo struct {
	Url, Title string
}

// インプットフィールド
type InputField struct {
	Name, Value string
}

// ログオン情報
type LogonInfo struct {
	Page   *PageInfo
	User   *InputField
	Pass   *InputField
	Submit string
	Button string
}

func (page *ExPage) Logon(info *LogonInfo) {
	//log.Debug("手前の画面へ遷移")
	page.Navigate(info.Page.Url)
	page.FindByButton(info.Button).Submit()

	//log.Debug("認証情報の入力")
	page.FindByID(info.User.Name).Fill(info.User.Value)
	page.FindByID(info.Pass.Name).Fill(info.Pass.Value)

	//log.Debug("ログオン実行")
	page.FindByID(info.Submit).Submit()
}

func main() {
	your_ID := `IDを入力してください`
	your_Pass := `Passwordを入力してください`
	// Chromeを利用することを宣言
	agoutiDriver := agouti.ChromeDriver()
	agoutiDriver.Start()
	defer agoutiDriver.Stop()
	page, _ := agoutiDriver.NewPage()
	expage := &ExPage{page}

	expage.Logon(&LogonInfo{
		Page: &PageInfo{
			Url: "https://utas.adm.u-tokyo.ac.jp/campusweb/campusportal.do",
		},
		User: &InputField{
			Name:  "userNameInput",
			Value:  your_ID,
		},
		Pass: &InputField{
			Name:  "passwordInput",
			Value:  your_Pass,
		},
		Submit: "submissionArea",
		Button: "ログイン",
	})
	page.Navigate("https://utas.adm.u-tokyo.ac.jp/campusweb/campusportal.do?page=main&tabId=si")
    page.FindByID("main-frame-if").SwitchToFrame()
    //少し待たないとチェックボックスをチェックする前に次に進む
    time.Sleep(1 * time.Second)
    page.FindByID("shozokuCd2").Click()
    page.FindByXPath(`//*[@id="rishuSeisekiReferListForm"]/table/tfoot/tr/td/input[1]`).Click()
    
    //contentにhtmlが入っている

	content, err := page.HTML()
	if err != nil {
		log.Printf("Failed to get html: %v", err)
    }
    
    //ここから煩雑　綺麗にしたい
	content = strings.Replace(content, "\n", "", -1)
	r := regexp.MustCompile(`(?im)電気電子情報実験・演習第一.*?</td>.*?</td>.*?</td>.*?</td>.*?</td>.*?<td align="center">(.*?)</td>`)
	r1 := regexp.MustCompile(`(?im)<td align="center">(.*?)</td>`)
	first := r.FindStringSubmatch(content)
	zero := strings.Replace(first[0], `<td align="center">`, "", 3)
	second := r1.FindStringSubmatch(zero)
	third := strings.Replace(second[0], `<td align="center">`, "", -1)
	forth := strings.Replace(third, `</td>`, "", -1)
	log.Println(`電気電子情報実験・演習第一は` + forth + `です。`)
}
