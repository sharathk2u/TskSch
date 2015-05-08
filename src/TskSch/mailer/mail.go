package mailer
import (
	"code.google.com/p/goconf/conf"
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
)
type Euser struct {
	User   string
	Pass   string
	Server string
	Port   int
}
var user string
var pass string
var tos string
var server string
var port int
func Mail(subject string, content string) bool {
	
	//Extracting conf
	c, err := conf.ReadConfigFile("../server.conf")
	if err != nil {
		fmt.Println("CAN'T READ CONF FIILE",err)
	}
	user, _ = c.GetString("email", "name")
	pass, _ = c.GetString("email", "password")
	tos, _ = c.GetString("email","to")
	server, _ = c.GetString("email", "server")
	port, _ = c.GetInt("email","port")
	emailUser := &Euser{user, pass, server, port}
	auth := smtp.PlainAuth("",
		emailUser.User,
		emailUser.Pass,
		emailUser.Server,
	)
	//for each recepient
	ml := strings.Split(tos, ",")
	for _, val := range ml {
		from := mail.Address{"Unbxd Task Monitoring Tool", emailUser.User}
		to := mail.Address{strings.Split(val,":")[0], strings.Split(val,":")[1]}
		title := subject
		body := content
		header := make(map[string]string)
		header["From"] = from.String()
		header["To"] = to.String()
		header["Subject"] = func(String string) string {
			// use mail's rfc2047 to encode any string
			addr := mail.Address{String, ""}
			return strings.Trim(addr.String(), " <>")
		}(title)
		header["MIME-Version"] = "1.0"
		header["Content-Type"] = "text/plain; charset=\"utf-8\""
		header["Content-Transfer-Encoding"] = "base64"

		message := ""
		for k, v := range header {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
		err := smtp.SendMail(
			emailUser.Server+":"+strconv.Itoa(emailUser.Port),
			auth,
			from.Address,
			[]string{to.Address},
			[]byte(message),
		)
		if err != nil {
			log.Print("Unable to send mail" + err.Error())
		}
	}
	return true
}
