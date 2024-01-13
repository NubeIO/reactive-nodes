package main

import (
	"github.com/NubeIO/reactive"
	"github.com/NubeIO/reactive-nodes/constants"
	"math/rand"
	"time"
)

// exports
var Trigger triggerFloat

type portDataType string

const (
	portTypeAny    portDataType = "any"
	portTypeFloat  portDataType = "float"
	portTypeString portDataType = "string"
	portTypeBool   portDataType = "bool"
)

// triggerFloat generates random values at regular intervals.
type triggerFloat struct {
	*reactive.BaseNode
	stop chan struct{}
}

// NewTriggerNode creates a new triggerFloat with the given ID, name, EventBus, and Flow.
func NewTriggerNode(nodeUUID, name string, bus *reactive.EventBus, settings *reactive.Settings) reactive.Node {
	node := reactive.NewBaseNode(reactive.NodeInfo(trigger, nodeUUID, name, pluginName), bus)
	node.NewOutputPort(constants.Output, constants.Output, "float")
	return &triggerFloat{
		BaseNode: node,
		stop:     make(chan struct{}),
	}

}

func (n *triggerFloat) New(nodeUUID, name string, bus *reactive.EventBus, settings *reactive.Settings) reactive.Node {
	newNode := NewTriggerNode(nodeUUID, name, bus, settings)
	return newNode
}

func (n *triggerFloat) Start() {
	go func() {
		ticker := time.NewTicker(2000 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case <-n.stop:
					return // Stop triggering when the stopTrigger channel is closed
				default:
					ranValue := randFloat()
					out := &reactive.Port{
						ID:        constants.Output,
						Name:      constants.Output,
						Value:     ranValue,
						Direction: "output",
						DataType:  "float",
					}
					n.PublishMessage(out, true)
				}
			}
		}
	}()
}

func (n *triggerFloat) Delete() {
	close(n.stop)
	n.RemoveNodeFromRuntime()
}

func (n *triggerFloat) BuildSchema() {

}

func randFloat() float64 {
	rand.NewSource(time.Now().UnixNano())
	randomFloat := rand.Float64()*9 + 1
	return float64(int(randomFloat))
}
