package hdwallet

import (
	"testing"
)

func TestGenerateAddress(t *testing.T) {
	address, err := GenerateAddress("xpub6C8Zf1BmXA11nWYYt7CLVU7VDEJayfRyYkM6TUnLzFgLbofCYYhUR9Fdbufhrp495YN3KBHaXETx2shPiWFSRqieJy3AQrw3R1KeYzqApDB", 0, 0)

	if err != nil {
		t.Fatalf("expected to see no error saw %v", err)
	}

	if address != "0x96B148BD4759256be05B0CdcA1401734B56fD9Af" {
		t.Fatalf("expected to see address %v saw address %v", "0x96B148BD4759256be05B0CdcA1401734B56fD9Af", address)
	}

	_, err = GenerateAddress("invalid PubK", 0, 0)

	if err == nil {
		t.Fatal("expected to see error saw no error")
	}
}

func TestGetPrivateKey(t *testing.T) {
	pvKey, err := GetPrivateKey("xprv9y9DFVesgnSia2U5n5fL8LAkfCU6aCi8BXRVf6NjRv9Mj1L411PDsLw9kc9Vm98qtAqg6bAazHewPrpTnwdoEMgTNJigfvuhUhJb6RiQCsb", 0, 0)
	if err != nil {
		t.Fatalf("expected to see no error saw %v", err)
	}

	if pvKey != "0x31fa0959ea766b401b72465e7d8e41918d74ffba94d7a43dffff1bfc866280cb" {
		t.Fatalf("expected to see address %v saw address %v", "0x31fa0959ea766b401b72465e7d8e41918d74ffba94d7a43dffff1bfc866280cb", pvKey)
	}

	_, err = GetPrivateKey("invalid PvK", 0, 0)
	if err == nil {
		t.Fatal("expected to see error saw no error")
	}
}
