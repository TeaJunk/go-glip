package gliphook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

const (
	icon    = "https://www.citypng.com/public/uploads/small/11596400565bgzbpdobxpy0poq6wg4xsyjywrtqovbjgpygoeiqhq9gwwmd3hj4q9lyhagu3dkekkwgsmfkiuta6euql9giz1v3hydfmlxj7f4b.png"
	hookurl = "https://hooks.glip.com/webhook/testhook"
)

func TestGlipNotification(t *testing.T) {
	message := GlipMessageSimple{
		Activity: "Reminder",
		Title:    "Building glipnotification",
		Body:     "In terraform repo",
		IconUrl:  icon,
	}
	err := GlipSendHook(hookurl, &message)
	if err != nil {
		fmt.Println(err)
		t.Error() // to indicate test failed
	}

}

func TestGlipNotificationCard(t *testing.T) {
	var message GlipMessageCard
	data := []byte(fmt.Sprintf(`{
    "attachments": [
    {
				"type": "Card",
				"author_name" : "aws:sns",
				"author_icon": "%s",
        "fallback": "Something bad happened",
        "color": "#00ff2a",
        "title": "I felt something...",
        "text": "...as if millions of voices suddenly cried out in terror and were suddenly silenced.",
        "fields": [
        {
            "title": "Where",
            "value": "Alderaan",
            "short": true
        },
        {
            "title": "What",
            "value": "Death Star",
            "short": true
        }
        ]
    }
    ]
}`, icon))

	err := json.Unmarshal(data, &message)
	if err != nil {
		fmt.Println(err)
		t.Error() // to indicate test failed
	}
	err = GlipSendHook(hookurl, &message)
	if err != nil {
		fmt.Println(err)
		t.Error() // to indicate test failed
	}

	message = GlipMessageCard{
		IconUrl: icon,
		Attachments: []GlipMessageCardAttachment{
			{
				Type:       "Card",
				AuthorName: "aws:sns",
				AuthorIcon: icon,
				Color:      "#FF0000",
				Fallback:   "Fallback text - An AWS Backup job was stopped. Resource ARN : arn:aws:dynamodb:eu-central-1:accountid:table/foo. BackupJob ID : randomid",
				Text:       "An AWS Backup job was stopped. Resource ARN : arn:aws:dynamodb:eu-central-1:accountid:table/foo. BackupJob ID : randomid",
				Title:      "**Notification from AWS Backup**",
			},
		},
	}
	err = GlipSendHook(hookurl, &message)
	if err != nil {
		fmt.Println(err)
		t.Error() // to indicate test failed
	}

}

func TestMultipleRecords(t *testing.T) {
	data, err := ioutil.ReadFile("testevents/samplemultipleevent.json")
	if err != nil {
		log.Fatalln(err)
	}
	messagetype := "Card"
	authoricon := icon
	color := "#FF0000"
	var snsEvent events.SNSEvent
	err = json.Unmarshal(data, &snsEvent)
	message := GlipMessageCard{
		Attachments: []GlipMessageCardAttachment{},
	}
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		message.Attachments = append(message.Attachments, GlipMessageCardAttachment{
			Type:       messagetype,
			AuthorName: record.EventSource,
			AuthorIcon: authoricon,
			Color:      color,
			Fallback:   snsRecord.Message,
			Text:       snsRecord.Message,
			Title:      fmt.Sprintf("**%s**", snsRecord.Subject),
		})
	}
	err = GlipSendHook(hookurl, &message)
	if err != nil {
		fmt.Println(err)
		t.Error() // to indicate test failed
	}
}

func TestCardStructConvertion(t *testing.T) {
	var testMessage GlipMessageCard
	teststring := fmt.Sprintf(`{
		"icon": "%[1]s",
    "attachments": [
			{
					"type": "Card",
					"author_name" : "aws:sns",
					"author_icon": "%[1]s",
					"fallback": "Fallback - Something bad happened",
					"color": "#FF0000",
					"title": "**Notification from AWS Backup**",
					"text": "An AWS Backup job was stopped. Resource ARN : arn:aws:dynamodb:eu-central-1:accountid:table/foo. BackupJob ID : randomid"
			}
    ]
}`, icon)
	err := json.Unmarshal([]byte(teststring), &testMessage)
	if err != nil {
		log.Fatalln(err)
	}
	message := GlipMessageCard{
		IconUrl: icon,
		Attachments: []GlipMessageCardAttachment{
			{
				Type:       "Card",
				AuthorName: "aws:sns",
				AuthorIcon: icon,
				Color:      "#FF0000",
				Fallback:   "Fallback - Something bad happened",
				Text:       "An AWS Backup job was stopped. Resource ARN : arn:aws:dynamodb:eu-central-1:accountid:table/foo. BackupJob ID : randomid",
				Title:      "**Notification from AWS Backup**",
			},
		},
	}

	// marshalMessage, err := message.glipMessageMarshal()
	if err != nil {
		t.Error()
	}
	assert.Equal(t, testMessage, message, "Should be equal")
	testMessageMarshal, err := json.Marshal(testMessage)
	if err != nil {
		log.Println(err)
		t.Error()
	}
	messageMarshal, err := message.glipMessageMarshal()
	if err != nil {
		log.Println(err)
		t.Error()
	}
	assert.Equal(t, testMessageMarshal, messageMarshal, "Should be equal")
}
