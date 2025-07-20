package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/willtrojniak/TabAppBackend/env"
	"github.com/willtrojniak/TabAppBackend/util"
)

type Token string

func (t Token) String() string {
	return string(t)
}

func (t *Token) ScanText(v pgtype.Text) error {
	if !v.Valid {
		*t = ""
		return nil
	}

	token, err := util.Decrypt(v.String, []byte(env.Envs.ENCRYPT_SECRET))
	if err != nil {
		return err
	}

	*t = Token(string(token))

	return nil
}

func (t Token) TextValue() (pgtype.Text, error) {
	v := pgtype.Text{}
	cipher, err := util.Encrypt([]byte(t), []byte(env.Envs.ENCRYPT_SECRET))
	if err != nil {
		return v, err
	}
	v.String = cipher
	v.Valid = true

	return v, nil
}
