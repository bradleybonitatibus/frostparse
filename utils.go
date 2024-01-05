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
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"
)

const combatLogTimestampFormat = "2006/1/_2 15:04:05.000"

func rowsInFile(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func mustParseTimestamp(t string) time.Time {
	ts, err := time.Parse(combatLogTimestampFormat, t)
	if err != nil {
		panic(err)
	}
	return ts
}

func mustParseInt(s string) int64 {
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return t
}

func mustParseSpellSchool(s string) SpellSchool {
	if strings.HasPrefix(s, "0x") {
		return SpellSchool(mustParseHexInt(s))
	}
	return SpellSchool(mustParseUint(s))
}

func mustParseHexInt(s string) uint64 {
	t := strings.Replace(s, "0x", "", -1)
	return mustParseUint(t)
}

func mustParseUint(s string) uint64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint64(i)
}

func removeQuoteString(s string) string {
	return strings.ReplaceAll(s, `"`, "")
}

func mustParseIntOrNil(s string) uint64 {
	if strings.Contains(s, "nil") {
		return 0
	}
	return mustParseUint(s)
}

func parseNilBool(s string) bool {
	if strings.Contains(s, "nil") {
		return false
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}

func sliceContains[T comparable](seq []T, v T) bool {
	for i := range seq {
		if seq[i] == v {
			return true
		}
	}
	return false
}

func isDamageEvent(c CombatLogRecord) bool {
	return sliceContains(DamageEvents, c.EventType)
}

func isHealingEvent(c CombatLogRecord) bool {
	return sliceContains(HealEvents, c.EventType)
}

func isOverlayEvent(c CombatLogRecord) bool {
	return sliceContains(OverlayEvents, c.EventType)
}

func isBossName(s string) bool {
	return sliceContains(BossNames, s)
}

func isBossID(v string) bool {
	return strings.HasPrefix(v, "0xF15")
}

func isNPCID(v string) bool {
	return strings.HasPrefix(v, "0xF13")
}

func isPlayerID(v string) bool {
	return strings.HasPrefix(v, "0x07")
}
