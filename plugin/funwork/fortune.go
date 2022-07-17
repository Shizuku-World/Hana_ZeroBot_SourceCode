// Package myfortune 简单的测人品
package funwork

import (
	"encoding/json"
	"fmt"
	"os"

	"math/rand"
	"strconv"
	"time"

	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

const (
	bed = "file:///root/Lucy_Project/tarot/"
)

type card struct {
	Name     string `json:"name"`
	Cardtype string `json:"cardtype"`
	Info     struct {
		Description        string `json:"description"`
		ReverseDescription string `json:"reverseDescription"`
		ImgURL             string `json:"imgUrl"`
	} `json:"info"`
}

type cardset = map[string]card

var (
	jrrpbk   string
	info     string
	uptime   string
	cardMap  = make(cardset, 256)
	reasons  = []string{" | "}
	position = []string{"正位", "逆位"}
	result   map[int64](int)
	egg      map[string](int)
	signTF   map[string](int)
)

func init() {
	signTF = make(map[string](int))
	egg = make(map[string](int))
	result = make(map[int64](int))
	engine.OnFullMatch("今日人品", ctxext.DoOnceOnSuccess(
		func(ctx *zero.Ctx) bool { // 检查 塔罗牌文件是否存在
			data, err := os.ReadFile(engine.DataFolder() + "tarots.json")
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return false
			}
			err = json.Unmarshal(data, &cardMap)
			if err != nil {
				panic(err)
			}
			return true
		},
	)).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			yiyan, err := web.RequestDataWith(web.NewDefaultClient(), "http://ovooa.com/API/yiyan/api.php", "GET", Referer, ua)
			if err != nil {
				return
			} // 获取一言

			p := rand.Intn(2)
			i := rand.Intn(78)
			card := cardMap[(strconv.Itoa(i))]
			name := card.Name
			cardtype := card.Cardtype
			cardurl := card.Info.ImgURL
			if p == 0 {
				info = card.Info.Description
			} else {
				info = card.Info.ReverseDescription
			} // 塔罗牌生成

			// 写的非常恶心 建议有时间赶紧重构x awa
			user := ctx.Event.UserID
			userS := strconv.FormatInt(user, 10)
			now := time.Now().Format("20060102")
			var si string = now + userS // 合成
			rand.Seed(time.Now().UnixNano())
			today := rand.Intn(100)
			dyn := time.Now().Hour()

			switch {
			case dyn <= 6 && dyn >= 0:
				uptime = "凌晨好~还没有睡觉呢 这样不是好孩子哦" // 计算是早上还是晚上
			case dyn <= 11 && dyn > 6:
				uptime = "上午好~呼呼~(*/ω＼*) "
			case dyn <= 14 && dyn > 11:
				uptime = "中午好~吃饭了嘛w 如果没有快去吃饭哦w"
			case dyn <= 18 && dyn > 14:
				uptime = "下午好ww~咱很高兴看到你精力充沛的样子w"
			case dyn <= 24 && dyn > 18:
				uptime = "晚上好吖w~今天过的开心嘛ww"
			}
			// CTRL C + CTRL V
			if signTF[si] == 0 {
				signTF[si] = (1)
				result[user] = (today)
				switch {
				case result[user] <= 20:
					jrrpbk = "[小凶]\n#Lucy抱了抱你~"
				case result[user] > 20 && result[user] < 50:
					jrrpbk = "[小吉]\n#Lucy偷瞄瞄~w"
				case result[user] >= 50 && result[user] < 90:
					jrrpbk = "[中吉]\n#Lucy捏了捏你的脸"
				case result[user] >= 90 && result[user] < 100:
					jrrpbk = "[吉]\n#Lucy摸了摸你的脸"
				case result[user] == 100:
					jrrpbk = "[大吉]\n#好诶~Lucy给你递了张彩票"
				}

				ctx.SendChain(message.At(user),
					message.Text(fmt.Sprintf("\n%s\nLucy正在帮你整理~\n", uptime)),
					message.Text("今日的人品值为", result[user]),
					message.Text(jrrpbk),
					message.Text("\n今日一言:\n"),
					message.Text(helper.BytesToString(yiyan), "\n"),
					message.Text("今日塔罗牌是: \n归类于", cardtype, reasons[rand.Intn(len(reasons))], position[p], " 的 ", name, "\n"),
					message.Image(bed+cardurl),
					message.Text("\n其意义为：\n", info))

			} else {
				ctx.SendChain(message.At(user), message.Text(" 今天已经测过了哦~今日的人品值为", result[user], "呢~"))
			}
			// special Time !
			if result[user] >= 90 && result[user] < 100 && egg[si] == 0 {
				egg[si] = (1)
				img, err := web.RequestDataWith(web.NewDefaultClient(), "http://iw233.fgimax2.fgnwctvip.com/API/Ghs.php?type=json", "GET", Referer, ua)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				picURL := gjson.Get(string(img), "pic").String()
				time.Sleep(time.Second * 2)
				deleteme := ctx.SendChain(message.At(user), message.Text("\n这是今日奖励哦"), message.Text(picURL))
				time.Sleep(time.Second * 20)
				ctx.DeleteMessage(deleteme)
			}
		})
}
