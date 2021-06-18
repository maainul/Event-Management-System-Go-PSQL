package postgres

import (
	"Event-Management-System-Go-PSQL/storage"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStorage_GetUser(t *testing.T) {
	dbString := newDBFromConfig()
	store, err := NewStorage(dbString)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		storage *Storage
		in      int32
		want    *storage.User
		wantErr bool
	}{
		{
			name:    "GET_USER_SUCCESS",
			storage: store,
			in:      1,
			want: &storage.User{
				ID:        1,
				FirstName: "Rahim",
				LastName:  "Karim",
				Username:  "rahim",
				Email:     "rahim@gmail.com",
				IsActive:  true,
				IsAdmin:   true,
			},
			wantErr: false,
		},
		{
			name:    "GET_USER_FAILED",
			storage: store,
			in:      2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				db: tt.storage.db,
			}
			got, err := s.GetUser(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tOps := []cmp.Option{
				cmpopts.IgnoreFields(storage.User{}, "Password", "CreatedAt", "UpdatedAt"),
			}

			if !tt.wantErr && !cmp.Equal(got, tt.want, tOps...) {
				t.Error(cmp.Diff(got, tt.want, tOps...))
			}

			if tt.wantErr && err == nil {
				t.Error("want error but got user")
			}
		})
	}
}

func newDBFromConfig() string {
	dbParams := " " + "user=postgres"
	dbParams += " " + "host=localhost"
	dbParams += " " + "port=5432"
	dbParams += " " + "dbname=practice"
	dbParams += " " + "password=0"
	dbParams += " " + "sslmode=disable"

	return dbParams
}
