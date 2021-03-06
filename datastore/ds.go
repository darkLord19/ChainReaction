package datastore

import (
	"strings"
	"sync"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/chainreaction/models"
)

var (
	allGameInstances    map[string]*models.Instance
	activeGameInstances map[string]*models.Instance
	mutex               sync.RWMutex
)

func init() {
	allGameInstances = make(map[string]*models.Instance)
	activeGameInstances = make(map[string]*models.Instance)
}

// GetGameInstance returns game instance from instance id
func GetGameInstance(iid string) (*models.Instance, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	val, ok := activeGameInstances[iid]
	return val, ok
}

// AddGameInstance adds game instance in a data store indexed by instance id
func AddGameInstance(gameInstance *models.Instance) {
	mutex.Lock()
	allGameInstances[gameInstance.RoomName] = gameInstance
	activeGameInstances[gameInstance.RoomName] = gameInstance
	mutex.Unlock()
}

// GetNewUniqueRoomName returns new random unique name for game room
func GetNewUniqueRoomName() string {
	mutex.Lock()
	defer mutex.Unlock()
	name := randomdata.SillyName()
	for activeGameInstances[name] != nil {
		name = randomdata.SillyName()
	}
	return strings.ToLower(name)
}

// Cleanup removes expired and over games from data store
func Cleanup() {
	for {
		time.Sleep(1 * time.Minute)
		mutex.Lock()
		for k, v := range activeGameInstances {
			if v.ExpiresOn.Sub(time.Now().UTC()) < 0 && v.IsOver {
				delete(activeGameInstances, k)
			}
		}
		mutex.Unlock()
	}
}
