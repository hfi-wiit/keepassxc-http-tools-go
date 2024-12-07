package keepassxc

import (
	"encoding/json"
	"keepassxc-http-tools-go/pkg/utils"
	"strings"

	"github.com/kevinburke/nacl"
)

/*
	entry representation
*/

// Password is a string type that prevents accidental prints.
type Password string

// String stringifies a Password to only asterisks.
func (p *Password) String() string {
	return "*****"
}

// Plaintext stringifies a Password to its actual value.
func (p *Password) Plaintext() string {
	return string(*p)
}

// StringFields represents the user defined additional fields as returned by the api.
// That is a list of maps with one element each, each key has to start with "KPH: " to be returned by the api.
type StringFields []map[string]Password

// ToMap converts the structure returned from the api (list of single entry maps)
// to a simple key value map (keys without the leading "KPH: ").
// The map uses Password values, since values could contain sensible data.
func (f StringFields) ToMap() map[string]Password {
	fMap := make(map[string]Password, len(f))
	for _, e := range f {
		for k, v := range e {
			fMap[strings.TrimSpace(strings.TrimPrefix(k, utils.StringFieldKeyPrefix))] = v
		}
	}
	return fMap
}

// String stringifies the contents of the StringFields as a json representation of the map returned by ToMap().
func (f StringFields) String() string {
	v, _ := json.Marshal(f.ToMap())
	return string(v)
}

// Entry represents a single password entry as returned by the keepassxc http api.
// Example entry as returned from the api, if every field has some value (incl. tags, expire date and totp):
// [{"group":"foo","login":"myname","name":"bar","password":"myPa$$w0rd","stringFields":[{"KPH: bar":"barval"},{"KPH: foo":"fooval"}],
// "totp":"175413","uuid":"92bfee4f24614ef9ac6e1f440eff3292"}]
// If expired, the entry will actually not be found/returned by the api.
type Entry struct {
	// The name/identifier of the password entry.
	Name string `json:"name"`
	// The user name of the password entry.
	Login string `json:"login"`
	// The password of the password entry.
	Password Password `json:"password"`
	// The current generated totp of the password entry, if a totp is set up.
	Totp string `json:"totp"`
	// The group/folder of the password entry inside the database.
	Group string `json:"group"`
	// The UUID of the password entry.
	Uuid string `json:"uuid"`
	// The user defined additional fields of the password entry.
	// See ToMap() for an actual usable representation.
	StringFields StringFields `json:"stringFields"`
}

// StringFieldsMap converts the StringFields structure returned from the api (list of single entry maps)
// to a simple key value map (keys without the leading "KPH: ").
// See also StringFields.ToMap().
func (e Entry) StringFieldsMap() map[string]Password {
	return e.StringFields.ToMap()
}

// Entries represents a list of Entry objects.
type Entries []*Entry

// TODO required?
func (e Entries) Names() []string {
	names := make([]string, len(e))
	for i, v := range e {
		names[i] = v.Name
	}
	return names
}

// TODO required?
func (e Entries) FilterByName(name string) Entries {
	newEntries := make(Entries, len(e))
	count := 0
	for _, entry := range e {
		if strings.Contains(entry.Name, name) {
			newEntries[count] = entry
			count += 1
		}
	}
	return newEntries[:count]
}

/*
	client helper structs
*/

// Implement this interface to provide the association data for the keepassxc http api connection.
type KeepassxcClientProfile interface {
	// GetAssocName returns the name of the profile.
	GetAssocName() string
	// GetAssocKey returns the nacl.Key of the profile.
	// If not yet associated, this is supposed to return nil.
	GetAssocKey() nacl.Key
	// SetAssoc saves the assoctiation name and nacl.Key in the profile.
	// If GetAssocKey() returns nil, a key will be generated and passed to this function.
	SetAssoc(string, nacl.Key) error
}

// Message represents the api input from the client.
type Message map[string]interface{}

// Response represents the api response to the client.
type Response map[string]interface{}

// entries tries to parse the entries from an api response.
func (r Response) entries() (Entries, error) {
	var data []byte
	if msg, ok := r["message"]; ok {
		if v, ok := msg.(map[string]interface{})["entries"]; ok {
			var err error
			data, err = json.Marshal(v)
			if err != nil {
				return nil, err
			}
		}
	}
	if len(data) == 0 {
		return nil, utils.ErrKeepassxcInvalidResponse
	}

	var entries Entries
	err := json.Unmarshal(data, &entries)
	return entries, err
}
