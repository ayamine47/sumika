package cloud

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ayamine47/sumika/lib/config"
	"github.com/ayamine47/sumika/lib/embed"
	"github.com/ayamine47/sumika/lib/utils"
)

type ResponseData struct {
	XMLName xml.Name `xml:"ocs"`
	Meta    struct {
		Status     string `xml:"status"`
		StatusCode int    `xml:"statuscode"`
		Message    string `xml:"message"`
	} `xml:"meta"`
	Data struct {
		URL string `xml:"url"`
	} `xml:"data"`
}

func UploadFileToNextCloud(fileName string, msgInfo *embed.MsgInfo) {
	filePath := config.CurrentConfig.NextCloud.Path + "/" + fileName + "_720.mp4"

	file, err := os.Open("./" + fileName + "_720.mp4")
	if err != nil {
		log.Print("Failed to opening file: ", err)
		return
	}

	defer file.Close()

	err = client.WriteStream(filePath, file, 0644)
	if err != nil {
		log.Print("Failed to upload file to NextCloud: ", err)
		return
	}

	nextCloudURL, err := url.Parse(config.CurrentConfig.NextCloud.Url)
	if err != nil {
		log.Print("Failed to parse NextCloud URL: ", err)
		return
	}

	expireDate := time.Now().Add(1 * time.Hour).Format(time.DateTime)
	apiURL := fmt.Sprintf("https://%s/ocs/v2.php/apps/files_sharing/api/v1/shares", nextCloudURL.Host)

	form := url.Values{}
	form.Add("path", filePath)
	form.Add("shareType", "3")
	form.Add("permissions", "1")
	form.Add("expireDate", expireDate)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Print("Failed to create request: ", err)
	}

	req.SetBasicAuth(config.CurrentConfig.NextCloud.Username, config.CurrentConfig.NextCloud.Password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("OCS-APIRequest", "true")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Print("Failed to send API request: ", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Failed to read response: ", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("API Error (Status Code: %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var respData ResponseData
	err = xml.Unmarshal(bodyBytes, &respData)
	if err != nil {
		log.Print("Failed to parse XML: ", err)
	}

	if respData.Meta.Message == "OK" {
		embed.SendSuccessEmbed(msgInfo, "Video URL: "+respData.Data.URL+"\nDue Date: "+expireDate, msgInfo.OrgMsg.Reference())

		err = msgInfo.Session.MessageReactionRemove(msgInfo.OrgMsg.ChannelID, msgInfo.OrgMsg.ID, "🤔", msgInfo.Session.State.User.ID)
		if err != nil {
			log.Print("Failed to add reaction: ", err)
		}

		err = msgInfo.Session.MessageReactionAdd(msgInfo.OrgMsg.ChannelID, msgInfo.OrgMsg.ID, "✅")
		if err != nil {
			log.Print("Failed to add reaction: ", err)
		}
	} else {
		log.Print("Failed to create public link: ", err)
		return
	}

	utils.CleanUpLocalVideo(fileName)
}
