package category

import "errors"

var (
	ErrQueryFailed       = errors.New("category: query execution failed")
	ErrInsertFailed      = errors.New("category: failed to insert category")
	ErrScanFailed        = errors.New("category: failed to scan row")
	ErrUUIDParseFailed   = errors.New("category: failed to parse UUID")
	ErrTransactionBegin  = errors.New("category: failed to begin transaction")
	ErrTransactionCommit = errors.New("category: failed to commit transaction")
	ErrPrepareStmtFailed = errors.New("category: failed to prepare statement")
	ErrInsertRelation    = errors.New("category: failed to insert category-post relation")
	ErrExistsCheckFailed = errors.New("category: failed to check existence")
)
