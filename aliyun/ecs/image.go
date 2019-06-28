package ecs

import (
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	ecssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// Image struct
type Image struct {
	ImageID   string
	ImageName string
	OSName    string
	OSType    string
	SizeGB    int
}

// GetImageByID returns image object with the given image id
func (c *Client) GetImageByID(id string) (Image, error) {
	req := ecssdk.CreateDescribeImagesRequest()
	req.PageSize = requests.NewInteger(50)
	req.ImageId = id

	resp, err := c.DescribeImages(req)
	if err != nil {
		return Image{}, err
	}

	if len(resp.Images.Image) <= 0 {
		return Image{}, errors.New("Image not found")
	}
	respImage := resp.Images.Image[0]

	img := Image{}
	img.ImageID = respImage.ImageId
	img.ImageName = respImage.ImageName
	img.OSName = respImage.OSName
	img.OSType = respImage.OSType
	img.SizeGB = respImage.Size

	return img, nil
}

// GetImageNameByID returns name of the image
func (c *Client) GetImageNameByID(id string) (string, error) {
	img, err := c.GetImageByID(id)
	if err != nil {
		return "", err
	}

	return img.ImageName, nil
}
