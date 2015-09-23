/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package protos

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/op/go-logging"

	"golang.org/x/crypto/sha3"
)

// NewBlock creates a new block with the specified proposer ID, list of,
// transactions, and hash of the state calculated by calling State.GetHash()
// after running all transactions in the block and updating the state.
//
// TODO Remove proposerID parameter. This should be fetched from the configuration
// TODO Remove the stateHash parameter. The transactions in this block should
// be run when blockchain.AddBlock() is called. This function will then update
// the stateHash in this block.
// func NewBlock(proposerID string, transactions []transaction.Transaction, stateHash []byte) *Block {
// 	block := new(Block)
// 	block.proposerID = proposerID
// 	block.transactions = transactions
// 	block.stateHash = stateHash
// 	return block
// }

var log = logging.MustGetLogger("protos")

func NewBlock(proposerID string, transactions []*Transaction, stateHash []byte) *Block {
	block := new(Block)
	block.ProposerID = proposerID
	block.Transactions = transactions
	block.StateHash = stateHash
	return block
}

// GetHash returns the hash of this block.
func (block *Block) GetHash() ([]byte, error) {
	hash := make([]byte, 64)
	data, err := block.Bytes()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not calculate has of block: %s", err))
	}
	sha3.ShakeSum256(hash, data)
	return hash, nil
}

// GetStateHash returns the stateHash stored in this block. The stateHash
// is the value returned by state.GetHash() after running all transactions in
// the block.
func (block *Block) GetStateHash() []byte {
	return block.StateHash
}

// SetPreviousBlockHash sets the hash of the previous block. This will be
// called by blockchain.AddBlock when then the block is added.
func (block *Block) SetPreviousBlockHash(previousBlockHash []byte) {
	block.PreviousBlockHash = previousBlockHash
}

// Bytes returns this block as an array of bytes
func (block *Block) Bytes() ([]byte, error) {
	data, err := proto.Marshal(block)
	if err != nil {
		log.Error("Error marshalling block: %s", err)
		return nil, errors.New(fmt.Sprintf("Could not marshal block: %s", err))
	}
	return data, nil
}
