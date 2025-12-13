package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbpkg "github.com/KrrMaxim/task-6/internal/db"
)

var (
	errQuery = errors.New("query error")
	errRow   = errors.New("row iteration error")
)

func TestNew(test *testing.T) {
	test.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(test, err)
	defer mockDB.Close()

	service := dbpkg.New(mockDB)
	require.NotNil(test, service)
	require.Equal(test, mockDB, service.DB)
}

func TestGetNames(test *testing.T) {
	test.Parallel()

	cases := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		expected    []string
		errContains string
	}{
		{
			name: "success multiple rows",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Maxim").
					AddRow("Artem")
				mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)
			},
			expected: []string{"Maxim", "Artem"},
		},
		{
			name: "success empty",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)
			},
			expected: nil,
		},
		{
			name: "query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT name FROM users`).WillReturnError(errQuery)
			},
			errContains: "db query",
		},
		{
			name: "scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)
			},
			errContains: "rows scanning",
		},
		{
			name: "rows iteration error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("first").
					AddRow("second")
				rows.RowError(1, errRow)
				mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)
			},
			errContains: "rows error",
		},
	}

	for _, tc := range cases {
		test.Run(tc.name, func(test *testing.T) {
			test.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(test, err)
			defer mockDB.Close()

			tc.setupMock(mock)

			service := dbpkg.New(mockDB)
			result, err := service.GetNames()

			require.NoError(test, mock.ExpectationsWereMet())

			if tc.errContains != "" {
				require.Error(test, err)
				assert.Contains(test, err.Error(), tc.errContains)
				assert.Nil(test, result)
			} else {
				require.NoError(test, err)
				assert.Equal(test, tc.expected, result)
			}
		})
	}
}

func TestGetUniqueNames(test *testing.T) {
	test.Parallel()

	runCase := func(
		test *testing.T,
		prepare func(sqlmock.Sqlmock),
		expected []string,
		errContains string,
	) {
		test.Helper()

		mockDB, mock, err := sqlmock.New()
		require.NoError(test, err)
		defer mockDB.Close()

		prepare(mock)

		service := dbpkg.New(mockDB)
		result, err := service.GetUniqueNames()

		require.NoError(test, mock.ExpectationsWereMet())

		if errContains != "" {
			require.Error(test, err)
			assert.Contains(test, err.Error(), errContains)
			assert.Nil(test, result)
			return
		}

		require.NoError(test, err)
		assert.Equal(test, expected, result)
	}

	test.Run("success distinct", func(test *testing.T) {
		test.Parallel()

		prepare := func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"name"}).
				AddRow("u1").
				AddRow("u2")
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		}

		runCase(test, prepare, []string{"u1", "u2"}, "")
	})

	test.Run("query error", func(test *testing.T) {
		test.Parallel()

		prepare := func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("SELECT DISTINCT name FROM users").
				WillReturnError(errors.New("fail"))
		}

		runCase(test, prepare, nil, "db query")
	})

	test.Run("scan error", func(test *testing.T) {
		test.Parallel()

		prepare := func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		}

		runCase(test, prepare, nil, "rows scanning")
	})

	test.Run("rows error", func(test *testing.T) {
		test.Parallel()

		prepare := func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"name"}).
				AddRow("ok").
				AddRow("bad")
			rows.RowError(1, errors.New("iter"))
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		}

		runCase(test, prepare, nil, "rows error")
	})
}
