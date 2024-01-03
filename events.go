/*
Copyright 2023 Bradley Bonitatibus.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package frostparse

// CombatLogRecordCallback is the callback function signature.
type CombatLogRecordCallback func(CombatLogRecord)

type EventListener interface {
	AddEventListener(event EventType, callback CombatLogRecordCallback)
	Get(EventType) (CombatLogRecordCallback, bool)
}

// listener stores eventtype and callbacks in a map.
type listener struct {
	cbs map[EventType]CombatLogRecordCallback
}

// AddEventListener registers a callback for a given event type.
func (e listener) AddEventListener(event EventType, cb CombatLogRecordCallback) {
	e.cbs[event] = cb
}

// Get returns the callback and an `ok` to indicate if the key existed in
// the event callback map.
func (e listener) Get(event EventType) (CombatLogRecordCallback, bool) {
	cb, ok := e.cbs[event]
	return cb, ok
}

// NewEventListener initializes and allocates an EventLisener implementation
// and returns it.
func NewEventListener() EventListener {
	return &listener{
		cbs: map[EventType]CombatLogRecordCallback{},
	}
}
