package Job

import (
	"encoding/json"
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

type HuzhouColletion struct {
	Pool []Huzhou
}

type Huzhou struct {
	Ip string `json:"ip"`
	Trafficin string `json:"trafficin"`
	Trafficout string `json:"trafficout"`
	Syn string `json:"syn"`
	Udp string `json:"udp"`
	Icmp string `json:"icmp"`
	Other string `json:"other"`
	//,'type',state,bypass,trafficin,trafficout,syn,udp,icmp,frag,other,dns,tcplinks,udplinks,tcpnew,udpnew string
}

func (hc *HuzhouColletion) FromJson(jsonStr string) error {
    var data = &hc.Pool
    b := []byte(jsonStr)
    return json.Unmarshal(b, data)
}

func (f *Firewall) Run() {
	for {	
		Utils.LogInfo("Firewall delay %d Second", Config.GetLoopTime().Firewall)
		time.Sleep(time.Second * Config.GetLoopTime().Firewall)

		go GetHuzhouFirewall(Config.GetUrl().HuzhouFW1)
		go GetHuzhouFirewall(Config.GetUrl().HuzhouFW2)

	}	
}

func ParseData(data []Huzhou) {
	arr := []string{}
	arrRealtime := []string{}
	arrRealtimeip := []string{}
	for _, v1 := range data {
			ip := v1.Ip;
			trafficin := strings.Split(v1.Trafficin, " / ");//>100M
			trafficout := strings.Split(v1.Trafficout, " / "); //>50M
			syn := strings.Split(v1.Syn, " / "); //>10M
			udp := strings.Split(v1.Udp, " / "); //>1M
			icmp := strings.Split(v1.Icmp, " / "); //>1M
			other := strings.Split(v1.Other, " / "); //>1M
			trafficin_1, _ := strconv.ParseFloat(trafficin[0], 64)
			trafficin_2, _ := strconv.ParseFloat(trafficin[1], 64)
			trafficout_1, _ := strconv.ParseFloat(trafficout[0], 64)
			trafficout_2, _ := strconv.ParseFloat(trafficout[1], 64)
			syn_1, _ := strconv.ParseFloat(syn[0], 64)
			syn_2, _ := strconv.ParseFloat(syn[1], 64)
			udp_1, _ := strconv.ParseFloat(udp[0], 64)
			udp_2, _ := strconv.ParseFloat(udp[1], 64)
			icmp_1, _ := strconv.ParseFloat(icmp[0], 64)
			icmp_2, _ := strconv.ParseFloat(icmp[1], 64)
			other_1, _ := strconv.ParseFloat(other[0], 64)
			other_2, _ := strconv.ParseFloat(other[1], 64)	
			//fmt.Printf("%q,%q,%f,%f", v1.Trafficin,trafficin,trafficin_1,trafficin_2)																	
			arrRealtime = append(arrRealtime, fmt.Sprintf("('%s','%f','%f','%f','%f', '%f','%f','%f','%f','%f','%f', '%f', '%f')", ip, trafficin_1,trafficin_2,trafficout_1,trafficout_2,syn_1,syn_2,udp_1,udp_2,icmp_1,icmp_2,other_1,other_2))
			arrRealtimeip = append(arrRealtimeip, fmt.Sprintf("'%s'", ip))
			if strings.HasPrefix(ip, "61.153.107") { //61开头的为大墙集群
				//trafficin_1 = trafficin_1 * 4
			}
			threshold,_ := strconv.ParseFloat("100", 64)
			if false || trafficin_1 > threshold || syn_1 > threshold || udp_1 > threshold || icmp_1 > threshold || other_1 > threshold {
				//fmt.Printf("trafficin_1:%f,trafficin_2:%f,trafficout_1:%f,trafficout_2:%f,syn_1:%f,syn_2:%f,udp_1:%f,udp_2:%f,icmp_1:%f,icmp_2:%f,other_1:%f,other_2:%f, =>db\n", trafficin_1, trafficin_2, trafficout_1, trafficout_2, syn_1, syn_2, udp_1, udp_2, icmp_1, icmp_2, other_1, other_2)
				arr = append(arr, fmt.Sprintf("('%s','%f','%f','%f','%f', '%f','%f','%f','%f','%f','%f', '%f', '%f')", ip, trafficin_1,trafficin_2,trafficout_1,trafficout_2,syn_1,syn_2,udp_1,udp_2,icmp_1,icmp_2,other_1,other_2))
			}
	}
	dbs, err := sql.Open("mysql", Config.GetDb().HzFirewall)
	defer dbs.Close()
	if err != nil {
		Utils.LogInfo("can't connect to db:%v", err)
		return
	}
	
	if len(arrRealtime) > 0 {
		sql := fmt.Sprintf("DELETE FROM `huzhou_realtime` WHERE ip in (%s)", strings.Join(arrRealtimeip, ","))
		dbs.Exec(sql)
		sql = fmt.Sprintf("insert into `huzhou_realtime` (ip,trafficin_1,trafficin_2,trafficout_1,trafficout_2,syn_1,syn_2,udp_1,udp_2,icmp_1,icmp_2,other_1,other_2) values %s", strings.Join(arrRealtime, ","))
		_, err = dbs.Exec(sql)
		if err != nil {
			Utils.LogInfo("can't query :%v, sql:\n%q\n", err, sql)
		}		 
	}	
	if len(arr) > 0 {
		sql := fmt.Sprintf("insert into `huzhou` (ip,trafficin_1,trafficin_2,trafficout_1,trafficout_2,syn_1,syn_2,udp_1,udp_2,icmp_1,icmp_2,other_1,other_2) values %s", strings.Join(arr, ","))
		//fmt.Println(sql)
		_, err = dbs.Exec(sql)
		if err != nil {
			Utils.LogInfo("can't query :%v, sql:\n%q\n", err, sql)
		}
	}	
}

func GetHuzhouFirewall(url string)(err error) {
	f := HuzhouColletion{}.Pool
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
		err = json.Unmarshal(res, &f)
		if err != nil {
			Utils.LogInfo("Parse json: %v", err)
			return err
		}
		//Utils.LogInfo("res:%s,json: %q", res,f)
		ParseData(f)
	}
	return  nil
}
