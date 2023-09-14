package api

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"time"

	"github.com/s12chung/text2anki/db/pkg/db"
	"github.com/s12chung/text2anki/pkg/anki"
	"github.com/s12chung/text2anki/pkg/util/archive/ziputil"
	"github.com/s12chung/text2anki/pkg/util/httptyped"
	"github.com/s12chung/text2anki/pkg/util/jhttp"
	"github.com/s12chung/text2anki/pkg/util/logg"
)

func init() {
	httptyped.RegisterType(db.Note{})
}

// NotesIndex shows lists all the notes
func (rs Routes) NotesIndex(_ *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	return jhttp.ReturnModelOr500(func() (any, error) {
		return txQs.NotesIndex(txQs.Ctx())
	})
}

// NoteCreate creates a new note
func (rs Routes) NoteCreate(r *http.Request, txQs db.TxQs) (any, *jhttp.HTTPError) {
	req := db.NoteCreateParams{}
	if httpErr := extractAndValidate(r, &req); httpErr != nil {
		return nil, httpErr
	}
	return jhttp.ReturnModelOr500(func() (any, error) {
		return txQs.NoteCreate(txQs.Ctx(), req)
	})
}

// NotesDownload downloads the not downloaded notes
func (rs Routes) NotesDownload(w http.ResponseWriter, r *http.Request) {
	httpErr := rs.runOr500(r, func(r *http.Request, tx db.TxQs) error {
		notes, err := rs.downloadedAnkiNotes(tx)
		if err != nil {
			return err
		}
		id, b, err := rs.exportNotes(notes)
		if err != nil {
			return err
		}

		updated, err := tx.NotesUpdateDownloaded(tx.Ctx())
		if err != nil {
			return err
		}
		rs.Log.LogAttrs(tx.Ctx(), slog.LevelInfo, fmt.Sprintf("updated count: %v", updated), logg.RequestAttrs(r)...)

		filename := "text2anki-" + id + ".zip"
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(b))
		return nil
	})
	if httpErr != nil {
		jhttp.LogAndRespondError(w, r, httpErr)
		return
	}
}

func (rs Routes) downloadedAnkiNotes(tx db.TxQs) ([]anki.Note, error) {
	notes, err := tx.NotesDownloaded(tx.Ctx())
	if err != nil {
		return nil, err
	}
	ankiNotes, err := db.AnkiNotes(notes)
	if err != nil {
		return nil, err
	}
	if err := rs.SoundSetter.SetSound(tx.Ctx(), ankiNotes); err != nil {
		return nil, err
	}
	return ankiNotes, nil
}

func (rs Routes) exportNotes(notes []anki.Note) (string, []byte, error) {
	id, err := rs.UUIDGenerator.Generate()
	if err != nil {
		return "", nil, err
	}
	exportDir := path.Join(rs.CacheDir, "NotesDownload", id)
	if err := anki.ExportFiles(exportDir, notes); err != nil {
		return "", nil, err
	}
	b, err := ziputil.ZipDir(exportDir)
	return id, b, err
}
