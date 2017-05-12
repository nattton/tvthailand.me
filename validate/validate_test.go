package validate

import (
	"testing"

	"github.com/code-mobi/tvthailand.me/utils"
	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	assert := assert.New(t)

	var url string
	var isValid bool
	var err error

	url = "https://docs.google.com/file/d/0B-Gdt0P5RlxcQWxXWVBXRkxEQzA/preview"
	isValid, err = URL(url)
	assert.NoError(err, url)
	assert.False(isValid, url)

	url = "https://drive.google.com/file/d/0B-lw9i7B9dS3NTFuTjZZLTlyU0U/view"
	isValid, err = URL(url)
	assert.NoError(err, url)
	assert.True(isValid, url)

	url = "http://goo.gl/OjfAja"
	isValid, err = URL(url)
	assert.NoError(err, url)
	assert.True(isValid, url)

}

func TestRunWebURL(t *testing.T) {
	db, _ := utils.OpenDB()
	RunWebURL(db, 0, 5, 0)
}
