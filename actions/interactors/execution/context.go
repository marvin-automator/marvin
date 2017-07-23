package execution

import (
	"github.com/bigblind/marvin/actions/domain"
	"context"
	appdomain "github.com/bigblind/marvin/app/domain"
	"net/http"
	idinteractors "github.com/bigblind/marvin/identityproviders/interactors"
	"reflect"
)

var (
	globalActionContext context.Context
	cancelAllActions    context.CancelFunc
	choreContexts        = map[string]choreContext{}
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
		cancel: cancel,
		logger: globalActionLogger.WithField("chore", cid),
	}

	go func() {
		<- ctx.Done()
		cc.logger.Info("Stopping chore %v", cid)
	}()

	return cc
}

type actionContext struct {
	exec	   *Executor
	context    context.Context
	Cancel     context.CancelFunc
	isTestCall bool
	logger     appdomain.Logger

	chore 	   domain.Chore
	instance   domain.ActionInstance
	actionMeta domain.ActionMeta
}

// newActionContext creates a new ActionContext. ch should be the chore this action is executing in.
func newActionContext(ex *Executor, ch domain.Chore, ac domain.BaseAction, inst domain.ActionInstance) *actionContext {
	cc := getChoreContext(ch.ID)
	ctx, cancel := context.WithCancel(cc)
	a := actionContext{
		exec: ex,
		context: ctx,
		Cancel: cancel,
		isTestCall: false,
		logger: cc.logger.WithFields(map[string]interface{}{
			"action": ac.Meta().Key,
			"action_instance": inst.ID,
		}),

		chore: ch,
		instance: inst,
		actionMeta: ac.Meta(),

	}
	return &a
}

// InvocationStore returns a store valid for the duration of this invocation
func (a *actionContext) InvocationStore() domain.KVStore {
	panic("implement me")
}

// InstanceStore returns a store with data that is specific to this instance.
func (a *actionContext) InstanceStore() domain.KVStore {
	panic("implement me")
}

// AccountGlobalStore stores data that's accessible to all actions in this provider,
// Data is partitioned by account.
func (a *actionContext) AccountGlobalStore() domain.KVStore {
	panic("implement me")
}

// GetCallbackURL returns a url that actions can use to receive information from other (web) applications
func (a *actionContext) GetCallbackURL(path string) string {
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

// GlobalConfig returns the global configuration object. This object can only be changed on the settingss page, so changes to the object will not be persisted.
func (a *actionContext) GlobalConfig() interface{} {
	panic("implement me")
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

// HTTPClient returns a http.Client that, if the action requires an identity, and everything is configured correctly,
// will automatically make requests with the correct credentials.
// In any other case, this returns http.DefaultClient.
func (a *actionContext) HTTPClient() *http.Client{
	ip := idinteractors.IdentityProvider{
		Provider: a.exec.registry.Provider(a.instance.ActionProvider),
		IdentityStore: a.exec.identityStore,
		Logger: a.logger,
	}

	gc := a.GlobalConfig()
	gcv := reflect.ValueOf(gc)
	cid := gcv.FieldByName("ClientID").String()
	csec := gcv.FieldByName("ClientSecret").String()

	return ip.GetHTTPClient(cid, csec, a.chore.Owner, a.instance.Identity)
}