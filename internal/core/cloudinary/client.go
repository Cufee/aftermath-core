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

		host:    "api.cloudinary.com",
		version: "v1_1",

		signatureExclude: []string{"file", "cloud_name", "resource_type", "api_key"},
	}
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

	host    string
	version string

	signatureExclude []string
}

func (c *Client) newUploadURL(preset, fileName, fileData string) (string, url.Values) {
	// Make signature
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	// Generate form
	form := url.Values{}
	form.Add("file", fileData)
	if fileName != "" {
		form.Add("public_id", fileName)
	}
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

	return fmt.Sprintf("https://%s/%s/%s", c.host, c.version, c.cloudName) + "/image/upload", form
}

func (c *Client) adminUrl(path string) (*url.URL, error) {
	link, err := url.Parse(fmt.Sprintf("https://%s/%s/%s", c.host, c.version, c.cloudName) + path)
	if err != nil {
		return nil, err
	}

	link.User = url.UserPassword(c.apiKey, c.apiSecret)
	return link, nil
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
func (c *Client) ManualUpload(image string) (string, error) {
	return c.uploadImage("manual-upload", "", image)
}

func (c *Client) uploadImage(preset, name, image string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Send post request
	url, form := c.newUploadURL(preset, name, image)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
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

func (c *Client) GetFolderImages(folder string, limit int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	requestUrl, err := c.adminUrl("/resources/image/upload")
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("prefix", folder)
	if limit > 0 {
		form.Add("max_results", strconv.Itoa(limit))
	}

	req, err := http.NewRequest("GET", requestUrl.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(res.Status)
	}

	var result struct {
		Resources []struct {
			SecureURL string `json:"secure_url"`
		} `json:"resources"`
	}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var images []string
	for _, asset := range result.Resources {
		images = append(images, asset.SecureURL)
	}

	return images, nil
}
