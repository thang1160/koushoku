// Code generated by SQLBoiler 4.11.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Submission is an object representing the database table.
type Submission struct {
	ID         int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	CreatedAt  time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt  time.Time   `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	Name       string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Submitter  null.String `boil:"submitter" json:"submitter,omitempty" toml:"submitter" yaml:"submitter,omitempty"`
	Content    string      `boil:"content" json:"content" toml:"content" yaml:"content"`
	Notes      null.String `boil:"notes" json:"notes,omitempty" toml:"notes" yaml:"notes,omitempty"`
	AcceptedAt null.Time   `boil:"accepted_at" json:"accepted_at,omitempty" toml:"accepted_at" yaml:"accepted_at,omitempty"`
	RejectedAt null.Time   `boil:"rejected_at" json:"rejected_at,omitempty" toml:"rejected_at" yaml:"rejected_at,omitempty"`
	Accepted   bool        `boil:"accepted" json:"accepted" toml:"accepted" yaml:"accepted"`
	Rejected   bool        `boil:"rejected" json:"rejected" toml:"rejected" yaml:"rejected"`

	R *submissionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L submissionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SubmissionColumns = struct {
	ID         string
	CreatedAt  string
	UpdatedAt  string
	Name       string
	Submitter  string
	Content    string
	Notes      string
	AcceptedAt string
	RejectedAt string
	Accepted   string
	Rejected   string
}{
	ID:         "id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	Name:       "name",
	Submitter:  "submitter",
	Content:    "content",
	Notes:      "notes",
	AcceptedAt: "accepted_at",
	RejectedAt: "rejected_at",
	Accepted:   "accepted",
	Rejected:   "rejected",
}

var SubmissionTableColumns = struct {
	ID         string
	CreatedAt  string
	UpdatedAt  string
	Name       string
	Submitter  string
	Content    string
	Notes      string
	AcceptedAt string
	RejectedAt string
	Accepted   string
	Rejected   string
}{
	ID:         "submission.id",
	CreatedAt:  "submission.created_at",
	UpdatedAt:  "submission.updated_at",
	Name:       "submission.name",
	Submitter:  "submission.submitter",
	Content:    "submission.content",
	Notes:      "submission.notes",
	AcceptedAt: "submission.accepted_at",
	RejectedAt: "submission.rejected_at",
	Accepted:   "submission.accepted",
	Rejected:   "submission.rejected",
}

// Generated where

type whereHelpernull_String struct{ field string }

func (w whereHelpernull_String) EQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_String) NEQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_String) LT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_String) LTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_String) GT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_String) GTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_String) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_String) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var SubmissionWhere = struct {
	ID         whereHelperint64
	CreatedAt  whereHelpertime_Time
	UpdatedAt  whereHelpertime_Time
	Name       whereHelperstring
	Submitter  whereHelpernull_String
	Content    whereHelperstring
	Notes      whereHelpernull_String
	AcceptedAt whereHelpernull_Time
	RejectedAt whereHelpernull_Time
	Accepted   whereHelperbool
	Rejected   whereHelperbool
}{
	ID:         whereHelperint64{field: "\"submission\".\"id\""},
	CreatedAt:  whereHelpertime_Time{field: "\"submission\".\"created_at\""},
	UpdatedAt:  whereHelpertime_Time{field: "\"submission\".\"updated_at\""},
	Name:       whereHelperstring{field: "\"submission\".\"name\""},
	Submitter:  whereHelpernull_String{field: "\"submission\".\"submitter\""},
	Content:    whereHelperstring{field: "\"submission\".\"content\""},
	Notes:      whereHelpernull_String{field: "\"submission\".\"notes\""},
	AcceptedAt: whereHelpernull_Time{field: "\"submission\".\"accepted_at\""},
	RejectedAt: whereHelpernull_Time{field: "\"submission\".\"rejected_at\""},
	Accepted:   whereHelperbool{field: "\"submission\".\"accepted\""},
	Rejected:   whereHelperbool{field: "\"submission\".\"rejected\""},
}

// SubmissionRels is where relationship names are stored.
var SubmissionRels = struct {
	Archives string
}{
	Archives: "Archives",
}

// submissionR is where relationships are stored.
type submissionR struct {
	Archives ArchiveSlice `boil:"Archives" json:"Archives" toml:"Archives" yaml:"Archives"`
}

// NewStruct creates a new relationship struct
func (*submissionR) NewStruct() *submissionR {
	return &submissionR{}
}

func (r *submissionR) GetArchives() ArchiveSlice {
	if r == nil {
		return nil
	}
	return r.Archives
}

// submissionL is where Load methods for each relationship are stored.
type submissionL struct{}

var (
	submissionAllColumns            = []string{"id", "created_at", "updated_at", "name", "submitter", "content", "notes", "accepted_at", "rejected_at", "accepted", "rejected"}
	submissionColumnsWithoutDefault = []string{}
	submissionColumnsWithDefault    = []string{"id", "created_at", "updated_at", "name", "submitter", "content", "notes", "accepted_at", "rejected_at", "accepted", "rejected"}
	submissionPrimaryKeyColumns     = []string{"id"}
	submissionGeneratedColumns      = []string{}
)

type (
	// SubmissionSlice is an alias for a slice of pointers to Submission.
	// This should almost always be used instead of []Submission.
	SubmissionSlice []*Submission
	// SubmissionHook is the signature for custom Submission hook methods
	SubmissionHook func(boil.Executor, *Submission) error

	submissionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	submissionType                 = reflect.TypeOf(&Submission{})
	submissionMapping              = queries.MakeStructMapping(submissionType)
	submissionPrimaryKeyMapping, _ = queries.BindMapping(submissionType, submissionMapping, submissionPrimaryKeyColumns)
	submissionInsertCacheMut       sync.RWMutex
	submissionInsertCache          = make(map[string]insertCache)
	submissionUpdateCacheMut       sync.RWMutex
	submissionUpdateCache          = make(map[string]updateCache)
	submissionUpsertCacheMut       sync.RWMutex
	submissionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var submissionAfterSelectHooks []SubmissionHook

var submissionBeforeInsertHooks []SubmissionHook
var submissionAfterInsertHooks []SubmissionHook

var submissionBeforeUpdateHooks []SubmissionHook
var submissionAfterUpdateHooks []SubmissionHook

var submissionBeforeDeleteHooks []SubmissionHook
var submissionAfterDeleteHooks []SubmissionHook

var submissionBeforeUpsertHooks []SubmissionHook
var submissionAfterUpsertHooks []SubmissionHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Submission) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range submissionAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Submission) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range submissionBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Submission) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range submissionAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Submission) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range submissionBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Submission) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range submissionAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Submission) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range submissionBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Submission) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range submissionAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Submission) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range submissionBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Submission) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range submissionAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSubmissionHook registers your hook function for all future operations.
func AddSubmissionHook(hookPoint boil.HookPoint, submissionHook SubmissionHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		submissionAfterSelectHooks = append(submissionAfterSelectHooks, submissionHook)
	case boil.BeforeInsertHook:
		submissionBeforeInsertHooks = append(submissionBeforeInsertHooks, submissionHook)
	case boil.AfterInsertHook:
		submissionAfterInsertHooks = append(submissionAfterInsertHooks, submissionHook)
	case boil.BeforeUpdateHook:
		submissionBeforeUpdateHooks = append(submissionBeforeUpdateHooks, submissionHook)
	case boil.AfterUpdateHook:
		submissionAfterUpdateHooks = append(submissionAfterUpdateHooks, submissionHook)
	case boil.BeforeDeleteHook:
		submissionBeforeDeleteHooks = append(submissionBeforeDeleteHooks, submissionHook)
	case boil.AfterDeleteHook:
		submissionAfterDeleteHooks = append(submissionAfterDeleteHooks, submissionHook)
	case boil.BeforeUpsertHook:
		submissionBeforeUpsertHooks = append(submissionBeforeUpsertHooks, submissionHook)
	case boil.AfterUpsertHook:
		submissionAfterUpsertHooks = append(submissionAfterUpsertHooks, submissionHook)
	}
}

// OneG returns a single submission record from the query using the global executor.
func (q submissionQuery) OneG() (*Submission, error) {
	return q.One(boil.GetDB())
}

// One returns a single submission record from the query.
func (q submissionQuery) One(exec boil.Executor) (*Submission, error) {
	o := &Submission{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for submission")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all Submission records from the query using the global executor.
func (q submissionQuery) AllG() (SubmissionSlice, error) {
	return q.All(boil.GetDB())
}

// All returns all Submission records from the query.
func (q submissionQuery) All(exec boil.Executor) (SubmissionSlice, error) {
	var o []*Submission

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Submission slice")
	}

	if len(submissionAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all Submission records in the query using the global executor
func (q submissionQuery) CountG() (int64, error) {
	return q.Count(boil.GetDB())
}

// Count returns the count of all Submission records in the query.
func (q submissionQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count submission rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q submissionQuery) ExistsG() (bool, error) {
	return q.Exists(boil.GetDB())
}

// Exists checks if the row exists in the table.
func (q submissionQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if submission exists")
	}

	return count > 0, nil
}

// Archives retrieves all the archive's Archives with an executor.
func (o *Submission) Archives(mods ...qm.QueryMod) archiveQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"archive\".\"submission_id\"=?", o.ID),
	)

	return Archives(queryMods...)
}

// LoadArchives allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (submissionL) LoadArchives(e boil.Executor, singular bool, maybeSubmission interface{}, mods queries.Applicator) error {
	var slice []*Submission
	var object *Submission

	if singular {
		object = maybeSubmission.(*Submission)
	} else {
		slice = *maybeSubmission.(*[]*Submission)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &submissionR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &submissionR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`archive`),
		qm.WhereIn(`archive.submission_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load archive")
	}

	var resultSlice []*Archive
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice archive")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on archive")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for archive")
	}

	if len(archiveAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Archives = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &archiveR{}
			}
			foreign.R.Submission = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.SubmissionID) {
				local.R.Archives = append(local.R.Archives, foreign)
				if foreign.R == nil {
					foreign.R = &archiveR{}
				}
				foreign.R.Submission = local
				break
			}
		}
	}

	return nil
}

// AddArchivesG adds the given related objects to the existing relationships
// of the submission, optionally inserting them as new records.
// Appends related to o.R.Archives.
// Sets related.R.Submission appropriately.
// Uses the global database handle.
func (o *Submission) AddArchivesG(insert bool, related ...*Archive) error {
	return o.AddArchives(boil.GetDB(), insert, related...)
}

// AddArchives adds the given related objects to the existing relationships
// of the submission, optionally inserting them as new records.
// Appends related to o.R.Archives.
// Sets related.R.Submission appropriately.
func (o *Submission) AddArchives(exec boil.Executor, insert bool, related ...*Archive) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.SubmissionID, o.ID)
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"archive\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"submission_id"}),
				strmangle.WhereClause("\"", "\"", 2, archivePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.SubmissionID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &submissionR{
			Archives: related,
		}
	} else {
		o.R.Archives = append(o.R.Archives, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &archiveR{
				Submission: o,
			}
		} else {
			rel.R.Submission = o
		}
	}
	return nil
}

// SetArchivesG removes all previously related items of the
// submission replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Submission's Archives accordingly.
// Replaces o.R.Archives with related.
// Sets related.R.Submission's Archives accordingly.
// Uses the global database handle.
func (o *Submission) SetArchivesG(insert bool, related ...*Archive) error {
	return o.SetArchives(boil.GetDB(), insert, related...)
}

// SetArchives removes all previously related items of the
// submission replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Submission's Archives accordingly.
// Replaces o.R.Archives with related.
// Sets related.R.Submission's Archives accordingly.
func (o *Submission) SetArchives(exec boil.Executor, insert bool, related ...*Archive) error {
	query := "update \"archive\" set \"submission_id\" = null where \"submission_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.Archives {
			queries.SetScanner(&rel.SubmissionID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Submission = nil
		}
		o.R.Archives = nil
	}

	return o.AddArchives(exec, insert, related...)
}

// RemoveArchivesG relationships from objects passed in.
// Removes related items from R.Archives (uses pointer comparison, removal does not keep order)
// Sets related.R.Submission.
// Uses the global database handle.
func (o *Submission) RemoveArchivesG(related ...*Archive) error {
	return o.RemoveArchives(boil.GetDB(), related...)
}

// RemoveArchives relationships from objects passed in.
// Removes related items from R.Archives (uses pointer comparison, removal does not keep order)
// Sets related.R.Submission.
func (o *Submission) RemoveArchives(exec boil.Executor, related ...*Archive) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.SubmissionID, nil)
		if rel.R != nil {
			rel.R.Submission = nil
		}
		if err = rel.Update(exec, boil.Whitelist("submission_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Archives {
			if rel != ri {
				continue
			}

			ln := len(o.R.Archives)
			if ln > 1 && i < ln-1 {
				o.R.Archives[i] = o.R.Archives[ln-1]
			}
			o.R.Archives = o.R.Archives[:ln-1]
			break
		}
	}

	return nil
}

// Submissions retrieves all the records using an executor.
func Submissions(mods ...qm.QueryMod) submissionQuery {
	mods = append(mods, qm.From("\"submission\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"submission\".*"})
	}

	return submissionQuery{q}
}

// FindSubmissionG retrieves a single record by ID.
func FindSubmissionG(iD int64, selectCols ...string) (*Submission, error) {
	return FindSubmission(boil.GetDB(), iD, selectCols...)
}

// FindSubmission retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSubmission(exec boil.Executor, iD int64, selectCols ...string) (*Submission, error) {
	submissionObj := &Submission{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"submission\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, submissionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from submission")
	}

	if err = submissionObj.doAfterSelectHooks(exec); err != nil {
		return submissionObj, err
	}

	return submissionObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Submission) InsertG(columns boil.Columns) error {
	return o.Insert(boil.GetDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Submission) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no submission provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(submissionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	submissionInsertCacheMut.RLock()
	cache, cached := submissionInsertCache[key]
	submissionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			submissionAllColumns,
			submissionColumnsWithDefault,
			submissionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(submissionType, submissionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(submissionType, submissionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"submission\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"submission\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into submission")
	}

	if !cached {
		submissionInsertCacheMut.Lock()
		submissionInsertCache[key] = cache
		submissionInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Submission record using the global executor.
// See Update for more documentation.
func (o *Submission) UpdateG(columns boil.Columns) error {
	return o.Update(boil.GetDB(), columns)
}

// Update uses an executor to update the Submission.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Submission) Update(exec boil.Executor, columns boil.Columns) error {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(columns, nil)
	submissionUpdateCacheMut.RLock()
	cache, cached := submissionUpdateCache[key]
	submissionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			submissionAllColumns,
			submissionPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update submission, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"submission\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, submissionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(submissionType, submissionMapping, append(wl, submissionPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update submission row")
	}

	if !cached {
		submissionUpdateCacheMut.Lock()
		submissionUpdateCache[key] = cache
		submissionUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q submissionQuery) UpdateAllG(cols M) error {
	return q.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q submissionQuery) UpdateAll(exec boil.Executor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for submission")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SubmissionSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SubmissionSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), submissionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"submission\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, submissionPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in submission slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Submission) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Submission) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no submission provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	o.UpdatedAt = currTime

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(submissionColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	submissionUpsertCacheMut.RLock()
	cache, cached := submissionUpsertCache[key]
	submissionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			submissionAllColumns,
			submissionColumnsWithDefault,
			submissionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			submissionAllColumns,
			submissionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert submission, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(submissionPrimaryKeyColumns))
			copy(conflict, submissionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"submission\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(submissionType, submissionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(submissionType, submissionMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert submission")
	}

	if !cached {
		submissionUpsertCacheMut.Lock()
		submissionUpsertCache[key] = cache
		submissionUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteG deletes a single Submission record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Submission) DeleteG() error {
	return o.Delete(boil.GetDB())
}

// Delete deletes a single Submission record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Submission) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Submission provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), submissionPrimaryKeyMapping)
	sql := "DELETE FROM \"submission\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from submission")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

func (q submissionQuery) DeleteAllG() error {
	return q.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all matching rows.
func (q submissionQuery) DeleteAll(exec boil.Executor) error {
	if q.Query == nil {
		return errors.New("models: no submissionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from submission")
	}

	return nil
}

// DeleteAllG deletes all rows in the slice.
func (o SubmissionSlice) DeleteAllG() error {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SubmissionSlice) DeleteAll(exec boil.Executor) error {
	if len(o) == 0 {
		return nil
	}

	if len(submissionBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), submissionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"submission\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, submissionPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from submission slice")
	}

	if len(submissionAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Submission) ReloadG() error {
	if o == nil {
		return errors.New("models: no Submission provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Submission) Reload(exec boil.Executor) error {
	ret, err := FindSubmission(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SubmissionSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty SubmissionSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SubmissionSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SubmissionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), submissionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"submission\".* FROM \"submission\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, submissionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SubmissionSlice")
	}

	*o = slice

	return nil
}

// SubmissionExistsG checks if the Submission row exists.
func SubmissionExistsG(iD int64) (bool, error) {
	return SubmissionExists(boil.GetDB(), iD)
}

// SubmissionExists checks if the Submission row exists.
func SubmissionExists(exec boil.Executor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"submission\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if submission exists")
	}

	return exists, nil
}
