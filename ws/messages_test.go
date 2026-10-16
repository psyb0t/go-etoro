package etorows_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/psyb0t/ctxerrors/commerr"
	"github.com/psyb0t/gen-asyncapi/asyncapirt"
	etorows "github.com/psyb0t/go-etoro/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	pushID       = "3f6c2d8e-0b7a-4c1e-9d2f-5a8b7c6d4e3f"
	pushTopic    = "instrument:1"
	subscribe    = "Subscribe"
	rateType     = "Trading.Instrument.Rate"
	candleType   = "Trading.Instrument.Candle"
	statusType   = "Trading.Instrument.MarketStatus"
	openedType   = "Trading.Position.Opened"
	closedType   = "Trading.Position.Closed"
	orderType    = "Trading.OrderForOpen.Update"
	closeType    = "Trading.OrderForCloseMultiple.Update"
	unlistedType = "Trading.Unlisted.Type"
)

func pushFrame(t *testing.T, pushType string, content any) []byte {
	t.Helper()

	encoded, err := json.Marshal(content)
	require.NoError(t, err)

	frame, err := json.Marshal(map[string]any{
		"messages": []map[string]any{{
			"topic":   pushTopic,
			"type":    pushType,
			"id":      pushID,
			"content": string(encoded),
		}},
	})
	require.NoError(t, err)

	return frame
}

func TestDecodeReceivedMessagePicksThePushType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		pushType string
		wantName string
	}{
		{rateType, etorows.MessageInstrumentRateMessage},
		{candleType, etorows.MessageInstrumentCandleMessage},
		{statusType, etorows.MessageInstrumentMarketStatusMessage},
		{openedType, etorows.MessagePositionOpenedMessage},
		{closedType, etorows.MessagePositionClosedMessage},
		{orderType, etorows.MessageOrderUpdateMessage},
		{closeType, etorows.MessageOrderUpdateMessage},
	}

	for _, tc := range testCases {
		t.Run(tc.pushType, func(t *testing.T) {
			t.Parallel()

			message, err := etorows.DecodeReceivedMessage(
				pushFrame(t, tc.pushType, map[string]any{}),
			)
			require.NoError(t, err)
			assert.Equal(t, tc.wantName, message.Name)
		})
	}
}

func TestDecodeReceivedMessageRejectsAnUnlistedPushType(t *testing.T) {
	t.Parallel()

	_, err := etorows.DecodeReceivedMessage(
		pushFrame(t, unlistedType, map[string]any{}),
	)
	require.ErrorIs(t, err, asyncapirt.ErrUnknownMessage)
}

func TestInstrumentRateContentKeepsWireStrings(t *testing.T) {
	t.Parallel()

	frame := pushFrame(t, rateType, map[string]any{
		"Ask":         "1.25",
		"Bid":         "1.2",
		"Date":        "2026-01-02T03:04:05Z",
		"PriceRateID": "1",
	})

	message, err := etorows.DecodeReceivedMessage(frame)
	require.NoError(t, err)

	envelope, ok := message.Payload.(etorows.InstrumentRateEnvelope)
	require.True(t, ok)
	require.Len(t, envelope.Messages, 1)

	rate, err := envelope.Messages[0].ContentValue()
	require.NoError(t, err)
	require.NotNil(t, rate.Ask)
	assert.Equal(t, "1.25", *rate.Ask)
	require.NotNil(t, rate.Bid)
	assert.Equal(t, "1.2", *rate.Bid)
}

func TestDecodeReceivedMessageOperationResponse(t *testing.T) {
	t.Parallel()

	frame, err := json.Marshal(map[string]any{
		"id":           pushID,
		"success":      false,
		"operation":    subscribe,
		"errorMessage": "Key is invalid",
	})
	require.NoError(t, err)

	message, err := etorows.DecodeReceivedMessage(frame)
	require.NoError(t, err)
	assert.Equal(t, etorows.MessageOperationResponse, message.Name)

	response, ok := message.Payload.(etorows.OperationResponse)
	require.True(t, ok)
	assert.False(t, response.Success)
	assert.Equal(
		t,
		etorows.OperationResponseOperationSubscribe,
		response.Operation,
	)
}

func TestDecodeReceivedMessageRejectsUnknownFrames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		frame   string
		wantErr error
	}{
		{
			name:    "not JSON",
			frame:   "{",
			wantErr: commerr.ErrCouldNotDecode,
		},
		{
			name:    "no known shape",
			frame:   `{"hello":"world"}`,
			wantErr: asyncapirt.ErrUnknownMessage,
		},
		{
			name:    "envelope item missing fields",
			frame:   `{"messages":[{}]}`,
			wantErr: asyncapirt.ErrUnknownMessage,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := etorows.DecodeReceivedMessage([]byte(tc.frame))
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

type recordingTransport struct {
	sent [][]byte
}

func (r *recordingTransport) WriteMessage(
	_ context.Context,
	data []byte,
) error {
	r.sent = append(r.sent, data)

	return nil
}

func (r *recordingTransport) ReadMessage(context.Context) ([]byte, error) {
	return nil, context.Canceled
}

func (r *recordingTransport) Close() error {
	return nil
}

func TestClientSubscribeWritesTheEnvelope(t *testing.T) {
	t.Parallel()

	transport := &recordingTransport{}
	client := etorows.NewClient(transport)

	err := client.Subscribe(context.Background(), etorows.SubscribeEnvelope{
		ID:        pushID,
		Operation: subscribe,
		Data:      etorows.SubscribeData{Topics: []string{pushTopic}},
	})
	require.NoError(t, err)
	require.Len(t, transport.sent, 1)

	var sent struct {
		Operation string `json:"operation"`
		Data      struct {
			Topics []string `json:"topics"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(transport.sent[0], &sent))
	assert.Equal(t, subscribe, sent.Operation)
	assert.Equal(t, []string{pushTopic}, sent.Data.Topics)
}
