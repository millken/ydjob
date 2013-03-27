package Job

import (
	"github.com/nutrun/lentil"
	"Utils"
	"Config"
	"net/smtp"
	"strings"
	"time"
	"fmt"
)

type Jobmail struct {
	Loop time.Duration
	Host, User, Pass, From string
}

//https://gist.github.com/andelf/5118732
type loginAuth struct {
	username, password string
}
 
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}
 
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}
 
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
			case "Username:":
			return []byte(a.username), nil
			case "Password:":
			return []byte(a.password), nil
			default:
			return nil, nil
		}
	}
	return nil, nil
}

func (j *Jobmail) Run() {
    conn, err := lentil.Dial(Config.GetBeanstalk().Server)
    if err != nil {
        Utils.LogPanicErr(err)
    }	
	for {

		Utils.LogInfo("Jobmail delay %d Second", Config.GetLoopTime().Mail)
		time.Sleep(time.Second * Config.GetLoopTime().Mail)
		if err != nil {
			conn, err = lentil.Dial(Config.GetBeanstalk().Server)
		}
		if err == nil {
		 	err = conn.Use(Config.GetBeanstalk().MailQueue)
		}
	    if err != nil {
	        Utils.LogPanicErr(err)
	    }else{
	    	for i := 0; i < 10; i++ {
				job,err := conn.PeekReady()
			    if err != nil {
			        Utils.LogPanicErr(err)
			        break
			    }else{
			    	fmt.Printf("Job id %d , body: %s \n", job.Id, job.Body)
					//SendMail("millken@test.cn", "test mail", "<b>"+ string(job.Body)+"</b>",  "html")	    	
			    	conn.Delete(job.Id)
			    }	    		
	    	}			
		}

	}	
}

/*http://www.oschina.net/code/snippet_173630_12032
 *	user : example@example.com login smtp server user
 *	password: xxxxx login smtp server password
 *	host: smtp.example.com:port   smtp.163.com:25
 *	to: example@example.com;example1@163.com;example2@sina.com.cn;...
 *  subject:The subject of mail
 *  body: The content of mail
 *  mailtyoe: mail type html or text
 */


func SendMail(to, subject, body, mailtype string) {
	hp := strings.Split(Config.GetSmtp().Host, ":")
	auth := LoginAuth(Config.GetSmtp().Username, Config.GetSmtp().Password)
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/"+ mailtype + "; charset=UTF-8"
	}else{
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	from := Config.GetSmtp().From

	msg := []byte("To: " + to + "\r\nFrom: " + from + "<"+ Config.GetSmtp().Username +"@"+hp[0]+">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(Config.GetSmtp().Host, auth, Config.GetSmtp().Username, send_to, msg)
    if err != nil {
            Utils.LogPanicErr(err)
    }	

}