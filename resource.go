package activity

import (
	"fmt"
	"strings"

	"github.com/moisespsena-go/aorm"
	"github.com/ecletus/admin"
)

func prepareGetActivitiesDB(context *admin.Context, result interface{}, types ...string) *aorm.DB {
	resourceID := getPrimaryKey(context, result)
	db := context.DB().Order("id asc").Where("resource_id = ? AND resource_type = ?", resourceID, context.Resource.ToParam())

	var inTypes, notInTypes []string
	for _, t := range types {
		if strings.HasPrefix(t, "-") {
			notInTypes = append(notInTypes, strings.TrimPrefix(t, "-"))
		} else {
			inTypes = append(inTypes, t)
		}
	}

	if len(inTypes) > 0 {
		db = db.Where("type IN (?)", inTypes)
	}

	if len(notInTypes) > 0 {
		db = db.Where("type NOT IN (?)", notInTypes)
	}

	return db
}

// GetActivities get activities for selected types
func GetActivities(context *admin.Context, result interface{}, types ...string) ([]QorActivity, error) {
	var activities []QorActivity
	db := prepareGetActivitiesDB(context, result, types...)
	err := db.Find(&activities).Error
	return activities, err
}

// GetActivitiesCount get activities's count for selected types
func GetActivitiesCount(context *admin.Context, result interface{}, types ...string) int {
	var count int
	prepareGetActivitiesDB(context, result, types...).Model(&QorActivity{}).Count(&count)
	return count
}

// CreateActivity creates an activity for this context
func CreateActivity(context *admin.Context, activity *QorActivity, result interface{}) error {
	var activityResource = context.Admin.GetResourceByID("QorActivity")

	// fill in necessary activity fields
	activity.ResourceParam = context.Resource.ToParam()
	activity.ResourceID = getPrimaryKey(context, result)
	if user := context.CurrentUser(); user != nil {
		activity.CreatorName = user.DisplayName()
	}

	return activityResource.Crud(context.Context).Create(activity)
}

func getPrimaryKey(context *admin.Context, record interface{}) string {
	db := context.DB()

	var primaryValues []string
	for _, field := range db.NewScope(record).PrimaryFields() {
		primaryValues = append(primaryValues, fmt.Sprint(field.Field.Interface()))
	}
	return strings.Join(primaryValues, "::")
}
