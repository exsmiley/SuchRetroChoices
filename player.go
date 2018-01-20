package main

import (
    // "log"
)

type Player struct {
    name string
    actions []int
}

type Game struct {
    id string // it will be a uuid
    players map[string]Player // maps player name to map
}

type GameMaster struct {
    playerToGame map[string]string // TODO periodically garbage collect disconnected players to free space
    games map[string]Game
    story *Story
}

func newGameMaster(story *Story) GameMaster {
    gm := GameMaster{}
    gm.playerToGame = make(map[string]string)
    gm.games = make(map[string]Game)
    gm.story = story
    return gm
}

// admin game stuff
func (gm *GameMaster) isInGame(playerName string) bool {
    _, ok := gm.playerToGame[playerName]
    return ok
}

func (gm *GameMaster) getGame(playerName string) Game {
    return gm.games[gm.playerToGame[playerName]]
}

func (gm *GameMaster) joinGame(playerName string, host string) bool {
    game, ok := gm.games[gm.playerToGame[host]]

    if !ok || len(game.players) != 1 {
        return false
    }

    game.players[playerName] = Player{playerName, []int{1}}
    return true
}

