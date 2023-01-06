package polygon

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/0xPolygon/polygon-edge/helper/keccak"
	"github.com/0xPolygon/polygon-edge/types"

	//"github.com/ethereum/go-ethereum/core/types"

	"github.com/umbracle/fastrlp"
)

func ibftConsensusHashEthereum(h *types.Header) ([]byte, error) {
	arena := fastrlp.DefaultArenaPool.Get()
	defer fastrlp.DefaultArenaPool.Put(arena)

	/*
		vv := arena.NewArray()
		vv.Set(arena.NewBytes(h.ParentHash.Bytes()))
		vv.Set(arena.NewBytes(h.UncleHash.Bytes()))
		vv.Set(arena.NewBytes(h.Coinbase.Bytes()))
		vv.Set(arena.NewBytes(h.Root.Bytes()))
		vv.Set(arena.NewBytes(h.TxHash.Bytes()))
		vv.Set(arena.NewBytes(h.ReceiptHash.Bytes()))
		vv.Set(arena.NewBytes(h.Bloom[:]))
		vv.Set(arena.NewBigInt(h.Difficulty))
		vv.Set(arena.NewBigInt(h.Number))
		vv.Set(arena.NewUint(h.GasLimit))
		vv.Set(arena.NewUint(h.GasUsed))
		vv.Set(arena.NewUint(h.Time))
		vv.Set(arena.NewCopyBytes(h.Extra))

		// buf := keccak.Keccak256Rlp(nil, vv)
		hash := crypto.Keccak256Hash([]byte("Balance(uint256)"))
		return buf, nil
	*/
	return nil, nil
}

func ibftConsensusHashEdgeFromJson(Data []byte) (string, error) {
	// fmt.Printf("--- Data: %v", string(Data))

	var b block
	err := json.Unmarshal(Data, &b)
	if err != nil {
		return "", err
	}

	// b1, _ := json.Marshal(b)
	// fmt.Printf("--- Block: %v", string(b1))

	return b.ComputeHash(), nil
}

// MarshalRLPWith marshals the header to RLP with a specific fastrlp.Arena
func (h *block) MarshalRLPWith(arena *fastrlp.Arena) *fastrlp.Value {
	vv := arena.NewArray()

	vv.Set(arena.NewBytes(h.ParentHash.Bytes()))
	vv.Set(arena.NewBytes(h.Sha3Uncles.Bytes()))
	vv.Set(arena.NewCopyBytes(h.Miner[:]))
	vv.Set(arena.NewBytes(h.StateRoot.Bytes()))
	vv.Set(arena.NewBytes(h.TxRoot.Bytes()))
	vv.Set(arena.NewBytes(h.ReceiptsRoot.Bytes()))
	vv.Set(arena.NewCopyBytes(h.LogsBloom[:]))

	vv.Set(arena.NewUint(uint64(h.Difficulty)))
	vv.Set(arena.NewUint(uint64(h.Number)))
	vv.Set(arena.NewUint(uint64(h.GasLimit)))
	vv.Set(arena.NewUint(uint64(h.GasUsed)))
	vv.Set(arena.NewUint(uint64(h.Timestamp)))

	vv.Set(arena.NewCopyBytes(h.ExtraData))
	vv.Set(arena.NewBytes(h.MixHash.Bytes()))
	vv.Set(arena.NewCopyBytes(h.Nonce[:]))

	return vv
}

var marshalArenaPool fastrlp.ArenaPool

func DefHeaderHash(h *block) (hash types.Hash) {
	// default header hashing
	ar := marshalArenaPool.Get()
	hasher := keccak.DefaultKeccakPool.Get()

	v := h.MarshalRLPWith(ar)
	hasher.WriteRlp(hash[:0], v)

	marshalArenaPool.Put(ar)
	keccak.DefaultKeccakPool.Put(hasher)

	return
}

// ComputeHash computes the hash of the header
func (h *block) ComputeHash() string {
	hash := DefHeaderHash(h)

	return hash.String()
}

func TestPolygonParser_GetBlockHash(t *testing.T) {
	type args struct {
		blockHeight int64
		blockJson   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "quicknode",
			args: args{blockHeight: 15974356, blockJson: "0xf3bfd4--quicknode.json"},
			want: "0x929e50eb1acd370284e2f4b5069cbde30d777222696d4d40d6fe52d61f1f4501",
		},
		{
			name: "polygon-rpc",
			args: args{blockHeight: 15974356, blockJson: "0xf3bfd4--polygon-rpc.json"},
			want: "0x929e50eb1acd370284e2f4b5069cbde30d777222696d4d40d6fe52d61f1f4501",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("testsuite", tt.args.blockJson)
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal("error reading source file:", err)
			}

			got, err := ibftConsensusHashEdgeFromJson(source)
			if err != nil {
				t.Fatal("error computing hash:", err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Polygon.GetBlockHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Polygon.GetBlockHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
