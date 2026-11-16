// SPDX-FileCopyrightText: 2025 DENIC eG and OpenPAQ contributors
// SPDX-License-Identifier: AGPL-3.0-only

package types

type DatabaseConfig struct {
	DbUserName     string
	DbUserPassword string
	DbHost         string
	DbPort         string
	DataBase       string
	Table          string
}
