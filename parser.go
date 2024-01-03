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
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type ParserFunc func(*Parser)

// Parser is responsible for loading the combat log file and parsing the data
// into the CombatLogRecord struct.
type Parser struct {
	LogFile       string
	EventListener EventListener
}

// WithLogFile is a ParserFunc that sets the parsers log file.
func WithLogFile(path string) ParserFunc {
	return func(p *Parser) {
		p.LogFile = path
	}
}

// WithEventListener sets the parsers EventListener.
func WithEventListener(listener EventListener) ParserFunc {
	return func(p *Parser) {
		p.EventListener = listener
	}
}

// New initializes and allocates a parser and applies any ParserFunc options
// and returns a pointer to the Parser.
func New(opts ...ParserFunc) *Parser {
	p := &Parser{
		LogFile:       os.Getenv("FROSTPARSE_LOG_FILE"),
		EventListener: NewEventListener(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Parse opens the combat log file and returns a slice of pointers to CombatLogRecords
// and an error if an error occurs during any part of the parsing.
func (p *Parser) Parse() ([]*CombatLogRecord, error) {
	empty := []*CombatLogRecord{}
	f, err := os.Open(p.LogFile)
	defer func() {
		f.Close()
	}()

	if err != nil {
		return empty, err
	}
	rows, err := rowsInFile(f)
	if err != nil {
		return empty, err
	}
	// pre-allocate based on the number of rows identified in the combat log file
	// to limit number of allocations during parsing
	out := make([]*CombatLogRecord, rows)
	// after rowsInFile is called, we need to seek back to beginning of file.
	_, err = f.Seek(0, 0)
	if err != nil {
		return empty, err
	}
	start := time.Now()
	s := bufio.NewScanner(f)
	i := 0
	for s.Scan() {
		v := parseRow(start, s.Text())
		out[i] = &v
		if cb, ok := p.EventListener.Get(v.EventType); ok {
			cb(v)
		}
		i++
	}
	return out, nil
}

// parseRow parses the string data from the combat log and stores it in a
// CombatLogRecord struct and returns it.
func parseRow(startTime time.Time, data string) CombatLogRecord {
	s := strings.Split(data, "  ")
	s[0] = fmt.Sprintf("%d/%s", startTime.Year(), s[0])
	t := mustParseTimestamp(s[0])
	eventParts := strings.Split(s[1], ",")
	eventType := EventType(eventParts[0])
	be := BaseCombatEvent{
		Timestamp:  t,
		EventType:  eventType,
		SourceID:   eventParts[1],
		SourceName: strings.ReplaceAll(eventParts[2], `"`, ""),
		TargetID:   eventParts[4],
		TargetName: strings.ReplaceAll(eventParts[5], `"`, ""),
	}
	prefix := Prefix{}
	suffix := Suffix{}
	switch eventType {
	case UnitDied:
		// can ignore
		break
	case SpellInstakill:
		// can ignore
		break
	case PartyKill:
		// can ignore
		break
	case SwingDamage:
		suffix.DamageSuffix = parseDamageSuffix(eventParts, 7)
	case SpellDamage:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.DamageSuffix = parseDamageSuffix(eventParts, 10)
	case SpellPeriodicDamage:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.DamageSuffix = parseDamageSuffix(eventParts, 10)
	case DamageShield:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.DamageSuffix = parseDamageSuffix(eventParts, 10)
	case DamageSplit:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.DamageSuffix = parseDamageSuffix(eventParts, 10)
	case SpellDrain:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.LeechOrDrainSuffix = parseLeachOrDrainSuffix(eventParts)
	case EnvironmentalDamage:
		prefix.EnvironmentalPrefix = parseEnvironmentalPrefix(eventParts)
		suffix.DamageSuffix = parseDamageSuffix(eventParts, 8)
	case RangeMissed:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.MissSuffix = parseMissSuffix(eventParts)
	case SpellAuraApplied:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.AuraSuffix = parseAuraSuffix(eventParts)
	case SpellHeal:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.HealSuffix = parseHealSuffix(eventParts)
	case SpellAuraRemoved:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.AuraSuffix = parseAuraSuffix(eventParts)
	case SpellCastStart:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
	case SpellCastFailed:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
	case SpellAuraRefresh:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.AuraSuffix = parseAuraSuffix(eventParts)
	case SpellEnergize:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.EnergizeSuffix = parseEnergizeSuffix(eventParts)
	case SwingMissed:
		suffix.MissSuffix = parseMissSuffix(eventParts)
	case SpellAuraAppliedDose:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.AuraSuffix = parseAuraSuffix(eventParts)
	case SpellPeriodicEnergize:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.EnergizeSuffix = parseEnergizeSuffix(eventParts)
	case SpellPeriodicHeal:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.HealSuffix = parseHealSuffix(eventParts)
	case SpellInterrupt:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.InterruptSuffix = parseInterruptSuffix(eventParts)
	case SpellMissed:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.MissSuffix = parseMissSuffix(eventParts)
	case SpellCreate:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
	case RangeDamage:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.DamageSuffix = parseDamageSuffix(eventParts, 10)
	case SpellExtraAttacks:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.ExtraAttacksSuffix = parseExtraAttackSuffix(eventParts)
	case SpellPeriodicMissed:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.MissSuffix = parseMissSuffix(eventParts)
	case SpellAuraRemovedDose:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
	case EnchantApplied:
		prefix.EnchantPrefix = parseEnchantPrefix(eventParts)
	case EnchantRemoved:
		prefix.EnchantPrefix = parseEnchantPrefix(eventParts)
	case SpellResurrect:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
	case SpellDispell:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.DispelOrStolenSuffix = parseDispellOrStolenSuffix(eventParts)
	case DamageShieldMissed:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.MissSuffix = parseMissSuffix(eventParts)
	case SpellPeriodicLeech:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
		suffix.LeechOrDrainSuffix = parseLeachOrDrainSuffix(eventParts)
	case SpellSummon:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
	case SpellCastSuccess:
		prefix.SpellAndRangePrefix = parseSpellPrefix(eventParts)
	default:
		fmt.Println("unknown eventType: ", eventType)
	}

	return CombatLogRecord{
		BaseCombatEvent: be,
		Prefix:          prefix,
		Suffix:          suffix,
	}
}

func parseSpellPrefix(eventParts []string) *SpellPrefix {
	return &SpellPrefix{
		SpellID:     mustParseUint(eventParts[7]),
		SpellName:   removeQuoteString(eventParts[8]),
		SpellSchool: mustParseSpellSchool(eventParts[9]),
	}
}

func parseDamageSuffix(eventParts []string, initialOffset int) *DamageSuffix {
	return &DamageSuffix{
		Amount:      mustParseUint(eventParts[initialOffset]),
		Overkill:    mustParseUint(eventParts[initialOffset+1]),
		SpellSchool: SpellSchool(mustParseInt(eventParts[initialOffset+2])),
		Resisted:    mustParseIntOrNil(eventParts[initialOffset+3]),
		Blocked:     mustParseIntOrNil(eventParts[initialOffset+4]),
		Absorbed:    mustParseIntOrNil(eventParts[initialOffset+5]),
		Critical:    parseNilBool(eventParts[initialOffset+6]),
	}
}

func parseAuraSuffix(eventParts []string) *AuraSuffix {
	return &AuraSuffix{
		AuraType: AuraType(removeQuoteString(eventParts[10])),
	}
}

func parseEnergizeSuffix(eventParts []string) *EnergizeSuffix {
	return &EnergizeSuffix{
		Amount:    mustParseInt(eventParts[10]),
		PowerType: PowerType(mustParseUint(eventParts[11])),
	}
}

func parseMissSuffix(eventParts []string) *MissSuffix {
	return &MissSuffix{
		MissType: eventParts[7],
	}
}

func parseHealSuffix(eventParts []string) *HealSuffix {
	return &HealSuffix{
		Amount:      mustParseUint(eventParts[10]),
		Overhealing: mustParseUint(eventParts[11]),
		Absorbed:    mustParseUint(eventParts[12]),
		Critical:    parseNilBool(eventParts[13]),
	}
}

func parseInterruptSuffix(eventParts []string) *InterruptSuffix {
	return &InterruptSuffix{
		ExtraSpellID:     mustParseUint(eventParts[10]),
		ExtraSpellName:   removeQuoteString(eventParts[11]),
		ExtraSpellSchool: SpellSchool(mustParseUint(eventParts[12])),
	}
}

func parseExtraAttackSuffix(eventParts []string) *ExtraAttacksSuffix {
	return &ExtraAttacksSuffix{
		Amount: mustParseUint(eventParts[10]),
	}
}

func parseEnchantPrefix(eventParts []string) *EnchantPrefix {
	return &EnchantPrefix{
		SpellName: removeQuoteString(eventParts[7]),
		ItemID:    mustParseUint(eventParts[8]),
		ItemName:  removeQuoteString(eventParts[9]),
	}
}

func parseDispellOrStolenSuffix(eventParts []string) *DispelOrStolenSuffix {
	return &DispelOrStolenSuffix{
		ExtraSpellID:     mustParseUint(eventParts[10]),
		ExtraSpellName:   removeQuoteString(eventParts[11]),
		ExtraSpellSchool: mustParseSpellSchool(eventParts[12]),
	}
}

func parseLeachOrDrainSuffix(eventParts []string) *LeechOrDrainSuffix {
	return &LeechOrDrainSuffix{
		Amount:      mustParseUint(eventParts[10]),
		PowerType:   PowerType(mustParseUint(eventParts[11])),
		ExtraAmount: mustParseUint(eventParts[12]),
	}
}

func parseEnvironmentalPrefix(eventParts []string) *EnvironmentalPrefix {
	return &EnvironmentalPrefix{
		EnvironmentalType: EnvironmentalType(eventParts[7]),
	}
}
