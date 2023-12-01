// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package merkledb

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/database/memdb"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/maybe"
)

func Test_History_Simple(t *testing.T) {
	require := require.New(t)

	db, err := New(
		context.Background(),
		memdb.New(),
		newDefaultConfig(),
	)
	require.NoError(err)
	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value")))
	require.NoError(batch.Write())

	val, err := db.Get([]byte("key"))
	require.NoError(err)
	require.Equal([]byte("value"), val)

	origProof, err := db.GetRangeProof(context.Background(), []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(origProof)
	origRootID := db.root.id
	require.NoError(origProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value0")))
	require.NoError(batch.Write())
	newProof, err := db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key1"), []byte("value1")))
	require.NoError(batch.Put([]byte("key8"), []byte("value8")))
	require.NoError(batch.Write())
	newProof, err = db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("k"), []byte("v")))
	require.NoError(batch.Write())
	newProof, err = db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))

	batch = db.NewBatch()
	require.NoError(batch.Delete([]byte("k")))
	require.NoError(batch.Delete([]byte("ke")))
	require.NoError(batch.Delete([]byte("key")))
	require.NoError(batch.Delete([]byte("key1")))
	require.NoError(batch.Put([]byte("key2"), []byte("value2")))
	require.NoError(batch.Delete([]byte("key3")))
	require.NoError(batch.Delete([]byte("key4")))
	require.NoError(batch.Delete([]byte("key5")))
	require.NoError(batch.Delete([]byte("key8")))
	require.NoError(batch.Write())
	newProof, err = db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))
}

func Test_History_Large(t *testing.T) {
	require := require.New(t)

	numIters := 250

	for i := 1; i < 10; i++ {
		config := newDefaultConfig()
		// History must be large enough to get the change proof
		// after this loop. Multiply by four because every loop
		// iteration we do two puts and up to two deletes.
		config.HistoryLength = 4 * numIters
		db, err := New(
			context.Background(),
			memdb.New(),
			config,
		)
		require.NoError(err)
		roots := []ids.ID{}

		now := time.Now().UnixNano()
		t.Logf("seed for iter %d: %d", i, now)
		r := rand.New(rand.NewSource(now)) // #nosec G404
		// make sure they stay in sync
		for x := 0; x < numIters; x++ {
			addkey := make([]byte, r.Intn(50))
			_, err := r.Read(addkey)
			require.NoError(err)
			val := make([]byte, r.Intn(50))
			_, err = r.Read(val)
			require.NoError(err)

			require.NoError(db.Put(addkey, val))

			addNilkey := make([]byte, r.Intn(50))
			_, err = r.Read(addNilkey)
			require.NoError(err)
			require.NoError(db.Put(addNilkey, nil))

			deleteKeyStart := make([]byte, r.Intn(50))
			_, err = r.Read(deleteKeyStart)
			require.NoError(err)

			it := db.NewIteratorWithStart(deleteKeyStart)
			if it.Next() {
				require.NoError(db.Delete(it.Key()))
			}
			require.NoError(it.Error())
			it.Release()

			root, err := db.GetMerkleRoot(context.Background())
			require.NoError(err)
			roots = append(roots, root)
		}
		proof, err := db.GetRangeProofAtRoot(context.Background(), roots[0], nil, maybe.Nothing[[]byte](), 10)
		require.NoError(err)
		require.NotNil(proof)

		require.NoError(proof.Verify(context.Background(), nil, maybe.Nothing[[]byte](), roots[0]))
	}
}

func Test_History_Bad_GetValueChanges_Input(t *testing.T) {
	require := require.New(t)

	config := newDefaultConfig()
	config.HistoryLength = 5

	db, err := New(
		context.Background(),
		memdb.New(),
		config,
	)
	require.NoError(err)

	// Do 5 puts (i.e. the history length)
	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value")))
	require.NoError(batch.Write())

	toBeDeletedRoot := db.getMerkleRoot()

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value0")))
	require.NoError(batch.Write())

	startRoot := db.getMerkleRoot()

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key1"), []byte("value0")))
	require.NoError(batch.Write())

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key1"), []byte("value1")))
	require.NoError(batch.Write())

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key2"), []byte("value3")))
	require.NoError(batch.Write())

	endRoot := db.getMerkleRoot()

	// ensure these start as valid calls
	_, err = db.history.getValueChanges(toBeDeletedRoot, endRoot, nil, maybe.Nothing[[]byte](), 1)
	require.NoError(err)
	_, err = db.history.getValueChanges(startRoot, endRoot, nil, maybe.Nothing[[]byte](), 1)
	require.NoError(err)

	_, err = db.history.getValueChanges(startRoot, endRoot, nil, maybe.Nothing[[]byte](), -1)
	require.ErrorIs(err, ErrInvalidMaxLength)

	_, err = db.history.getValueChanges(endRoot, startRoot, nil, maybe.Nothing[[]byte](), 1)
	require.ErrorIs(err, ErrInsufficientHistory)

	// trigger the first root to be deleted by exiting the lookback window
	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key2"), []byte("value4")))
	require.NoError(batch.Write())

	// now this root should no longer be present
	_, err = db.history.getValueChanges(toBeDeletedRoot, endRoot, nil, maybe.Nothing[[]byte](), 1)
	require.ErrorIs(err, ErrInsufficientHistory)

	// same start/end roots should yield an empty changelist
	changes, err := db.history.getValueChanges(endRoot, endRoot, nil, maybe.Nothing[[]byte](), 10)
	require.NoError(err)
	require.Empty(changes.values)
}

func Test_History_Trigger_History_Queue_Looping(t *testing.T) {
	require := require.New(t)

	config := newDefaultConfig()
	config.HistoryLength = 2

	db, err := New(
		context.Background(),
		memdb.New(),
		config,
	)
	require.NoError(err)

	// Do 2 puts (i.e. the history length)
	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value")))
	require.NoError(batch.Write())
	origRootID := db.getMerkleRoot()

	origProof, err := db.GetRangeProof(context.Background(), []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(origProof)
	require.NoError(origProof.Verify(
		context.Background(),
		[]byte("k"),
		maybe.Some([]byte("key3")),
		origRootID,
	))

	// write a new value into the db, now there should be 2 roots in the history
	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value0")))
	require.NoError(batch.Write())

	// ensure that previous root is still present and generates a valid proof
	newProof, err := db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(
		context.Background(),
		[]byte("k"),
		maybe.Some([]byte("key3")),
		origRootID,
	))

	// trigger a new root to be added to the history, which should cause rollover since there can only be 2
	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key1"), []byte("value1")))
	require.NoError(batch.Write())

	// proof from first root shouldn't be generatable since it should have been removed from the history
	_, err = db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.ErrorIs(err, ErrInsufficientHistory)
}

func Test_History_Values_Lookup_Over_Queue_Break(t *testing.T) {
	require := require.New(t)

	config := newDefaultConfig()
	config.HistoryLength = 4
	db, err := New(
		context.Background(),
		memdb.New(),
		config,
	)
	require.NoError(err)

	// Do 4 puts (i.e. the history length)
	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value")))
	require.NoError(batch.Write())

	// write a new value into the db
	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value0")))
	require.NoError(batch.Write())

	startRoot := db.getMerkleRoot()

	// write a new value into the db
	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key1"), []byte("value0")))
	require.NoError(batch.Write())

	// write a new value into the db that overwrites key1
	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key1"), []byte("value1")))
	require.NoError(batch.Write())

	// trigger a new root to be added to the history, which should cause rollover since there can only be 3
	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key2"), []byte("value3")))
	require.NoError(batch.Write())

	endRoot := db.getMerkleRoot()

	// changes should still be collectable even though the history has had to loop due to hitting max size
	changes, err := db.history.getValueChanges(startRoot, endRoot, nil, maybe.Nothing[[]byte](), 10)
	require.NoError(err)
	require.Contains(changes.values, newPath([]byte("key1")))
	require.Equal([]byte("value1"), changes.values[newPath([]byte("key1"))].after.Value())
	require.Contains(changes.values, newPath([]byte("key2")))
	require.Equal([]byte("value3"), changes.values[newPath([]byte("key2"))].after.Value())
}

func Test_History_RepeatedRoot(t *testing.T) {
	require := require.New(t)

	db, err := New(
		context.Background(),
		memdb.New(),
		newDefaultConfig(),
	)
	require.NoError(err)
	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key1"), []byte("value1")))
	require.NoError(batch.Put([]byte("key2"), []byte("value2")))
	require.NoError(batch.Put([]byte("key3"), []byte("value3")))
	require.NoError(batch.Write())

	origProof, err := db.GetRangeProof(context.Background(), []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(origProof)
	origRootID := db.root.id
	require.NoError(origProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key1"), []byte("other")))
	require.NoError(batch.Put([]byte("key2"), []byte("other")))
	require.NoError(batch.Put([]byte("key3"), []byte("other")))
	require.NoError(batch.Write())
	newProof, err := db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))

	// revert state to be the same as in orig proof
	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key1"), []byte("value1")))
	require.NoError(batch.Put([]byte("key2"), []byte("value2")))
	require.NoError(batch.Put([]byte("key3"), []byte("value3")))
	require.NoError(batch.Write())

	newProof, err = db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))
}

func Test_History_ExcessDeletes(t *testing.T) {
	require := require.New(t)

	db, err := New(
		context.Background(),
		memdb.New(),
		newDefaultConfig(),
	)
	require.NoError(err)
	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value")))
	require.NoError(batch.Write())

	origProof, err := db.GetRangeProof(context.Background(), []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(origProof)
	origRootID := db.root.id
	require.NoError(origProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))

	batch = db.NewBatch()
	require.NoError(batch.Delete([]byte("key1")))
	require.NoError(batch.Delete([]byte("key2")))
	require.NoError(batch.Delete([]byte("key3")))
	require.NoError(batch.Delete([]byte("key4")))
	require.NoError(batch.Delete([]byte("key5")))
	require.NoError(batch.Write())
	newProof, err := db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))
}

func Test_History_DontIncludeAllNodes(t *testing.T) {
	require := require.New(t)

	db, err := New(
		context.Background(),
		memdb.New(),
		newDefaultConfig(),
	)
	require.NoError(err)
	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value")))
	require.NoError(batch.Write())

	origProof, err := db.GetRangeProof(context.Background(), []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(origProof)
	origRootID := db.root.id
	require.NoError(origProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("z"), []byte("z")))
	require.NoError(batch.Write())
	newProof, err := db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))
}

func Test_History_Branching2Nodes(t *testing.T) {
	require := require.New(t)

	db, err := New(
		context.Background(),
		memdb.New(),
		newDefaultConfig(),
	)
	require.NoError(err)
	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value")))
	require.NoError(batch.Write())

	origProof, err := db.GetRangeProof(context.Background(), []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(origProof)
	origRootID := db.root.id
	require.NoError(origProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("k"), []byte("v")))
	require.NoError(batch.Write())
	newProof, err := db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))
}

func Test_History_Branching3Nodes(t *testing.T) {
	require := require.New(t)

	db, err := New(
		context.Background(),
		memdb.New(),
		newDefaultConfig(),
	)
	require.NoError(err)
	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key123"), []byte("value123")))
	require.NoError(batch.Write())

	origProof, err := db.GetRangeProof(context.Background(), []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(origProof)
	origRootID := db.root.id
	require.NoError(origProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key321"), []byte("value321")))
	require.NoError(batch.Write())
	newProof, err := db.GetRangeProofAtRoot(context.Background(), origRootID, []byte("k"), maybe.Some([]byte("key3")), 10)
	require.NoError(err)
	require.NotNil(newProof)
	require.NoError(newProof.Verify(context.Background(), []byte("k"), maybe.Some([]byte("key3")), origRootID))
}

func Test_History_MaxLength(t *testing.T) {
	require := require.New(t)

	config := newDefaultConfig()
	config.HistoryLength = 2
	db, err := New(
		context.Background(),
		memdb.New(),
		config,
	)
	require.NoError(err)

	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key"), []byte("value")))
	require.NoError(batch.Write())

	oldRoot, err := db.GetMerkleRoot(context.Background())
	require.NoError(err)

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("k"), []byte("v")))
	require.NoError(batch.Write())

	require.Contains(db.history.lastChanges, oldRoot)

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("k1"), []byte("v2"))) // Overwrites oldest element in history
	require.NoError(batch.Write())

	require.NotContains(db.history.lastChanges, oldRoot)
}

func Test_Change_List(t *testing.T) {
	require := require.New(t)

	db, err := New(
		context.Background(),
		memdb.New(),
		newDefaultConfig(),
	)
	require.NoError(err)
	batch := db.NewBatch()
	require.NoError(batch.Put([]byte("key20"), []byte("value20")))
	require.NoError(batch.Put([]byte("key21"), []byte("value21")))
	require.NoError(batch.Put([]byte("key22"), []byte("value22")))
	require.NoError(batch.Put([]byte("key23"), []byte("value23")))
	require.NoError(batch.Put([]byte("key24"), []byte("value24")))
	require.NoError(batch.Write())
	startRoot, err := db.GetMerkleRoot(context.Background())
	require.NoError(err)

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key25"), []byte("value25")))
	require.NoError(batch.Put([]byte("key26"), []byte("value26")))
	require.NoError(batch.Put([]byte("key27"), []byte("value27")))
	require.NoError(batch.Put([]byte("key28"), []byte("value28")))
	require.NoError(batch.Put([]byte("key29"), []byte("value29")))
	require.NoError(batch.Write())

	batch = db.NewBatch()
	require.NoError(batch.Put([]byte("key30"), []byte("value30")))
	require.NoError(batch.Put([]byte("key31"), []byte("value31")))
	require.NoError(batch.Put([]byte("key32"), []byte("value32")))
	require.NoError(batch.Delete([]byte("key21")))
	require.NoError(batch.Delete([]byte("key22")))
	require.NoError(batch.Write())

	endRoot, err := db.GetMerkleRoot(context.Background())
	require.NoError(err)

	changes, err := db.history.getValueChanges(startRoot, endRoot, nil, maybe.Nothing[[]byte](), 8)
	require.NoError(err)
	require.Len(changes.values, 8)
}

func TestHistoryRecord(t *testing.T) {
	require := require.New(t)

	maxHistoryLen := 3
	th := newTrieHistory(maxHistoryLen)

	changes := []*changeSummary{}
	for i := 0; i < maxHistoryLen; i++ { // Fill the history
		changes = append(changes, &changeSummary{rootID: ids.GenerateTestID()})

		th.record(changes[i])
		require.Equal(uint64(i+1), th.nextInsertNumber)
		require.Equal(i+1, th.history.Len())
		require.Len(th.lastChanges, i+1)
		require.Contains(th.lastChanges, changes[i].rootID)
		changeAndIndex := th.lastChanges[changes[i].rootID]
		require.Equal(uint64(i), changeAndIndex.insertNumber)
		got, ok := th.history.Index(int(changeAndIndex.insertNumber))
		require.True(ok)
		require.Equal(changes[i], got.changeSummary)
	}
	// history is [changes[0], changes[1], changes[2]]

	// Add a new change
	change3 := &changeSummary{rootID: ids.GenerateTestID()}
	th.record(change3)
	// history is [changes[1], changes[2], change3]
	require.Equal(uint64(maxHistoryLen+1), th.nextInsertNumber)
	require.Equal(maxHistoryLen, th.history.Len())
	require.Len(th.lastChanges, maxHistoryLen)
	require.Contains(th.lastChanges, change3.rootID)
	changeAndIndex := th.lastChanges[change3.rootID]
	require.Equal(uint64(maxHistoryLen), changeAndIndex.insertNumber)
	got, ok := th.history.PeekRight()
	require.True(ok)
	require.Equal(change3, got.changeSummary)

	// // Make sure the oldest change was evicted
	require.NotContains(th.lastChanges, changes[0].rootID)
	oldestChange, ok := th.history.PeekLeft()
	require.True(ok)
	require.Equal(uint64(1), oldestChange.insertNumber)

	// Add another change which was the same root ID as changes[2]
	change4 := &changeSummary{rootID: changes[2].rootID}
	th.record(change4)
	// history is [changes[2], change3, change4]

	change5 := &changeSummary{rootID: ids.GenerateTestID()}
	th.record(change5)
	// history is [change3, change4, change5]

	// Make sure that even though changes[2] was evicted, we still remember
	// that the most recent change resulting in that change's root ID.
	require.Len(th.lastChanges, maxHistoryLen)
	require.Contains(th.lastChanges, changes[2].rootID)
	changeAndIndex = th.lastChanges[changes[2].rootID]
	require.Equal(uint64(maxHistoryLen+1), changeAndIndex.insertNumber)

	// Make sure [t.history] is right.
	require.Equal(maxHistoryLen, th.history.Len())
	got, ok = th.history.PopLeft()
	require.True(ok)
	require.Equal(uint64(maxHistoryLen), got.insertNumber)
	require.Equal(change3.rootID, got.rootID)
	got, ok = th.history.PopLeft()
	require.True(ok)
	require.Equal(uint64(maxHistoryLen+1), got.insertNumber)
	require.Equal(change4.rootID, got.rootID)
	got, ok = th.history.PopLeft()
	require.True(ok)
	require.Equal(uint64(maxHistoryLen+2), got.insertNumber)
	require.Equal(change5.rootID, got.rootID)
}

func TestHistoryGetChangesToRoot(t *testing.T) {
	maxHistoryLen := 3
	history := newTrieHistory(maxHistoryLen)

	changes := []*changeSummary{}
	for i := 0; i < maxHistoryLen; i++ { // Fill the history
		changes = append(changes, &changeSummary{
			rootID: ids.GenerateTestID(),
			nodes: map[path]*change[*node]{
				newPath([]byte{byte(i)}): {
					before: &node{id: ids.GenerateTestID()},
					after:  &node{id: ids.GenerateTestID()},
				},
			},
			values: map[path]*change[maybe.Maybe[[]byte]]{
				newPath([]byte{byte(i)}): {
					before: maybe.Some([]byte{byte(i)}),
					after:  maybe.Some([]byte{byte(i + 1)}),
				},
			},
		})
		history.record(changes[i])
	}

	type test struct {
		name         string
		rootID       ids.ID
		start        []byte
		end          maybe.Maybe[[]byte]
		validateFunc func(*require.Assertions, *changeSummary)
		expectedErr  error
	}

	tests := []test{
		{
			name:        "unknown root ID",
			rootID:      ids.GenerateTestID(),
			expectedErr: ErrInsufficientHistory,
		},
		{
			name:   "most recent change",
			rootID: changes[maxHistoryLen-1].rootID,
			validateFunc: func(require *require.Assertions, got *changeSummary) {
				require.Equal(newChangeSummary(defaultPreallocationSize), got)
			},
		},
		{
			name:   "second most recent change",
			rootID: changes[maxHistoryLen-2].rootID,
			validateFunc: func(require *require.Assertions, got *changeSummary) {
				// Ensure this is the reverse of the most recent change
				require.Len(got.nodes, 1)
				require.Len(got.values, 1)
				reversedChanges := changes[maxHistoryLen-1]
				removedKey := newPath([]byte{byte(maxHistoryLen - 1)})
				require.Equal(reversedChanges.nodes[removedKey].before, got.nodes[removedKey].after)
				require.Equal(reversedChanges.values[removedKey].before, got.values[removedKey].after)
				require.Equal(reversedChanges.values[removedKey].after, got.values[removedKey].before)
			},
		},
		{
			name:   "third most recent change",
			rootID: changes[maxHistoryLen-3].rootID,
			validateFunc: func(require *require.Assertions, got *changeSummary) {
				require.Len(got.nodes, 2)
				require.Len(got.values, 2)
				reversedChanges1 := changes[maxHistoryLen-1]
				removedKey1 := newPath([]byte{byte(maxHistoryLen - 1)})
				require.Equal(reversedChanges1.nodes[removedKey1].before, got.nodes[removedKey1].after)
				require.Equal(reversedChanges1.values[removedKey1].before, got.values[removedKey1].after)
				require.Equal(reversedChanges1.values[removedKey1].after, got.values[removedKey1].before)
				reversedChanges2 := changes[maxHistoryLen-2]
				removedKey2 := newPath([]byte{byte(maxHistoryLen - 2)})
				require.Equal(reversedChanges2.nodes[removedKey2].before, got.nodes[removedKey2].after)
				require.Equal(reversedChanges2.values[removedKey2].before, got.values[removedKey2].after)
				require.Equal(reversedChanges2.values[removedKey2].after, got.values[removedKey2].before)
			},
		},
		{
			name:   "third most recent change with start filter",
			rootID: changes[maxHistoryLen-3].rootID,
			start:  []byte{byte(maxHistoryLen - 1)}, // Omit values from second most recent change
			validateFunc: func(require *require.Assertions, got *changeSummary) {
				require.Len(got.nodes, 2)
				require.Len(got.values, 1)
				reversedChanges1 := changes[maxHistoryLen-1]
				removedKey1 := newPath([]byte{byte(maxHistoryLen - 1)})
				require.Equal(reversedChanges1.nodes[removedKey1].before, got.nodes[removedKey1].after)
				require.Equal(reversedChanges1.values[removedKey1].before, got.values[removedKey1].after)
				require.Equal(reversedChanges1.values[removedKey1].after, got.values[removedKey1].before)
				reversedChanges2 := changes[maxHistoryLen-2]
				removedKey2 := newPath([]byte{byte(maxHistoryLen - 2)})
				require.Equal(reversedChanges2.nodes[removedKey2].before, got.nodes[removedKey2].after)
			},
		},
		{
			name:   "third most recent change with end filter",
			rootID: changes[maxHistoryLen-3].rootID,
			end:    maybe.Some([]byte{byte(maxHistoryLen - 2)}), // Omit values from most recent change
			validateFunc: func(require *require.Assertions, got *changeSummary) {
				require.Len(got.nodes, 2)
				require.Len(got.values, 1)
				reversedChanges1 := changes[maxHistoryLen-1]
				removedKey1 := newPath([]byte{byte(maxHistoryLen - 1)})
				require.Equal(reversedChanges1.nodes[removedKey1].before, got.nodes[removedKey1].after)
				reversedChanges2 := changes[maxHistoryLen-2]
				removedKey2 := newPath([]byte{byte(maxHistoryLen - 2)})
				require.Equal(reversedChanges2.nodes[removedKey2].before, got.nodes[removedKey2].after)
				require.Equal(reversedChanges2.values[removedKey2].before, got.values[removedKey2].after)
				require.Equal(reversedChanges2.values[removedKey2].after, got.values[removedKey2].before)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			got, err := history.getChangesToGetToRoot(tt.rootID, tt.start, tt.end)
			require.ErrorIs(err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}
			tt.validateFunc(require, got)
		})
	}
}
