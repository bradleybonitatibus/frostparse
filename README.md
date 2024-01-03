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
        frostparse.WithLogFile("C:\\Program Files (x86)\\World of Warcraft 3.3.5a\\Logs\\WoWCombatLog"),
    )
    _data, err := p.Parse()
    if err != nil {
        log.Fatal("failed to parse combatlog: ", err)
    }
    // handle []*frostparse.CombatLogRecord how you please
}

```