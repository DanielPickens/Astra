package naming

import (
	"github.com/daniel-pickens/service-binding-operator/apis"
	"github.com/daniel-pickens/service-binding-operator/pkg/naming"
	"github.com/daniel-pickens/service-binding-operator/pkg/reconcile/pipeline"
)

const StrategyError = "NamingStrategyError"

func Handle(ctx pipeline.Context) {
	for _, item := range ctx.BindingItems() {
		if item.Source != nil {
			template, err := naming.NewTemplate(ctx.NamingTemplate(), templateData(item.Source))
			if err != nil {
				stop(ctx, err)
				return
			}
			item.Name, err = template.GetBindingName(item.Name)
			if err != nil {
				stop(ctx, err)
				return
			}
		}
	}
}

func templateData(service pipeline.Service) map[string]interface{} {
	res := service.Resource()
	return map[string]interface{}{
		"kind": res.GetKind(),
		"name": res.GetName(),
	}
}

func stop(ctx pipeline.Context, err error) {
	ctx.Error(err)
	ctx.StopProcessing()
	ctx.SetCondition(apis.Conditions().NotCollectionReady().Reason(StrategyError).Msg(err.Error()).Build())
}
