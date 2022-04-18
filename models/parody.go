// Code generated by SQLBoiler 4.8.6 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Parody is an object representing the database table.
type Parody struct {
	ID   int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	Slug string `boil:"slug" json:"slug" toml:"slug" yaml:"slug"`
	Name string `boil:"name" json:"name" toml:"name" yaml:"name"`

	R *parodyR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L parodyL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ParodyColumns = struct {
	ID   string
	Slug string
	Name string
}{
	ID:   "id",
	Slug: "slug",
	Name: "name",
}

var ParodyTableColumns = struct {
	ID   string
	Slug string
	Name string
}{
	ID:   "parody.id",
	Slug: "parody.slug",
	Name: "parody.name",
}

// Generated where

var ParodyWhere = struct {
	ID   whereHelperint64
	Slug whereHelperstring
	Name whereHelperstring
}{
	ID:   whereHelperint64{field: "\"parody\".\"id\""},
	Slug: whereHelperstring{field: "\"parody\".\"slug\""},
	Name: whereHelperstring{field: "\"parody\".\"name\""},
}

// ParodyRels is where relationship names are stored.
var ParodyRels = struct {
	Archives string
}{
	Archives: "Archives",
}

// parodyR is where relationships are stored.
type parodyR struct {
	Archives ArchiveSlice `boil:"Archives" json:"Archives" toml:"Archives" yaml:"Archives"`
}

// NewStruct creates a new relationship struct
func (*parodyR) NewStruct() *parodyR {
	return &parodyR{}
}

// parodyL is where Load methods for each relationship are stored.
type parodyL struct{}

var (
	parodyAllColumns            = []string{"id", "slug", "name"}
	parodyColumnsWithoutDefault = []string{}
	parodyColumnsWithDefault    = []string{"id", "slug", "name"}
	parodyPrimaryKeyColumns     = []string{"id"}
	parodyGeneratedColumns      = []string{}
)

type (
	// ParodySlice is an alias for a slice of pointers to Parody.
	// This should almost always be used instead of []Parody.
	ParodySlice []*Parody
	// ParodyHook is the signature for custom Parody hook methods
	ParodyHook func(boil.Executor, *Parody) error

	parodyQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	parodyType                 = reflect.TypeOf(&Parody{})
	parodyMapping              = queries.MakeStructMapping(parodyType)
	parodyPrimaryKeyMapping, _ = queries.BindMapping(parodyType, parodyMapping, parodyPrimaryKeyColumns)
	parodyInsertCacheMut       sync.RWMutex
	parodyInsertCache          = make(map[string]insertCache)
	parodyUpdateCacheMut       sync.RWMutex
	parodyUpdateCache          = make(map[string]updateCache)
	parodyUpsertCacheMut       sync.RWMutex
	parodyUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var parodyAfterSelectHooks []ParodyHook

var parodyBeforeInsertHooks []ParodyHook
var parodyAfterInsertHooks []ParodyHook

var parodyBeforeUpdateHooks []ParodyHook
var parodyAfterUpdateHooks []ParodyHook

var parodyBeforeDeleteHooks []ParodyHook
var parodyAfterDeleteHooks []ParodyHook

var parodyBeforeUpsertHooks []ParodyHook
var parodyAfterUpsertHooks []ParodyHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Parody) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range parodyAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Parody) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range parodyBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Parody) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range parodyAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Parody) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range parodyBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Parody) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range parodyAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Parody) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range parodyBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Parody) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range parodyAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Parody) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range parodyBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Parody) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range parodyAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddParodyHook registers your hook function for all future operations.
func AddParodyHook(hookPoint boil.HookPoint, parodyHook ParodyHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		parodyAfterSelectHooks = append(parodyAfterSelectHooks, parodyHook)
	case boil.BeforeInsertHook:
		parodyBeforeInsertHooks = append(parodyBeforeInsertHooks, parodyHook)
	case boil.AfterInsertHook:
		parodyAfterInsertHooks = append(parodyAfterInsertHooks, parodyHook)
	case boil.BeforeUpdateHook:
		parodyBeforeUpdateHooks = append(parodyBeforeUpdateHooks, parodyHook)
	case boil.AfterUpdateHook:
		parodyAfterUpdateHooks = append(parodyAfterUpdateHooks, parodyHook)
	case boil.BeforeDeleteHook:
		parodyBeforeDeleteHooks = append(parodyBeforeDeleteHooks, parodyHook)
	case boil.AfterDeleteHook:
		parodyAfterDeleteHooks = append(parodyAfterDeleteHooks, parodyHook)
	case boil.BeforeUpsertHook:
		parodyBeforeUpsertHooks = append(parodyBeforeUpsertHooks, parodyHook)
	case boil.AfterUpsertHook:
		parodyAfterUpsertHooks = append(parodyAfterUpsertHooks, parodyHook)
	}
}

// OneG returns a single parody record from the query using the global executor.
func (q parodyQuery) OneG() (*Parody, error) {
	return q.One(boil.GetDB())
}

// One returns a single parody record from the query.
func (q parodyQuery) One(exec boil.Executor) (*Parody, error) {
	o := &Parody{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for parody")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all Parody records from the query using the global executor.
func (q parodyQuery) AllG() (ParodySlice, error) {
	return q.All(boil.GetDB())
}

// All returns all Parody records from the query.
func (q parodyQuery) All(exec boil.Executor) (ParodySlice, error) {
	var o []*Parody

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Parody slice")
	}

	if len(parodyAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all Parody records in the query, and panics on error.
func (q parodyQuery) CountG() (int64, error) {
	return q.Count(boil.GetDB())
}

// Count returns the count of all Parody records in the query.
func (q parodyQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count parody rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q parodyQuery) ExistsG() (bool, error) {
	return q.Exists(boil.GetDB())
}

// Exists checks if the row exists in the table.
func (q parodyQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if parody exists")
	}

	return count > 0, nil
}

// Archives retrieves all the archive's Archives with an executor.
func (o *Parody) Archives(mods ...qm.QueryMod) archiveQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.InnerJoin("\"archive_parodies\" on \"archive\".\"id\" = \"archive_parodies\".\"archive_id\""),
		qm.Where("\"archive_parodies\".\"parody_id\"=?", o.ID),
	)

	query := Archives(queryMods...)
	queries.SetFrom(query.Query, "\"archive\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"archive\".*"})
	}

	return query
}

// LoadArchives allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (parodyL) LoadArchives(e boil.Executor, singular bool, maybeParody interface{}, mods queries.Applicator) error {
	var slice []*Parody
	var object *Parody

	if singular {
		object = maybeParody.(*Parody)
	} else {
		slice = *maybeParody.(*[]*Parody)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &parodyR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &parodyR{}
			}

			for _, a := range args {
				if a == obj.ID {
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
		qm.Select("\"archive\".id, \"archive\".path, \"archive\".created_at, \"archive\".updated_at, \"archive\".published_at, \"archive\".title, \"archive\".slug, \"archive\".pages, \"archive\".size, \"archive\".submission_id, \"a\".\"parody_id\""),
		qm.From("\"archive\""),
		qm.InnerJoin("\"archive_parodies\" as \"a\" on \"archive\".\"id\" = \"a\".\"archive_id\""),
		qm.WhereIn("\"a\".\"parody_id\" in ?", args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load archive")
	}

	var resultSlice []*Archive

	var localJoinCols []int64
	for results.Next() {
		one := new(Archive)
		var localJoinCol int64

		err = results.Scan(&one.ID, &one.Path, &one.CreatedAt, &one.UpdatedAt, &one.PublishedAt, &one.Title, &one.Slug, &one.Pages, &one.Size, &one.SubmissionID, &localJoinCol)
		if err != nil {
			return errors.Wrap(err, "failed to scan eager loaded results for archive")
		}
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice archive")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
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
			foreign.R.Parodies = append(foreign.R.Parodies, object)
		}
		return nil
	}

	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			if local.ID == localJoinCol {
				local.R.Archives = append(local.R.Archives, foreign)
				if foreign.R == nil {
					foreign.R = &archiveR{}
				}
				foreign.R.Parodies = append(foreign.R.Parodies, local)
				break
			}
		}
	}

	return nil
}

// AddArchivesG adds the given related objects to the existing relationships
// of the parody, optionally inserting them as new records.
// Appends related to o.R.Archives.
// Sets related.R.Parodies appropriately.
// Uses the global database handle.
func (o *Parody) AddArchivesG(insert bool, related ...*Archive) error {
	return o.AddArchives(boil.GetDB(), insert, related...)
}

// AddArchives adds the given related objects to the existing relationships
// of the parody, optionally inserting them as new records.
// Appends related to o.R.Archives.
// Sets related.R.Parodies appropriately.
func (o *Parody) AddArchives(exec boil.Executor, insert bool, related ...*Archive) error {
	var err error
	for _, rel := range related {
		if insert {
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}
	}

	for _, rel := range related {
		query := "insert into \"archive_parodies\" (\"parody_id\", \"archive_id\") values ($1, $2)"
		values := []interface{}{o.ID, rel.ID}

		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, query)
			fmt.Fprintln(boil.DebugWriter, values)
		}
		_, err = exec.Exec(query, values...)
		if err != nil {
			return errors.Wrap(err, "failed to insert into join table")
		}
	}
	if o.R == nil {
		o.R = &parodyR{
			Archives: related,
		}
	} else {
		o.R.Archives = append(o.R.Archives, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &archiveR{
				Parodies: ParodySlice{o},
			}
		} else {
			rel.R.Parodies = append(rel.R.Parodies, o)
		}
	}
	return nil
}

// SetArchivesG removes all previously related items of the
// parody replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parodies's Archives accordingly.
// Replaces o.R.Archives with related.
// Sets related.R.Parodies's Archives accordingly.
// Uses the global database handle.
func (o *Parody) SetArchivesG(insert bool, related ...*Archive) error {
	return o.SetArchives(boil.GetDB(), insert, related...)
}

// SetArchives removes all previously related items of the
// parody replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parodies's Archives accordingly.
// Replaces o.R.Archives with related.
// Sets related.R.Parodies's Archives accordingly.
func (o *Parody) SetArchives(exec boil.Executor, insert bool, related ...*Archive) error {
	query := "delete from \"archive_parodies\" where \"parody_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	removeArchivesFromParodiesSlice(o, related)
	if o.R != nil {
		o.R.Archives = nil
	}
	return o.AddArchives(exec, insert, related...)
}

// RemoveArchivesG relationships from objects passed in.
// Removes related items from R.Archives (uses pointer comparison, removal does not keep order)
// Sets related.R.Parodies.
// Uses the global database handle.
func (o *Parody) RemoveArchivesG(related ...*Archive) error {
	return o.RemoveArchives(boil.GetDB(), related...)
}

// RemoveArchives relationships from objects passed in.
// Removes related items from R.Archives (uses pointer comparison, removal does not keep order)
// Sets related.R.Parodies.
func (o *Parody) RemoveArchives(exec boil.Executor, related ...*Archive) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	query := fmt.Sprintf(
		"delete from \"archive_parodies\" where \"parody_id\" = $1 and \"archive_id\" in (%s)",
		strmangle.Placeholders(dialect.UseIndexPlaceholders, len(related), 2, 1),
	)
	values := []interface{}{o.ID}
	for _, rel := range related {
		values = append(values, rel.ID)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err = exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}
	removeArchivesFromParodiesSlice(o, related)
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

func removeArchivesFromParodiesSlice(o *Parody, related []*Archive) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.Parodies {
			if o.ID != ri.ID {
				continue
			}

			ln := len(rel.R.Parodies)
			if ln > 1 && i < ln-1 {
				rel.R.Parodies[i] = rel.R.Parodies[ln-1]
			}
			rel.R.Parodies = rel.R.Parodies[:ln-1]
			break
		}
	}
}

// Parodies retrieves all the records using an executor.
func Parodies(mods ...qm.QueryMod) parodyQuery {
	mods = append(mods, qm.From("\"parody\""))
	return parodyQuery{NewQuery(mods...)}
}

// FindParodyG retrieves a single record by ID.
func FindParodyG(iD int64, selectCols ...string) (*Parody, error) {
	return FindParody(boil.GetDB(), iD, selectCols...)
}

// FindParody retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindParody(exec boil.Executor, iD int64, selectCols ...string) (*Parody, error) {
	parodyObj := &Parody{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"parody\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, parodyObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from parody")
	}

	if err = parodyObj.doAfterSelectHooks(exec); err != nil {
		return parodyObj, err
	}

	return parodyObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Parody) InsertG(columns boil.Columns) error {
	return o.Insert(boil.GetDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Parody) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no parody provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(parodyColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	parodyInsertCacheMut.RLock()
	cache, cached := parodyInsertCache[key]
	parodyInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			parodyAllColumns,
			parodyColumnsWithDefault,
			parodyColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(parodyType, parodyMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(parodyType, parodyMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"parody\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"parody\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into parody")
	}

	if !cached {
		parodyInsertCacheMut.Lock()
		parodyInsertCache[key] = cache
		parodyInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Parody record using the global executor.
// See Update for more documentation.
func (o *Parody) UpdateG(columns boil.Columns) error {
	return o.Update(boil.GetDB(), columns)
}

// Update uses an executor to update the Parody.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Parody) Update(exec boil.Executor, columns boil.Columns) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(columns, nil)
	parodyUpdateCacheMut.RLock()
	cache, cached := parodyUpdateCache[key]
	parodyUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			parodyAllColumns,
			parodyPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update parody, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"parody\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, parodyPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(parodyType, parodyMapping, append(wl, parodyPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update parody row")
	}

	if !cached {
		parodyUpdateCacheMut.Lock()
		parodyUpdateCache[key] = cache
		parodyUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q parodyQuery) UpdateAllG(cols M) error {
	return q.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q parodyQuery) UpdateAll(exec boil.Executor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for parody")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ParodySlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ParodySlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), parodyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"parody\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, parodyPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in parody slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Parody) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Parody) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no parody provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(parodyColumnsWithDefault, o)

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

	parodyUpsertCacheMut.RLock()
	cache, cached := parodyUpsertCache[key]
	parodyUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			parodyAllColumns,
			parodyColumnsWithDefault,
			parodyColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			parodyAllColumns,
			parodyPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert parody, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(parodyPrimaryKeyColumns))
			copy(conflict, parodyPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"parody\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(parodyType, parodyMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(parodyType, parodyMapping, ret)
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
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert parody")
	}

	if !cached {
		parodyUpsertCacheMut.Lock()
		parodyUpsertCache[key] = cache
		parodyUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteG deletes a single Parody record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Parody) DeleteG() error {
	return o.Delete(boil.GetDB())
}

// Delete deletes a single Parody record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Parody) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Parody provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), parodyPrimaryKeyMapping)
	sql := "DELETE FROM \"parody\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from parody")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

func (q parodyQuery) DeleteAllG() error {
	return q.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all matching rows.
func (q parodyQuery) DeleteAll(exec boil.Executor) error {
	if q.Query == nil {
		return errors.New("models: no parodyQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from parody")
	}

	return nil
}

// DeleteAllG deletes all rows in the slice.
func (o ParodySlice) DeleteAllG() error {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ParodySlice) DeleteAll(exec boil.Executor) error {
	if len(o) == 0 {
		return nil
	}

	if len(parodyBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), parodyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"parody\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, parodyPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from parody slice")
	}

	if len(parodyAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Parody) ReloadG() error {
	if o == nil {
		return errors.New("models: no Parody provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Parody) Reload(exec boil.Executor) error {
	ret, err := FindParody(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ParodySlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty ParodySlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ParodySlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ParodySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), parodyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"parody\".* FROM \"parody\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, parodyPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ParodySlice")
	}

	*o = slice

	return nil
}

// ParodyExistsG checks if the Parody row exists.
func ParodyExistsG(iD int64) (bool, error) {
	return ParodyExists(boil.GetDB(), iD)
}

// ParodyExists checks if the Parody row exists.
func ParodyExists(exec boil.Executor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"parody\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if parody exists")
	}

	return exists, nil
}
