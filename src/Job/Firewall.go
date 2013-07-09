package Job

import (
	"encoding/xml"
	"Utils"
	"Config"
	"net/http"
	"io/ioutil"
	"time"
	"strings"
	"strconv"
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
	arrRealtime := []string{}
	arrRealtimeip := []string{}
	for _, v1 := range data {
			wptotbps := strings.Split(v1.WPTotBps, " / ");//>100M
			wptcpbps := strings.Split(v1.WPTcpBps, " / "); //>50M
			wpsynbps := strings.Split(v1.WPSynBps, " / "); //>10M
			wpudpbps := strings.Split(v1.WPUdpBps, " / "); //>1M
			wpicmpbps := strings.Split(v1.WPIcmpBps, " / "); //>1M
			wpothbps := strings.Split(v1.WPOthBps, " / "); //>1M
			in_totbps, _ := strconv.ParseFloat(wptotbps[0], 64)
			in_tcpbps, _ := strconv.ParseFloat(wptcpbps[0], 64)
			in_synbps, _ := strconv.ParseFloat(wpsynbps[0], 64)
			in_udpbps, _ := strconv.ParseFloat(wpudpbps[0], 64)
			in_icmpbps, _ := strconv.ParseFloat(wpicmpbps[0], 64)
			in_othbps, _ := strconv.ParseFloat(wpothbps[0], 64)
			arrRealtime = append(arrRealtime, fmt.Sprintf("('%s','%s','%s','%s','%s', '%s','%s','%s','%s','%s','%s', '%s', '%s','%s','%s','%s','%s','%s', '%s', '%s','%s','%s','%s','%s','%s', '%s', '%s','%s','%s','%s','%s','%s','%s')", v1.Sip, v1.WDConn,v1.WDTcpConn,v1.WDUdpConn,v1.WDOthConn,v1.WPTotBps,v1.WPTcpBps,v1.WPSynBps,v1.WPUdpBps,v1.WPIcmpBps,v1.WPOthBps,v1.WPFragBps,v1.DPTotBps,v1.DPTcpBps,v1.DPSynBps,v1.DPUdpBps,v1.DPIcmpBps,v1.DPOthBps,v1.DPFragBps,v1.WPTotPps,v1.WPTcpPps,v1.WPSynPps,v1.WPUdpPps,v1.WPIcmpPps,v1.WPOthPps,v1.WPFragPps,v1.DPTotPps,v1.DPTcpPps,v1.DPSynPps,v1.DPUdpPps,v1.DPIcmpPps,v1.DPOthPps,v1.DPFragPps))
			arrRealtimeip = append(arrRealtimeip, fmt.Sprintf("'%s'", v1.Sip))
			if true || in_totbps > 100 || in_tcpbps >50 || in_synbps > 10 || in_udpbps > 1 || in_icmpbps > 1 || in_othbps > 1 {
				//fmt.Printf("in_totbps:%+v,in_tcpbps:%+v,in_synbps:%+v,in_udpbps:%+v,in_icmpbps:%+v =>db\n", in_totbps,in_tcpbps,in_synbps,in_udpbps,in_icmpbps)
				arr = append(arr, fmt.Sprintf("('%s','%s','%s','%s','%s', '%s','%s','%s','%s','%s','%s', '%s', '%s','%s','%s','%s','%s','%s', '%s', '%s','%s','%s','%s','%s','%s', '%s', '%s','%s','%s','%s','%s','%s','%s')", v1.Sip, v1.WDConn,v1.WDTcpConn,v1.WDUdpConn,v1.WDOthConn,v1.WPTotBps,v1.WPTcpBps,v1.WPSynBps,v1.WPUdpBps,v1.WPIcmpBps,v1.WPOthBps,v1.WPFragBps,v1.DPTotBps,v1.DPTcpBps,v1.DPSynBps,v1.DPUdpBps,v1.DPIcmpBps,v1.DPOthBps,v1.DPFragBps,v1.WPTotPps,v1.WPTcpPps,v1.WPSynPps,v1.WPUdpPps,v1.WPIcmpPps,v1.WPOthPps,v1.WPFragPps,v1.DPTotPps,v1.DPTcpPps,v1.DPSynPps,v1.DPUdpPps,v1.DPIcmpPps,v1.DPOthPps,v1.DPFragPps))
			}
	}
	dbs, err := sql.Open("mysql", Config.GetDb().HzFirewall)
	defer dbs.Close()
	if err != nil {
		Utils.LogInfo("can't connect to db:%v", err)
		return
	}
	
	if len(arrRealtime) > 0 {
		sql := fmt.Sprintf("DELETE FROM `firewall_hz_realtime` WHERE sip in (%s)", strings.Join(arrRealtimeip, ","))
		dbs.Exec(sql)
		sql = fmt.Sprintf("insert into `firewall_hz_realtime` (sip,wdconn,wdtcpconn,wdudpconn,wdothconn,wptotbps,wptcpbps,wpsynbps,wpudpbps,wpicmpbps,wpothbps,wpfragbps,dptotbps,dptcpbps,dpsynbps,dpudpbps,dpicmpbps,dpothbps,dpfragbps,wptotpps,wptcppps,wpsynpps,wpudppps,wpicmppps,wpothpps,wpfragpps,dptotpps,dptcppps,dpsynpps,dpudppps,dpicmppps,dpothpps,dpfragpps) values %s", strings.Join(arrRealtime, ","))
		_, err = dbs.Exec(sql)
		if err != nil {
			Utils.LogInfo("can't query :%v, sql:\n%q\n", err, sql)
		}		 
	}	
	if len(arr) > 0 {
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
