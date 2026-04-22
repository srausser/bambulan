package bambu3mf

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestParseMetadata(t *testing.T) {
	filename := writeTestProjectFile(t)

	r, err := Open(filename)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer r.Close()

	md, err := r.ParseMetadata()
	if err != nil {
		t.Fatalf("ParseMetadata failed: %v", err)
	}

	if md.Title != "Fixture Cube" {
		t.Fatalf("Title = %q, want %q", md.Title, "Fixture Cube")
	}
	if md.Designer != "Codex" {
		t.Fatalf("Designer = %q, want %q", md.Designer, "Codex")
	}
	if md.Description != "Self-contained 3MF test fixture" {
		t.Fatalf("Description = %q, want %q", md.Description, "Self-contained 3MF test fixture")
	}
	if md.ThumbnailPath != "Metadata/thumbnail.png" {
		t.Fatalf("ThumbnailPath = %q, want %q", md.ThumbnailPath, "Metadata/thumbnail.png")
	}

	if len(md.Plates) != 1 {
		t.Fatalf("len(Plates) = %d, want 1", len(md.Plates))
	}
	plate := md.Plates[0]
	if plate.ID != 1 {
		t.Fatalf("Plate ID = %d, want 1", plate.ID)
	}
	if plate.Name != "Fixture Plate" {
		t.Fatalf("Plate Name = %q, want %q", plate.Name, "Fixture Plate")
	}
	if plate.ThumbnailPath != "Metadata/plate_1.png" {
		t.Fatalf("Plate ThumbnailPath = %q, want %q", plate.ThumbnailPath, "Metadata/plate_1.png")
	}
	if plate.ThumbnailSmall != "Metadata/plate_1_small.png" {
		t.Fatalf("Plate ThumbnailSmall = %q, want %q", plate.ThumbnailSmall, "Metadata/plate_1_small.png")
	}

	if len(md.Filaments) != 1 {
		t.Fatalf("len(Filaments) = %d, want 1", len(md.Filaments))
	}
	filament := md.Filaments[0]
	if filament.ID != 1 {
		t.Fatalf("Filament ID = %d, want 1", filament.ID)
	}
	if filament.Type != "PLA" {
		t.Fatalf("Filament Type = %q, want %q", filament.Type, "PLA")
	}
	if filament.Color != "#00FF00" {
		t.Fatalf("Filament Color = %q, want %q", filament.Color, "#00FF00")
	}
	if filament.UsedGrams != 12.5 {
		t.Fatalf("Filament UsedGrams = %v, want 12.5", filament.UsedGrams)
	}

	thumbnail, err := r.GetThumbnail(plate.ID)
	if err != nil {
		t.Fatalf("GetThumbnail failed: %v", err)
	}
	if string(thumbnail) != "plate-thumb" {
		t.Fatalf("thumbnail = %q, want %q", string(thumbnail), "plate-thumb")
	}

	smallThumbnail, err := r.GetThumbnailSmall(plate.ID)
	if err != nil {
		t.Fatalf("GetThumbnailSmall failed: %v", err)
	}
	if string(smallThumbnail) != "plate-thumb-small" {
		t.Fatalf("small thumbnail = %q, want %q", string(smallThumbnail), "plate-thumb-small")
	}

	packageThumbnail, err := r.GetPackageThumbnail()
	if err != nil {
		t.Fatalf("GetPackageThumbnail failed: %v", err)
	}
	if string(packageThumbnail) != "package-thumb" {
		t.Fatalf("package thumbnail = %q, want %q", string(packageThumbnail), "package-thumb")
	}
}

func writeTestProjectFile(t *testing.T) string {
	t.Helper()

	filename := filepath.Join(t.TempDir(), "fixture.gcode.3mf")
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	zipWriter := zip.NewWriter(file)
	writeZipEntry(t, zipWriter, "Metadata/slice_info.config", sliceInfoFixture)
	writeZipEntry(t, zipWriter, "Metadata/plate_1.json", `{"plate_name":"Fixture Plate"}`)
	writeZipEntry(t, zipWriter, "Metadata/plate_1.png", "plate-thumb")
	writeZipEntry(t, zipWriter, "Metadata/plate_1_small.png", "plate-thumb-small")
	writeZipEntry(t, zipWriter, "Metadata/thumbnail.png", "package-thumb")
	writeZipEntry(t, zipWriter, "3D/3dmodel.model", modelFixture)

	if err := zipWriter.Close(); err != nil {
		t.Fatalf("zip close failed: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("file close failed: %v", err)
	}

	return filename
}

func writeZipEntry(t *testing.T, zipWriter *zip.Writer, name, contents string) {
	t.Helper()

	writer, err := zipWriter.Create(name)
	if err != nil {
		t.Fatalf("Create(%q) failed: %v", name, err)
	}
	if _, err := writer.Write([]byte(contents)); err != nil {
		t.Fatalf("Write(%q) failed: %v", name, err)
	}
}

const sliceInfoFixture = `<?xml version="1.0" encoding="UTF-8"?>
<config>
  <plate>
    <metadata key="index" value="1"></metadata>
    <filament id="1" type="PLA" color="#00FF00" used_g="12.5"></filament>
  </plate>
</config>
`

const modelFixture = `<?xml version="1.0" encoding="UTF-8"?>
<model unit="millimeter" xml:lang="en-US" xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02">
  <metadata name="Title">Fixture Cube</metadata>
  <metadata name="Designer">Codex</metadata>
  <metadata name="Description">Self-contained 3MF test fixture</metadata>
</model>
`
