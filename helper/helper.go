package helper

import (
	"fmt"
	"net"
	"time"
	"io/ioutil"
	"github.com/spf13/viper"
)

func HandleError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}


func sendToSyslogServer(msg string, settings map[string]interface{}) {
	fmt.Println("Starting to send to remote syslogserver")
	servAddr := settings["remote_syslog_server"].(string)+":"+settings["port"].(string)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", servAddr)
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
    conn.SetNoDelay(false)
    conn.SetWriteBuffer(10000)
    start := time.Now()
    conn.Write([]byte(msg))
    fmt.Println("took:", time.Since(start))
}

func ReadFile(fileName string) (string){
    data, err := ioutil.ReadFile(fileName)
    if err != nil {
        return "error"
    }
	return string(data)
}

func WriteFile(fileName string,content []byte) (bool){
	err := ioutil.WriteFile(fileName, content, 0644)
    if err != nil {
        return false
    }
    return true
}

func Raw_connect()(bool){
		settings := viper.AllSettings()
	
        timeout := time.Second
        conn, err := net.DialTimeout("tcp", net.JoinHostPort(settings["remote_syslog_server"].(string), settings["port"].(string)), timeout)
        if err != nil {
            return false
        }
        if conn != nil {
            defer conn.Close()
            fmt.Println("Opened connection")
        }
		return true
}