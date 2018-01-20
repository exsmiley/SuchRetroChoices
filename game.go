package main

import (
    "log"
)

type Player struct {
    cookie string
    score int
    actions []string
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
    cookie string
    score int
    waiting bool
    actions []int
}


func newGameMaster(story *Story) GameMaster {
    gm := GameMaster{}
    gm.playerToGame = make(map[string]string) // cookie name to string
    gm.games = make(map[string]Game)
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

    game.players[playerCookie] = Player{playerCookie, 0, []string{"start"}}
    gm.playerToGame[playerCookie] = game.Id

    return true
}

func (gm *GameMaster) hostGame(playerCookie string) string {
    // gm.games[gm.playerToGame]
    game := Game{genId(), playerCookie, map[string]Player{}, ""}
    game.players[playerCookie] = Player{playerCookie, 0, []string{"start"}}

    gm.playerToGame[playerCookie] = game.Id
    gm.games[game.Id] = game
    return game.Id
}

func (gm *GameMaster) getRooms(playerCookie string) GameRooms {
    grs := GameRooms{[]GameRoom{}}

    for _, game := range gm.games {
        hostId := game.host
        if hostId == playerCookie {
            continue
        }

        hostName := gm.getName(hostId)
        room := GameRoom{hostId, hostName}
        grs.Rooms = append(grs.Rooms, room)
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
    actions := player.actions
    lastAction := actions[len(actions)-1]

    return gm.story.hasEnded(lastAction)
}


func (gm *GameMaster) doAction(playerCookie string, action string) string {
    player := gm.getPlayer(playerCookie)
    actions := player.actions
    // lastAction := actions[len(actions)-1]
    // path := gm.story.getPath(lastAction)
    // nextElement := gm.story.getElement(action)

    // basic util stuff
    player.actions = append(actions, action)
    // player.score += nextElement.points

    // if intInSlice(action, path.normal) {
    //     // TODO normal action
    // } else if intInSlice(action, path.help) {
    //     // TODO help action
    // } else if intInSlice(action, path.force) {
    //     // TODO force action
    // }
    return ""
}


func (gm *GameMaster) getState(playerCookie string) GameState {
    // TODO and maybe change return type
    return GameState{}
}

