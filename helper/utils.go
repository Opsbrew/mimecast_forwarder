package helper

import(
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/nu7hatch/gouuid"
	"github.com/lestrrat-go/strftime"
	"time"
	"strings"
    "encoding/base64"
    "os"
    "crypto/sha1"
    "crypto/hmac"
	"encoding/json"
)


func Get_base_url(email string)(string){
	settings := viper.AllSettings()
	fmt.Println("Getting the base URL")
    fmt.Println(settings)

	base_url:= "no_url"

	EMAIL := make(map[string]interface{})
	EMAIL["emailAddress"] = settings["mm_email_address"].(string)
	var emails []interface{}
	emails = append(emails, EMAIL)

	REQ_BODY := make(map[string]interface{})
	REQ_BODY["data"]= new([]interface{})
	REQ_BODY["data"]= emails
	REQ_ID, _ := uuid.NewV4()
	REQ_DATE := get_hdr_date()

	// body := fmt.Sprint(REQ_BODY)

	client := resty.New()
	resp, err := client.R().
		SetBody(REQ_BODY).
		SetHeader("x-mc-app-id", settings["mm_app_id"].(string)).
		SetHeader("x-mc-req-id", REQ_ID.String()).
		SetHeader("x-mc-date", REQ_DATE).
		SetHeader("Accept", "application/json").
		Post("https://api.mimecast.com/api/login/discover-authentication")
	if(resp.StatusCode() == 200 && err == nil){
		var data map[string]interface{}
		json.Unmarshal(resp.Body(), &data)
		main_data := data["data"].([]interface{})
		if(len(main_data)> 0){
			region :=main_data[0].(map[string]interface{})["region"]
			if(region !=nil){
				temp_url := region.(map[string]interface{})["api"].(string)
				base_url = strings.Split(temp_url, "//")[1]
			}
		}	
	}

	return base_url;
}


func get_hdr_date() (string){
	then := time.Now().UTC()
    f, _ := strftime.New("%a, %d %b %Y %H:%M:%S UTC")
    out := f.FormatString(then)
	return out
}

func Get_mta_siem_logs(base_url string)(bool){
    settings := viper.AllSettings()
    mc_token := ""
    
    checkpoint_path := "./checkpoint"
    if _, err := os.Stat(checkpoint_path); os.IsNotExist(err) {
        os.Mkdir(checkpoint_path,0700)
    }

    file_path := checkpoint_path+"/checkpoint.ops"
    fmt.Println("Reading checkpoint paths ", file_path)
    if _, err := os.Stat(file_path); err == nil {
        fmt.Println("Checkpoint file found")
        out := ReadFile(file_path)
        mc_token = out
    }
    

    REQ_BODY := `{"data":[{"compress":false, "type":"MTA"}]}`
    if(mc_token != ""){
        REQ_BODY = `{"data":[{"compress":false, "type":"MTA", "token":"`+mc_token+`"}]}`
    }
    result,headers,body := post_request(base_url,REQ_BODY)
    if(result=="success"){
        fmt.Println("Updating token in checkpoint")
        wrote := WriteFile(file_path,[]byte(headers["mc-siem-token"]))
        if(wrote){
            fmt.Println("Updating token success")
        }
        if(headers["Content-Type"] == "application/json"){
            fmt.Println("No more logs to search")
        }else if(headers["Content-Type"] == "application/octet-stream"){
            fmt.Println("Log forwarding started")
            file_name := strings.Split(headers["Content-Disposition"],"=\"")
            filen := file_name[1][:len(file_name[1])-1]
            event := strings.Split(filen,"_")
            linex := strings.Split(body,"\r")
            
            for _,log  := range linex {
                if(len(log)>0){
                    
                    furlinex := strings.Split(log,"datetime")
                    for _,furlog  := range furlinex{
                        
                        if(len(furlog)>0){
                            go sendToSyslogServer("datetime"+furlog+"|event="+event[0],settings)
                        }
                    }   
                }
            }
        }

    }else{
        return false
    }

    return true
}
	


func post_request(base_url string,REQ_BODY string) (string,map[string]string,string){
 
    settings := viper.AllSettings()
    REQ_ID, _ := uuid.NewV4()
	REQ_DATE := get_hdr_date()
    
    unsigned_auth_header := fmt.Sprintf("%s:%s:%s:%s", REQ_DATE,REQ_ID.String(),settings["mm_uri"],settings["mm_app_key"]) 

    fmt.Println("generating signature")
    signature := hmacB64(settings["mm_secret_key"].(string),unsigned_auth_header)
    
    headers := make(map[string]string)

    if(signature != "error"){
            
            client := resty.New()
            resp, err := client.R().
            SetBody(REQ_BODY).
            SetHeader("Authorization", "MC "+settings["mm_access_key"].(string)+":"+signature).
            SetHeader("x-mc-req-id", REQ_ID.String()).
            SetHeader("x-mc-date", REQ_DATE).
            SetHeader("x-mc-app-id", settings["mm_app_id"].(string)).
            SetHeader("Content-Type", "application/json").
            SetHeader("Accept","application/octet-stream").
            Post("https://"+base_url+settings["mm_uri"].(string))
            if(resp.StatusCode() == 200 && err == nil){
                fmt.Println("Logs found")
                headers["Content-Disposition"] = resp.Header().Get("Content-Disposition")
                headers["Content-Type"] = resp.Header().Get("Content-Type")
                headers["mc-siem-token"] = resp.Header().Get("mc-siem-token")
                return "success",headers,string(resp.Body())
            }else{
                fmt.Println(err)
            }
    }

    return "failed",headers,""

}


 


func hmacB64(secret_key string, unsigned_auth_header string)(string){
    key,err := base64.StdEncoding.DecodeString(secret_key)
    if(err == nil){
        key_for_sign := []byte(key)
        h := hmac.New(sha1.New, key_for_sign)
        h.Write([]byte(unsigned_auth_header))
        b64encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))
        return b64encoded
    }
    return "error"
    
}