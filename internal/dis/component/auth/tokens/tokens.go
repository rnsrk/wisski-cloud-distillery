package tokens

import (
	"context"
	"crypto/rand"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/password"
	"github.com/tkw1536/pkglib/reflectx"
	"gorm.io/gorm"
)

// Tokens implements Tokens
type Tokens struct {
	component.Base

	Dependencies struct {
		SQL *sql.SQL
	}
}

var (
	_ component.UserDeleteHook = (*Tokens)(nil)
	_ component.Table          = (*Tokens)(nil)
)

func (tok *Tokens) TableInfo() component.TableInfo {
	return component.TableInfo{
		Name:  models.TokensTable,
		Model: reflectx.MakeType[models.Token](),
	}
}

func (tok *Tokens) table(ctx context.Context) (*gorm.DB, error) {
	return tok.Dependencies.SQL.QueryTable(ctx, tok)
}

func (tok *Tokens) OnUserDelete(ctx context.Context, user *models.User) error {
	table, err := tok.table(ctx)
	if err != nil {
		return err
	}
	return table.Delete(&models.Token{}, &models.Token{User: user.User}).Error
}

// Tokens returns a list of tokens for the given user
func (tok *Tokens) Tokens(ctx context.Context, user string) ([]models.Token, error) {
	// the empty user has no tokens
	if user == "" {
		return nil, nil
	}

	// get the table
	table, err := tok.table(ctx)
	if err != nil {
		return nil, err
	}

	var tokens []models.Token

	// make a query to find all keys (in the underlying model)
	query := table.Find(&tokens, &models.Token{User: user})
	if query.Error != nil {
		return nil, query.Error
	}

	return tokens, nil
}

const (
	tokenGroupLength                  = 8
	tokenGroupCount                   = 8
	tokenSeparator                    = "-"
	tokenCharset     password.Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// NewToken generates a new token
func NewToken() (string, error) {
	// generate a new password
	token, err := password.Generate(rand.Reader, tokenGroupCount*tokenGroupLength, tokenCharset)
	if err != nil {
		return "", err
	}

	// insert the token group separators
	var result strings.Builder
	result.Grow(len(token) + (tokenGroupCount-1)*len(tokenSeparator))

	for i := 0; i < tokenGroupCount; i++ {
		if i != 0 {
			result.WriteString(tokenSeparator)
		}

		start := i * tokenGroupLength
		result.WriteString(token[start : start+tokenGroupLength])
	}

	return result.String(), nil

}

// Add adds a new token, unless it already exists.
// The token is granted scopes with .SetScopes(scopes).
func (tok *Tokens) Add(ctx context.Context, user string, description string, scopes []string) (*models.Token, error) {

	// create a new token and set the scopes
	mk := models.Token{
		User:        user,
		Description: description,
	}
	mk.SetScopes(scopes)

	// generate a new random password
	var err error
	mk.Token, err = NewToken()
	if err != nil {
		return nil, err
	}

	// get the table
	table, err := tok.table(ctx)
	if err != nil {
		return nil, err
	}

	// create the token instance
	if err := table.Create(&mk).Error; err != nil {
		return nil, err
	}

	// and return
	return &mk, nil
}

// Remove removes a token with the given token from the user
func (tok *Tokens) Remove(ctx context.Context, user, token string) error {
	// get the table
	table, err := tok.table(ctx)
	if err != nil {
		return err
	}

	// and do the delete
	return table.Where("user = ? AND token = ?", user, token).Delete(&models.Token{}).Error
}
