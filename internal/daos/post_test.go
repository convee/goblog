package daos

import (
	"reflect"
	"testing"
)

func TestPostQueryBuilderBindsEachSearchID(t *testing.T) {
	qb := NewPostQueryBuilder().
		WithCategory("4").
		WithKeyword("go", map[string][]string{
			"ids":          {"1", "2"},
			"category_ids": {"3", "5"},
			"tag_ids":      {"8"},
		}).
		WithStatus(1).
		WithPagination(2, 20)

	wantSQL := "select id,title from post where category_id=? and (title like ? escape '|' or description like ? escape '|' or id in (?,?) or category_id in (?,?) or JSON_CONTAINS(tag_ids,?)) and status=? order by is_top desc,updated_at desc limit 20,20"
	if got := qb.BuildQuery("id,title"); got != wantSQL {
		t.Fatalf("BuildQuery()\n got: %s\nwant: %s", got, wantSQL)
	}
	wantArgs := []interface{}{"4", "%go%", "%go%", "1", "2", "3", "5", "8", 1}
	if got := qb.GetArgs(); !reflect.DeepEqual(got, wantArgs) {
		t.Fatalf("GetArgs() = %#v, want %#v", got, wantArgs)
	}
}

func TestPostQueryBuilderDoesNotMutateReturnedArgs(t *testing.T) {
	qb := NewPostQueryBuilder().WithStatus(1)
	args := qb.GetArgs()
	args[0] = 0
	if got := qb.GetArgs()[0]; got != 1 {
		t.Fatalf("GetArgs exposed internal state: %v", got)
	}
}

func TestLikePatternEscapesUserWildcards(t *testing.T) {
	got := likePattern("100%_go|安全")
	want := "%100|%|_go||安全%"
	if got != want {
		t.Fatalf("likePattern() = %q, want %q", got, want)
	}
}
