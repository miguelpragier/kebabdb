package kebabdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	// database/sql uses this package.
	_ "github.com/go-sql-driver/mysql"
)

// KebabDB represents the current state of database connection pool
type kebabDB struct {
	db           *sql.DB
	lastSQLQuery string
}

var (
	kdb              kebabDB
	connectionString = ""
)

func init() {
	if os.Getenv("DatabaseConnectionString") != "" {
		connectionString = os.Getenv("DatabaseConnectionString")
	}

	connectDB()
}

// SetConnectionString defines or redefines KebabDB connection string
func SetConnectionString(cs string) {
	connectionString = cs
}

// ConnectDB Assures that database could be reached.
func connectDB() {
	log.Println("Connecting database...")

	if connectionString == "" {
		log.Fatalf("Empty Connection String]!")
	}

	if sqldb, err := sql.Open("mysql", connectionString); err == nil {
		kdb.db = sqldb

		if err := kdb.db.Ping(); err != nil {
			log.Fatalf("database.init db.ping %v", err)
		}
	} else {
		log.Fatalf("database.init sql.open %v", err)
	}
}

// Execute sends a sql query to database. tipically update, delete and insert, and return the number of affected rows
func Execute(sqlQuery string, params ...interface{}) (int64, error) {
	if kdb.db == nil {
		connectDB()
	}

	kdb.lastSQLQuery = sqlQuery

	rs, err := kdb.db.Exec(sqlQuery, params...)

	if err != nil {
		log.Printf(`kebabdb.Execute db.Query The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)
		return 0, err
	}

	return rs.RowsAffected()
}

// Insert executes "insert" sql queries against database an try to return lastInsertedID
func Insert(sqlQuery string, params ...interface{}) (int64, error) {
	if kdb.db == nil {
		connectDB()
	}

	kdb.lastSQLQuery = sqlQuery

	rs, err := kdb.db.Exec(sqlQuery, params...)

	if err != nil {
		log.Printf(`kebabdb.Insert db.Exec The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)
		return 0, err
	}

	return rs.LastInsertId()
}

// GetString returns the text of the first row and column, or nil.
func GetString(sqlQuery string, params ...interface{}) string {
	if kdb.db == nil {
		connectDB()
	}

	kdb.lastSQLQuery = sqlQuery

	var s string

	err := kdb.db.QueryRow(sqlQuery, params...).Scan(&s)

	if err != nil {
		log.Printf(`kebabdb.GetString db.QueryRow().Scan() The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)
	}

	return s
}

// GetInt executes the query, and get the first rol/col data, and it converts to int.
// It returns an int anyway.
func GetInt(sqlQuery string, params ...interface{}) int {
	if kdb.db == nil {
		connectDB()
	}

	kdb.lastSQLQuery = sqlQuery

	var i int

	err := kdb.db.QueryRow(sqlQuery, params...).Scan(&i)

	if err != nil {
		log.Printf(`kebabdb.GetString db.QueryRow().Scan() The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)
	}

	return i
}

// GetCount returns the number of rows in the given table
func GetCount(tableName string) uint64 {
	if kdb.db == nil {
		connectDB()
	}

	sqlQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, tableName)

	kdb.lastSQLQuery = sqlQuery

	var i uint64

	err := kdb.db.QueryRow(sqlQuery).Scan(&i)

	if err != nil {
		log.Printf(`kebabdb.GetCount db.QueryRow().Scan() Query "%s" has failed with "%v"\n`, sqlQuery, err)
	}

	return i
}

// GetOne returns query result as a map[string]string
func GetOne(sqlQuery string, params ...interface{}) (map[string]string, error) {
	if kdb.db == nil {
		connectDB()
	}

	kdb.lastSQLQuery = sqlQuery

	rows, err := kdb.db.Query(sqlQuery, params...)

	if err != nil {
		log.Printf(`kebabdb.GetOne db.Query The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)
		return nil, err
	}

	// If there's no rows, return an empty resultset
	if !rows.Next() {
		if rows.Err() != nil {
			log.Printf(`kebabdb.GetOne rows.Next The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)

			return nil, rows.Err()
		}

		return map[string]string{}, nil
	}

	cols, err := rows.Columns()

	if err != nil {
		log.Printf(`kebabdb.GetOne rows.Columns The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)
		return nil, err
	}

	howManyCols := len(cols)

	// Byte array for store values
	realDest := make([][]byte, howManyCols)

	// Once Rows.Scan() expects []interface{} ( slice of interfaces ), we need one of these as a temporary buffer
	fakeDest := make([]interface{}, howManyCols)

	// Here we set a byteArrayPointer for each interface instance
	for i := range realDest {
		fakeDest[i] = &realDest[i]
	}

	err = rows.Scan(fakeDest...)

	if err != nil {
		log.Printf(`kebabdb.GetOne rows.Scan The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)
		return nil, err
	}

	// array for the resulting values of each row
	valuesArray := make(map[string]string, howManyCols)

	for i, raw := range realDest {
		if raw == nil {
			valuesArray[cols[i]] = ""
		} else {
			valuesArray[cols[i]] = string(raw)
		}
	}

	return valuesArray, nil
}

// GetMany returns sql query result as an array of map[string]string
func GetMany(sqlQuery string, params ...interface{}) ([]map[string]string, error) {
	if kdb.db == nil {
		connectDB()
	}

	kdb.lastSQLQuery = sqlQuery

	rows, err := kdb.db.Query(sqlQuery, params...)

	if err != nil {
		log.Printf(`kebabdb.GetMany db.Query The query "%s" with params "%v" has failed with "%v"\n`, sqlQuery, params, err)
		return nil, err
	}

	cols, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	howManyCols := len(cols)

	// Byte array for store values
	realDest := make([][]byte, howManyCols)

	// Once Rows.Scan() expects []interface{} ( slice of interfaces ), we need one of these as a temporary buffer
	fakeDest := make([]interface{}, howManyCols)

	// Here we set a byteArrayPointer for each interface instance
	for i := range realDest {
		fakeDest[i] = &realDest[i]
	}

	result := make([]map[string]string, 0)

	for rows.Next() {
		err = rows.Scan(fakeDest...)

		if err != nil {
			log.Printf(`kebabdb.GetMany rows.Scan The query "%s" has failed with "%v"\n`, sqlQuery, err)
			return nil, err
		}

		// Array for the resulting values of each row
		valuesArray := make(map[string]string, howManyCols)

		for i, raw := range realDest {
			if raw == nil {
				valuesArray[cols[i]] = ""
			} else {
				valuesArray[cols[i]] = string(raw)
			}
		}

		result = append(result, valuesArray)
	}

	return result, nil
}

// Exists returns true if the criteria is satisfied by database table
// Params are "table" and criteria
// @arg criteria example: "id=15 AND name='mike'"
func Exists(table, criteria string) bool {
	if kdb.db == nil {
		connectDB()
	}

	if strings.ToLower(string(criteria[0:6])) == "where " {
		criteria = strings.Replace(criteria, "WHERE ", "", 1)
	}

	sqlQuery := fmt.Sprintf("SELECT CASE WHEN EXISTS ( SELECT * FROM %s WHERE %s ) THEN 1 ELSE 0 END", table, criteria)

	kdb.lastSQLQuery = sqlQuery

	i := uint64(0)

	err := kdb.db.QueryRow(sqlQuery).Scan(&i)

	if err != nil {
		log.Printf(`kebabdb.Exists db.QueryRow().Scan() The query "%s" has failed with "%v"\n`, sqlQuery, err)
	}

	return i == 1
}

// GetLastSQLQuery is a draft way to return the last issued sql command.
// It's not working now, because I couldn't retrive the last sql from mysql driver.
func GetLastSQLQuery() string {
	return kdb.lastSQLQuery
}
