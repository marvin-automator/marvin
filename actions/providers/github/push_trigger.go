package github

import (
	"github.com/bigblind/marvin/actions/domain"
	"github.com/satori/go.uuid"
	"net/http"
	"fmt"
	"encoding/json"
	"crypto/hmac"
	"crypto/sha1"
	"bytes"
	"encoding/hex"
)

// PushTrigger is a trigger action that gets triggered when commits get pushed to a Github repository.
type PushTrigger struct {
	domain.ActionMeta
}

func newPushTrigger() PushTrigger {
	pt := PushTrigger{}
	pt.SetMeta("onpush", "Git Push", "Triggered when someone pushes commits to a repository", true, false, false)
	pt.SetupStepPath = "/setup"
	return pt
}

// PushOutput is the type that is output by this action.
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

// OutputType returns a struct of the type that this action will output.
func (p PushTrigger) OutputType(c domain.ActionContext) interface{} {
	return PushOutput{}
}

// Start initializes the trigger
func (p PushTrigger) Start(c domain.ActionContext) {
	// We don't need any setup
}

// Callback gets called when the trigger receives a URL request to a URL it registered.
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
	bodybuf := bytes.NewBuffer([]byte{})
	bodybuf.ReadFrom(req.Body)
	mac, err := hex.DecodeString(req.Header.Get("HTTP_X_HUB_SIGNATURE"))
	if err != nil {
		rw.Text(500, err.Error())
	}
	if !p.checkVerifier(bodybuf.Bytes(), mac, c) {
		rw.Text(400,"Incorrect hmac")
	}
	bodybuf.Reset()
	d := json.NewDecoder(bodybuf)
	out := PushOutput{}
	d.Decode(&out)
	c.Output(out)
}

func (p PushTrigger) handleSetup(req *http.Request, rw domain.ActionResponseWriter, c domain.ActionContext) {
	verifier, err := p.makeVerifier(c)
	if err != nil {
		c.Logger().Error(err)
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

func (p PushTrigger) makeVerifier(c domain.ActionContext) (string, error) {
	verifier := uuid.NewV4().String()
	err := c.InstanceStore().Put("verifier", verifier)
	return verifier, err
}

func (p PushTrigger) checkVerifier(body []byte, mac []byte, c domain.ActionContext) bool {
	ver, err := c.InstanceStore().Get("verifier")
	if err != nil {
		return false
	}
	h := hmac.New(sha1.New, []byte(ver.(string)))
	expected := h.Sum(nil)
	return hmac.Equal(expected, mac)
}





