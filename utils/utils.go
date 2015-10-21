package utils

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/mssola/user_agent"
	"strings"
)

func IsMobile(userAgent string) bool {
	ua := user_agent.New(userAgent)
	isiPad := strings.Contains(userAgent, "iPad")
	return ua.Mobile() && !isiPad
}
