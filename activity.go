package activity

import (
	"fmt"
	"strings"

	"github.com/ecletus/admin"
	"github.com/ecletus/core"
	"github.com/ecletus/core/resource"
	"github.com/ecletus/core/utils"
	"github.com/ecletus/media/asset_manager"
	"github.com/ecletus/validations"
	"github.com/moisespsena-go/aorm"
)

// QorActivity default model used to save resource's activities
type QorActivity struct {
	aorm.AuditedModel
	Action        string
	Content       string `sql:"size:5000"`
	Note          string `sql:"size:2000"`
	Type          string
	ResourceParam string
	ResourceID    string
	CreatorName   string
}

// Register register activity feature for an qor resource
func Register(res *admin.Resource) {
	var (
		qorAdmin = res.GetAdmin()
	)

	activityResource := res.AddResource(&admin.SubConfig{
		RawFieldFilter:  map[string]interface{}{"ResourceParam": res.Param},
		ParentFieldName: "ResourceID",
	}, &QorActivity{}, &admin.Config{Invisible: false})
	activityResource.Meta(&admin.Meta{Name: "Action", Type: "hidden", Valuer: func(value interface{}, ctx *core.Context) interface{} {
		act := value.(*QorActivity)
		if act.Action == "" {
			act.Action = "comment on"
		}
		return activityResource.GetAdmin().T(ctx, "activity."+act.Action, act.Action)
	}})
	activityResource.Meta(&admin.Meta{Name: "UpdatedAt", Type: "hidden", Valuer: func(value interface{}, ctx *core.Context) interface{} {
		return utils.FormatTime(value.(*QorActivity).UpdatedAt, "Jan 2 15:04", ctx)
	}})
	activityResource.Meta(&admin.Meta{Name: "URL", Valuer: func(value interface{}, ctx *core.Context) interface{} {
		return strings.Join([]string{res.GetContextURI(ctx, nil), "!" + activityResource.ToParam(), value.(*QorActivity).ID.String(), "edit"}, "/")
	}})

	assetManager := res.GetResourceByID("AssetManager")
	if assetManager == nil {
		assetManager = qorAdmin.AddResource(&asset_manager.AssetManager{}, &admin.Config{Invisible: true})
	}

	activityResource.Meta(&admin.Meta{Name: "Content", Type: "rich_editor", Resource: assetManager})
	activityResource.Meta(&admin.Meta{Name: "Note", Type: "string", Resource: assetManager})
	activityResource.EditAttrs("Action", "Content", "Note")
	activityResource.ShowAttrs("ID", "Action", "Content", "Note", "URL", "UpdatedAt", "CreatorName")
	activityResource.AddValidator(func(record interface{}, metaValues *resource.MetaValues, context *core.Context) error {
		if meta := metaValues.Get("Content"); meta != nil {
			if name := utils.ToString(meta.Value); strings.TrimSpace(name) == "" {
				return validations.Failed(record, "Content", "Content can't be blank")
			}
		}
		return nil
	})

	res.UseTheme("activity")

	qorAdmin.RegisterFuncMap("activity_resource", func() *admin.Resource {
		return qorAdmin.GetResourceByID("QorActivity")
	})

	qorAdmin.RegisterFuncMap("get_activities_count", func(context *admin.Context) int {
		return GetActivitiesCount(context, context.Result)
	})

	router := res.ItemRouter
	ctrl := controller{ActivityResource: activityResource}
	router.Get("/!qor_activities", ctrl.GetActivity)
	router.Post("/!qor_activities", ctrl.CreateActivity)
	router.Post(fmt.Sprintf("/!qor_activities/%v/edit", activityResource.ParamIDPattern()), ctrl.UpdateActivity)
}
