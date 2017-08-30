package execution

import (
	"context"
	"github.com/marvin-automator/marvin/actions/domain"
	appdomain "github.com/marvin-automator/marvin/app/domain"
	idinteractors "github.com/marvin-automator/marvin/identityproviders/interactors"
	"net/http"
	"reflect"
)

var (
	globalActionContext context.Context
	cancelAllActions    context.CancelFunc
	choreContexts       = map[string]choreContext{}
	globalActionLogger  appdomain.Logger
)

type choreContext struct {
	context.Context
	cancel context.CancelFunc
	logger appdomain.Logger
}

func getChoreContext(cid string) choreContext {
	if c, ok := choreContexts[cid]; ok {
		return c
	}

	ctx, cancel := context.WithCancel(globalActionContext)

	cc := choreContext{
		Context: ctx,
		cancel:  cancel,
		logger:  globalActionLogger.WithField("chore", cid),
	}

	choreContexts[cid] = cc

	go func() {
		<-ctx.Done()
		cc.logger.Info("Stopping chore %v", cid)
		delete(choreContexts, cid)
	}()

	return cc
}

type configurationContext struct {
	context 	context.Context
	logger		appdomain.Logger
	exec		*Executor

	account		string
	instance	domain.ActionInstance
	actionMeta	domain.ActionMeta
}

func newConfigurationContext(ex *Executor, account string, log appdomain.Logger, ac domain.BaseAction, inst domain.ActionInstance, ctx context.Context) *configurationContext {
	c := configurationContext{
		context: ctx,
		logger: log,
		exec: ex,
		account: account,
		instance: inst,
		actionMeta: ac.Meta(),
	}
	return &c
}

func (c *configurationContext) InstanceStore() domain.KVStore {
	panic("implement me")
}

func (c *configurationContext) GetCallbackURL(path string) string {
	panic("implement me")
}

func (c *configurationContext) GlobalConfig() interface{} {
	panic("implement me")
}

func (c *configurationContext) Logger() appdomain.Logger {
	return c.logger
}

// HTTPClient returns a http.Client that, if the action requires an identity, and everything is configured correctly,
// will automatically make requests with the correct credentials.
// In any other case, this returns http.DefaultClient.
func (c *configurationContext) HTTPClient() *http.Client {
	ip := idinteractors.IdentityProvider{
		Provider:      c.exec.registry.Provider(c.instance.ActionProvider),
		IdentityStore: c.exec.identityStore,
		Logger:        c.logger,
	}

	gc := c.GlobalConfig()
	gcv := reflect.ValueOf(gc)
	cid := gcv.FieldByName("ClientID").String()
	csec := gcv.FieldByName("ClientSecret").String()

	return ip.GetHTTPClient(cid, csec, c.account, c.instance.Identity)
}

type actionContext struct {
	*configurationContext

	Cancel		context.CancelFunc
	isTestCall	bool
	chore		domain.Chore
}

// newActionContext creates a new ActionContext. ch should be the chore this action is executing in.
func newActionContext(ex *Executor, ch domain.Chore, ac domain.BaseAction, inst domain.ActionInstance) *actionContext {
	cc := getChoreContext(ch.ID)
	ctx, cancel := context.WithCancel(cc)
	log := cc.logger.WithFields(map[string]interface{}{
		"action":          ac.Meta().Key,
		"action_instance": inst.ID,
	})

	configCtx := newConfigurationContext(ex, ch.Owner, log, ac, inst, ctx)
	a := actionContext{
		configurationContext: configCtx,

		Cancel:     cancel,
		isTestCall: false,
		chore:      ch,
	}
	return &a
}

// InvocationStore returns a store valid for the duration of this invocation
func (a *actionContext) InvocationStore() domain.KVStore {
	panic("implement me")
}

// IsTestCall returns whether the current call is a test call
func (a *actionContext) IsTestCall() bool {
	return a.isTestCall
}

// MarkTestCall marks the call that this context will be given as a test call
func (a *actionContext) MarkTestCall() {
	a.isTestCall = true
}

// Output should be called by the action when it wants to send output to the next action
func (a *actionContext) Output(interface{}) {
	panic("implement me")
}

// Done should be called by the action when it will stop sending outputs.
// Triggers should never call this
func (a *actionContext) Done() {
	panic("implement me")
}

// The action should call Error to log en error, in cases where it isn't possible or appropriate to return an error.
func (a *actionContext) Logger() appdomain.Logger {
	return a.logger
}
