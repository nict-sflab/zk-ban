package load

import (
	"os"
	"testing"

	zkban "github.com/akakou/zk-ban"
	"github.com/akakou/zk-ban/circuit"
	"github.com/akakou/zk-ban/dump"
	"github.com/akakou/zk-ban/snark"
)

func TestReDumpUserKey(t *testing.T) {
	oldKeyPath := dump.KeyPath
	dump.KeyPath = t.TempDir() + string(os.PathSeparator)
	t.Cleanup(func() {
		dump.KeyPath = oldKeyPath
	})

	const (
		name     = "redump-user-key"
		protocol = "join"
	)

	params, err := snark.InitSNARK(&circuit.JoinRequestCircuit{})
	if err != nil {
		t.Fatalf("init snark: %v", err)
	}

	dump.DumpSafeKeys(name, protocol, params)
	ReDumpUserKey(name, protocol)

	prover, err := LoadUserKey(name, protocol)
	if err != nil {
		t.Fatalf("load redumped user key: %v", err)
	}

	req, _, err := zkban.RequestJoin(20240101, prover)
	if err != nil {
		t.Fatalf("prove with redumped user key: %v", err)
	}

	if err := req.Verify(20240101, params.VerifyKey); err != nil {
		t.Fatalf("verify proof from redumped user key: %v", err)
	}
}
