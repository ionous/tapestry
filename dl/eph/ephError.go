package eph

import "github.com/ionous/errutil"

// fix? consider moving domain error to catalog processing internals ( and removing explicit external use )
// would need to add "domain" to conflict for that. ( which it should probably have anyway )
type domainError struct {
	Domain string
	Err    error
}

func (n domainError) Error() string {
	return errutil.Sprintf("%v in domain %q", n.Err, n.Domain)
}
func (n domainError) Unwrap() error {
	return n.Err
}
