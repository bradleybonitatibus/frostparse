# `frostparse`

World of Warcraft 3.3.5a Combat Log Parser library.

## Overview

This combat log parser is wrintten in `go` with the purpose of concurrently processing
a combat log file and providing a means of structuring the 3.3.5a combat log format
into well-structured data.

## Usage

```go
package main

import (   
    "log"
    "github.com/bradleybonitatibus/frostparse"
)

func main() {
    p := frostparse.New(
        frostparse.WithLogFile("C:\\Program Files (x86)\\World of Warcraft 3.3.5a\\Logs\\WoWCombatLog.txt"),
    )
    _data, err := p.Parse()
    if err != nil {
        log.Fatal("failed to parse combatlog: ", err)
    }
    // handle []*frostparse.CombatLogRecord how you please
}
```

If you want basic summary statistics from the combat log, you can use the `Collector` struct:
```go
package main

import (   
    "log"
    "github.com/bradleybonitatibus/frostparse"
)

func main() {
    p := frostparse.New(
        frostparse.WithLogFile("C:\\Program Files (x86)\\World of Warcraft 3.3.5a\\Logs\\WoWCombatLog.txt"),
    )
    data, err := p.Parse()
    if err != nil {
        log.Fatal("failed to parse combatlog: ", err)
    }
    coll := NewCollector()
    coll.Run(data)
    fmt.Println("DamageBySource: ", coll.DamageBySource)
	fmt.Println("HealingBySource: ", coll.HealingBySource)
	fmt.Println("DamageTakenBySource: ", coll.DamageTakenBySource)
	fmt.Println("DamageTakenBySpell: ", coll.DamageTakenBySpell)
}
```

## Data Model

The `CombatLogRecord` struct aggregates a `BaseCombatEvent`, a `Prefix` and a `Suffix`.
The `Prefix` struct is an aggregate to various prefixes, `SpellAndRangePrefix`, `EnchantPrefix`, and `EnvironmentalPrefix`. The member fields of the struct
are pointers because some properties are not populated based on the `BaseCombatEvent.EventType` field of the log record.
