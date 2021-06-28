package types

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres

	uuid "github.com/kthomas/go.uuid"
	"github.com/unibrightio/baseledger/dbutil"
)

//Put here our Types needed for the proxy elements?
type OffchainProcessMessageReferenceType string

type OffchainProcessMessage struct {
	Id                                   uuid.UUID
	SenderId                             uuid.UUID
	ReceiverId                           uuid.UUID
	Topic                                string
	ReferencedOffchainProcessMessageId   uuid.UUID
	BaseledgerSyncTreeJson               string
	WorkstepType                         string
	BusinessObjectProof                  string
	BusinessObjectType                   string
	TendermintTransactionIdOfStoredProof uuid.UUID
	BaseledgerTransactionIdOfStoredProof uuid.UUID
	BaseledgerBusinessObjectId           uuid.UUID
	ReferencedBaseledgerBusinessObjectId uuid.UUID
	StatusTextMessage                    string
	BaseledgerTransactionType            string
	ReferencedBaseledgerTransactionId    uuid.UUID
	EntryType                            string
}

// TODO rename after clean up
type SynchronizationRequest struct {
	WorkgroupId                          uuid.UUID
	Recipient                            string
	WorkstepType                         string
	BusinessObjectType                   string
	BaseledgerBusinessObjectId           string
	BusinessObjectJson                   string
	ReferencedBaseledgerBusinessObjectId string
	ReferencedBaseledgerTransactionId    string
}

type SynchronizationFeedback struct {
	WorkgroupId                                uuid.UUID
	BaseledgerProvenBusinessObjectJson         string
	BusinessObjectType                         string
	Recipient                                  string
	Approved                                   bool
	BaseledgerBusinessObjectIdOfApprovedObject string
	HashOfObjectToApprove                      string
	OriginalBaseledgerTransactionId            string
	OriginalOffchainProcessMessageId           string
	FeedbackMessage                            string
}

type BaseledgerTransactionPayload struct {
	SenderId                             string
	TransactionType                      string
	OffchainMessageId                    string
	ReferencedOffchainMessageId          string
	ReferencedBaseledgerTransactionId    string
	BaseledgerTransactionId              string
	Proof                                string
	BaseledgerBusinessObjectId           string
	ReferencedBaseledgerBusinessObjectId string
}

func (o *OffchainProcessMessage) Create() bool {
	db := dbutil.Db.GetConn()
	if db.NewRecord(o) {
		result := db.Create(&o)
		rowsAffected := result.RowsAffected
		errors := result.GetErrors()
		if len(errors) > 0 {
			fmt.Printf("errors while creating new offchain process msg entry %v\n", errors)
			return false
		}
		return rowsAffected > 0
	}

	return false
}

func GetOffchainMsgById(id uuid.UUID) (msg *OffchainProcessMessage, err error) {
	db := dbutil.Db.GetConn()
	var offchainMsg OffchainProcessMessage
	res := db.First(&offchainMsg, "id = ?", id.String())

	if res.Error != nil {
		fmt.Printf("error when getting offchain msg from db %v\n", err)
		return nil, res.Error
	}

	return &offchainMsg, nil
}

// all other types for hasing, privacy, off-chain messaging
