package cloudinary

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cufee/aftermath-core/internal/core/utils"
)

var DefaultClient *Client

func init() {
	DefaultClient = &Client{
		cloudName: utils.MustGetEnv("CLOUDINARY_API_CLOUD_NAME"),
		apiSecret: utils.MustGetEnv("CLOUDINARY_API_SECRET"),
		apiKey:    utils.MustGetEnv("CLOUDINARY_API_KEY"),
		preset:    utils.MustGetEnv("CLOUDINARY_API_PRESET_NAME"),

		signatureExclude: []string{"file", "cloud_name", "resource_type", "api_key"},
	}
	DefaultClient.baseURL = fmt.Sprintf("https://api.cloudinary.com/v1_1/%s", DefaultClient.cloudName)
}

type uploadResponse struct {
	PublicID  string    `json:"public_id"`
	Format    string    `json:"format"`
	CreatedAt time.Time `json:"created_at"`
	URL       string    `json:"url"`
}

type Client struct {
	cloudName string
	apiSecret string
	apiKey    string

	baseURL          string
	preset           string
	signatureExclude []string
}

func (c *Client) newUploadURL(preset, fileName, fileData string) (string, url.Values) {
	// Make signature
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	// Generate form
	form := url.Values{}
	form.Add("file", fileData)
	form.Add("public_id", fileName)
	form.Add("upload_preset", preset)

	form.Add("api_key", c.apiKey)
	form.Add("timestamp", timestamp)

	signatureValues := url.Values{}
	for key, value := range form {
		if len(value) != 1 {
			continue
		}
		if slices.Contains(c.signatureExclude, key) {
			continue
		}
		signatureValues.Add(key, value[0])
	}

	// Encode and add signature
	h := sha1.New()
	h.Write([]byte(signatureValues.Encode() + c.apiSecret))
	signature := hex.EncodeToString(h.Sum(nil))
	form.Add("signature", signature)

	return c.baseURL + "/image/upload", form
}

/*
Uploads a base64 encoded image to Cloudinary with moderation using AWS Rekognition
*/
func (c *Client) UploadWithModeration(name, image string) (string, error) {
	return c.uploadImage("user-upload", name, image)
}

/*
Uploads a base64 encoded image to Cloudinary without moderation
*/
func (c *Client) ManualUpload(name, image string) (string, error) {
	return c.uploadImage("manual-upload", name, image)
}

func (c *Client) uploadImage(preset, name, image string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Send post request
	url, form := c.newUploadURL(preset, name, image)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf(res.Status)
	}

	// Return image URL and err
	var response uploadResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}
	if response.URL == "" {
		return "", fmt.Errorf("failed to upload image")
	}

	return response.URL, nil
}
