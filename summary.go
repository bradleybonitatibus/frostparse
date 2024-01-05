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

import "time"


type Encounter struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// SummaryStats is responsible for listening to the parser.CombatLogRecord stream
// and aggregating the events into well-known raid metrics.
type SummaryStats struct {
	DamageDoneOverTime   map[time.Time]uint64  `json:"damage_done"`
	HealingpDoneOverTime map[time.Time]uint64  `json:"healing_done"`
	DamageTakenOverTime  map[time.Time]uint64  `json:"damage_taken"`
	EncounterOverlays    map[string]Encounter `json:"encounter_overlays"`
	DamageBySource       map[string]uint64     `json:"damage_by_source"`
	HealingBySource      map[string]uint64     `json:"healing_by_source"`
	DamageTakenBySource  map[string]uint64     `json:"damage_taken_by_source"`
	DamageTakenBySpell   map[string]uint64     `json:"damage_taken_by_spell"`
	InterruptsBySource   map[string]uint64     `json:"interrupts_by_source"`
	DispellsBySource     map[string]uint64     `json:"dispells_by_source"`
}

type Collector struct {
	TimeResolution time.Duration
}

type CollectorFunc func(*Collector)

func WithTimeresolution(res time.Duration) CollectorFunc {
	return func(c *Collector) {
		c.TimeResolution=res
	}
}

// NewCollector initializes, allocates and returns a pointer to a Collector struct.
func NewCollector(opts ...CollectorFunc) *Collector {
	t := &Collector{
		TimeResolution: time.Second * 30,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Run consumes the input channel of parser.CombatLogRecord and processes
// each event in the event handler.
func (c *Collector) Run(data []*CombatLogRecord) *SummaryStats {
	s := &SummaryStats{
		DamageDoneOverTime:   map[time.Time]uint64{},
		HealingpDoneOverTime: map[time.Time]uint64{},
		DamageTakenOverTime:  map[time.Time]uint64{},
		DamageBySource:       map[string]uint64{},
		HealingBySource:      map[string]uint64{},
		DamageTakenBySource:  map[string]uint64{},
		DamageTakenBySpell:   map[string]uint64{},
		InterruptsBySource:   map[string]uint64{},
		DispellsBySource:     map[string]uint64{},
		EncounterOverlays:    map[string]Encounter{},
	}
	for i := range data {
		s.handleEvent(*data[i], c.TimeResolution)
	}
	return s
}

// handleEvent is responsible for aggregating the event based on event type
// and source-> target directionality.
func (c *SummaryStats) handleEvent(row CombatLogRecord, resolution time.Duration) {
	if isDamageEvent(row) {
		var amount uint64 = 0
		if row.ExtraAttacksSuffix != nil {
			amount = row.ExtraAttacksSuffix.Amount
		} else if row.DamageSuffix != nil {
			amount = row.DamageSuffix.Amount
		}
		if isBossName(row.TargetName) {
			encounter, ok := c.EncounterOverlays[row.TargetName]
			now := row.Timestamp.Truncate(resolution)
			if !ok {
				encounter = Encounter{
					StartTime: now,
					EndTime:   now,
				}
			} else {
				encounter.EndTime = now
			}
			c.EncounterOverlays[row.TargetName] = encounter
		}
		if (isBossID(row.SourceID) || isNPCID(row.SourceID)) && isPlayerID(row.TargetID) {
			// NPC -> player, accumulate damage taken
			c.DamageTakenBySource[row.SourceName] += amount
			c.DamageTakenOverTime[row.Timestamp.Truncate(resolution)] += amount
			if row.SpellAndRangePrefix != nil {
				c.DamageTakenBySpell[row.SpellAndRangePrefix.SpellName] += amount
			}
			return
		}
		if isPlayerID(row.SourceID) && isNPCID(row.TargetID) || isBossID(row.TargetID) {
			// player -> npc, accumulate damage done
			c.DamageBySource[row.SourceName] += amount
			c.DamageDoneOverTime[row.Timestamp.Truncate(resolution)] += amount
			return
		}
		return
	}
	if isHealingEvent(row) {
		if isPlayerID(row.SourceID) {
			c.HealingBySource[row.SourceName] += row.HealSuffix.Amount
			c.HealingpDoneOverTime[row.Timestamp.Truncate(resolution)] += row.HealSuffix.Amount
		}
		return
	}
	if isOverlayEvent(row) {
		c.DispellsBySource[row.SourceName] += 1
		c.InterruptsBySource[row.SourceName] += 1
		return
	}
}
