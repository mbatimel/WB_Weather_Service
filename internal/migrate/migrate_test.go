package migrate

import ("testing"
"github.com/mbatimel/WB_Weather_Service/internal/repo"
)

func TestApplyMigrations(t *testing.T) {
	db, err := repo.SetConfigs("../../config/config.yaml")
	want :=error(nil)
    if err != nil {
		t.Errorf("got %q, wanted %q", err, want)
    }
    defer db.Close()

    err = db.ConnectToDataBase()
    if err != nil {
		t.Errorf("got %q, wanted %q", err, want)
    }

    // Example migration application
    err = ApplyMigrations(db, "../../migrations/migrate.sql")
    if err != nil {
		t.Errorf("got %q, wanted %q", err, want)
    }
}