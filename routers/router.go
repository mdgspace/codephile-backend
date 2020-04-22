// @APIVersion 1.0.0
// @Title Codephile Official API
// @Description  Documentation for Codephile API
// @SecurityDefinition token_auth apiKey Authorization header "Enter the token here with bearer keyword Eg: Bearer {token}"
// @Contact mdg@iitr.ac.in
// @TermsOfServiceUrl
// @License
// @LicenseUrl
package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/controllers"
	"github.com/mdg-iitr/Codephile/middleware"
	"net/http"
	"path"
)

func init() {
	beego.InsertFilter("/v1/*", beego.BeforeRouter, middleware.Authenticate)
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/contests",
			beego.NSInclude(
				&controllers.ContestController{},
			),
		),
		beego.NSNamespace("/submission",
			beego.NSInclude(
				&controllers.SubmissionController{},
			),
		),
		beego.NSNamespace("/friends",
			beego.NSInclude(
				&controllers.FriendsController{},
			),
		),
		beego.NSNamespace("/feed",
			beego.NSInclude(
				&controllers.FeedController{},
			),
		),
	)

	ns2 := beego.NewNamespace("/institutes", beego.NSGet("/", func(context *context.Context) {
		http.ServeFile(context.ResponseWriter, context.Request, path.Join(conf.AppRootDir, "conf/institute_list.json"))
	}))
	beego.AddNamespace(ns, ns2)
}
