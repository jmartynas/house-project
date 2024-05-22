// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Breziny struct {
	ID           int32
	Brezinys     string
	FkSutartisID int16
}

type BrezinysIDSeq struct {
	LastValue pgtype.Int4
	LogCnt    pgtype.Int4
	IsCalled  pgtype.Bool
}

type Leidima struct {
	ID           int64
	Leidimas     string
	FkBrezinysID int32
}

type LeidimasIDSeq struct {
	LastValue pgtype.Int8
	LogCnt    pgtype.Int4
	IsCalled  pgtype.Bool
}

type Sutarti struct {
	ID       int16
	Sutartis string
	Kaina    string
}

type SutartisIDSeq struct {
	LastValue pgtype.Int2
	LogCnt    pgtype.Int4
	IsCalled  pgtype.Bool
}
