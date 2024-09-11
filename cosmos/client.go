package cosmos

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typetx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/google"
)

func NewGrpcClient() (grpcClient *grpc.ClientConn, err error) {
	// 创建grpc连接
	target := "cosmos-grpc.publicnode.com:443"
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithCredentialsBundle(google.NewDefaultCredentials()))
	grpcConn, err := grpc.NewClient(target, opts...)
	if err != nil {
		panic(err)
	}
	txClient := typetx.NewServiceClient(grpcConn)

	req := &typetx.GetBlockWithTxsRequest{Height: 20758858}
	resp, err := txClient.GetBlockWithTxs(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Txs[0].Body)
	return grpcConn, nil
}

func transaction() {
	// 构造交易
	pk, _ := Import("you-prikeyHex")
	var privKeyAccAddr sdk.AccAddress = pk.PubKey().Address().Bytes()

	toAddr := sdk.MustAccAddressFromBech32("cosmos1f4wkdv7qdj9wrgqdgtnw50nt9a36663x93229g")
	coin := sdk.NewInt64Coin("uatom", 100000)
	amount := sdk.NewCoins(coin)
	msg := banktypes.NewMsgSend(privKeyAccAddr, toAddr, amount)
	msgs := []sdk.Msg{msg}
	fmt.Println(msg)

	grpcConn, _ := NewGrpcClient()

	authClient := authtypes.NewQueryClient(grpcConn)
	fromAddrInfoAny, err := authClient.Account(context.Background(), &authtypes.QueryAccountRequest{Address: privKeyAccAddr.String()})
	if err != nil {
		fmt.Println(err)
	}
	// 尝试将Account字段中的Any转换为BaseAccount类型
	var f authtypes.BaseAccount
	if err := proto.Unmarshal(fromAddrInfoAny.Account.Value, &f); err != nil {
		panic(err)
	}

	fee := sdk.NewInt64Coin("uatom", 2500)
	gasLimit := int64(100000)
	// 签名
	txRaw, err := BuildTxV2("cosmoshub-4", f.Sequence,
		f.AccountNumber, pk, fee, gasLimit, msgs)
	if err != nil {
		panic(err)
	}
	txBytes, err := proto.Marshal(txRaw)
	if err != nil {
		panic(err)
	}

	// 广播
	txClient := typetx.NewServiceClient(grpcConn)
	req := &typetx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    typetx.BroadcastMode_BROADCAST_MODE_SYNC,
	}
	txResp, err := txClient.BroadcastTx(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(txResp.TxResponse)
}

func BuildTxV2(chainId string, sequence, accountNumber uint64, privKey *secp256k1.PrivKey, fee sdk.Coin, gaslimit int64, msgs []sdk.Msg) (*typetx.TxRaw, error) {
	txBodyMessage := make([]*types.Any, 0)
	for i := 0; i < len(msgs); i++ {
		msgAnyValue, err := types.NewAnyWithValue(msgs[i])
		if err != nil {
			return nil, err
		}
		txBodyMessage = append(txBodyMessage, msgAnyValue)
	}
	txBody := &typetx.TxBody{
		Messages:                    txBodyMessage,
		Memo:                        "",
		TimeoutHeight:               0,
		ExtensionOptions:            nil,
		NonCriticalExtensionOptions: nil,
	}
	txBodyBytes, err := proto.Marshal(txBody)
	if err != nil {
		return nil, err
	}
	pubAny, err := types.NewAnyWithValue(privKey.PubKey())
	if err != nil {
		return nil, err
	}
	authInfo := &typetx.AuthInfo{
		SignerInfos: []*typetx.SignerInfo{
			{
				PublicKey: pubAny,
				ModeInfo: &typetx.ModeInfo{
					Sum: &typetx.ModeInfo_Single_{
						Single: &typetx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
					},
				},
				Sequence: sequence,
			},
		},
		Fee: &typetx.Fee{
			Amount:   sdk.NewCoins(fee),
			GasLimit: uint64(gaslimit),
			Payer:    "",
			Granter:  "",
		},
	}

	txAuthInfoBytes, err := proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	signDoc := &typetx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       chainId,
		AccountNumber: accountNumber,
	}
	signatures, err := proto.Marshal(signDoc)
	if err != nil {
		return nil, err
	}
	sign, err := privKey.Sign(signatures)
	if err != nil {
		return nil, err
	}
	return &typetx.TxRaw{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	}, nil
}
