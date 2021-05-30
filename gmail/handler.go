package gmail

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

type GmailHandler struct {
	sync.Mutex
	Messages map[string]string
	S *gmail.Service
}

func NewGmailHandler() *GmailHandler {
	return &GmailHandler {
		Messages: make(map[string]string),
	}
}

func (handler *GmailHandler) Login() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	handler.S, err = gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
}

func (handler *GmailHandler) AddEmail(id, code string) {
	handler.Lock()
	defer handler.Unlock()

	handler.Messages[id] = code
}

func (handler *GmailHandler) RefreshEmails() (*[]string, error) {

	var newMessages []string

	user := "me"
	r, err := handler.S.Users.Messages.List(user).Q(fmt.Sprintf("from:orders@oe.target.com is:unread after:%s", time.Now().Add(-24 * time.Hour).Format("2006/01/02"))).Do()
	if err != nil {
		return nil, err
	}
	for _, l := range r.Messages {
		msg, err := handler.S.Users.Messages.Get(user, l.Id).Do()
		if err != nil {
			return nil, err
		}

		unEncoded, err := b64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		if err != nil {
			return nil, err
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(unEncoded))
		if err != nil {
			return nil, err
		}

		code := doc.Find("h2").Text()
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}

		if _, ok := handler.Messages[l.Id]; !ok {
			handler.Messages[l.Id] = code
			newMessages = append(newMessages, code)
		}
	}

	return &newMessages, nil
}
