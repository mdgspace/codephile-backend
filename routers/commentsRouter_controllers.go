package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ContestController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ContestController"],
        beego.ControllerComments{
            Method: "GetContests",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ContestController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ContestController"],
        beego.ControllerComments{
            Method: "GetSpecificContests",
            Router: `/:site`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:FeedController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:FeedController"],
        beego.ControllerComments{
            Method: "ContestsFeed",
            Router: `/contests`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:FeedController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:FeedController"],
        beego.ControllerComments{
            Method: "FriendsFeed",
            Router: `/friend-activity`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:FollowController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:FollowController"],
        beego.ControllerComments{
            Method: "CompareUser",
            Router: `/compare`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:FollowController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:FollowController"],
        beego.ControllerComments{
            Method: "FollowUser",
            Router: `/following`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:FollowController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:FollowController"],
        beego.ControllerComments{
            Method: "GetFollowing",
            Router: `/following`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:SubmissionController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:SubmissionController"],
        beego.ControllerComments{
            Method: "GetSubmission",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:SubmissionController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:SubmissionController"],
        beego.ControllerComments{
            Method: "SaveSubmission",
            Router: `/:site`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:SubmissionController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:SubmissionController"],
        beego.ControllerComments{
            Method: "FilterSubmission",
            Router: `/:site/:uid/filter`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:SubmissionController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:SubmissionController"],
        beego.ControllerComments{
            Method: "FilterSubmission",
            Router: `/:site/filter`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:SubmissionController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:SubmissionController"],
        beego.ControllerComments{
            Method: "GetSubmission",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/all`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "IsAvailable",
            Router: `/available`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "ReturnAllProfiles",
            Router: `/fetch/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "Fetch",
            Router: `/fetch/:site`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "ReturnAllProfiles",
            Router: `/fetch/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/logout`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "ProfilePic",
            Router: `/picture`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "Search",
            Router: `/search`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "CreateUser",
            Router: `/signup`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/mdg-iitr/Codephile/controllers:UserController"],
        beego.ControllerComments{
            Method: "Verify",
            Router: `/verify/:site`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
