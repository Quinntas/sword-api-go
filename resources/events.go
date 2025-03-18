package resources

import (
	"github.com/quinntas/go-fiber-template/eventEmitter"
	"github.com/quinntas/go-fiber-template/resources/task"
)

func SetupEvents(manager *eventEmitter.ChannelManager) {
	task.SetupEvents(manager)
}
