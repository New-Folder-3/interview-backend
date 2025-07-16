package op

import (
	"github.com/pkg/errors"
	"interview/internal/db"
	"interview/internal/model"
	"interview/util"
)

func CreateContent(MessageID string, c Content) error {
	id := util.GenerateToken(16)
	var videoPtr *string
	if c.Video != nil {
		videos := util.StringListToDB(*c.Video)
		videoPtr = &videos
	} else {
		videoPtr = nil
	}
	var audio, image *string
	switch {
	case c.InputAudio != nil:
		audio = &c.InputAudio.Data
	case c.ImageURL != nil:
		image = &c.ImageURL.URL
	}
	content := model.Content{
		ID:        id,
		MessageID: MessageID,
		Text:      c.Text,
		Audio:     audio,
		Image:     image,
		Video:     videoPtr,
	}
	err := db.CreateContent(&content)
	if err != nil {
		return errors.WithStack(err)
	}
	messageDB, err := db.GetMessage(MessageID)
	if err != nil {
		return errors.WithStack(err)
	}
	messageDB.Contents = util.StringListToDB(append(util.DBToStringList(messageDB.Contents), id))
	err = db.UpdateMessage(messageDB)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func DeleteContent(ContentID string, internal bool) error {
	if !internal { // remove in father message
		content, err := db.GetContent(ContentID)
		if err != nil {
			return errors.WithStack(err)
		}
		message, err := db.GetMessage(content.MessageID)
		if err != nil {
			return errors.WithStack(err)
		}
		message.Contents = util.RemoveFromDBList(message.Contents, ContentID)
		if err := db.ReplaceMessage(message); err != nil {
			return errors.WithStack(err)
		}
	}
	if err := db.DeleteContent(ContentID); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func GetContent(contentIDs []string) (*[]Content, error) {
	var ret []Content
	for _, contentID := range contentIDs {
		contentDB, err := db.GetContent(contentID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		var videoPtr *[]string
		if contentDB.Video != nil {
			videos := util.DBToStringList(*contentDB.Video)
			videoPtr = &videos
		} else {
			videoPtr = nil
		}

		content := Content{}
		switch {
		case contentDB.Text != nil:
			content.Type = "text"
			content.Text = contentDB.Text
		case contentDB.Audio != nil:
			content.Type = "input_audio"
			content.InputAudio = &InputAudio{
				Format: "wav",
				Data:   *contentDB.Audio,
			}
		case contentDB.Image != nil:
			content.Type = "input_image"
			content.ImageURL = &ImageURL{
				URL: *contentDB.Image,
			}
		case videoPtr != nil:
			content.Type = "input_video"
			content.Video = videoPtr
		}
		ret = append(ret, content)
	}
	return &ret, nil
}
