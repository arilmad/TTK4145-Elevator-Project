package nodestates

import (
	"../datatypes"
	"../fsm"
)

// NodeStateMsg ...
// Used for broadcasting the local node state and for receiving remote node states
type NodeStateMsg struct {
	ID    datatypes.NodeID
	State fsm.NodeState
}

// Channels ...
// Used for communication between this module and other modules
type Channels struct {
	LocalNodeStateChan chan fsm.NodeState
	AllNodeStatesChan  chan map[datatypes.NodeID]fsm.NodeState
	NodeLostChan       chan datatypes.NodeID
}

// deepcopyNodeStates ...
// @return: A pointer to a deep copied map of allNodeStates
func deepcopyNodeStates(m map[datatypes.NodeID]fsm.NodeState) map[datatypes.NodeID]fsm.NodeState {
	cpy := make(map[datatypes.NodeID]fsm.NodeState)

	for currID := range m {
		temp := fsm.NodeState{
			Behaviour: m[currID].Behaviour,
			Floor:     m[currID].Floor,
			Dir:       m[currID].Dir,
		}
		cpy[currID] = temp
	}

	return cpy
}

// Handler ...
// The nodestates handler keeps an updated state on all nodes currently in the system
// (that is, nodes that are in peerlist).
// Lost nodes will be deleted from the collection of states, and new nodes will
// be added to the collection of states immediately.
func Handler(
	localID datatypes.NodeID,
	FsmLocalNodeStateChan <-chan fsm.NodeState,
	NetworkAllNodeStatesChan chan<- map[datatypes.NodeID]fsm.NodeState,
	NodeLost <-chan datatypes.NodeID,
	NetworkLocalNodeStateChan chan<- fsm.NodeState,
	RemoteNodeStatesChan <-chan NodeStateMsg) {

	var allNodeStates = make(map[datatypes.NodeID]fsm.NodeState)

	for {
		select {

		// Send received localState from FSM to the network module
		case a := <-FsmLocalNodeStateChan:
			NetworkLocalNodeStateChan <- a

		// Update allNodeStates with the received node state, and
		// update the network module
		case a := <-RemoteNodeStatesChan:
			allNodeStates[a.ID] = a.State
			NetworkAllNodeStatesChan <- deepcopyNodeStates(allNodeStates)

		// Remove lost nodes from allNodeStates
		case a := <-NodeLost:
			delete(allNodeStates, a)
		}

	}
}
