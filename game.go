package main

import (
    // "log"
)

type Player struct {
    cookie string
    score int
    actions []int
}

type Game struct {
    id string // it will be a uuid
    players map[string]Player // maps player name to map
}

type GameMaster struct {
    playerToGame map[string]string // TODO periodically garbage collect disconnected players to free space
    cookieToName map[string]string
    games map[string]Game
    story *Story
}

func newGameMaster(story *Story) GameMaster {
    gm := GameMaster{}
    gm.playerToGame = make(map[string]string) // cookie name to string
    gm.games = make(map[string]Game)
    gm.story = story
    return gm
}

// admin game stuff
func (gm *GameMaster) isInGame(playerCookie string) bool {
    _, ok := gm.playerToGame[playerCookie]
    return ok
}

func (gm *GameMaster) getGame(playerCookie string) Game {
    return gm.games[gm.playerToGame[playerCookie]]
}

// TODO maybe map cookie to something else for hosts? too lazy for now
func (gm *GameMaster) joinGame(playerCookie string, host string) bool {
    game, ok := gm.games[gm.playerToGame[host]]

    if !ok || len(game.players) != 1 {
        return false
    }

    game.players[playerCookie] = Player{playerCookie, 0, []int{1}}
    gm.playerToGame[playerCookie] = game.id

    return true
}

func (gm *GameMaster) hostGame(playerCookie string) {
    // gm.games[gm.playerToGame]
    game := Game{genId(), map[string]Player{}}
    game.players[playerCookie] = Player{playerCookie, 0, []int{1}}

    gm.playerToGame[playerCookie] = game.id
    gm.games[game.id] = game
}

// game utils
func (gm *GameMaster) getPlayer(playerCookie string) Player {
    return gm.games[gm.playerToGame[playerCookie]].players[playerCookie]
}

func (gm *GameMaster) getName(playerCookie string) string {
    return gm.cookieToName[playerCookie]
}

func (gm *GameMaster) setName(playerCookie string, name string) {
    gm.cookieToName[playerCookie] = name
}


// game mechanics
// check if the player's game has ended
func (gm *GameMaster) hasEnded(playerCookie string) bool {
    player := gm.getPlayer(playerCookie)
    actions := player.actions
    lastAction := actions[len(actions)-1]

    return gm.story.hasEnded(lastAction)
}


func (gm *GameMaster) doAction(playerCookie string, action int) {
    player := gm.getPlayer(playerCookie)
    actions := player.actions
    lastAction := actions[len(actions)-1]
    path := gm.story.getPath(lastAction)
    nextElement := gm.story.getElement(action)

    // basic util stuff
    player.actions = append(actions, action)
    player.score += nextElement.points

    if intInSlice(action, path.normal) {
        // TODO normal action
    } else if intInSlice(action, path.help) {
        // TODO help action
    } else if intInSlice(action, path.force) {
        // TODO force action
    }

}



