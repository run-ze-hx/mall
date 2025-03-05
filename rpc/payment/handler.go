package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/google/uuid"
	payment "mall/kitex_gen/payment"
	"mall/model"
	"mall/rpc/payment/paymentdao/mysql"
	"time"
)

// PaymentServiceImpl implements the last service interface defined in the IDL.
type PaymentServiceImpl struct{}

// Charge implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) Charge(ctx context.Context, req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {

	var TransactionId uuid.UUID
	TransactionId, err = uuid.NewRandom()
	if err != nil {
		return nil, kerrors.NewBizStatusError(400501, err.Error())
	}

	err = mysql.CreatePayment(ctx, &model.Payment{
		UserId:        req.UserId,
		OrderId:       req.OrderId,
		TransactionId: TransactionId.String(),
		Amount:        req.Amount,
		PayAt:         time.Now(),
	})
	if err != nil {
		return nil, kerrors.NewBizStatusError(400502, err.Error())
	}

	return &payment.ChargeResp{TransactionId: TransactionId.String()}, nil
}
