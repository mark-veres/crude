# crude
Turning **GO**RM **to** an **CRUD** API.

```sh
go get github.com/mark-veres/crude
```

## Usage
```go
package main

type Post struct {
	gorm.Model
	Name    string
	Content string
}

func main() {
    r := gin.Default()
    api := r.Group("/api")
    db, _ := gorm.Open(...)
    db.AutoMigrate(&Post{})

    config := crude.Config{DB: db}
    crude.Register[Post](&config, api, "posts")

    r.Run(":80")
}
```

## Routes
The `crude.Register` function attaches the following routes to the given `*gin.RouterGroup`.

### `POST` /{name}/new
This route creates a new record in the database by marshalling the request body to the given model.
```sh
curl -X POST -H "Content-Type: application/json" \
    -d '{"name": "new post", "content": ""}' \
    http://localhost:80/.../{name}/new
```

### `POST` /{name}/update
This route updates the record that has the same [`ID`](https://gorm.io/docs/models.html#gorm-Model) as the JSON parsed request body with the request body.  
Updating is accomplished with the gorm [`(*DB)Save`](https://pkg.go.dev/gorm.io/gorm@v1.25.4#DB.Save) function.
```json
// The request body must take this form
{
    "ID": 1,
    ...
}
```
```sh
curl -X POST -H "Content-Type: application/json" \
    -d '{"id": 0, "content": "the content was updated"}' \
    http://localhost:80/.../{name}/updates
```

### `GET` /{name}/list
This route returns all elements in the table.
```sh
curl http://localhost:80/.../{name}/list
```

### `GET` /{name}/by/{property}?value
This route retuns all records that match the SQL query:  
`SELECT * FROM {name} WHERE {property} = {value}` (+ other stuff gorm automatically generates)

### `GET` /{name}/where/{property}/{operator}
This route allows simple queries using the WHERE clause.  
Possible [operators](https://www.w3schools.com/sql/sql_where.asp):
- `=`, `>`, `<`, `>=`, `<=`, `!=`
- `between`
- `like`

#### Usage of `=`, `>`, `<`, `>=`, `<=`, `!=` operators
```sh
curl http://localhost:80/.../{name}/where/{property}/{operator}?value={value}
```

#### Usage of `between` operator
```sh
curl http://localhost:80/.../{name}/where/{property}/{operator}?from={from}&to={to}
```

#### Usage of `like` operator
```sh
curl http://localhost:80/.../{name}/where/{property}/{operator}?pattern={pattern}
```

## Custom Middleware
The `Config` object has a list of middleware that can be run before every one of the 4 CRUD operations.
- `Config.CreateMiddleware`
    + `POST` /{name}/new
- `Config.ReadMiddleware`
    + `GET` /{name}/list
    + `GET` /{name}/by/{property}
    + `GET` /{name}/where/{property}/{operator}
- `Config.UpdateMiddleware`
    + `POST` /{name}/update
- `Config.DeleteMiddleware`
    + `GET` /{name}/delete