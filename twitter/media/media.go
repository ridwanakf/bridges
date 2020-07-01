package media

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/sling"
)

/*
	This code is part of this PR: https://github.com/dghubble/go-twitter/pull/152/
	which I've been given permission to edit and refactor it myself.

	The main reason is because https://github.com/dghubble/go-twitter/ seems to be abandoned
	and there's no Twitter Media API Implementation.
 */

const twitterUploadAPI = "https://upload.twitter.com/1.1/"

//Service ... provides methods for accessing Twitter status API endpoints.
type Service struct {
	sling *sling.Sling
}

// NewStatusService returns a new StatusService.
func NewService(httpClient *http.Client) *Service {
	upload := sling.New().Client(httpClient).Base(twitterUploadAPI)
	return &Service{
		sling: upload.Path("media/"),
	}
}

//UploadParams ... are the parameters for StatusService.Update
type UploadParams struct {
	File     []byte
	MimeType string
}

//Media ... response of uploaded file
type Media struct {
	MediaID          int64  `json:"media_id"`
	MediaIDString    string `json:"media_id_string"`
	ExpiresAfterSecs uint64 `json:"expires_after_secs"`
}

type mediaUploadCommand struct {
	Command      string `url:"command,omitempty"`
	MediaID      string `url:"media_id,omitempty"`
	MediaType    string `url:"media_type,omitempty"`
	MediaData    string `url:"media_data,omitempty"`
	SegmentIndex string `url:"segment_index,omitempty"`
	TotalBytes   string `url:"total_bytes,omitempty"`
}

func (m UploadParams) getTotalBytes() int {
	if m.File != nil {
		return len(m.File)
	}

	return 0
}

// Upload media file
// Requires a user auth context.
// https://dev.twitter.com/rest/reference/post/media/upload
func (s *Service) Upload(params *UploadParams) (*Media, *http.Response, error) {
	var resp *http.Response
	var err error
	var twitterMediaID *Media

	twitterMediaID, resp, err = s.mediaInit(params)
	if err != nil {
		return nil, resp, err
	}

	resp, err = s.mediaAppend(twitterMediaID, params)
	if err != nil {
		return nil, resp, err
	}

	resp, err = s.mediaFinalize(twitterMediaID.MediaID)
	if err != nil {
		return nil, resp, err
	}

	return twitterMediaID, resp, nil
}

func (s *Service) mediaInit(p *UploadParams) (*Media, *http.Response, error) {
	paramsBody := mediaUploadCommand{
		Command:    "INIT",
		MediaType:  p.MimeType,
		TotalBytes: fmt.Sprintf("%d", p.getTotalBytes()),
	}

	twitterMediaID := new(Media)
	apiError := new(twitter.APIError)
	resp, err := s.sling.New().Post("upload.json").Add("Content-Type", "application/x-www-form-urlencoded").BodyForm(paramsBody).Receive(twitterMediaID, apiError)
	return twitterMediaID, resp, relevantError(err, *apiError)
}

func (s *Service) mediaAppend(twitterMediaID *Media, params *UploadParams) (*http.Response, error) {
	media := params.File
	mediaID := twitterMediaID.MediaIDString
	mediaBase64 := b64.StdEncoding.EncodeToString(media)

	step := 500 * 1024
	for i := 0; i*step < len(mediaBase64); i++ {
		rangeBeginning := i * step
		rangeEnd := (i + 1) * step
		if rangeEnd > len(mediaBase64) {
			rangeEnd = len(mediaBase64)
		}
		_ = rangeBeginning
		params := mediaUploadCommand{
			Command:      "APPEND",
			MediaID:      mediaID,
			MediaData:    mediaBase64[rangeBeginning:rangeEnd],
			SegmentIndex: fmt.Sprint(i),
		}

		apiError := new(twitter.APIError)
		resp, err := s.sling.New().Post("upload.json").Add("Content-Type", "application/x-www-form-urlencoded").BodyForm(params).Receive(nil, apiError)
		if err != nil {
			return resp, relevantError(err, *apiError)
		}
	}

	return nil, nil
}

func (s *Service) mediaFinalize(mediaID int64) (*http.Response, error) {
	params := mediaUploadCommand{
		Command: "FINALIZE",
		MediaID: fmt.Sprint(mediaID),
	}

	apiError := new(twitter.APIError)
	resp, err := s.sling.New().Post("upload.json").Add("Content-Type", "application/x-www-form-urlencoded").BodyForm(params).Receive(nil, apiError)
	if err != nil {
		return resp, relevantError(err, *apiError)
	}

	return resp, nil
}

func relevantError(httpError error, apiError twitter.APIError) error {
	if httpError != nil {
		return httpError
	}
	if apiError.Empty() {
		return nil
	}
	return apiError
}
