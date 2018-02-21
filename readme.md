# RQL - Rio Query Language

This is a toy DBMS built purely for educational purposes.  It implements a small
subset of MySQL/Postgresql.  Definitely DO NOT use this for anything real.  100%
the worst database ever built.

## RQL Grammar

```sh
<Field> := IDENT
<Constant> := STRING_CONST | INT_CONST
<Expression> := <Field> | <Constant>
<Term> := <Expression> = <Expression>
<Predicate> := <Term> [ AND <Predicate> ]
<SelectList> := <Field> [ , <SelectList> ]
<TableList> := IdTok [ , <TableList> ]
<UpdateCmd> 
```

## Major Components Todo

- [ ] Token/Grammar definitions
- [ ] Lexer
- [ ] Parser (Recursive descent)
- [ ] Planner
- [ ] Query Optimizer
- [ ] Server
- [ ] Driver

## Parts of the Server

* Remote
* Planner
* Parse
* Query
* Metadata
* Record
* Transaction
* Buffer
* Log
* File
