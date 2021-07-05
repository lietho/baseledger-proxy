package types

import (
	"database/sql"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres

	uuid "github.com/kthomas/go.uuid"
	common "github.com/unibrightio/baseledger/common"
	"github.com/unibrightio/baseledger/dbutil"
	"github.com/unibrightio/baseledger/logger"
)

type TrustmeshEntry struct {
	TendermintBlockId                    sql.NullString
	TendermintTransactionId              uuid.UUID
	TendermintTransactionTimestamp       sql.NullTime
	EntryType                            string
	SenderOrgId                          uuid.UUID
	ReceiverOrgId                        uuid.UUID
	WorkgroupId                          uuid.UUID
	WorkstepType                         string
	BaseledgerTransactionType            string
	BaseledgerTransactionId              uuid.UUID
	ReferencedBaseledgerTransactionId    uuid.UUID
	BusinessObjectType                   string
	BaseledgerBusinessObjectId           uuid.UUID
	ReferencedBaseledgerBusinessObjectId uuid.UUID
	OffchainProcessMessageId             uuid.UUID
	ReferencedProcessMessageId           uuid.UUID
	CommitmentState                      string
	TransactionHash                      string
	TrustmeshId                          uuid.UUID
}

type Trustmesh struct {
	Id                  uuid.UUID
	CreatedAt           time.Time
	StartTime           time.Time
	EndTime             time.Time
	Participants        string
	BusinessObjectTypes string
	Finalized           bool
	ContainsRejections  bool
	Entries             []TrustmeshEntry
}

func (t *TrustmeshEntry) Create() bool {
	t.CommitmentState = common.UncommittedCommitmentState
	t.TendermintBlockId = sql.NullString{Valid: false}
	t.TendermintTransactionTimestamp = sql.NullTime{Valid: false}
	if dbutil.Db.GetConn().NewRecord(t) {
		result := dbutil.Db.GetConn().Create(&t)
		rowsAffected := result.RowsAffected
		errors := result.GetErrors()
		if len(errors) > 0 {
			logger.Errorf("errors while creating new entry %v\n", errors)
			return false
		}
		return rowsAffected > 0
	}

	return false
}
