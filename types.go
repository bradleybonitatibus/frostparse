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
	"time"
)

// AuraType represents the enumeration of Buff or Debuf auras.
type AuraType string

// EnvironmentalType represents environmental damage types.
type EnvironmentalType string

// EventType represents the valid event names that appear in the combat log.
type EventType string

// PowerType represents the various power types for characters.
type PowerType int

// SpellSchool is an integer-based representation of schools of magic for a spell.
type SpellSchool int

const (
	DamageShield          EventType = "DAMAGE_SHIELD"
	DamageShieldMissed    EventType = "DAMAGE_SHIELD_MISSED"
	DamageSplit           EventType = "DAMAGE_SPLIT"
	EnchantApplied        EventType = "ENCHANT_APPLIED"
	EnchantRemoved        EventType = "ENCHANT_REMOVED"
	EnvironmentalDamage   EventType = "ENVIRONMENTAL_DAMAGE"
	PartyKill             EventType = "PARTY_KILL"
	RangeDamage           EventType = "RANGE_DAMAGE"
	RangeMissed           EventType = "RANGE_MISSED"
	SpellAuraApplied      EventType = "SPELL_AURA_APPLIED"
	SpellAuraAppliedDose  EventType = "SPELL_AURA_APPLIED_DOSE"
	SpellAuraRefresh      EventType = "SPELL_AURA_REFRESH"
	SpellAuraRemoved      EventType = "SPELL_AURA_REMOVED"
	SpellAuraRemovedDose  EventType = "SPELL_AURA_REMOVED_DOSE"
	SpellCastFailed       EventType = "SPELL_CAST_FAILED"
	SpellCastStart        EventType = "SPELL_CAST_START"
	SpellCastSuccess      EventType = "SPELL_CAST_SUCCESS"
	SpellCreate           EventType = "SPELL_CREATE"
	SpellDamage           EventType = "SPELL_DAMAGE"
	SpellDispell          EventType = "SPELL_DISPEL"
	SpellDrain            EventType = "SPELL_DRAIN"
	SpellEnergize         EventType = "SPELL_ENERGIZE"
	SpellExtraAttacks     EventType = "SPELL_EXTRA_ATTACKS"
	SpellHeal             EventType = "SPELL_HEAL"
	SpellInterrupt        EventType = "SPELL_INTERRUPT"
	SpellInstakill        EventType = "SPELL_INSTAKILL"
	SpellMissed           EventType = "SPELL_MISSED"
	SpellPeriodicDamage   EventType = "SPELL_PERIODIC_DAMAGE"
	SpellPeriodicEnergize EventType = "SPELL_PERIODIC_ENERGIZE"
	SpellPeriodicHeal     EventType = "SPELL_PERIODIC_HEAL"
	SpellPeriodicLeech    EventType = "SPELL_PERIODIC_LEECH"
	SpellPeriodicMissed   EventType = "SPELL_PERIODIC_MISSED"
	SpellResurrect        EventType = "SPELL_RESURRECT"
	SpellSummon           EventType = "SPELL_SUMMON"
	SwingDamage           EventType = "SWING_DAMAGE"
	SwingMissed           EventType = "SWING_MISSED"
	UnitDied              EventType = "UNIT_DIED"
)

// DamageEvents contains the events that dealt damage.
var DamageEvents []EventType = []EventType{
	DamageShield,
	DamageSplit,
	RangeDamage,
	SpellDamage,
	SpellDrain,
	SpellExtraAttacks,
	SpellInstakill,
	SpellPeriodicDamage,
	SpellPeriodicLeech,
	SwingDamage,
}

// HealEvents contain the event types where healing was applied to a target.
var HealEvents []EventType = []EventType{
	SpellHeal,
	SpellPeriodicHeal,
	SpellPeriodicLeech,
}

// OverlayEvents are non-damage/healing events of interest.
var OverlayEvents []EventType = []EventType{
	SpellAuraApplied,
	SpellAuraAppliedDose,
	SpellAuraRemoved,
	SpellAuraRefresh,
	SpellAuraRemovedDose,
	SpellDispell,
	SpellInterrupt,
	UnitDied,
}

// BossNames is the string enumeration containing the ICC Boss names.
var BossNames []string = []string{
	"Lord Marrowgar",
	"Lady Deathwhisper",
	"The Skybreaker",
	"Orgrim's Hammer",
	"Deathbringer Saurfang",
	"Rotface",
	"Festergut",
	"Professor Putricide",
	"Valithria Dreamwalker",
	"Sindragosa",
	"The Lich King",
}

const (
	// BuffAura is when a buff is applied to a target.
	BuffAura AuraType = "BUFF"
	// DebufAura is when a debuf is applied to a target.
	DebufAura AuraType = "DEBUF"
)

const (
	Physical     SpellSchool = 1
	Holy         SpellSchool = 2
	Fire         SpellSchool = 4
	Nature       SpellSchool = 8
	Frost        SpellSchool = 16
	Shadow       SpellSchool = 32
	Arcane       SpellSchool = 64
	Holystrike   SpellSchool = 3
	Flamestrike  SpellSchool = 5
	Radiant      SpellSchool = 6
	Stormstrike  SpellSchool = 9
	Holystorm    SpellSchool = 10
	Volcanic     SpellSchool = 12
	Froststrike  SpellSchool = 17
	Holyfrost    SpellSchool = 18
	Frostfire    SpellSchool = 20
	Froststorm   SpellSchool = 24
	Shadowstrike SpellSchool = 33
	Twilight     SpellSchool = 34
	Shadowflame  SpellSchool = 36
	Plague       SpellSchool = 40
	Shadowfrost  SpellSchool = 48
	Spellstrike  SpellSchool = 65
	Divine       SpellSchool = 66
	Spellfire    SpellSchool = 68
	Astral       SpellSchool = 72
	Spellfrost   SpellSchool = 80
	Spellshadow  SpellSchool = 96
	Elemental    SpellSchool = 28
	Chromatic    SpellSchool = 62
	Cosmic       SpellSchool = 106
	Chaos        SpellSchool = 124
	Magic        SpellSchool = 126
	Fel          SpellSchool = 127
)

const (
	Drowning EnvironmentalType = "DROWNING"
	Falling  EnvironmentalType = "FALLING"
	Fatigue  EnvironmentalType = "FATIGUE"
	FireET   EnvironmentalType = "FIRE"
	Lava     EnvironmentalType = "LAVA"
	Slime    EnvironmentalType = "SLIME"
)

// String implementation of SpellSchool.
func (s SpellSchool) String() string {
	switch s {
	case Physical:
		return "Physical"
	case Holy:
		return "Holy"
	case Fire:
		return "Fire"
	case Nature:
		return "Nature"
	case Frost:
		return "Frost"
	case Shadow:
		return "Shadow"
	case Arcane:
		return "Arcane"
	case Holystrike:
		return "Holystrike"
	case Flamestrike:
		return "Flamestrike"
	case Radiant:
		return "Radiant"
	case Stormstrike:
		return "Stormstrike"
	case Holystorm:
		return "Holystorm"
	case Volcanic:
		return "Volcanic"
	case Froststrike:
		return "Froststrike"
	case Holyfrost:
		return "Holyfrost"
	case Frostfire:
		return "Frostfire"
	case Froststorm:
		return "Froststorm"
	case Shadowstrike:
		return "Shadowstrike"
	case Twilight:
		return "Twilight"
	case Shadowflame:
		return "Shadowflame"
	case Plague:
		return "Plague"
	case Shadowfrost:
		return "Shadowfrost"
	case Spellstrike:
		return "Spellstrike"
	case Divine:
		return "Divine"
	case Spellfire:
		return "Spellfire"
	case Astral:
		return "Astral"
	case Spellfrost:
		return "Spellfrost"
	case Spellshadow:
		return "Spellshadow"
	case Elemental:
		return "Elemental"
	case Chromatic:
		return "Chromatic"
	case Cosmic:
		return "Cosmic"
	case Chaos:
		return "Chaos"
	case Magic:
		return "Magic"
	case Fel:
		return "Fel"
	default:
		return "unknown"
	}
}

// String implementation for PowerType.
func (pt PowerType) String() string {
	switch pt {
	case -2:
		return "Health cost"
	case -1:
		return "None"
	case 0:
		return "Mana"
	case 1:
		return "Rage"
	case 2:
		return "Focus"
	case 3:
		return "Energy"
	case 4:
		return "Combo Points"
	case 5:
		return "Runes"
	case 6:
		return "Runic Power"
	case 7:
		return "Soul Shards"
	default:
		return "N/A"
	}
}

// SwingPrefix is an empty prefix for SWING_ prefixed event types.
type SwingPrefix struct{}

// EnvironmentalPrefix is the prefix for ENVIRONMENTAL_ prefixed events.
// Usually for things like taking falling damage etc.
type EnvironmentalPrefix struct {
	EnvironmentalType EnvironmentalType
}

// SpellAndRangePrefix is the most common prefix containing spell metadata for SPELL_
// and RANGE_ prefixed event types.
type SpellAndRangePrefix struct {
	SpellID     uint64
	SpellName   string
	SpellSchool SpellSchool
}

// DamageSuffix contains damage related metadata.
type DamageSuffix struct {
	Amount      uint64
	Overkill    uint64
	SpellSchool SpellSchool
	Resisted    uint64
	Blocked     uint64
	Absorbed    uint64
	Critical    bool
}

// AuraSuffix contains aura related metadata.
type AuraSuffix struct {
	AuraType AuraType
}

// EnergizeSuffix contains metadata related to a unit getting their power energized
// (either positively, with something like a Mana Potion, or negatively with "Burnout".)
type EnergizeSuffix struct {
	Amount    int64
	PowerType PowerType
}

// MissSuffix is used to identify why some spell / swing missed.
type MissSuffix struct {
	MissType string
}

// HealSuffix contains info about a _HEAL event, usual spell heal, or spell periodic heal.
type HealSuffix struct {
	Amount      uint64
	Overhealing uint64
	Absorbed    uint64
	Critical    bool
}

// InterruptSuffix contains metadata about which spell was interrupted.
type InterruptSuffix struct {
	ExtraSpellID     uint64
	ExtraSpellName   string
	ExtraSpellSchool SpellSchool
}

// DispellOrStolenSuffix provides what spell was debuffed.
type DispelOrStolenSuffix struct {
	ExtraSpellID     uint64
	ExtraSpellName   string
	ExtraSpellSchool SpellSchool
	AuraType         AuraType
}

// LeechOrDrainSuffix provides the amount of power that was leeched or drained from
// a given target.
type LeechOrDrainSuffix struct {
	Amount      uint64
	PowerType   PowerType
	ExtraAmount uint64
}

// EnchantPrefix is used for ENCHANT_ events and providing what item provided
// a spell name.
type EnchantPrefix struct {
	SpellName string
	ItemID    uint64
	ItemName  string
}

// Prefix aggregates all the prefix types. The sub-prefixes will be `nil` if the
// event type does not match the prefix.
type Prefix struct {
	*SpellAndRangePrefix
	*EnchantPrefix
	*EnvironmentalPrefix
}

// ExtraAttacksSuffix provides metadata for how much an extra-attack hit for.
type ExtraAttacksSuffix struct {
	Amount uint64
}

// Suffix aggregates all the suffixes into pointers. Pointers will be `nil` when
// the `BaseCombatEvent.EventType` matches a given suffix.
// For example, if `BaseCombatEvent.EventType == "SPELL_DAMAGE"`, the DamageSuffix
// struct will not be `nil`, but the others will be.
type Suffix struct {
	*DamageSuffix
	*AuraSuffix
	*EnergizeSuffix
	*MissSuffix
	*HealSuffix
	*InterruptSuffix
	*ExtraAttacksSuffix
	*DispelOrStolenSuffix
	*LeechOrDrainSuffix
}

// BaseCombatEvent is the common properties across all combat log lines.
type BaseCombatEvent struct {
	Timestamp  time.Time
	EventType  EventType
	SourceName string
	SourceID   string
	TargetName string
	TargetID   string
}

// CombatLogRecord composes the `BaseCombatEvent`, `Prefix`, and `Suffix` structs
// into a single struct to represent a single line in a combat log file.
type CombatLogRecord struct {
	BaseCombatEvent
	Prefix
	Suffix
}
