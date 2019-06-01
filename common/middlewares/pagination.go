package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dheerajgopi/todo-api/common"
	todoErr "github.com/dheerajgopi/todo-api/common/error"
)

// Paginator gets pagination data from request and set it in request context
func Paginator(fieldTypeMapping map[string]common.FieldType) MiddlewareFunc {
	return func(f common.HandlerFunc) common.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request, reqCtx *common.RequestContext) (int, interface{}, *todoErr.APIError) {
			reqParams := req.URL.Query()

			var pagesize int
			var pagenumber int
			var err error

			// set page size if provided, else set default value
			pageSizeParams, ok := reqParams["limit"]

			if ok {
				pagesize, err = strconv.Atoi(pageSizeParams[0])

				if err != nil || pagesize <= 0 {
					apiError := todoErr.NewAPIError("Invalid pagination", &todoErr.APIErrorBody{
						Message: "Limit should be a number greater than zero",
						Target:  "limit",
					})

					return http.StatusBadRequest, nil, apiError
				}
			} else {
				pagesize = 20
			}

			// set page number if provided, else set default value
			pageNumberParams, ok := reqParams["page"]

			if ok {
				pagenumber, err = strconv.Atoi(pageNumberParams[0])

				if err != nil || pagenumber < 0 {
					apiError := todoErr.NewAPIError("Invalid pagination", &todoErr.APIErrorBody{
						Message: "Page should be non-negative number",
						Target:  "page",
					})

					return http.StatusBadRequest, nil, apiError
				}
			} else {
				pagenumber = 0
			}

			page := &common.Page{
				Limit:  pagesize,
				Offset: int64(pagesize * pagenumber),
			}

			cursorFields := make([]*common.Sort, 0)

			// set cursor if provided, else set default value
			cursorParams, cursorOk := reqParams["cursor"]
			sortParams, sortOk := reqParams["sort"]

			if cursorOk && sortOk {
				apiError := todoErr.NewAPIError("Cannot provide both sort and cursor", &todoErr.APIErrorBody{
					Message: "Not required",
					Target:  "sort",
				})

				return http.StatusBadRequest, nil, apiError
			}

			if cursorOk {
				cursor := cursorParams[0]
				decodedBytes, err := base64.URLEncoding.DecodeString(cursor)

				if err != nil {
					return invalidCursorError()
				}

				if err := json.Unmarshal(decodedBytes, &cursorFields); err != nil {
					return invalidCursorError()
				}

				for _, eachSort := range cursorFields {
					fieldType, ok := fieldTypeMapping[eachSort.Field]

					if ok && !eachSort.ValidateLastVal(fieldType) {
						return invalidCursorError()
					}
				}

				page.Cursor = cursorFields

				a, _ := json.Marshal(&cursorFields)
				fmt.Println(string(a))
			}

			sortFields, apiError := parseSortParams(sortParams)

			if apiError != nil {
				return http.StatusBadRequest, nil, apiError
			}

			page.Sort = sortFields

			reqCtx.Page = page

			return f(res, req, reqCtx)
		}
	}
}

func parseSortParams(sortParams []string) ([]*common.Sort, *todoErr.APIError) {
	sortFields := make([]*common.Sort, 0)

	if len(sortParams) == 0 {
		return sortFields, nil
	}

	for _, eachSortParam := range sortParams {
		eachSort, apiError := parseSort(eachSortParam)

		if apiError != nil {
			return nil, apiError
		}

		sortFields = append(sortFields, eachSort...)
	}

	if apiError := hasDuplicateSortFields(sortFields); apiError != nil {
		return nil, apiError
	}

	return sortFields, nil
}

func parseSort(sort string) ([]*common.Sort, *todoErr.APIError) {
	sort = strings.ToLower(strings.TrimSpace(sort))
	sortFields := make([]*common.Sort, 0)

	if len(sort) == 0 {
		return sortFields, nil
	}

	sorts := strings.Split(sort, ",")

	for _, eachSortParam := range sorts {
		eachSortParam = strings.TrimSpace(eachSortParam)

		if len(eachSortParam) == 0 {
			continue
		}

		fieldAndDirectionPair := strings.Split(eachSortParam, ":")
		fieldAndDirectionLen := len(fieldAndDirectionPair)

		if fieldAndDirectionLen > 2 {
			apiError := todoErr.NewAPIError("Invalid sort", &todoErr.APIErrorBody{
				Message: "Invalid sort parameter",
				Target:  "sort",
			})

			return nil, apiError
		}

		if fieldAndDirectionLen == 2 {
			direction := strings.TrimSpace(fieldAndDirectionPair[1])

			if !(direction == "asc" || direction == "desc") {
				apiError := todoErr.NewAPIError("Invalid sort parameter", &todoErr.APIErrorBody{
					Message: "Invalid sort direction",
					Target:  "sort",
				})

				return nil, apiError
			}
		}

		eachSort := &common.Sort{}
		eachSort.Field = strings.TrimSpace(fieldAndDirectionPair[0])

		if fieldAndDirectionLen == 1 {
			eachSort.SetDescending()
		} else {
			eachSort.Direction = strings.TrimSpace(fieldAndDirectionPair[1])
		}

		sortFields = append(sortFields, eachSort)
	}

	return sortFields, nil
}

func hasDuplicateSortFields(sorts []*common.Sort) *todoErr.APIError {
	seen := make(map[string]bool, len(sorts))
	var apiError *todoErr.APIError

	for _, each := range sorts {
		if _, ok := seen[each.Field]; ok {
			apiError = todoErr.NewAPIError("Invalid sort parameter", &todoErr.APIErrorBody{
				Message: "Duplicate sort fields",
				Target:  "sort",
			})

			break
		}

		seen[each.Field] = false
	}

	return apiError
}

func invalidCursorError() (int, interface{}, *todoErr.APIError) {
	apiError := todoErr.NewAPIError("Invalid cursor parameter", &todoErr.APIErrorBody{
		Message: "Invalid cursor",
		Target:  "cursor",
	})

	return http.StatusBadRequest, nil, apiError
}
