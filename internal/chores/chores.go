package chores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/augustoroman/v8"
	"github.com/augustoroman/v8/v8console"
	"github.com/marvin-automator/marvin/actions"
	"github.com/marvin-automator/marvin/internal"
	"github.com/marvin-automator/marvin/internal/db"
	"os"
	"reflect"
)

// choreTrigger is a trigger bound to a store.
type choreTrigger struct {
	RegisteredTrigger
	Input interface{} `json:"-"`
}

// start calls the trigger function, and listens for triggered events. When an event is fired, it sends a callback to the store
// so it can run the JavaScript function that was registered.
func (ct *choreTrigger) start(c *Chore, index int, ctx context.Context) error {
	t, err := actions.Registry.GetAction(ct.Provider, ct.Group, ct.Action)
	if err != nil {
		return err
	}

	out, err := t.Run(ct.Input, ctx)
	if err != nil {
		return err
	}

	outInterfaces, err := receiveValues(out, ctx)

	go func() {
		for {
			select {
			case v := <-outInterfaces:
				fmt.Printf("Received value %v\n", v)
				c.triggerCallback(index, v, ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// receiveValues is a helper function that takes a receiving channel of unknown type, and outputs all the values to a new channel
func receiveValues(in interface{}, ctx context.Context) (<-chan interface{}, error) {
	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Chan {
		return nil, errors.New("Output from trigger is non-channel")
	}

	out := make(chan interface{}, 20)
	go func() {
		for {
			i, outv, _ := reflect.Select([]reflect.SelectCase{
				{Dir: reflect.SelectRecv, Chan: v},
				{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())},
			})
			fmt.Printf("Received event %v, %v", i, outv)
			if i == 1 {
				return
			}
			out <- outv.Interface()
		}
	}()

	return out, nil
}

// UnmarshalJSON lets choreTrigger implement the JSONUnmarshaler interface.
// This is necessary because we need to convert the trigger inputs we got to the correct type for the trigger function.
func (ct *choreTrigger) UnmarshalJSON(data []byte) error {
	var rt RegisteredTrigger
	err := json.Unmarshal(data, &rt)
	if err != nil {
		return err
	}

	ct.RegisteredTrigger = rt

	a, err := actions.Registry.GetAction(rt.Provider, rt.Group, rt.Action)
	if err != nil {
		return err
	}

	inputHolder := struct{ Input interface{} }{}
	inputHolder.Input = reflect.New(a.Info().InputType).Interface()
	err = json.Unmarshal(data, &inputHolder)
	if err != nil {
		return err
	}

	ct.Input = reflect.ValueOf(inputHolder.Input).Elem().Interface()

	return nil
}

// holds configuration data for the chore.
type choreConfig struct {
	// configurable parameters that are used by the script specify behavior.
	Inputs   map[string]string `json:"inputs"`

	// The registered triggers, associated with their parameters.
	Triggers []choreTrigger    `json:"triggers"`
}

// A chore is a workflow for Marvin to execute. It consists of a number of triggers with callbacks that specify what
// should happen when a trigger fires.
type Chore struct {
	Name     string        `json:"name"`
	Id       string        `json:"id"`
	Active   bool          `json:"active"`
	Template ChoreTemplate `json:"template"`
	Config   choreConfig   `json:"choreSettings"`
	Snapshot []byte        `json:"-"`
}

// FromTemplate creates a new Chore based on a template.
func FromTemplate(ct *ChoreTemplate, name string, inputs map[string]string) (*Chore, error) {
	conf, err := ct.GenerateChoreConfig(inputs)
	if err != nil {
		return nil, err
	}

	id, err := internal.NewId()

	return &Chore{
		Name:     name,
		Config:   *conf,
		Template: *ct,
		Snapshot: ct.GetChoreSnapshot(inputs),
		Id:       id,
	}, nil
}

// Start activates the chore, and starts all triggers.
func (c *Chore) Start(ctx context.Context) {
	c.Active = true

	for i, ct := range c.Config.Triggers {
		//TODO handle any errors returned by the trigger
		ct.start(c, i, ctx)
	}
}

var choreContexts = make(map[string]*v8.Context)

// TriggerCallback is called by a trigger when it has fired an event.
// This will call the JavaScript callback that was registered for this trigger.
func (c *Chore) triggerCallback(index int, value interface{}, ctx context.Context) {
	jsCtx, ok := choreContexts[c.Id]
	if !ok {
		fmt.Printf("Creating context for chore %v\n", c.Name)
		jsCtx = c.createContext(ctx)
		fmt.Printf("Created context for chore %v\n", c.Name)
		choreContexts[c.Id] = jsCtx
	}

	eventValue, err := jsCtx.Create(value)
	if err != nil {
		fmt.Printf("Received malformed value from a trigger in chore%v", c.Name)
		return
	}

	b, _ := eventValue.MarshalJSON()
	fmt.Printf("Set event input value %v\n", string(b))
	err = jsCtx.Global().Set("__triggeredEvent", eventValue)
	if err != nil {
		panic(err)
	}
	fmt.Println("Running JS for trigger...")

	go func() {
		code := c.Template.combineScriptWithBoilerplate(c.Config.Inputs) + fmt.Sprintf("marvin.isSetup=false;marvin._triggers[%v].callback(__triggeredEvent)", index)
		fmt.Println("CODE:\n\n", code, "\n\n_______")
		res, err := jsCtx.Eval(code, "name.js")
		fmt.Println("result:", res, " error:", err)
	}()

}

// createContext creates a JavaScript context for this chore.
func (c *Chore) createContext(ctx context.Context) *v8.Context {
	is := v8.NewIsolate()
	jsCtx := is.NewContext()

	runAction := jsCtx.Bind("_runAction", func(args v8.CallbackArgs) (*v8.Value, error) {
		provider := args.Arg(0).String()
		group := args.Arg(1).String()
		action := args.Arg(2).String()
		undefined, _ := jsCtx.Create("undefined")

		a, err := actions.Registry.GetAction(provider, group, action)
		if err != nil {
			return undefined, err
		}

		inBytes, err := args.Arg(3).MarshalJSON()
		if err != nil {
			return undefined, err
		}

		in := reflect.New(a.Info().InputType).Interface()
		err = json.Unmarshal(inBytes, in)
		if err != nil {
			return undefined, err
		}

		out, err := a.Run(in, ctx)
		if err != nil {
			return undefined, err
		}

		return jsCtx.Create(out)
	})
	fmt.Println("Assigning _runAction")
	if err := jsCtx.Global().Set("_runAction", runAction); err != nil {
		panic(err)
	}

	fmt.Println("Injecting console")
	cons := v8console.Config{Stdout: os.Stdout, Stderr: os.Stderr}
	cons.Inject(jsCtx)

	return jsCtx
}

const choreStoreName = "chores"

var choreCache = map[string]*Chore{}
var choresLoaded = false

// Save saves a chore to the database.
func (c *Chore) Save() error {
	s := db.GetStore(choreStoreName)
	choreCache[c.Id] = c
	return s.Set(c.Id, c)
}

// Delete removes a chore from the database.
func (c *Chore) Delete() error {
	s := db.GetStore(choreStoreName)
	err := s.Delete(c.Id)
	if err != nil {
		return err
	}

	delete(choreCache, c.Id)
	return nil
}

// GetChore gets a single chore from the database.
func GetChore(id string) (*Chore, error) {
	if c, ok := choreCache[id]; ok {
		return c, nil
	}

	s := db.GetStore(choreStoreName)
	c := new(Chore)
	err := s.Get(id, s)
	if err != nil {
		return nil, err
	}

	choreCache[id] = c
	return c, nil
}

// GetChores gets all chores from the database.
func GetChores() ([]*Chore, error) {
	res := make([]*Chore, len(choreCache))
	if choresLoaded {
		for _, c := range choreCache {
			res = append(res, c)
		}
	} else {
		s := db.GetStore(choreStoreName)
		c := new(Chore)
		err := s.EachKeyWithPrefix("", c, func(key string) error {
			ccopy := *c
			choreCache[key] = &ccopy

			res = append(res, choreCache[key])
			return nil
		})

		if err != nil {
			return nil, err
		}

		choresLoaded = true
	}

	return res, nil
}
