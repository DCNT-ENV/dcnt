// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package identity_test

import (
	"testing"

	"bytes"
	"encoding/json"

	"github.com/DCNT-Hammer/dcnt/common/constants"
	"github.com/DCNT-Hammer/dcnt/common/factoid"
	. "github.com/DCNT-Hammer/dcnt/common/identity"
	"github.com/DCNT-Hammer/dcnt/common/interfaces"
	"github.com/DCNT-Hammer/dcnt/common/messages"
	"github.com/DCNT-Hammer/dcnt/common/primitives"
	"github.com/DCNT-Hammer/dcnt/testHelper"
)

func TestAuthorityType(t *testing.T) {
	auth := new(Authority)
	if auth.Type() != -1 {
		t.Errorf("Invalid type returned - %v", auth.Type())
	}

	auth.Status = constants.IDENTITY_FEDERATED_SERVER
	if auth.Type() != 1 {
		t.Errorf("Invalid type returned - %v", auth.Type())
	}

	auth.Status = constants.IDENTITY_AUDIT_SERVER
	if auth.Type() != 0 {
		t.Errorf("Invalid type returned - %v", auth.Type())
	}
}

//func TestAuthoritySignature(t *testing.T) {
//	s := testHelper.CreateAndPopulateTestState()
//	idindex := s.CreateBlankFactomIdentity(primitives.NewZeroHash())
//	s.Identities[idindex].ManagementChainID = primitives.NewZeroHash()
//
//	index := s.AddAuthorityFromChainID(primitives.NewZeroHash())
//	s.Authorities[index].SigningKey = *(s.GetServerPublicKey())
//	s.Authorities[index].Status = 1
//
//	ack := new(messages.Ack)
//	ack.DBHeight = s.LLeaderHeight
//	ack.VMIndex = 1
//	ack.Minute = byte(5)
//	ack.Timestamp = s.GetTimestamp()
//	ack.MessageHash = primitives.NewZeroHash()
//	ack.LeaderChainID = s.IdentityChainID
//	ack.SerialHash = primitives.NewZeroHash()
//
//	err := ack.Sign(s)
//	if err != nil {
//		t.Error("Authority Test Failed when signing message")
//	}
//
//	msg, err := ack.MarshalForSignature()
//	if err != nil {
//		t.Error("Authority Test Failed when marshalling for sig")
//	}
//
//	sig := ack.GetSignature()
//	server, err := s.Authorities[0].VerifySignature(msg, sig.GetSignature())
//	if !server || err != nil {
//		t.Error("Authority Test Failed when checking sigs")
//	}
//}

//func TestMarshalJSON(t *testing.T) {
//	s := testHelper.CreateAndPopulateTestState()
//	idindex := s.CreateBlankFactomIdentity(primitives.NewZeroHash())
//	s.Identities[idindex].ManagementChainID = primitives.NewZeroHash()
//
//	index := s.AddAuthorityFromChainID(primitives.NewZeroHash())
//	s.Authorities[index].SigningKey = *(s.GetServerPublicKey())
//	s.Authorities[index].Status = 1
//
//	j, err := s.Authorities[index].MarshalJSON()
//	if err != nil {
//		t.Errorf("%v", err)
//	}
//
//	expected := `{"chainid":"0000000000000000000000000000000000000000000000000000000000000000","manageid":"0000000000000000000000000000000000000000000000000000000000000000","matroyshka":null,"signingkey":"cc1985cdfae4e32b5a454dfda8ce5e1361558482684f3367649c3ad852c8e31a","status":"federated","anchorkeys":null}`
//	if string(j) != expected {
//		t.Errorf("Invalid json returned - %v vs %v", string(j), expected)
//	}
//}

func TestAuthorityMarshalUnmarshal(t *testing.T) {
	for i := 0; i < 1000; i++ {
		a := RandomAuthority()

		h, err := a.MarshalBinary()
		if err != nil {
			t.Errorf("%v", err)
		}
		a2 := new(Authority)
		err = a2.UnmarshalBinary(h)
		if err != nil {
			t.Errorf("%v", err)
		}
		if a.IsSameAs(a2) == false {
			t.Errorf("Authorities are not identical")
		}
	}
}

func TestAuthorityClone(t *testing.T) {
	for i := 0; i < 1000; i++ {
		auth := RandomAuthority()
		auth2 := auth.Clone()
		if auth.IsSameAs(auth2) == false {
			t.Errorf("Authorities are not the same")
		}

		// Check their marshalled values
		d1, err := auth.MarshalBinary()
		if err != nil {
			t.Errorf("%v", err)
		}
		d2, err := auth2.MarshalBinary()
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Compare(d1, d2) != 0 {
			t.Errorf("Authorities are not the same when marshalled")
		}
	}
}

func TestVerify(t *testing.T) {
	s := testHelper.CreateEmptyTestState()
	pl := s.ProcessLists.Get(10)

	var privs []*primitives.PrivateKey
	var ids []interfaces.IHash
	for i := 0; i < 10; i++ {
		p := primitives.RandomPrivateKey()
		id := primitives.RandomHash()
		privs = append(privs, p)
		ids = append(ids, id)

		auth := new(Authority)
		if i%2 == 0 {
			pl.AddAuditServer(id)
			auth.Status = constants.IDENTITY_AUDIT_SERVER
		} else {
			pl.AddFedServer(id)
			auth.Status = constants.IDENTITY_FEDERATED_SERVER
		}
		auth.AuthorityChainID = id
		auth.SigningKey = *(p.Pub)

		s.IdentityControl.SetAuthority(auth.AuthorityChainID, auth)
	}

	for i := 0; i < len(ids); i++ {
		msg := newAck(ids[i], s.GetTimestamp())
		msg.Sign(privs[i])
		b, _ := msg.MarshalForSignature()
		v, err := s.VerifyAuthoritySignature(b, msg.GetSignature().GetSignature(), 10)
		if err != nil {
			t.Error(err)
		}
		if i%2 == 0 {
			if v != 0 {
				t.Errorf("Should be 0 for audit, found %d", v)
			}
		} else {
			if v != 1 {
				t.Errorf("Should be 1 for fed, found %d", v)
			}
		}

		v2, err := s.FastVerifyAuthoritySignature(b, msg.GetSignature(), 10)
		if err != nil {
			t.Error(err)
		}

		if v != v2 {
			t.Error("Should be equal validates")
		}
	}

}

func TestSameAuth(t *testing.T) {
	a := RandomAuthority()
	d, _ := a.MarshalBinary()
	b := new(Authority)
	b.UnmarshalBinary(d)

	if !a.IsSameAs(b) {
		t.Error("Should be same, both empty")
	}

	a.AuthorityChainID = primitives.RandomHash()
	if a.IsSameAs(b) {
		t.Error("Diff auth chains, should be different")
	}
	a.AuthorityChainID = b.AuthorityChainID

	a.ManagementChainID = primitives.RandomHash()
	if a.IsSameAs(b) {
		t.Error("Diff ManagementChainID chains, should be different")
	}
	a.ManagementChainID = b.ManagementChainID

	a.MatryoshkaHash = primitives.RandomHash()
	if a.IsSameAs(b) {
		t.Error("Diff MatryoshkaHash chains, should be different")
	}
	a.MatryoshkaHash = b.MatryoshkaHash

	a.SigningKey = *(primitives.RandomPrivateKey().Pub)
	if a.IsSameAs(b) {
		t.Error("Diff SigningKey chains, should be different")
	}
	a.SigningKey = b.SigningKey

	a.Status = b.Status + 1
	if a.IsSameAs(b) {
		t.Error("Diff Status chains, should be different")
	}
	a.Status = b.Status

	a.AnchorKeys = append(a.AnchorKeys, AnchorSigningKey{})
	if a.IsSameAs(b) {
		t.Error("Diff Status AnchorKeys, should be different")
	}
	a.AnchorKeys = b.AnchorKeys

	a.KeyHistory = append(a.KeyHistory, HistoricKey{})
	if a.IsSameAs(b) {
		t.Error("Diff Status KeyHistory, should be different")
	}
	a.KeyHistory = b.KeyHistory

}

func newAck(id interfaces.IHash, ts interfaces.Timestamp) *messages.Ack {
	ack := new(messages.Ack)
	ack.DBHeight = 0
	ack.VMIndex = 1
	ack.Minute = byte(5)
	ack.Timestamp = ts
	ack.MessageHash = primitives.NewZeroHash()
	ack.LeaderChainID = id
	ack.SerialHash = primitives.NewZeroHash()

	return ack
}

func TestAuthorityJsonMarshal(t *testing.T) {
	// Testing Human readable json marshal
	a := NewAuthority()
	a.CoinbaseAddress = factoid.NewAddress(make([]byte, 32))
	a.Efficiency = 100

	data, err := a.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	var dst bytes.Buffer
	exp := `
		{
			"chainid": "0000000000000000000000000000000000000000000000000000000000000000",
			"manageid": "0000000000000000000000000000000000000000000000000000000000000000",
			"matroyshka": "0000000000000000000000000000000000000000000000000000000000000000",
			"signingkey": "0000000000000000000000000000000000000000000000000000000000000000",
			"status": "none",
			"anchorkeys": null,
			"efficiency": 100,
			"coinbaseaddress": "FA1y5ZGuHSLmf2TqNf6hVMkPiNGyQpQDTFJvDLRkKQaoPo4bmbgu"
		}
	`
	json.Compact(&dst, []byte(exp))
	if bytes.Compare(dst.Bytes(), data) != 0 {
		t.Errorf("Does not match expected")
	}
}
