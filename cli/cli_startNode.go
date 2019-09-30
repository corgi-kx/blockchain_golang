package cli

import "myCode/public_blockchain/part7-network/network"

func (cli *Cli) startNode() {
	network.StartNode(nodeID)
}
