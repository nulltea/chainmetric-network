package intention

import (
	"context"
	"encoding/hex"
	"hash/fnv"

	"github.com/cnf/structhash"
)

// EventConcernBase partially implements EventConcern with base data structure and functionality.
type EventConcernBase struct {
	HashSum  string      `bson:"hash"`
	Kind     EventKind   `json:"kind" bson:"event_kind"`
	Contract string      `json:"contract" bson:"source_contract"`
	Args     interface{} `json:"args"`
}

// NewEventConcernBase constructs new EventConcernBase instance and calculates structural hash.
func NewEventConcernBase(kind EventKind, contract string, args interface{}) EventConcernBase {
	var ecb = EventConcernBase{
		Kind: kind,
		Contract: contract,
		Args: args,
	}

	ecb.HashSum = hex.EncodeToString(fnv.New64().Sum(structhash.Dump(ecb, 1)))

	return ecb
}

func (rv EventConcernBase) SourceContract() string {
	return rv.Contract
}

func (rv EventConcernBase) Context(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

func (rv EventConcernBase) OfKind() EventKind {
	return rv.Kind
}

func (rv EventConcernBase) Hash() string {
	return rv.HashSum
}

func (rv EventConcernBase) IsEqual(concern EventConcern) bool {
	return rv.HashSum == concern.Hash()
}
