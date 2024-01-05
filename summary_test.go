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

import (
	"fmt"
	"testing"
)


func newTestParser() *Parser {
	return New(
		WithLogFile("./testdata/test.txt"),
	)
}

func TestCollectorRun(t *testing.T) {
	p := newTestParser()
	data, err := p.Parse()
	if err != nil {
		t.Error(err)
	}
	coll := NewCollector()
	stats := coll.Run(data)
	fmt.Println("DamageBySource: ", stats.DamageBySource)
	fmt.Println("HealingBySource: ", stats.HealingBySource)
	fmt.Println("DamageTakenBySource: ", stats.DamageTakenBySource)
	fmt.Println("DamageTakenBySpell: ", stats.DamageTakenBySpell)
}
