# Goven 🧑‍🍳

Goven (go-oven) is a go library that allows you to have a drop-in query language for your database schema. 

* Take any gorm database object and with a few lines of code you can make it searchable.
* [Safe query language not exposing SQL to your users](https://imgs.xkcd.com/comics/exploits_of_a_mom.png).
* Easily extensible to support more advanced queries for your own schema.
* Basic grammar that allows for powerful queries.

Like a real life oven, it takes something raw (your database struct + a query input) and produces something specific (the schema specific parser + SQL output). We call the adaptors "recipes" that goven can make. Currently Goven only supports a SQL adaptor, but the AST produced by the lexer/parser can easily be extended to other query languages.

## Recipes

### Basic Example

TODO: here and in examples folder

### Extension Example

TODO: here and in examples folder

## Grammar

Goven has a simple syntax that allows for powerful queries.

Fields can be compared using the following operators: 

`=`, `!=`, `>=`, `<=`, `<`, `>`

Multiple queries can be combined using `AND`, `OR`.

Together this means you can build up a  query like this:

`model_name=iris AND version>=2.0`

More advanced queries can be built up using bracketed expressions:

`(model_name=iris AND version>=2.0) OR artifact_type=TENSORFLOW`
