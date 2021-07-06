/*
Package reputation - location for reputation record structs.
*/
package reputation

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

// Record represents an event that will affect peer's reputation.
type Record struct {
	reason    string
	point     int64
	violation bool
}

// Reason gets the reason of this record.
func (r *Record) Reason() string {
	return r.reason
}

// Point gets the point update of this record.
func (r *Record) Point() int64 {
	return r.point
}

// Violation gets whether this record is a violation.
func (r *Record) Violation() bool {
	return r.violation
}

// Copy gets a copy of this record.
func (r *Record) Copy() *Record {
	return &Record{reason: r.reason, point: r.point, violation: r.violation}
}

// A list of global variables representing list records.
var MockGoodRecord = Record{
	reason:    "Mock good record",
	point:     1000,
	violation: false,
}

var MockBadRecord = Record{
	reason:    "Mock bad record",
	point:     -1,
	violation: true,
}

var StandardOfferRetrieved = Record{
	reason:    "Retrieved one offer from standard discovery",
	point:     5,
	violation: false,
}

var DHTOfferRetrieved = Record{
	reason:    "Retrived one offer from dht discovery",
	point:     10,
	violation: false,
}

var ContentRetrieved = Record{
	reason:    "Retrived a content from a pre-signed offer",
	point:     10,
	violation: false,
}

var NetworkError = Record{
	reason:    "Network error",
	point:     -1,
	violation: false,
}

var InvalidResponse = Record{
	reason:    "Received an invalid response",
	point:     -2,
	violation: false,
}

var NetworkErrorAfterPayment = Record{
	reason:    "Network error after a payment is made",
	point:     -10,
	violation: true,
}

var InvalidResponseAfterPayment = Record{
	reason:    "Received an invalid response after a payment is made",
	point:     -50,
	violation: true,
}

var InvalidRefund = Record{
	reason:    "Received an invalid refund",
	point:     -50,
	violation: true,
}
