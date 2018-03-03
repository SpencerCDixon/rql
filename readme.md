# RQL - Rio Query Language

This is a toy DBMS built purely for educational purposes.  It implements a small
subset of MySQL/Postgresql.  Definitely DO NOT use this for anything real.  100%
the worst database ever built.

## RQL Grammar

Below is the basic grammer used in the RQL DBMS.  The characters `[`, `]`, and
`|` are not used as delimeters in the language so we use them as punctuation to
help define patterns.

The `|` character represents a boolean `or`.  So:

```sh
<Constant> := STRING_TOK | INT_TOK
```

Means a constant can be either a string or an int.

The `[`/`]` characters are used to denote an _optional_ aspect of the grammar.
So:

```sh
<Predicate> := <Term> [ AND <Predicate> ]
```

Can be composed of a `<Term>` optionally followed by a `<Predicate>`.  The
following are all valid `<Predicate>`'s:

```sh
UserName = 'Spencer'
Color = 'red' AND Category = 'flower'
Language = 'english' AND Age = 25 AND EyeColor = 'blue'
```

With that, the entire grammer for the RQL language (small subset of SQL) can be
found below:

```sh
<Field>       := IDENT
<Constant>    := STRING_TOK | INT_TOK
<Expression>  := <Field> | <Constant>
<Term>        := <Expression> = <Expression>
<Predicate>   := <Term> [ AND <Predicate> ]
<SelectList>  := <Field> [ , <SelectList> ]
<TableList>   := IDENT [ , <TableList> ]
<UpdateCmd>   := <Insert> | <Delete> | <Modify> | <Create>
<Create>      := <CreateTable> | <CreateIndex>
<Insert>      := INSERT INTO IDENT ( <FieldList> ) VALUES ( <ConstList> )
<FieldList>   := <Field> [ , <FieldList> ]
<ConstList>   := <Constant> [ , <Constant> ]
<Delete>      := DELETE FROM IDENT [ WHERE <Predicate> ]
<Modify>      := UPDATE IDENT SET <Field> = <Expression> [ WHERE <Predicate> ]
<CreateTable> := CREATE TABLE IDENT ( <FieldDefs> )
<FieldDefs>   := <FieldDef> [ , <FieldDefs> ]
<FieldDef>    := IDENT <TypeDef>
<TypeDef>     := INT | VARCHAR ( INT_TOK )
<CreateIndex> := CREATE INDEX IDENT ON IDENT ( <Field> )
``` 

## Major Components

* [ ] CLI
* [ ] Remote
* [ ] Planner
* [ ] Parse
* [x] Lexer
* [ ] Query
* [ ] Metadata
* [ ] Record
* [ ] Transaction
* [ ] Buffer
* [ ] Tracer/Stats
* [x] Log
* [x] File

## Milestones

### Create/Insert/Select

```sql
CREATE TABLE users (
  id int,
  name varchar(200),
  company varchar(100)
);
INSERT INTO users VALUES (1, 'Spencer Dixon', 'Rio');
INSERT INTO users VALUES (1, 'Stefan VanBuren', 'Rio');
SELECT name FROM users;
```
