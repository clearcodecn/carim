package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	d = `{"orderDetail":{"activityId":"5f731149e4b0804576a6d512","ticketId":"1dc1138c94124b76aca88fbb7bb62dee","ticketType":1,"num":"1","goodsType":1,"price":80},"areaCode":"86_CN","telephone":"18529570908","orderPerformers":"[{\"documentNumber\":\"46000319930808741X\",\"documentType\":1,\"documentTypeStr\":\"身份证\",\"id\":1704172,\"isSelf\":1,\"name\":\"陈康\",\"showDocumentNumber\":\"460003199******41X\",\"updateAuditStatus\":0,\"userId\":3242612,\"selected\":true}]","customerName":"冼登辉","provinceName":"广东","cityName":"深圳","address":"宝安西乡盐田新二村191栋406","teamId":"","couponId":"","checkCode":"","formToken":"1603704061046569wS63FHtPGfBqNqc0a","payPlatName":"alipaywap","st_flpv":"1603703684243kNw5RJH38Vyvxd3JYrg0","sign":"61c59cf9e1cf003feb1b0e926795bb35","trackPath":"","terminal":"wap"}`
)

func main() {

	x := `Connection: keep-alive
Pragma: no-cache
Cache-Control: no-cache
st_flpv: 1603703684243kNw5RJH38Vyvxd3JYrg0
r: 1603704065982
s: 42979aa544b99ec2fcabbb8a1c966111
terminal: wap
User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1
sign: 61c59cf9e1cf003feb1b0e926795bb35
Content-Type: application/json
Accept: */*
Origin: https://wap.showstart.com
Sec-Fetch-Site: same-origin
Sec-Fetch-Mode: cors
Sec-Fetch-Dest: empty
Referer: https://wap.showstart.com/pages/order/activity/confirm/confirm?sequence=120265&ticketId=1dc1138c94124b76aca88fbb7bb62dee&ticketNum=1
Accept-Encoding: gzip, deflate, br
Accept-Language: zh-CN,zh;q=0.9
Cookie: Hm_lvt_da038bae565bb601b53cc9cb25cdca74=1603703684; o_s=https://wap.showstart.com/pages/passport/login/login; Hm_lpvt_da038bae565bb601b53cc9cb25cdca74=1603703727; u_s=https://wap.showstart.com/pages/order/activity/confirm/confirm?sequence=120265&ticketId=1dc1138c94124b76aca88fbb7bb62dee&ticketNum=1`
	req, _ := http.NewRequest(http.MethodPost, "https://wap.showstart.com/api/wap/order/order.json?sign=61c59cf9e1cf003feb1b0e926795bb35", strings.NewReader(d))

	h := strings.Split(x, "\n")
	for _, v := range h {
		arr := strings.Split(v, ": ")
		req.Header.Add(arr[0], arr[1])
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	xx, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(xx))
}
