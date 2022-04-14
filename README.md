[![Go Reference](https://pkg.go.dev/badge/github.com/seldonio/goven.svg)](https://pkg.go.dev/github.com/seldonio/goven)
[![Go Report Card](https://goreportcard.com/badge/github.com/seldonio/goven)](https://goreportcard.com/report/github.com/seldonio/goven)
[![codecov](https://codecov.io/gh/seldonio/goven/branch/master/graph/badge.svg?token=ZBCTOI896Y)](https://codecov.io/gh/seldonio/goven)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  

# Goven 🧑‍🍳

Goven (go-oven) is a go library that allows you to have a drop-in query language for your database schema. 

* Take any gorm database object and with a few lines of code you can make it searchable.
* [Safe query language not exposing SQL to your users](https://imgs.xkcd.com/comics/exploits_of_a_mom.png).
* Easily extensible to support more advanced queries for your own schema.
* Basic grammar that allows for powerful queries.

Like a real life oven, it takes something raw (your database struct + a query input) and produces something specific (the schema specific parser + SQL output). We call the adaptors "recipes" that Goven can make. Currently Goven only supports a SQL adaptor, but the AST produced by the lexer/parser can easily be extended to other query languages.

## Recipes

### Basic Example

You can make a basic query using gorm against your database, something like this: 

```go
reflection := reflect.ValueOf(&User{})
queryAdaptor, err := sql_adaptor.NewDefaultAdaptorFromStruct(reflection)
if err != nil {
    return nil, err
}

dbQuery := db.WithContext(ctx)
parsedQuery, err := queryAdaptor.Parse("(name=james AND age > 11) OR email=fred@gmail.com")
if err != nil {
	return nil, err
}
dbQuery = query.Model(User{}).Where(parsedQuery.Raw, sql_adaptor.StringSliceToInterfaceSlice(parsedQuery.Values)...)
err = dbQuery.Find(&users).Error
```

The values are interpolated to prevent injection attacks.

### Extension Example

You can also extend the basic query language with regex matchers. An example would be having a Tag struct on your User schema.

```go
type User struct {
	gorm.Model
	name string
	tags []Tag
}

type Tag struct {
	gorm.Model
	Key string
	Value string
}
```

You can make this searchable my defining a regex and a custom matcher.

e.g if we want `tags[key]=value` to work then we can add the following matcher when creating the adaptor.

```go
KeyValueRegex = `(tags)\[(.+)\]`

// keyValueMatcher is a custom matcher for and tags[y].
func keyValueMatcher(ex *goven_parser.Expression) (*goven_sql.SqlResponse, error) {
	reg, err := regexp.Compile(KeyValueRegex)
	if err != nil {
		return nil, err
	}
	slice := reg.FindStringSubmatch(ex.Field)
	if slice == nil {
		return nil, errors.New("didn't match regex expression")
	}
	if len(slice) < 3 {
		return nil, errors.New("regex match slice is too short")
	}
	sq := goven_sql.SqlResponse{
		Raw:    fmt.Sprintf("id IN (SELECT user_id FROM %s WHERE key=? AND value%s?)", slice[1], ex.Comparator),
		Values: []string{slice[2], ex.Value},
	}
	return &sq, nil
}
```

### Protecting Fields

Sometimes we may not want particular fields to be searchable by end users. You can protect them by removing them from the fields mapping when creating your adaptor.

```go
defaultFields := goven_sql.FieldParseValidatorFromStruct(gorm)
delete(defaultFields, "fieldname")
```

## Grammar

Goven has a simple syntax that allows for powerful queries.

Fields can be compared using the following operators: 

`=`, `!=`, `>=`, `<=`, `<`, `>`, `%`

The `%` operator allows you to do partial string matching using LIKE.

Multiple queries can be combined using `AND`, `OR`.

Together this means you can build up a  query like this:

`model_name=iris AND version>=2.0`

More advanced queries can be built up using bracketed expressions:

`(model_name=iris AND version>=2.0) OR artifact_type=TENSORFLOW`
