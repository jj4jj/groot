package feedback

import (
	"groot/service/feedback/model"
	"groot/sfw/app"
	"groot/sfw/db"
)

//logic state of FeedBack service
type FeedbackService struct {
}


func (logic *FeedbackService)SyncDb() error {
	return sfw_db.AutoMigrateTable(FeedbackDbName, &model.DbFeedback{})
}

func (logic *FeedbackService) Init(svc *app.ServiceCtx) error {
	//init db
	//using
	if e:= InitDbContext(); e!= nil {
		return e
	}
	//post
	svc.BindRpc("create", FeedbackCreate)

	return nil
}

func Initialize() {
	sfw_db.AddDbUseMap(FeedbackDbName, "prod-misc")
	app.RegisterService("feedback", &FeedbackService{}, "main", "need_auth")
}
