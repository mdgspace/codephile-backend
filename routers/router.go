// @APIVersion 1.0.0
// @Title Codephile Official API
// @Description  Documentation for Codephile API
// @Contact mdg@iitr.ac.in
// @TermsOfServiceUrl
// @License
// @LicenseUrl
package routers

import (
	"github.com/mdg-iitr/Codephile/controllers"

	"github.com/astaxie/beego"
)

func init() {
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
	)
	beego.AddNamespace(ns)
}
