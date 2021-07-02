/*
Package fcrdatabase - handles the database inside the system.
*/
package fcrdatabase

import "errors"

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

// FCRDatabaseImplV1 implements the FCRDatabase interface, it is a placeholder.
type FCRDatabaseImplV1 struct {
	start bool
}

func NewFCRDatabaseImplV1() FCRDatabase {
	return &FCRDatabaseImplV1{start: false}
}

func (db *FCRDatabaseImplV1) Start() error {
	if db.start {
		return errors.New("FCRDatabase has already started")
	}
	db.start = true
	return nil
}

func (db *FCRDatabaseImplV1) Shutdown() {
	if !db.start {
		return
	}
	db.start = false
}
