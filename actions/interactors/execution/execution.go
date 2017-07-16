package execution

import (
	"github.com/bigblind/marvin/actions/domain"
	accountsdomain "github.com/bigblind/marvin/accounts/domain"
	"fmt"
	"context"
	"github.com/bigblind/marvin/app/domain"
)

// SetupExecutionEnvironment should be called by the main function to set up some global variables that will be used
// when executing actions.
func SetupExecutionEnvironment(c context.Context, l domain.Logger) {
	globalActionContext, cancelAllActions = context.WithCancel(c)
	globalActionLogger = l
}

// StartChore starts the triggers for chores
type Executor struct {
	accountStore accountsdomain.AccountStore
	choreStore domain.ChoreStore
	registry domain.ProviderRegistry
}

// All calls the triggers for all saved chores
func (e *Executor) SAll() error {
	return e.accountStore.EachAccount(func(a accountsdomain.Account) error {
		chores, err := e.choreStore.GetAccountChores(a.ID)
		if err != nil {
			return err
		}

		for _, c := range chores {
			err = e.startChore(c)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

// One starts one specific chore, owned by account with ID aid, and chore ID cid
func (e *Executor) One(aid, cid string) error {
	c, err := e.choreStore.GetChore(aid, cid)
	if err != nil {
		return nil
	}

	return e.startChore(c)
}

func (e *Executor) startChore(c domain.Chore) error {
	inst := c.Actions[0]
	act := e.registry.Provider(inst.ActionProvider).Action(inst.Action)
	t, err := e.actionToTrigger(act)
	if err != nil {
		return err
	}

	ctx := newActionContext(c, act, inst)
	ctx.logger.Infof("Starting Chore %v", c.ID)
	go t.Start(ctx)
	return nil
}

func (e *Executor) actionToTrigger(a domain.BaseAction) (domain.Trigger, error) {
	if !a.Meta().IsTrigger {
		return nil, fmt.Errorf("Action %v is first of a chore, but meta.isTrigger is false.", a.Meta().Key)
	}
	if t, ok := a.(domain.Trigger); ok {
		return t, nil
	}

	return nil, fmt.Errorf("Action %v is first of a chore, but doesn't implement the Trigger interface.", a.Meta().Key)
}

// StopAllActions shuts down all the currently running actions.
func (e *Executor) StopAllActions() {
	cancelAllActions()
}

// CancelChore cancels any actions runing for the chore with the given chore ID.
func (e *Executor) CancelChore(cid string) {
	getChoreContext(cid).cancel()
}




