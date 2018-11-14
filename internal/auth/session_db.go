package auth

import (
	"github.com/kataras/iris/sessions"
	"github.com/marvin-automator/marvin/internal/db"
	"time"
)

type sessDB struct {
	store db.Store
	cache map[string]map[string]interface{}
}

func (s sessDB) getSession(sid string) (map[string]interface{}, error) {
	session, ok := s.cache[sid]
	if ok {
		return session, nil
	}

	session = make(map[string]interface{})
	err := s.store.Get(sid, &session)
	if _, ok := err.(db.KeyNotFoundError); !ok && err != nil {
		return nil, err
	}

	s.cache[sid] = session
	return session, nil
}

func (s sessDB) Acquire(sid string, expires time.Duration) sessions.LifeTime {
	s.cache[sid] = make(map[string]interface{})
	s.store.SetWithExpiration(sid, make(map[string]interface{}), expires)
	return sessions.LifeTime{}
}

func (s sessDB) OnUpdateExpiration(sid string, newExpires time.Duration) error {
	session, err := s.getSession(sid)
	if err != nil {
		return err
	}

	return s.store.SetWithExpiration(sid, session, newExpires)
}

func (s sessDB) Set(sid string, lifetime sessions.LifeTime, key string, value interface{}, immutable bool) {
	session, _ := s.getSession(sid)
	session[key] = value
	s.store.SetWithExpiration(sid, session, lifetime.DurationUntilExpiration())
}

func (s sessDB) Get(sid string, key string) interface{} {
	session, err := s.getSession(sid)
	if err != nil {
		return err
	}

	return session[key]
}

func (s sessDB) Visit(sid string, cb func(key string, value interface{})) {
	session, _ := s.getSession(sid)
	for k, v := range session {
		cb(k, v)
	}
}

func (s sessDB) Len(sid string) int {
	session, _ := s.getSession(sid)
	return len(session)
}

func (s sessDB) Delete(sid string, key string) (deleted bool) {
	session, err := s.getSession(sid)
	if err != nil {
		return false
	}

	delete(session, key)
	return true
}

func (s sessDB) Clear(sid string) {
	s.cache[sid] = make(map[string]interface{})
	exp, err := s.store.GetExpiration(sid)
	if err != nil {
		s.store.Set(sid, s.cache[sid])
	}

	s.store.SetWithExpiration(sid, s.cache[sid], exp)
}

func (s sessDB) Release(sid string) {
	delete(s.cache, sid)
	s.store.Delete(sid)
}

func getSessionDB() sessions.Database {
	return sessDB{db.GetStore("sessions"), make(map[string]map[string]interface{})}
}
