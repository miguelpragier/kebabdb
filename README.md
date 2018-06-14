# KebabDB
Database helper for go

For now, only **MySql** helper is written.

## Dependencies:
KebabDB depends on this nice work: [github.com/go-sql-driver/mysql]("github.com/go-sql-driver/mysql")

## API or How to Use
I bet it's simpler than you expected

First, you need to download kebabdb.go and put it in a directory called kebabdb.
Now you can call methods qualifying like this: kebabdb.GetCount("tablecustomers")


## Connection String
**Set an environment variable named KebabDBConnectionString**

### Examples
```go
/*
* All methods can be called with query + separated ( variadic ) params or with a ready query only.
*/

// Execute() sends a sql query to database. tipically update, delete and insert, and return the number of affected rows
affectedRows,err := kebabdb.Execute("UPDATE tablecustomers SET age=? WHERE id=?", age, id)

// Insert() executes "insert" sql queries against database an try to return lastInsertedID
newID, err := kebabdb.Insert("INSERT INTO tablecustomers (name,age) VALUES (?,?)")

// GetString returns the text of the first row and column, or nil.
customerName := kebabdb.GetString("SELECT name FROM tablecustomers WHERE id=?", 15)

// GetInt executes the query, and get the first rol/col data, and - if necessary - converts to int.
customerAge := kebabdb.GetInt("SELECT age FROM tablecustomers WHERE name='adelle' LIMIT 1")

// GetCount() returns the rows count of the given table
howManyRows := kebabdb.GetCount("tablecustomers")

// GetOne returns query result as a map[string]string
// If something goes wrong, it returns an empty map[string]string or nil, depending on the case
row, err := kebabdb.GetOne("SELECT id, name, age FROM tablecustomers WHERE id=8")

fmt.Printf("Customer's id: %s\n", row["id"])
fmt.Printf("Customer's name: %s\n", row["name"])
fmt.Printf("Customer's age: %s\n", row["age"])


// GetMany returns query result as a []map[string]string
// If something goes wrong, it returns an empty []map[string]string or nil, depending on the case
resultset, err := kebabdb.GetMany("SELECT id, name, age FROM tablecustomers ORDER BY name")

for _,row := range resultset {
    fmt.Printf("Customer's id: %s\n", row["id"])
    fmt.Printf("Customer's name: %s\n", row["name"])
    fmt.Printf("Customer's age: %s\n", row["age"])
}

```
