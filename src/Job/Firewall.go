package Job

import (
	"encoding/xml"
	"Utils"
	"Config"
	"net/http"
	"io/ioutil"
	"time"
	"strings"
	"fmt"
	"database/sql"
	_"github.com/go-sql-driver/mysql"	
)

type Firewall struct {
}

//XML结构
type NlMsgDb struct {
    Count []count `xml:"count"`
}

type count struct {
    Sip string `xml:"sip"`
    WDConn,WDTcpConn,WDUdpConn,WDOthConn,WPTotBps,WPTcpBps,WPSynBps,WPUdpBps,WPIcmpBps,WPOthBps,WPFragBps,DPTotBps,DPTcpBps,DPSynBps,DPUdpBps,DPIcmpBps,DPOthBps,DPFragBps,WPTotPps,WPTcpPps,WPSynPps,WPUdpPps,WPIcmpPps,WPOthPps,WPFragPps,DPTotPps,DPTcpPps,DPSynPps,DPUdpPps,DPIcmpPps,DPOthPps,DPFragPps string
}

func (f *Firewall) Run() {
	for {	
		Utils.LogInfo("Firewall delay %d Second", Config.GetLoopTime().Firewall)
		time.Sleep(time.Second * Config.GetLoopTime().Firewall)

		go GetHuzhouFirewall(Config.GetUrl().HuzhouFW1)
		go GetHuzhouFirewall(Config.GetUrl().HuzhouFW2)

	}	
}

func ParseData(data []count) {
	arr := []string{}
	for _, v1 := range data {
			arr = append(arr, fmt.Sprintf("('%s','%s','%s','%s','%s', '%s','%s','%s','%s','%s','%s', '%s', '%s','%s','%s','%s','%s','%s', '%s', '%s','%s','%s','%s','%s','%s', '%s', '%s','%s','%s','%s','%s','%s','%s')", v1.Sip, v1.WDConn,v1.WDTcpConn,v1.WDUdpConn,v1.WDOthConn,v1.WPTotBps,v1.WPTcpBps,v1.WPSynBps,v1.WPUdpBps,v1.WPIcmpBps,v1.WPOthBps,v1.WPFragBps,v1.DPTotBps,v1.DPTcpBps,v1.DPSynBps,v1.DPUdpBps,v1.DPIcmpBps,v1.DPOthBps,v1.DPFragBps,v1.WPTotPps,v1.WPTcpPps,v1.WPSynPps,v1.WPUdpPps,v1.WPIcmpPps,v1.WPOthPps,v1.WPFragPps,v1.DPTotPps,v1.DPTcpPps,v1.DPSynPps,v1.DPUdpPps,v1.DPIcmpPps,v1.DPOthPps,v1.DPFragPps))
	}	
	if len(arr) > 0 {
		dbs, err := sql.Open("mysql", Config.GetDb().HzFirewall)
		if err != nil {
			Utils.LogInfo("can't connect to db:%v", err)
		}
		defer dbs.Close()		
		sql := fmt.Sprintf("insert into `firewall_hz` (sip,wdconn,wdtcpconn,wdudpconn,wdothconn,wptotbps,wptcpbps,wpsynbps,wpudpbps,wpicmpbps,wpothbps,wpfragbps,dptotbps,dptcpbps,dpsynbps,dpudpbps,dpicmpbps,dpothbps,dpfragbps,wptotpps,wptcppps,wpsynpps,wpudppps,wpicmppps,wpothpps,wpfragpps,dptotpps,dptcppps,dpsynpps,dpudppps,dpicmppps,dpothpps,dpfragpps) values %s", strings.Join(arr, ","))
		//fmt.Println(sql)
		_, err = dbs.Exec(sql)
		if err != nil {
			Utils.LogInfo("can't query :%v, sql:\n%q\n", err, sql)
		}
	}	
}
func GetHuzhouFirewall(url string)(err error) {
	f := NlMsgDb{}
	var myTransport http.RoundTripper = &http.Transport{
		    Proxy:                 http.ProxyFromEnvironment,
		    ResponseHeaderTimeout: time.Second * 3,
	}

	var myClient = &http.Client{Transport: myTransport}

	resp, err := myClient.Get(url)
	if err != nil {
		Utils.LogInfo("http Client: %v", err)
		return err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Utils.LogInfo("Firewall fetch url %s => %s", url, err)
	}else{
		err = xml.Unmarshal(res, &f)
		if err != nil {
			Utils.LogInfo("Parse xml: %v", err)
			return err
		}
		ParseData(f.Count)
	}
	return  nil
}
