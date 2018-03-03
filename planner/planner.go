// Package planner uses the Lexer and Parser to generate query/update plans.  It
// is responsible for checking the deeper level semantics.  Does the table exist
// or not? etc.
package planner

// type Transaction struct {}
// type QueryData struct {}

// type Planner struct {
// qp QueryPlanner
// up UpdatePlanner
// }

// type QueryPlanner interface {
// CreatePlan(data QueryData, tx Transaction)
// }

// type UpdatePlanner interface {
// ExecuteInsert(data InsertData, tx Transaction)
// ExecuteCreateTable(data CreateTableData, tx Transaction)
// // TODO: ... DELETE/MODIFY/INDEX ...
// }
