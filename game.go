package main

import (
    // "log"
)

type Player struct {
    cookie string
    score int
    states []string
    next string // used after a wait occurs
}

type Game struct {
    Id string // it will be a uuid
    host string
    players map[string]Player // maps player name to map
    waiting string // cookie of player that is waiting for the other one
}

type GameMaster struct {
    playerToGame map[string]string // TODO periodically garbage collect disconnected players to free space
    cookieToName map[string]string
    games map[string]Game
    story *Story
}

type GameRoom struct {
    HostId string
    HostName string
}


type GameRooms struct {
    Rooms []GameRoom
}

// struct to send back and forth
type GameState struct {
    Waiting bool
    Text string
    Special string
    Actions []string
}


func newGameMaster(story *Story) GameMaster {
    gm := GameMaster{}
    gm.playerToGame = make(map[string]string) // cookie name to string
    gm.games = make(map[string]Game)
    gm.cookieToName = make(map[string]string)
    gm.story = story
    return gm
}

// admin game stuff
func (gm *GameMaster) isInActiveGame(playerCookie string) bool {
    game := gm.getGame(playerCookie)
    return len(game.players) == 2
}

func (gm *GameMaster) getGame(playerCookie string) Game {
    return gm.games[gm.playerToGame[playerCookie]]
}

// TODO maybe map cookie to something else for hosts? too lazy for now
// TODO stop race conditions
func (gm *GameMaster) joinGame(playerCookie string, host string) bool {
    game, ok := gm.games[gm.playerToGame[host]]

    if !ok || len(game.players) != 1 {
        return false
    }

    // remove unhosted game
    delete(gm.games, gm.playerToGame[playerCookie])

    game.players[playerCookie] = Player{playerCookie, 0, []string{"start"}, ""}
    gm.playerToGame[playerCookie] = game.Id

    return true
}

func (gm *GameMaster) hostGame(playerCookie string) string {
    // gm.games[gm.playerToGame]
    game := Game{genId(), playerCookie, map[string]Player{}, ""}
    game.players[playerCookie] = Player{playerCookie, 0, []string{"start"}, ""}

    gm.playerToGame[playerCookie] = game.Id
    gm.games[game.Id] = game
    return game.Id
}

func (gm *GameMaster) getRooms(playerCookie string) GameRooms {
    grs := GameRooms{[]GameRoom{}}

    hosts := map[string]bool{}

    for _, game := range gm.games {
        hostId := game.host
        if hostId == playerCookie || hosts[hostId] {
            continue
        }

        hostName := gm.getName(hostId)
        room := GameRoom{hostId, hostName}
        grs.Rooms = append(grs.Rooms, room)
        hosts[hostId] = true
    }

    return grs
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
    states := player.states
    lastState := states[len(states)-1]
    return gm.story.hasEnded(lastState)
}


func (gm *GameMaster) doAction(playerCookie string, action string) string {
    player := gm.getPlayer(playerCookie)
    states := player.states
    game := gm.games[gm.playerToGame[playerCookie]]
    lastState := states[len(states)-1]

    next, points := gm.story.makeChoice(lastState, action)
    player.score += points

    // handle if other player is waiting
    if game.waiting != "" {
        otherCookie := game.waiting
        otherPlayer := gm.getPlayer(otherCookie)
        otherLastState := otherPlayer.states[len(otherPlayer.states)-1]
        otherNext := gm.story.checkConditions(otherLastState, otherLastState, next)
        otherPlayer.next = otherNext

        nextNext := gm.story.checkConditions(lastState, lastState, otherLastState)
        player.next = nextNext
        game.waiting = ""
    } else if gm.story.needsToWait(next) {
        game.waiting = playerCookie
    } else {
        player.next = ""
    }


    actionString := ""
    if gm.story.triggersHelp(lastState, action) {
        actionString = "help"
        game.waiting = playerCookie
    } else if gm.story.triggersForce(lastState, action) {
        actionString = "force"
    }

    player.states = append(states, next)

    return actionString
}


func (gm *GameMaster) getState(playerCookie string) GameState {
    story := gm.story
    player := gm.getPlayer(playerCookie)
    states := player.states
    game := gm.games[gm.playerToGame[playerCookie]]
    lastState := states[len(states)-1]
    waiting := game.waiting == playerCookie

    actions := story.getActions(lastState, player.next)
    text := story.getText(lastState)

    return GameState{waiting, text, "", actions}
}

