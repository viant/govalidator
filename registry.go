package govalidator

import (
	"strings"
	"sync"
)

type registry struct {
	fn    map[string]NewIsValid
	alias map[string][]string
	sync.RWMutex
}

//Register register tag
func (r *registry) Register(tag string, fn NewIsValid) {
	r.RWMutex.Lock()
	r.fn[tag] = fn
	r.RWMutex.Unlock()
}

//Lookup returns tag NewIsValid function
func (r *registry) Lookup(tag string) NewIsValid {
	r.RWMutex.RLock()
	ret := r.fn[tag]
	r.RWMutex.RUnlock()
	return ret
}

//Alias returns tag  aliases
func (r *registry) Alias(tag string) []string {
	r.RWMutex.RLock()
	ret := r.alias[tag]
	r.RWMutex.RUnlock()
	return ret
}

//RegisterAlias register tag alias
func (r *registry) RegisterAlias(tag string, tags ...string) {
	r.RWMutex.Lock()
	r.alias[tag] = tags
	r.RWMutex.Unlock()
}

var _register = &registry{fn: map[string]NewIsValid{}, alias: map[string][]string{}}

//Register register tag
func Register(check string, fn NewIsValid) {
	_register.Register(strings.ToLower(check), fn)
}

//RegisterAlias register tag alias
func RegisterAlias(check string, checks ...string) {
	_register.RegisterAlias(strings.ToLower(check), checks...)
}

//Lookup returns tag NewIsValid
func Lookup(check string) NewIsValid {
	return _register.Lookup(strings.ToLower(check))
}

//Alias returns tag aliases
func Alias(check string) []string {
	return _register.Alias(strings.ToLower(check))
}

//LookupAll returns tag NewIsValids
func LookupAll(check string) NewIsValid {
	newIsValid := _register.Lookup(strings.ToLower(check))
	if newIsValid != nil {
		return newIsValid
	}
	var result []NewIsValid
	for _, alias := range Alias(check) {
		result = append(result, LookupAll(alias))
	}
	return atListOneValid(result...)
}
