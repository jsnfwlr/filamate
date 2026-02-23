package db_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/etc/db/migrations"
	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/test/containers"
	. "github.com/jsnfwlr/filamate/internal/types"
	"github.com/jsnfwlr/filamate/static"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	colorCount    = 32
	brandCount    = 11
	locationCount = 13
	materialCount = 15
	storeCount    = 29
	spoolCount    = 59
)

func TestDatabase(t *testing.T) {
	gCfg, err := go11y.LoadConfig()
	if err != nil {
		t.Fatalf("could not load go11y config: %v", err)
	}

	ctx, _, err := go11y.Initialise(context.Background(), gCfg, os.Stdout)
	if err != nil {
		t.Fatalf("could not initialise go11y: %v", err)
	}

	_, err = db.Connect(ctx, db.Config{Database: "invalid?description_cache_capacity=three", Port: "1234"})
	if err == nil {
		t.Fatal("expected error when connecting with invalid config")
	}

	ctr, cfg, err := containers.Postgres(t, ctx, "db_version", "17", PointerOf("full-db-test"))
	if err != nil {
		t.Fatalf("could not start the Postgres container: %v", err)
	}

	t.Cleanup(
		func() {
			ctr.Cleanup(t)
		},
	)

	dbClient, err := db.Connect(ctx, cfg)
	if err != nil {
		t.Fatalf("could not connect to the database: %v", err)
	}

	t.Cleanup(
		func() {
			dbClient.Close()
		},
	)

	testBeforeMigrations(t, ctx, dbClient)

	testMigrations(t, ctx, cfg, dbClient)

	t.Run("transactions", func(t *testing.T) {
		testBrandsTX(t, ctx, dbClient)
		testStoresTX(t, ctx, dbClient)
		testLocationsTX(t, ctx, dbClient)
		testMaterialsTX(t, ctx, dbClient)
		testColorsTX(t, ctx, dbClient)
		testSpoolsTX(t, ctx, dbClient)
	})

	t.Run("raw", func(t *testing.T) {
		testBrands(t, ctx, dbClient)
		testStores(t, ctx, dbClient)
		testLocations(t, ctx, dbClient)
		testMaterials(t, ctx, dbClient)
		testColors(t, ctx, dbClient)
		testSpools(t, ctx, dbClient)
	})
}

func getQuerier(t *testing.T, ctx context.Context, dbClient *db.Client) (querier *db.Queries, tx pgx.Tx) {
	trx, err := dbClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		t.Fatalf("could not begin transaction: %v", err)
	}
	return dbClient.Queries.WithTx(trx), trx
}

func testBeforeMigrations(t *testing.T, ctx context.Context, dbClient *db.Client) {
	if _, err := dbClient.Queries.FindBrands(ctx); err == nil {
		t.Fatal("should not be able to find brands")
	}

	if _, err := dbClient.Queries.FindColors(ctx); err == nil {
		t.Fatal("should not be able to find colors")
	}

	if _, err := dbClient.Queries.FindLocations(ctx); err == nil {
		t.Fatal("should not be able to find locations")
	}

	if _, err := dbClient.Queries.FindMaterials(ctx); err == nil {
		t.Fatal("should not be able to find materials")
	}

	if _, err := dbClient.Queries.FindSpools(ctx); err == nil {
		t.Fatal("should not be able to find spools")
	}

	if _, err := dbClient.Queries.FindStores(ctx); err == nil {
		t.Fatal("should not be able to find stores")
	}

	if _, err := dbClient.Queries.GetSpoolColors(ctx, 1); err == nil {
		t.Fatal("should not be able to get spool colors")
	}
}

func testMigrations(t *testing.T, ctx context.Context, dbConfig db.Config, dbClient *db.Client) {
	t.Run("bad migrate", func(t *testing.T) {
		badConfig := db.Config{
			Host:         "invalid_host",
			Port:         "9999",
			Database:     "nonexistent_db",
			Username:     "wrong_user",
			Password:     "wrong_password",
			VersionTable: "x.y.z",
		}

		t.Run("bad connection", func(t *testing.T) {
			_, err := db.NewMigrator(ctx, badConfig, nil)
			switch me := err.(type) {
			case db.MigratorError:
				if !me.Connect {
					t.Fatalf("expected Connect=true, got false")
				}
				if me.URL != badConfig.GetRedactedURI() {
					t.Fatalf("expected URL %q, got %q", badConfig.GetRedactedURI(), me.URL)
				}
				if me.Err == nil {
					t.Fatalf("expected non-nil Err")
				}
				if !strings.Contains(err.Error(), "could not connect to database") {
					t.Fatalf("expected connection error, got: %v", me.Err)
				}
			default:
				t.Fatalf("expected MigratorError, got %[1]T: %[1]v", err)
			}
		})

		t.Run("bad version table", func(t *testing.T) {
			badConfig = db.Config{
				Host:         dbConfig.Host,
				Port:         dbConfig.Port,
				Database:     dbConfig.Database,
				Username:     dbConfig.Username,
				Password:     dbConfig.Password,
				VersionTable: "invalid table name; --",
			}
			_, err := db.NewMigrator(ctx, badConfig, nil)
			switch me := err.(type) {
			case db.MigratorError:
				if !me.Create {
					t.Fatalf("expected Create=true, got false")
				}
				if me.Err == nil {
					t.Fatalf("expected non-nil Err")
				}
				if !strings.Contains(err.Error(), "could not create migrator") {
					t.Fatalf("expected connection error, got: %v", me.Err)
				}
			default:
				t.Fatalf("expected MigratorError, got %[1]T: %[1]v", err)
			}
		})

		t.Run("nil filesystem", func(t *testing.T) {
			_, err := db.NewMigrator(ctx, dbConfig, nil)
			switch me := err.(type) {
			case db.MigratorError:
				if !me.Filesystem {
					t.Fatalf("expected Filesystem=true, got false")
				}
				if me.Err == nil {
					t.Fatalf("expected non-nil Err")
				}
				if !strings.Contains(err.Error(), "could not load migrations") {
					t.Fatalf("expected connection error, got: %v", me.Err)
				}
			default:
				t.Fatalf("expected MigratorError, got %[1]T: %[1]v", err)
			}
		})

		t.Run("wrong filesystem", func(t *testing.T) {
			fs := migrations.Collection{
				Filesystem: static.Files,
			}

			_, err := db.NewMigrator(ctx, dbConfig, fs)
			switch me := err.(type) {
			case db.MigratorError:
				if !me.Filesystem {
					t.Fatalf("expected Filesystem=true, got false")
				}
				if me.Err == nil {
					t.Fatalf("expected non-nil Err")
				}
				if !strings.Contains(err.Error(), "could not load migrations") {
					t.Fatalf("expected connection error, got: %v", me.Err)
				}
			default:
				t.Fatalf("expected MigratorError, got %[1]T: %[1]v", err)
			}
		})

		t.Run("unknown error", func(t *testing.T) {
			err := db.MigratorError{}

			if !strings.Contains(err.Error(), "unknown migration error") {
				t.Fatalf("expected unknown migration error, got: %v", err)
			}
		})
	})

	var col migrations.Collection
	var err error

	t.Run("good migrate", func(t *testing.T) {
		col, err = migrations.New()
		if err != nil {
			t.Fatalf("could not create the embedded filesystem: %v", err)
		}

		m, err := db.NewMigrator(ctx, dbConfig, col)
		if err != nil {
			t.Fatalf("could not create the migrator: %v", err)
		}

		i, err := m.Info(ctx, -1)
		if err != nil {
			t.Fatalf("could not get the migration info: %v", err)
		}

		t.Logf("host: %s:%s, database: %s, currentVersion: %d, targetVersion: %d\n%s", i.Hostname, i.Port, i.Database, i.Migrations.CurrentVersion, i.Migrations.TargetVersion, i.Migrations.Summary)

		_, err = m.MigrateTo(ctx, col.Steps())
		if err != nil {
			t.Fatalf("could not run the migrations: %v", err)
		}

		_, err = m.MigrateTo(ctx, 0)
		if err != nil {
			t.Fatalf("could not run the migrations: %v", err)
		}

		_, err = m.Migrate(ctx)
		if err != nil {
			t.Fatalf("could not run the migrations: %v", err)
		}
	})

	mig, err := dbClient.Queries.CheckMigration(ctx)
	if err != nil {
		t.Fatalf("could not check migration: %v", err)
	}

	if mig != col.Steps() {
		t.Fatalf("migration steps do not match: expected %d, got %d", col.Steps(), mig)
	}
}

func testBrands(t *testing.T, ctx context.Context, dbClient *db.Client) {
	t.Run("brands", func(t *testing.T) {
		s, err := dbClient.Queries.CreateBrand(ctx, db.CreateBrandParams{Label: "Test Brand", Active: true, StoreID: nil})
		if err != nil {
			t.Fatalf("could not create brand: %v", err)
		}

		_, err = dbClient.Queries.CreateBrand(ctx, db.CreateBrandParams{Label: "Test Brand", Active: true, StoreID: nil})
		if err == nil {
			t.Fatal("should not be able create duplicate brand")
		}

		s, err = dbClient.Queries.GetBrandByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get brand by ID: %v", err)
		}

		s, err = dbClient.Queries.UpdateBrand(ctx, db.UpdateBrandParams{Label: "Test Brand2", ID: s.ID})
		if err != nil {
			t.Fatalf("could not update brand: %v", err)
		}

		_, err = dbClient.Queries.CreateBrand(ctx, db.CreateBrandParams{Label: "Test Brand", Active: true, StoreID: nil})
		if err != nil {
			t.Fatalf("should be able create brand: %v", err)
		}

		m, err := dbClient.Queries.FindBrands(ctx)
		if err != nil {
			t.Fatalf("could not find brands: %v", err)
		}

		if len(m) != (brandCount + 2) {
			t.Fatalf("expected %d brands, got %d", (brandCount + 2), len(m))
		}

		err = dbClient.Queries.DeleteBrand(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete brand: %v", err)
		}

		m, err = dbClient.Queries.FindBrands(ctx)
		if err != nil {
			t.Fatalf("could not find brands: %v", err)
		}

		if len(m) != (brandCount + 1) {
			t.Fatalf("expected %d brands, got %d", (brandCount + 1), len(m))
		}
	})
}

func testStores(t *testing.T, ctx context.Context, dbClient *db.Client) {
	t.Run("stores", func(t *testing.T) {
		s, err := dbClient.Queries.CreateStore(ctx, "Test Store", nil)
		if err != nil {
			t.Fatalf("could not create store: %v", err)
		}

		_, err = dbClient.Queries.CreateStore(ctx, "Test Store", nil)
		if err == nil {
			t.Fatal("should not be able to create duplicate store")
		}

		s, err = dbClient.Queries.GetStoreByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get store by ID: %v", err)
		}

		s, err = dbClient.Queries.UpdateStore(ctx, db.UpdateStoreParams{
			ID:    s.ID,
			Label: "Test Store2",
			URL:   PointerOf("https://example.com"),
		})
		if err != nil {
			t.Fatalf("could not update store: %v", err)
		}

		_, err = dbClient.Queries.CreateStore(ctx, "Test Store", nil)
		if err != nil {
			t.Fatalf("should be able create store: %v", err)
		}

		m, err := dbClient.Queries.FindStores(ctx)
		if err != nil {
			t.Fatalf("could not find stores: %v", err)
		}

		if len(m) != (storeCount + 2) {
			t.Fatalf("expected %d stores, got %d", (storeCount + 2), len(m))
		}

		err = dbClient.Queries.DeleteStore(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete store: %v", err)
		}

		m, err = dbClient.Queries.FindStores(ctx)
		if err != nil {
			t.Fatalf("could not find stores: %v", err)
		}

		if len(m) != (storeCount + 1) {
			t.Fatalf("expected %d stores, got %d", (storeCount + 1), len(m))
		}
	})
}

func testLocations(t *testing.T, ctx context.Context, dbClient *db.Client) {
	t.Run("locations", func(t *testing.T) {
		s, err := dbClient.Queries.CreateLocation(ctx, db.CreateLocationParams{
			Label:       "Test Location",
			Description: "Desc",
			Capacity:    10,
			Printable:   false,
			Tally:       false,
		})
		if err != nil {
			t.Fatalf("could not create location: %v", err)
		}

		_, err = dbClient.Queries.CreateLocation(ctx, db.CreateLocationParams{
			Label:       "Test Location",
			Description: "Desc",
			Capacity:    10,
			Printable:   false,
			Tally:       false,
		})
		if err == nil {
			t.Fatal("should not be able to create duplicate location")
		}

		s, err = dbClient.Queries.GetLocationByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get location by ID: %v", err)
		}

		s, err = dbClient.Queries.UpdateLocation(ctx, db.UpdateLocationParams{
			ID:          s.ID,
			Label:       "Test Location2",
			Description: "A location for testing",
			Capacity:    0,
			Printable:   true,
			Tally:       true,
		})
		if err != nil {
			t.Fatalf("could not update location: %v", err)
		}

		_, err = dbClient.Queries.CreateLocation(ctx, db.CreateLocationParams{
			Label:       "Test Location",
			Description: "Desc",
			Capacity:    10,
			Printable:   false,
			Tally:       false,
		})
		if err != nil {
			t.Fatalf("should be able create location: %v", err)
		}

		m, err := dbClient.Queries.FindLocations(ctx)
		if err != nil {
			t.Fatalf("could not find locations: %v", err)
		}

		if len(m) != (locationCount + 2) {
			t.Fatalf("expected %d locations, got %d", (locationCount + 2), len(m))
		}

		err = dbClient.Queries.DeleteLocation(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete location: %v", err)
		}

		m, err = dbClient.Queries.FindLocations(ctx)
		if err != nil {
			t.Fatalf("could not find locations: %v", err)
		}

		if len(m) != (locationCount + 1) {
			t.Fatalf("expected %d locations, got %d", (locationCount + 1), len(m))
		}
	})
}

func testMaterials(t *testing.T, ctx context.Context, dbClient *db.Client) {
	t.Run("materials", func(t *testing.T) {
		s, err := dbClient.Queries.CreateMaterial(ctx, db.CreateMaterialParams{
			Label:   "Test Material",
			Class:   "Desc",
			Special: false,
		})
		if err != nil {
			t.Fatalf("could not create material: %v", err)
		}

		_, err = dbClient.Queries.CreateMaterial(ctx, db.CreateMaterialParams{
			Label:   "Test Material",
			Class:   "Desc",
			Special: false,
		})
		if err == nil {
			t.Fatal("should not be able to create duplicate material")
		}

		s, err = dbClient.Queries.GetMaterialByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get material by ID: %v", err)
		}

		s, err = dbClient.Queries.UpdateMaterial(ctx, db.UpdateMaterialParams{
			ID:      s.ID,
			Label:   "Test Material2",
			Class:   "A material for testing",
			Special: true,
		})
		if err != nil {
			t.Fatalf("could not update material: %v", err)
		}

		_, err = dbClient.Queries.CreateMaterial(ctx, db.CreateMaterialParams{
			Label:   "Test Material",
			Class:   "Desc",
			Special: false,
		})
		if err != nil {
			t.Fatalf("should be able create material: %v", err)
		}

		m, err := dbClient.Queries.FindMaterials(ctx)
		if err != nil {
			t.Fatalf("could not find materials: %v", err)
		}

		if len(m) != (materialCount + 2) {
			t.Fatalf("expected %d materials, got %d", (materialCount + 2), len(m))
		}

		err = dbClient.Queries.DeleteMaterial(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete material: %v", err)
		}

		m, err = dbClient.Queries.FindMaterials(ctx)
		if err != nil {
			t.Fatalf("could not find materials: %v", err)
		}

		if len(m) != (materialCount + 1) {
			t.Fatalf("expected %d materials, got %d", (materialCount + 1), len(m))
		}
	})
}

func testColors(t *testing.T, ctx context.Context, dbClient *db.Client) {
	t.Run("colors", func(t *testing.T) {
		s, err := dbClient.Queries.CreateColor(ctx, db.CreateColorParams{
			Label:   "Test Color",
			HexCode: "#FFFFFF",
			Alias:   nil,
		})
		if err != nil {
			t.Fatalf("could not create color: %v", err)
		}

		_, err = dbClient.Queries.CreateColor(ctx, db.CreateColorParams{
			Label:   "Test Color",
			HexCode: "#FFFFFF",
			Alias:   nil,
		})
		if err == nil {
			t.Fatal("should not be able to create duplicate color")
		}

		s, err = dbClient.Queries.GetColorByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get color by ID: %v", err)
		}

		s, err = dbClient.Queries.UpdateColor(ctx, db.UpdateColorParams{
			ID:      s.ID,
			Label:   "Test Color2",
			HexCode: "#FFFFFF",
			Alias:   PointerOf("White"),
		})
		if err != nil {
			t.Fatalf("could not update color: %v", err)
		}

		_, err = dbClient.Queries.CreateColor(ctx, db.CreateColorParams{
			Label:   "Test Color",
			HexCode: "#FF0000",
			Alias:   nil,
		})
		if err != nil {
			t.Fatalf("should be able create color: %v", err)
		}

		m, err := dbClient.Queries.FindColors(ctx)
		if err != nil {
			t.Fatalf("could not find colors: %v", err)
		}

		if len(m) != (colorCount + 2) {
			t.Fatalf("expected %d colors, got %d", (colorCount + 2), len(m))
		}

		err = dbClient.Queries.DeleteColor(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete color: %v", err)
		}

		m, err = dbClient.Queries.FindColors(ctx)
		if err != nil {
			t.Fatalf("could not find colors: %v", err)
		}

		if len(m) != (colorCount + 1) {
			t.Fatalf("expected %d colors, got %d", colorCount, len(m))
		}
	})
}

func testSpools(t *testing.T, ctx context.Context, dbClient *db.Client) {
	t.Run("spools", func(t *testing.T) {
		locations, err := dbClient.Queries.FindLocations(ctx)
		if err != nil {
			t.Fatalf("could not find locations: %v", err)
		}

		if len(locations) == 0 {
			t.Fatalf("no locations found to assign to spool")
		}

		brands, err := dbClient.Queries.FindBrands(ctx)
		if err != nil {
			t.Fatalf("could not find brands: %v", err)
		}

		if len(brands) == 0 {
			t.Fatalf("no brands found to assign to spool")
		}

		materials, err := dbClient.Queries.FindMaterials(ctx)
		if err != nil {
			t.Fatalf("could not find materials: %v", err)
		}

		if len(materials) == 0 {
			t.Fatalf("no materials found to assign to spool")
		}

		stores, err := dbClient.Queries.FindStores(ctx)
		if err != nil {
			t.Fatalf("could not find stores: %v", err)
		}

		if len(stores) == 0 {
			t.Fatalf("no stores found to assign to spool")
		}

		colors, err := dbClient.Queries.FindColors(ctx)
		if err != nil {
			t.Fatalf("could not find colors: %v", err)
		}

		if len(colors) == 0 {
			t.Fatalf("no colors found to assign to spool")
		}

		price := pgtype.Numeric{}
		price.Scan("29.99")
		weight := pgtype.Numeric{}
		weight.Scan("1000.00")
		currentWeight := pgtype.Numeric{}
		currentWeight.Scan("900.00")
		combinedWeight := pgtype.Numeric{}
		combinedWeight.Scan("1167.00")

		s, err := dbClient.Queries.CreateSpool(ctx, db.CreateSpoolParams{
			Location:       locations[locationCount-1].ID,
			Brand:          brands[brandCount-1].ID,
			Material:       materials[materialCount-1].ID,
			Store:          stores[storeCount-1].ID,
			Weight:         weight,
			CurrentWeight:  currentWeight,
			CombinedWeight: combinedWeight,
			Price:          price,
			Empty:          false,
		})
		if err != nil {
			t.Fatalf("could not create spool: %v", err)
		}

		s, err = dbClient.Queries.GetSpoolByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get spool by ID: %v", err)
		}

		newWeight := pgtype.Numeric{}
		newWeight.Scan("2900.00")

		newCurrentWeight := pgtype.Numeric{}
		newCurrentWeight.Scan("2800.00")

		newCombinedWeight := pgtype.Numeric{}
		newCombinedWeight.Scan("2900.00")

		s, err = dbClient.Queries.UpdateSpool(ctx, db.UpdateSpoolParams{
			ID:             s.ID,
			Location:       locations[locationCount-1].ID,
			Brand:          brands[brandCount-1].ID,
			Material:       materials[materialCount-1].ID,
			Store:          stores[storeCount-1].ID,
			Weight:         newWeight,
			CurrentWeight:  newCurrentWeight,
			CombinedWeight: newCombinedWeight,
			Price:          price,
			Empty:          false,
		})
		if err != nil {
			t.Fatalf("could not update spool: %v", err)
		}

		_, err = dbClient.Queries.CreateSpool(ctx, db.CreateSpoolParams{
			Location:       locations[locationCount-1].ID,
			Brand:          brands[brandCount-1].ID,
			Material:       materials[materialCount-1].ID,
			Store:          stores[storeCount-1].ID,
			Weight:         newWeight,
			CurrentWeight:  newCurrentWeight,
			CombinedWeight: newCombinedWeight,
			Price:          price,
			Empty:          false,
		})
		if err != nil {
			t.Fatalf("should be able create spool: %v", err)
		}

		m, err := dbClient.Queries.FindSpools(ctx)
		if err != nil {
			t.Fatalf("could not find spools: %v", err)
		}

		if len(m) != (spoolCount + 2) {
			t.Fatalf("expected %d spools, got %d", (spoolCount + 2), len(m))
		}

		err = dbClient.Queries.DeleteSpool(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete spool: %v", err)
		}

		m, err = dbClient.Queries.FindSpools(ctx)
		if err != nil {
			t.Fatalf("could not find spools: %v", err)
		}

		if len(m) != (spoolCount + 1) {
			t.Fatalf("expected %d spools, got %d", (spoolCount + 1), len(m))
		}

		deleted := 0
		for _, s := range m {
			if s.DeletedAt != nil {
				deleted++
			}
		}

		if deleted > 0 {
			t.Fatalf("unexpected deleted spool, found %d deleted spools", deleted)
		}

		colorIDs := []int64{
			colors[colorCount-1].ID,
		}

		err = dbClient.Queries.AddSpoolColors(ctx, s.ID, colorIDs)
		if err != nil {
			t.Fatalf("could not add spool color: %v", err)
		}

		c, err := dbClient.Queries.GetSpoolColors(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get spool colors: %v", err)
		}

		if len(c) != 1 {
			t.Fatalf("expected 1 spool color, got %d", len(c))
		}

		err = dbClient.Queries.ResetSpoolColor(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not remove spool color: %v", err)
		}

		_, err = dbClient.Queries.GetStorageStats(ctx)
		if err != nil {
			t.Fatalf("could not get storage stats: %v", err)
		}

		_, err = dbClient.Queries.GetUsageStats(ctx)
		if err != nil {
			t.Fatalf("could not get usage stats: %v", err)
		}
	})
}

func testBrandsTX(t *testing.T, ctx context.Context, dbClient *db.Client) {
	querier, tx := getQuerier(t, ctx, dbClient)

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	t.Run("brands", func(t *testing.T) {
		s, err := querier.CreateBrand(ctx, db.CreateBrandParams{Label: "Test Brand TX", Active: true, StoreID: nil})
		if err != nil {
			t.Fatalf("could not create brand: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_, err = querier.CreateBrand(ctx, db.CreateBrandParams{Label: "Test Brand TX", Active: true, StoreID: nil})
		if err == nil {
			t.Fatal("should not be able create duplicate brand")
		}

		err = tx.Commit(ctx)
		if err == nil {
			t.Fatal("should not be able create duplicate brand")
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		s, err = querier.GetBrandByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get brand by ID: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		s, err = querier.UpdateBrand(ctx, db.UpdateBrandParams{Label: "Test Brand2 TX", ID: s.ID})
		if err != nil {
			t.Fatalf("could not update brand: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_, err = querier.CreateBrand(ctx, db.CreateBrandParams{Label: "Test Brand TX", Active: true, StoreID: nil})
		if err != nil {
			t.Fatalf("should be able create brand: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		m, err := querier.FindBrands(ctx)
		if err != nil {
			t.Fatalf("could not find brands: %v", err)
		}

		if len(m) != (brandCount + 1) {
			t.Fatalf("expected %d brands, got %d", (brandCount + 1), len(m))
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		err = querier.DeleteBrand(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete brand: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		m, err = querier.FindBrands(ctx)
		if err != nil {
			t.Fatalf("could not find brands: %v", err)
		}

		if len(m) != brandCount {
			t.Fatalf("expected %d brands, got %d", brandCount, len(m))
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
	})
}

func testStoresTX(t *testing.T, ctx context.Context, dbClient *db.Client) {
	querier, tx := getQuerier(t, ctx, dbClient)

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	t.Run("stores", func(t *testing.T) {
		s, err := querier.CreateStore(ctx, "Test Store TX", nil)
		if err != nil {
			t.Fatalf("could not create store: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_, err = querier.CreateStore(ctx, "Test Store TX", nil)
		if err == nil {
			t.Fatal("should not be able to create duplicate store")
		}

		err = tx.Commit(ctx)
		if err == nil {
			t.Fatal("should not be able to create duplicate store")
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		s, err = querier.GetStoreByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get store by ID: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		s, err = querier.UpdateStore(ctx, db.UpdateStoreParams{
			ID:    s.ID,
			Label: "Test Store2 TX",
			URL:   PointerOf("https://example.com"),
		})
		if err != nil {
			t.Fatalf("could not update store: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_, err = querier.CreateStore(ctx, "Test Store TX", nil)
		if err != nil {
			t.Fatalf("should be able create store: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		m, err := querier.FindStores(ctx)
		if err != nil {
			t.Fatalf("could not find stores: %v", err)
		}

		if len(m) != (storeCount + 1) {
			t.Fatalf("expected %d stores, got %d", (storeCount + 1), len(m))
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		err = querier.DeleteStore(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete store: %v", err)
		}

		m, err = querier.FindStores(ctx)
		if err != nil {
			t.Fatalf("could not find stores: %v", err)
		}

		if len(m) != storeCount {
			t.Fatalf("expected %d stores, got %d", storeCount, len(m))
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
	})
}

func testLocationsTX(t *testing.T, ctx context.Context, dbClient *db.Client) {
	querier, tx := getQuerier(t, ctx, dbClient)

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	t.Run("locations", func(t *testing.T) {
		s, err := querier.CreateLocation(ctx, db.CreateLocationParams{
			Label:       "Test Location TX",
			Description: "Desc",
			Capacity:    10,
			Printable:   false,
			Tally:       false,
		})
		if err != nil {
			t.Fatalf("could not create location: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_, err = querier.CreateLocation(ctx, db.CreateLocationParams{
			Label:       "Test Location TX",
			Description: "Desc",
			Capacity:    10,
			Printable:   false,
			Tally:       false,
		})
		if err == nil {
			t.Fatal("should not be able to create duplicate location")
		}

		err = tx.Commit(ctx)
		if err == nil {
			t.Fatal("should not be able to create duplicate location")
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		s, err = querier.GetLocationByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get location by ID: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		s, err = querier.UpdateLocation(ctx, db.UpdateLocationParams{
			ID:          s.ID,
			Label:       "Test Location2 TX",
			Description: "A location for testing",
			Capacity:    0,
			Printable:   true,
			Tally:       true,
		})
		if err != nil {
			t.Fatalf("could not update location: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_, err = querier.CreateLocation(ctx, db.CreateLocationParams{
			Label:       "Test Location TX",
			Description: "Desc",
			Capacity:    10,
			Printable:   false,
			Tally:       false,
		})
		if err != nil {
			t.Fatalf("should be able create location: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		m, err := querier.FindLocations(ctx)
		if err != nil {
			t.Fatalf("could not find locations: %v", err)
		}

		if len(m) != (locationCount + 1) {
			t.Fatalf("expected %d locations, got %d", (locationCount + 1), len(m))
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		err = querier.DeleteLocation(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete location: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		m, err = querier.FindLocations(ctx)
		if err != nil {
			t.Fatalf("could not find locations: %v", err)
		}

		if len(m) != locationCount {
			t.Fatalf("expected %d locations, got %d", locationCount, len(m))
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
	})
}

func testMaterialsTX(t *testing.T, ctx context.Context, dbClient *db.Client) {
	querier, tx := getQuerier(t, ctx, dbClient)

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	t.Run("materials", func(t *testing.T) {
		s, err := querier.CreateMaterial(ctx, db.CreateMaterialParams{
			Label:   "Test Material TX",
			Class:   "Desc",
			Special: false,
		})
		if err != nil {
			t.Fatalf("could not create material: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_, err = querier.CreateMaterial(ctx, db.CreateMaterialParams{
			Label:   "Test Material TX",
			Class:   "Desc",
			Special: false,
		})
		if err == nil {
			t.Fatal("should not be able to create duplicate material")
		}

		err = tx.Commit(ctx)
		if err == nil {
			t.Fatal("should not be able to create duplicate material")
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		s, err = querier.GetMaterialByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get material by ID: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		s, err = querier.UpdateMaterial(ctx, db.UpdateMaterialParams{
			ID:      s.ID,
			Label:   "Test Material2 TX",
			Class:   "A material for testing",
			Special: true,
		})
		if err != nil {
			t.Fatalf("could not update material: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_, err = querier.CreateMaterial(ctx, db.CreateMaterialParams{
			Label:   "Test Material TX",
			Class:   "Desc",
			Special: false,
		})
		if err != nil {
			t.Fatalf("should be able create material: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		m, err := querier.FindMaterials(ctx)
		if err != nil {
			t.Fatalf("could not find materials: %v", err)
		}

		if len(m) != (materialCount + 1) {
			t.Fatalf("%d materials, got %d", (materialCount + 1), len(m))
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		err = querier.DeleteMaterial(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete material: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		m, err = querier.FindMaterials(ctx)
		if err != nil {
			t.Fatalf("could not find materials: %v", err)
		}

		if len(m) != materialCount {
			t.Fatalf("expected %d materials, got %d", materialCount, len(m))
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
	})
}

func testColorsTX(t *testing.T, ctx context.Context, dbClient *db.Client) {
	querier, tx := getQuerier(t, ctx, dbClient)

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	t.Run("colors", func(t *testing.T) {
		s, err := querier.CreateColor(ctx, db.CreateColorParams{
			Label:   "Test Color TX",
			HexCode: "#FFFFFF",
			Alias:   nil,
		})
		if err != nil {
			t.Fatalf("could not create color: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_, err = querier.CreateColor(ctx, db.CreateColorParams{
			Label:   "Test Color TX",
			HexCode: "#FFFFFF",
			Alias:   nil,
		})
		if err == nil {
			t.Fatal("should not be able to create duplicate color")
		}

		err = tx.Commit(ctx)
		if err == nil {
			t.Fatal("should not be able to create duplicate color")
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		s, err = querier.GetColorByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get color by ID: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		s, err = querier.UpdateColor(ctx, db.UpdateColorParams{
			ID:      s.ID,
			Label:   "Test Color2 TX",
			HexCode: "#FFFFFF",
			Alias:   PointerOf("White"),
		})
		if err != nil {
			t.Fatalf("could not update color: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		_, err = querier.CreateColor(ctx, db.CreateColorParams{
			Label:   "Test Color TX",
			HexCode: "#FF0000",
			Alias:   nil,
		})
		if err != nil {
			t.Fatalf("should be able create color: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		m, err := querier.FindColors(ctx)
		if err != nil {
			t.Fatalf("could not find colors: %v", err)
		}

		if len(m) != (colorCount + 1) {
			t.Fatalf("expected %d colors, got %d", (colorCount + 1), len(m))
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		err = querier.DeleteColor(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete color: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
		querier, tx = getQuerier(t, ctx, dbClient)
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		m, err = querier.FindColors(ctx)
		if err != nil {
			t.Fatalf("could not find colors: %v", err)
		}

		if len(m) != colorCount {
			t.Fatalf("expected %d colors, got %d", colorCount, len(m))
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
	})
}

func testSpoolsTX(t *testing.T, ctx context.Context, dbClient *db.Client) {
	querier, tx := getQuerier(t, ctx, dbClient)

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	t.Run("spools", func(t *testing.T) {
		locations, err := querier.FindLocations(ctx)
		if err != nil {
			t.Fatalf("could not find locations: %v", err)
		}

		if len(locations) == 0 {
			t.Fatalf("no locations found to assign to spool")
		}

		brands, err := querier.FindBrands(ctx)
		if err != nil {
			t.Fatalf("could not find brands: %v", err)
		}

		if len(brands) == 0 {
			t.Fatalf("no brands found to assign to spool")
		}

		materials, err := querier.FindMaterials(ctx)
		if err != nil {
			t.Fatalf("could not find materials: %v", err)
		}

		if len(materials) == 0 {
			t.Fatalf("no materials found to assign to spool")
		}

		stores, err := querier.FindStores(ctx)
		if err != nil {
			t.Fatalf("could not find stores: %v", err)
		}

		if len(stores) == 0 {
			t.Fatalf("no stores found to assign to spool")
		}

		colors, err := querier.FindColors(ctx)
		if err != nil {
			t.Fatalf("could not find colors: %v", err)
		}

		if len(colors) == 0 {
			t.Fatalf("no colors found to assign to spool")
		}

		price := pgtype.Numeric{}
		price.Scan("29.99")
		weight := pgtype.Numeric{}
		weight.Scan("1000.00")
		currentWeight := pgtype.Numeric{}
		currentWeight.Scan("900.00")
		combinedWeight := pgtype.Numeric{}
		combinedWeight.Scan("1167.00")

		s, err := querier.CreateSpool(ctx, db.CreateSpoolParams{
			Location:       locations[locationCount-1].ID,
			Brand:          brands[brandCount-1].ID,
			Material:       materials[materialCount-1].ID,
			Store:          stores[storeCount-1].ID,
			Weight:         weight,
			CurrentWeight:  currentWeight,
			CombinedWeight: combinedWeight,
			Price:          price,
			Empty:          false,
		})
		if err != nil {
			t.Fatalf("could not create spool: %v", err)
		}

		s, err = querier.GetSpoolByID(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get spool by ID: %v", err)
		}

		newWeight := pgtype.Numeric{}
		newWeight.Scan("2900.00")

		newCurrentWeight := pgtype.Numeric{}
		newCurrentWeight.Scan("2800.00")

		newCombinedWeight := pgtype.Numeric{}
		newCombinedWeight.Scan("2900.00")

		s, err = querier.UpdateSpool(ctx, db.UpdateSpoolParams{
			ID:             s.ID,
			Location:       locations[locationCount-1].ID,
			Brand:          brands[brandCount-1].ID,
			Material:       materials[materialCount-1].ID,
			Store:          stores[storeCount-1].ID,
			Weight:         newWeight,
			CurrentWeight:  newCurrentWeight,
			CombinedWeight: newCombinedWeight,
			Price:          price,
			Empty:          false,
		})
		if err != nil {
			t.Fatalf("could not update spool: %v", err)
		}

		_, err = querier.CreateSpool(ctx, db.CreateSpoolParams{
			Location:       locations[locationCount-1].ID,
			Brand:          brands[brandCount-1].ID,
			Material:       materials[materialCount-1].ID,
			Store:          stores[storeCount-1].ID,
			Weight:         newWeight,
			CurrentWeight:  newCurrentWeight,
			CombinedWeight: newCombinedWeight,
			Price:          price,
			Empty:          false,
		})
		if err != nil {
			t.Fatalf("should be able create spool: %v", err)
		}

		m, err := querier.FindSpools(ctx)
		if err != nil {
			t.Fatalf("could not find spools: %v", err)
		}

		if len(m) != (spoolCount + 1) {
			t.Fatalf("expected %d spools, got %d", (spoolCount + 1), len(m))
		}

		err = querier.DeleteSpool(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not delete spool: %v", err)
		}

		m, err = querier.FindSpools(ctx)
		if err != nil {
			t.Fatalf("could not find spools: %v", err)
		}

		if len(m) != spoolCount {
			t.Fatalf("expected %d spools, got %d", spoolCount, len(m))
		}

		deleted := 0
		for _, s := range m {
			if s.DeletedAt != nil {
				deleted++
			}
		}

		if deleted > 0 {
			t.Fatalf("unexpected deleted spool, found %d deleted spools", deleted)
		}

		colorIDs := []int64{
			colors[colorCount-1].ID,
		}

		err = querier.AddSpoolColors(ctx, s.ID, colorIDs)
		if err != nil {
			t.Fatalf("could not add spool color: %v", err)
		}

		c, err := querier.GetSpoolColors(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not get spool colors: %v", err)
		}

		if len(c) != 1 {
			t.Fatalf("expected 1 spool color, got %d", len(c))
		}

		err = querier.ResetSpoolColor(ctx, s.ID)
		if err != nil {
			t.Fatalf("could not remove spool color: %v", err)
		}

		_, err = querier.GetStorageStats(ctx)
		if err != nil {
			t.Fatalf("could not get storage stats: %v", err)
		}

		_, err = querier.GetUsageStats(ctx)
		if err != nil {
			t.Fatalf("could not get usage stats: %v", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			t.Fatalf("could not commit transaction: %v", err)
		}
	})
}
