package rpchandlers

import (
	"github.com/c4ei/c4exd/app/appmessage"
	"github.com/c4ei/c4exd/app/rpc/rpccontext"
	"github.com/c4ei/c4exd/domain/consensus/model/externalapi"
	"github.com/c4ei/c4exd/domain/consensus/utils/transactionid"
	"github.com/c4ei/c4exd/infrastructure/network/netadapter/router"
)

// HandleGetMempoolEntry handles the respectively named RPC command
func HandleGetMempoolEntry(context *rpccontext.Context, _ *router.Router, request appmessage.Message) (appmessage.Message, error) {

	transaction := &externalapi.DomainTransaction{}
	var found bool
	var isOrphan bool

	getMempoolEntryRequest := request.(*appmessage.GetMempoolEntryRequestMessage)

	transactionID, err := transactionid.FromString(getMempoolEntryRequest.TxID)
	if err != nil {
		errorMessage := &appmessage.GetMempoolEntryResponseMessage{}
		errorMessage.Error = appmessage.RPCErrorf("Transaction ID could not be parsed: %s", err)
		return errorMessage, nil
	}

	mempoolTransaction, isOrphan, found := context.Domain.MiningManager().GetTransaction(transactionID, !getMempoolEntryRequest.FilterTransactionPool, getMempoolEntryRequest.IncludeOrphanPool)

	if !found {
		errorMessage := &appmessage.GetMempoolEntryResponseMessage{}
		errorMessage.Error = appmessage.RPCErrorf("Transaction %s was not found", transactionID)
		return errorMessage, nil
	}

	rpcTransaction := appmessage.DomainTransactionToRPCTransaction(mempoolTransaction)
	err = context.PopulateTransactionWithVerboseData(rpcTransaction, nil)
	if err != nil {
		return nil, err
	}
	return appmessage.NewGetMempoolEntryResponseMessage(transaction.Fee, rpcTransaction, isOrphan), nil
}
