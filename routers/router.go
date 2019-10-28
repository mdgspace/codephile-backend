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
	"github.com/mdg-iitr/Codephile/controllers"
	"github.com/mdg-iitr/Codephile/middleware"

	"github.com/astaxie/beego"
)

func init() {
	beego.InsertFilter("/v1/*",beego.BeforeRouter,middleware.Authenticate)
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
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
		beego.NSNamespace("/follow",
	        beego.NSInclude(
				&controllers.FollowController{},
			),
		),
		beego.NSNamespace("/feed",
	        beego.NSInclude(
				&controllers.FeedController{},
			),
	    ),
	)
	beego.AddNamespace(ns)
}
