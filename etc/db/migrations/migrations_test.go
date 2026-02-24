package migrations_test

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/jsnfwlr/filamate/etc/db/migrations"
)

func TestValidateMigrations(t *testing.T) {
	files, err := os.ReadDir("./")
	if err != nil {
		t.Fatalf("could not read files from the directory this test is in: %s", err)
	}

	col, err := migrations.New()
	if err != nil {
		t.Fatalf("could not read files from the embeddedFS: %s", err)
	}

	if col.Steps() == 0 {
		t.Fatalf("no migration files found in the embeddedFS")
	}

	sqlFiles := []string{}

	for _, f := range files {
		if migrations.AllPattern.MatchString(f.Name()) {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}

	t.Logf("Found %d regular SQL files in the local directory", len(sqlFiles))

	slices.Sort(sqlFiles)

	migs := []string{}
	fEntries := col.Migrations()
	for _, f := range fEntries {
		if migrations.AllPattern.MatchString(f.Name()) {
			migs = append(migs, f.Name())
		}
	}

	slices.Sort(migs)
	t.Logf("Found %d regular SQL files in the embeddedFS", len(migs))

	compareMigrations(t, sqlFiles, migs)
	checkCompliance(t, col, sqlFiles)
}

type infoLine struct {
	Sequence int
	Name     string
}

func makeInfo(t *testing.T, files []string) (targetVersion int32, information string, lines []infoLine, fault error) {
	t.Helper()

	var tv int32
	var l []infoLine

	i := ""

	fPartReg := regexp.MustCompile(`([0-9]+)_(.*)\.sql`)
	numLeadZero := regexp.MustCompile(`^0+`)
	for _, file := range files {
		matches := fPartReg.FindAllStringSubmatch(file, -1)

		indicator := "  "

		num := numLeadZero.ReplaceAllString(matches[0][1], "")
		i += fmt.Sprintf("%2s %3s %s\n", indicator, num, matches[0][2])
		if num == "" {
			num = "0"
		}

		t.Logf("Processing file %s - (%s %s %s)", file, matches[0][1], num, matches[0][2])

		n, err := strconv.Atoi(num)
		if err != nil {
			return 0, "", []infoLine{}, fmt.Errorf("could not convert %s to int: %w", num, err)
		}

		l = append(l, infoLine{Sequence: n, Name: matches[0][2]})
		if int32(n) > tv {
			tv = int32(n)
		}
	}

	return tv, i, l, nil
}

func compareMigrations(t *testing.T, sqlFiles, migs []string) {
	lTV, lI, lL, err := makeInfo(t, sqlFiles)
	if err != nil {
		t.Fatalf("could not get local information: %s", err)
	}

	t.Logf("Found %d SQL files in the embeddedFS", len(migs))
	mTV, mI, mL, err := makeInfo(t, migs)
	if err != nil {
		t.Fatalf("could not get migration information: %s", err)
	}

	t.Run("match_target_versions", func(t *testing.T) {
		if mTV != lTV {
			t.Errorf("Target version should be %d, not %d", lTV, mTV)
		}
	})

	t.Run("compare_migrations", func(t *testing.T) {
		if mI != lI {
			if len(mL) != len(lL) {
				t.Errorf("length of info differs: migrations - %d vs local - %d", len(mL), len(lL))
			}
			for i := range mL {
				if mL[i].Sequence != lL[i].Sequence {
					t.Errorf("%d - %d vs %d", i, mL[i].Sequence, lL[i].Sequence)
				}
				if mL[i].Name != lL[i].Name {
					t.Errorf("%d - %s vs %s", i, mL[i].Name, lL[i].Name)
				}
			}
		}
	})
}

func checkCompliance(t *testing.T, col migrations.Collection, sqlFiles []string) {
	t.Helper()

	t.Run("check_collection_compliance", func(t *testing.T) {
		files, _ := col.ReadDir(".")
		fileNames := []string{}
		for _, f := range files {
			fileNames = append(fileNames, f.Name())
		}

		if len(files) != len(sqlFiles) {
			t.Logf("Collection: [%s], Local: [%s]", strings.Join(fileNames, ", "), strings.Join(sqlFiles, ", "))
			t.Errorf("number of files in the collection (%d) does not match the number of files in the local directory (%d)", len(files), len(sqlFiles))
		}

		_, err := col.ReadDir("nonexistent")
		if err == nil {
			t.Errorf("expected an error when reading a nonexistent directory, but got none")
		}

		globs, err := col.Glob("*.sql")
		if err != nil {
			t.Errorf("could not get globs: %s", err)
		}

		if len(globs) != len(sqlFiles) {
			t.Logf("Glob: [%s], Local: [%s]", strings.Join(globs, ", "), strings.Join(sqlFiles, ", "))
			t.Errorf("number of glob matches (%d) does not match the number of files in the local directory (%d)", len(globs), len(sqlFiles))
		}

		b, err := col.ReadFile(globs[0])
		if err != nil {
			t.Errorf("could not read the first file in the glob list: %s", err)
		}
		if len(b) == 0 {
			t.Error("the first file in the glob list is empty")
		}

		_, err = col.Open("does_not_exist.sql")
		if err == nil {
			t.Error("should get an error here")
		}

		fi, err := col.Open(globs[0])
		if err != nil {
			t.Errorf("should not get an error here: %s", err)
		}

		stat, err := fi.Stat()
		if err != nil {
			t.Errorf("should not get an error here: %s", err)
		}

		if stat.Name() != globs[0] {
			t.Errorf("expected the name to be %s, but got %s", globs[0], stat.Name())
		}

		_, err = col.ReadFile("does_not_exist.sql")
		if err == nil {
			t.Error("should get an error here")
		}

		globs, err = col.Glob("nonexistent/*.sql")
		if err != nil {
			t.Errorf("should not get error here: %s", err)
		}

		if len(globs) != 0 {
			t.Errorf("expected 0 glob matches for a nonexistent directory, but got %d [%s]", len(globs), strings.Join(globs, ", "))
		}

		_, err = col.Glob("[]a]")
		if err == nil {
			t.Error("should get a pattern error here")
		}
	})
}
