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

type AttackHost struct {
	OnAttackingHosts []OnAttackingHost `json:"onAttackingHosts"`
}

type OnAttackingHost struct {
 	Ip string `json:"name"`
	Trafficin string `json:"trafficin"`
	Trafficout string `json:"trafficout"`
	Synin string `json:"synin"`
	Tcpin string `json:"tcpin"`
	Udpin string `json:"udpin"`
	Icmpin string `json:"icmpin"`
	Otherin string `json:"otherin"`
	Packetsin string `json:"packetsin"`
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

		for _,i := range Config.GetUrl().AttackFireWall {
			go GetAttackData(i)
		}

	}	
}

func ParseAttackHost(data []OnAttackingHost) {
	arr := []string{}
	arrRealtime := []string{}
	arrRealtimeip := []string{}
	for _, v1 := range data {
			ip := v1.Ip;
			if strings.HasPrefix(ip, "117.27.250") == false && strings.HasPrefix(ip, "61.153.107") == false {
				Utils.LogInfo("filter ip : %s", ip)
				continue;
			}

			trafficin := strings.Split(v1.Trafficin, " / ");//>100M
			trafficout := strings.Split(v1.Trafficout, " / "); //>50M
			synin := strings.Split(v1.Synin, " /"); //>10M
			tcpin := strings.Split(v1.Tcpin, " /"); //>10M
			udpin := strings.Split(v1.Udpin, " /"); //>1M
			icmpin := strings.Split(v1.Icmpin, " /"); //>1M
			otherin := strings.Split(v1.Otherin, " /"); //>1M
			packetsin := strings.Split(v1.Packetsin, " / ");
			trafficin_1, _ := strconv.ParseFloat(trafficin[0], 64)
			trafficin_2, _ := strconv.ParseFloat(trafficin[1], 64)
			trafficout_1, _ := strconv.ParseFloat(trafficout[0], 64)
			trafficout_2, _ := strconv.ParseFloat(trafficout[1], 64)
			synin_1, _ := strconv.ParseFloat(synin[0], 64)
			synin_2, _ := strconv.ParseFloat(synin[1], 64)
			tcpin_1, _ := strconv.ParseFloat(tcpin[0], 64)
			tcpin_2, _ := strconv.ParseFloat(tcpin[1], 64)			
			udpin_1, _ := strconv.ParseFloat(udpin[0], 64)
			udpin_2, _ := strconv.ParseFloat(udpin[1], 64)
			icmpin_1, _ := strconv.ParseFloat(icmpin[0], 64)
			icmpin_2, _ := strconv.ParseFloat(icmpin[1], 64)
			otherin_1, _ := strconv.ParseFloat(otherin[0], 64)
			otherin_2, _ := strconv.ParseFloat(otherin[1], 64)
			packetsin_1, _ := strconv.ParseFloat(packetsin[0], 64)	
			packetsin_2, _ := strconv.ParseFloat(packetsin[1], 64)																					
			arrRealtime = append(arrRealtime, fmt.Sprintf("('%s','%f','%f','%f','%f','%f','%f', '%f','%f','%f','%f','%f','%f', '%f', '%f', '%f', '%f')", ip, trafficin_1,trafficin_2,trafficout_1,trafficout_2,synin_1,synin_2,tcpin_1,tcpin_2,udpin_1,udpin_2,icmpin_1,icmpin_2,otherin_1,otherin_2,packetsin_1,packetsin_2))
			arrRealtimeip = append(arrRealtimeip, fmt.Sprintf("'%s'", ip))
			if strings.HasPrefix(ip, "117.27.250") {
				trafficin_1 = trafficin_1 * 10;
				synin_1 = synin_1 * 10;
				udpin_1 = udpin_1 * 10;
				icmpin_1 = icmpin_1 * 10;
				otherin_1 = otherin_1 * 10;
				packetsin_1 = packetsin_1 * 10;
			}
			if strings.HasPrefix(ip, "61.153.107") {
				trafficin_1 = trafficin_1 * 4;
				synin_1 = synin_1 * 4;
				udpin_1 = udpin_1 * 4;
				icmpin_1 = icmpin_1 * 4;
				otherin_1 = otherin_1 * 4;
				packetsin_1 = packetsin_1 * 4;
			}
			
			threshold,_ := strconv.ParseFloat("30", 64)
			if false || trafficin_1 > threshold || synin_1 > threshold || tcpin_1 > threshold || udpin_1 > threshold || icmpin_1 > threshold || otherin_1 > threshold {
				//fmt.Printf("trafficin_1:%f,trafficin_2:%f,trafficout_1:%f,trafficout_2:%f,syn_1:%f,syn_2:%f,udp_1:%f,udp_2:%f,icmp_1:%f,icmp_2:%f,other_1:%f,other_2:%f, =>db\n", trafficin_1, trafficin_2, trafficout_1, trafficout_2, synin_1, synin_2, udpin_1, udpin_2, icmpin_1, icmpin_2, otherin_1, otherin_2)
				arr = append(arr, fmt.Sprintf("('%s','%f','%f','%f','%f', '%f','%f','%f','%f','%f','%f','%f','%f', '%f', '%f', '%f', '%f')", ip, trafficin_1,trafficin_2,trafficout_1,trafficout_2,synin_1,synin_2,tcpin_1,tcpin_2,udpin_1,udpin_2,icmpin_1,icmpin_2,otherin_1,otherin_2,packetsin_1,packetsin_2))
			}			

	}
	dbs, err := sql.Open("mysql", Config.GetDb().Firewall)
	defer dbs.Close()
	if err != nil {
		Utils.LogInfo("can't connect to db:%v", err)
		return
	}
	
	if len(arrRealtime) > 0 {
		sql := fmt.Sprintf("DELETE FROM `attacking_realtime` WHERE ip in (%s)", strings.Join(arrRealtimeip, ","))
		dbs.Exec(sql)
		sql = fmt.Sprintf("insert into `attacking_realtime` (ip,trafficin_1,trafficin_2,trafficout_1,trafficout_2,synin_1,synin_2,tcpin_1,tcpin_2,udpin_1,udpin_2,icmpin_1,icmpin_2,otherin_1,otherin_2,packetsin_1,packetsin_2) values %s", strings.Join(arrRealtime, ","))
		_, err = dbs.Exec(sql)
		if err != nil {
			Utils.LogInfo("can't query :%v, sql:\n%q\n", err, sql)
		}		 
	}	
	if len(arr) > 0 {
		sql := fmt.Sprintf("insert into `attacking` (ip,trafficin_1,trafficin_2,trafficout_1,trafficout_2,synin_1,synin_2,tcpin_1,tcpin_2,udpin_1,udpin_2,icmpin_1,icmpin_2,otherin_1,otherin_2,packetsin_1,packetsin_2) values %s", strings.Join(arr, ","))
		//fmt.Println(sql)
		_, err = dbs.Exec(sql)
		if err != nil {
			Utils.LogInfo("can't query :%v, sql:\n%q\n", err, sql)
		}
	}	
}

//only get attack host data
func GetAttackData(url string)(err error) {
	f := AttackHost{}
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
		if(len(f.OnAttackingHosts) > 0) {
			ParseAttackHost(f.OnAttackingHosts)
		}
	}
	return  nil
}
