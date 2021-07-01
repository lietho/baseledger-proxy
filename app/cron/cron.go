package cron

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"

	businesslogic "github.com/unibrightio/baseledger/app/business_logic"
	"github.com/unibrightio/baseledger/app/types"
	proxytypes "github.com/unibrightio/baseledger/x/proxy/types"

	common "github.com/unibrightio/baseledger/common"
	"github.com/unibrightio/baseledger/dbutil"
	"github.com/unibrightio/baseledger/logger"
)

func queryTrustmeshes() {
	logger.Info("query trustmesh start")

	var trustmeshEntries []proxytypes.TrustmeshEntry
	dbutil.Db.GetConn().Where("commitment_state=?", common.UncommittedCommitmentState).Find(&trustmeshEntries)

	logger.Infof("found %v trustmesh entries\n", len(trustmeshEntries))
	var jobs = make(chan types.Job, len(trustmeshEntries))
	var results = make(chan types.Result, len(trustmeshEntries))
	createWorkerPool(1, jobs, results)

	for _, trustmeshEntry := range trustmeshEntries {
		logger.Infof("creating job for %v\n", trustmeshEntry.TransactionHash)
		job := types.Job{TrustmeshEntry: trustmeshEntry}
		jobs <- job
	}
	close(jobs)

	for result := range results {
		logger.Infof("Tx hash %v, height %v, timestamp %v\n", result.Job.TrustmeshEntry.TransactionHash, result.TxInfo.TxHeight, result.TxInfo.TxTimestamp)
	}

	logger.Info("query trustmesh end")
}

func getTxInfo(txHash string) (txInfo *types.TxInfo, err error) {
	// fetching tx details
	// TODO: BAS-33
	str := "http://localhost:26657/tx?hash=0x" + txHash
	httpRes, err := http.Get(str)
	if err != nil {
		logger.Errorf("error during http tx req %v\n", err)
		return &types.TxInfo{}, errors.New("get tx info error")
	}

	// if transaction is not found it is not yet committed
	if httpRes.StatusCode != 200 {
		logger.Error("tx not committed yet")
		return &types.TxInfo{TxCommitted: false}, errors.New("get tx info error")
	}

	// if it's found should be committed at this point, decode
	var committedTx types.TxResp
	err = json.NewDecoder(httpRes.Body).Decode(&committedTx)
	if err != nil {
		logger.Errorf("error decoding tx %v\n", err)
		return &types.TxInfo{}, errors.New("error decoding tx")
	}
	// query for block at specific height to find timestamp
	// TODO: BAS-33
	str = "http://localhost:26657/block?height" + committedTx.TxResult.Height
	httpRes, err = http.Get(str)
	if err != nil {
		logger.Errorf("error during http block req %v\n", err)
		return &types.TxInfo{}, errors.New("get blcok info error")
	}
	var commitedBlock types.BlockResp
	err = json.NewDecoder(httpRes.Body).Decode(&commitedBlock)
	if err != nil {
		logger.Errorf("error decoding block %v\n", err)
		return &types.TxInfo{}, errors.New("error decoding block")
	}
	logger.Infof("DECODED COMMITTED TX HEIGHT %v AND TIMESTAMP %v\n", committedTx.TxResult.Height, commitedBlock.BlockResult.Block.Header.Time)
	return &types.TxInfo{
		TxHeight:    committedTx.TxResult.Height,
		TxTimestamp: commitedBlock.BlockResult.Block.Header.Time,
		TxCommitted: true,
	}, nil
}

func worker(jobs chan types.Job, results chan types.Result) {
	defer close(results)
	for job := range jobs {
		txInfo, err := getTxInfo(job.TrustmeshEntry.TransactionHash)
		if err != nil {
			// here it would be http error
			// it seems that we can just let it go through result channel
			logger.Warnf("result error %v\n", err)
		}
		logger.Infof("result tx %v transaction type %v\n", txInfo, job.TrustmeshEntry.BaseledgerTransactionType)
		output := types.Result{Job: job, TxInfo: *txInfo}
		businesslogic.ExecuteBusinessLogic(output)
		results <- output
	}
}

func createWorkerPool(noOfWorkers int, jobs chan types.Job, results chan types.Result) {
	for i := 0; i < noOfWorkers; i++ {
		go worker(jobs, results)
	}
}

func StartCron() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Seconds().SingletonMode().Do(queryTrustmeshes)

	s.StartAsync()
}
