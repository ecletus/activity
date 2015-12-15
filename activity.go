package activity

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor/admin"
	"github.com/qor/qor/audited"
)

type QorActivity struct {
	gorm.Model
	Subject      string
	Context      string
	Note         string
	Type         string
	ResourceType string
	ResourceID   string
	audited.AuditedModel
}

func RegisterActivityMeta(res *admin.Resource) {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		admin.RegisterViewPath(path.Join(gopath, "src/github.com/qor/activity/views"))
	}

	res.UseTheme("activities")

	res.GetAdmin().RegisterFuncMap("get_activities", func(context *admin.Context, types ...string) []QorActivity {
		activities, _ := GetActivities(context, types...)
		return activities
	})

	router := res.GetAdmin().GetRouter()
	router.Post(fmt.Sprintf("/%v/(.*?)/!activity", res.ToParam()), CreateActivityHandler)
}
