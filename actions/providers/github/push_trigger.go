package github

import (
	"github.com/bigblind/marvin/actions/domain"
	"github.com/satori/go.uuid"
	"net/http"
	"fmt"
	"encoding/json"
)

type PushTrigger struct {
	domain.ActionMeta
}

func newPushTrigger() PushTrigger {
	pt := PushTrigger{}
	pt.SetMeta("onpush", "Git Push", "Triggered when someone pushes commits to a repository", true, false, false)
	pt.SetupStepPath = "/setup"
	return pt
}

type PushOutput struct {
	Ref string `json:"ref" description:"The full Git ref that was pushed. Example: 'refs/heads/master'.",json:"ref"`
	Head string `json:"head" description:"The SHA of the most recent commit on ref after the push."`
	Before string `json:"before" description:"The SHA of the most recent commit on ref before the push."`
	Size int `json:"size" description:"The number of commits in the push."`
	DistinctSize int `json:"distinct_size" description:"The number of distinct commits in the push."`
	Commits []Commit `json:"commits" description:"The commits that were pushed."`
	Sender   User    `json:"sender" description:"The GitHub user who pushed the changes"`
	Pusher   CommitAuthor `json:"pusher" description:"The Git author who pushed the changes."`
}

func (p PushTrigger) OutputType(c domain.ActionContext) interface{} {
	return PushOutput{}
}

func (p PushTrigger) Start(c domain.ActionContext) {
	v := uuid.NewV4().String()
	c.InstanceStore().Put("verifier", v)
}

func (p PushTrigger) Callback(req *http.Request, rw domain.ActionResponseWriter, path string, c domain.ActionContext) {
	switch path {
	case "/setup":
		p.handleSetup(req, rw, c)
		return
	case "/pushed":
		p.handlePush(req, rw, c)
		return
	default:
		p.handleUnknownPath(rw, path)
	}
}

func (p PushTrigger) handlePush(req *http.Request, rw domain.ActionResponseWriter, c domain.ActionContext) {
	d := json.NewDecoder(req.Body)
	out := PushOutput{}
	d.Decode(&out)
	c.Output(out)
}

func (p PushTrigger) handleSetup(req *http.Request, rw domain.ActionResponseWriter, c domain.ActionContext) {
	verifier, err := c.InstanceStore().Get("verifier")
	if err != nil {
		rw.Text(500, "Something went wrong:\n" + err.Error())
		return
	}

	url := c.GetCallbackURL("/pushed")

	resp := `
	<h2>Set Up a Webhook on your GitHub Repository</h1>
	<ol>
		<li>Go to your repository's page, and click on "Settings".</li>
		<li>Then, click on "Webhooks".</li>
		<li>Click "Add Webhook", and enter your GitHub password.</li>
		<li>In the <b>Payload URL</b> field, enter: <pre>%v</pre>.</li>
		<li>In the <b>Content Type</b> dropdown, select "Application/JSON".</li>
		<li>In the <b>Secret</b> field, enter: <pre>%v</pre>.</li>
		</li>Leave the setting to just send push events on, and leave the active checkbox checked. Click "Add Webhook".</li>
	</ul>
	`
	resp = fmt.Sprintf(resp, url, verifier)
	rw.HTML(200, resp)
}

func (p PushTrigger) handleUnknownPath(rw domain.ActionResponseWriter, path string) {
	rw.Text(404, "Unknown path: " + path)
}




