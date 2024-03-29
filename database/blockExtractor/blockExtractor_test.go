package blockExtractor_test

import (
	"testing"

	. "github.com/DCNT-Hammer/dcnt/database/blockExtractor"
	"github.com/DCNT-Hammer/dcnt/testHelper"
)

func TestTest(t *testing.T) {
	dbo := testHelper.CreateAndPopulateTestDatabaseOverlay()

	be := new(BlockExtractor)

	err := be.ExportDChain(dbo)
	if err != nil {
		t.Error(err)
	}
	err = be.ExportECChain(dbo)
	if err != nil {
		t.Error(err)
	}
	err = be.ExportAChain(dbo)
	if err != nil {
		t.Error(err)
	}
	err = be.ExportFctChain(dbo)
	if err != nil {
		t.Error(err)
	}
	err = be.ExportDirBlockInfo(dbo)
	if err != nil {
		t.Error(err)
	}
	err = be.ExportEChain(testHelper.GetChainID().String(), dbo)
	if err != nil {
		t.Error(err)
	}
	err = be.ExportEChain(testHelper.GetAnchorChainID().String(), dbo)
	if err != nil {
		t.Error(err)
	}
}
