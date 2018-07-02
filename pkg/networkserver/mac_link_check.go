// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package networkserver

import (
	"context"

	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

func handleLinkCheckReq(ctx context.Context, dev *ttnpb.EndDevice, msg *ttnpb.UplinkMessage) error {
	floor, ok := demodulationFloor[msg.Settings.SpreadingFactor][msg.Settings.Bandwidth]
	if !ok {
		return ErrInvalidDataRate.New(nil)
	}

	gtws := make(map[string]struct{}, len(msg.RxMetadata))

	maxSNR := msg.RxMetadata[0].SNR
	for _, md := range msg.RxMetadata {
		gtws[md.GatewayIdentifiers.UniqueID(ctx)] = struct{}{}

		snr := md.SNR
		if snr <= maxSNR {
			continue
		}
		maxSNR = snr
	}

	dev.QueuedMACCommands = append(
		dev.QueuedMACCommands,
		(&ttnpb.MACCommand_LinkCheckAns{
			Margin:       uint32(uint8(maxSNR - floor)),
			GatewayCount: uint32(len(gtws)),
		}).MACCommand(),
	)
	return nil
}