package domain

import "github.com/vyolayer/vyolayer/internal/shared/auth"

type Password struct {
	Hash string
}

func NewPassword(p string) *Password {
	hash, _ := auth.GenerateHash(p)

	return &Password{
		Hash: hash,
	}
}

func ReconstructPassword(hash string) *Password {
	return &Password{
		Hash: hash,
	}
}

func (p *Password) VerifyPassword(s string) bool {
	return auth.CheckHash(s, p.Hash)
}

func (p *Password) IsSamePassword(s string) bool {
	return auth.CheckHash(s, p.Hash)
}

func (p *Password) ChangePassword(s string) error {
	newHash, e := auth.GenerateHash(s)
	if e != nil {
		return e
	}
	p.Hash = newHash
	return nil
}

func (p *Password) String() string {
	return p.Hash
}
