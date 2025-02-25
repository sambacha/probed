// Copyright © 2021, 2022 Weald Technology Trading.
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

// Package rest provides a REST implementation of the probe daemon.
package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wealdtech/probed/services/daemon/rest/types"
	"github.com/wealdtech/probed/services/probedb"
)

func (s *Service) postAggregateAttestation(w http.ResponseWriter, r *http.Request) {
	var aggregateAttestation types.AggregateAttestation
	if err := json.NewDecoder(r.Body).Decode(&aggregateAttestation); err != nil {
		log.Debug().Err(err).Msg("Supplied with invalid data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sourceIP, err := sourceIP(r)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to obtain source IP")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.aggregateAttestationsSetter.SetAggregateAttestation(context.Background(), &probedb.AggregateAttestation{
		IPAddr:          sourceIP,
		Source:          aggregateAttestation.Source,
		Method:          aggregateAttestation.Method,
		Slot:            aggregateAttestation.Slot,
		CommitteeIndex:  aggregateAttestation.CommitteeIndex,
		AggregationBits: aggregateAttestation.AggregationBits,
		BeaconBlockRoot: aggregateAttestation.BeaconBlockRoot,
		SourceRoot:      aggregateAttestation.SourceRoot,
		TargetRoot:      aggregateAttestation.TargetRoot,
		DelayMS:         aggregateAttestation.DelayMS,
	}); err != nil {
		log.Warn().Err(err).Msg("Failed to set aggregate attestation")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//	log.Trace().
	//		Str("ip_addr", sourceIP.String()).
	//		Str("source", blockDelay.Source).
	//		Str("method", blockDelay.Method).
	//		Uint32("slot", blockDelay.Slot).
	//		Uint32("delay_ms", blockDelay.DelayMS).
	//		Msg("Metric accepted")
	w.WriteHeader(http.StatusCreated)
}
