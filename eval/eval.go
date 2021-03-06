package eval

import (
	"monkey/ast"
	"os"
)

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL  = &Null{}
	BREAK = &Break{}
)

var includeScope *Scope

func Eval(node ast.Node, scope *Scope) Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, scope)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, scope)
	case *ast.IncludeStatement:
		return evalIncludeStatement(node, scope)
	case *ast.LetStatement:
		return evalLetStatement(node, scope)
	case *ast.ReturnStatement:
		return evalReturnStatment(node, scope)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	case *ast.StringLiteral:
		return evalStringLiteral(node)
	case *ast.InterpolatedString:
		return evalInterpolatedString(node, scope)
	case *ast.Identifier:
		return evalIdentifier(node, scope)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, scope)
	case *ast.HashLiteral:
		return evalHashLiteral(node, scope)
	case *ast.StructLiteral:
		return evalStructLiteral(node, scope)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, scope)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, scope)
	case *ast.InfixExpression:
		return evalInfixExpression(node, scope)
	case *ast.IfExpression:
		return evalIfExpression(node, scope)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, scope)
	case *ast.CallExpression:
		return evalFunctionCall(node, scope)
	case *ast.MethodCallExpression:
		return evalMethodCallExpression(node, scope)
	case *ast.IndexExpression:
		return evalIndexExpression(node, scope)
	case *ast.DoLoop:
		return evalDoLoopExpression(node, scope)
	case *ast.BreakExpression:
		return BREAK
	case *ast.AssignExpression:
		return evalAssignStatement(node, scope)
	}
	return nil
}

// Program Evaluation Entry Point Functions, and Helpers:

func evalProgram(program *ast.Program, scope *Scope) (results Object) {
	loadIncludes(program.Includes, scope)
	for _, statement := range program.Statements {
		results = Eval(statement, scope)
		switch s := results.(type) {
		case *ReturnValue:
			return s.Value
		case *Error:
			return s
		}
	}
	if results == nil {
		return NULL
	}
	return results
}

func loadIncludes(includes map[string]*ast.IncludeStatement, s *Scope) {
	if includeScope == nil {
		includeScope = NewScope(nil)
	}
	for _, p := range includes {
		Eval(p, s)
	}
}

// Statements...
func evalIncludeStatement(i *ast.IncludeStatement, s *Scope) Object {
	if i.IsModule {
		return evalProgram(i.Program, s)
	}
	imported := &IncludedObject{Name: i.IncludePath.String(), Scope: NewScope(nil)}

	// capture stdout to suppress output during evaluating import
	so := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	if _, ok := includeScope.Get(i.IncludePath.String()); !ok {
		evalProgram(i.Program, imported.Scope)
		includeScope.Set(i.IncludePath.String(), imported)
	}

	// restore stdout
	w.Close()
	os.Stdout = so

	return imported
}

func evalLetStatement(l *ast.LetStatement, scope *Scope) (val Object) {
	if val = Eval(l.Value, scope); val.Type() != ERROR_OBJ {
		return scope.Set(l.Name.String(), val)
	}
	return
}

func evalAssignStatement(a *ast.AssignExpression, scope *Scope) (val Object) {
	if val = Eval(a.Value, scope); val.Type() != ERROR_OBJ {
		v, ok := scope.Reset(a.Name.String(), val)
		if ok {
			return v
		}
		return newError(UNKNOWNIDENT, a.Name.String())
	}
	return
}

func evalReturnStatment(r *ast.ReturnStatement, scope *Scope) Object {
	if value := Eval(r.ReturnValue, scope); value != nil {
		return &ReturnValue{Value: value}
	}
	return NULL
}

// Booleans
func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// Literals
func evalIntegerLiteral(i *ast.IntegerLiteral) Object {
	return &Integer{Value: i.Value}
}

func evalStringLiteral(s *ast.StringLiteral) Object {
	return &String{Value: s.Value}
}

func evalInterpolatedString(is *ast.InterpolatedString, scope *Scope) Object {
	s := &InterpolatedString{Value: &String{}, RawValue: is.Value, Expressions: is.ExprMap}
	s.Interpolate(scope)
	return s
}

func evalArrayLiteral(a *ast.ArrayLiteral, scope *Scope) Object {
	return &Array{Members: evalArgs(a.Members, scope)}
}

// Identifier not a literal, but felt logicially like it belonged here.. Literal expressions continue below
func evalIdentifier(i *ast.Identifier, scope *Scope) Object {
	val, ok := scope.Get(i.String())
	if !ok {
		if val, ok = includeScope.Get(i.String()); !ok {
			return newError(UNKNOWNIDENT, i.String())
		}
	}
	if i, ok := val.(*InterpolatedString); ok {
		i.Interpolate(scope)
		return i
	}
	return val
}

func evalHashLiteral(hl *ast.HashLiteral, scope *Scope) Object {
	hashMap := make(map[HashKey]HashPair)
	for key, value := range hl.Pairs {
		key := Eval(key, scope)
		if hashable, ok := key.(Hashable); ok {
			hashMap[hashable.HashKey()] = HashPair{Key: key, Value: Eval(value, scope)}
		} else {
			return newError(KEYERROR, key.Type())
		}
	}
	return &Hash{Pairs: hashMap}
}

func evalStructLiteral(s *ast.StructLiteral, scope *Scope) Object {
	structScope := NewScope(nil)
	for key, value := range s.Pairs {
		if ident, ok := key.(*ast.Identifier); ok {
			structScope.Set(ident.String(), Eval(value, scope))
		} else {
			return newError(KEYERROR, "IDENT")
		}
	}
	return &Struct{Scope: structScope, methods: make(map[string]*Function)}
}

func evalFunctionLiteral(fl *ast.FunctionLiteral, scope *Scope) Object {
	return &Function{Literal: fl, Scope: scope}
}

// Prefix expressions, e.g. `!true, -5`
func evalPrefixExpression(p *ast.PrefixExpression, s *Scope) Object {
	right := Eval(p.Right, s)
	if right.Type() == ERROR_OBJ {
		return right
	}
	switch p.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		if i, ok := right.(*Integer); ok {
			i.Value = -i.Value
			return i
		}
	}
	return newError(PREFIXOP, p.Operator, right.Type())
}

// Helper for evaluating Bang(!) expressions. Coerces truthyness based on object presence.
func evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// Evaluate infix expressions, e.g 1 + 2, a = 5, true == true, etc...
func evalInfixExpression(i *ast.InfixExpression, s *Scope) Object {
	left := Eval(i.Left, s)
	right := Eval(i.Right, s)
	if left.Type() == ERROR_OBJ {
		return left
	} else if right.Type() == ERROR_OBJ {
		return right
	}

	switch {
	case i.Operator == "and":
		return nativeBoolToBooleanObject(objectToNativeBoolean(left) && objectToNativeBoolean(right))
	case i.Operator == "or":
		return nativeBoolToBooleanObject(objectToNativeBoolean(left) || objectToNativeBoolean(right))
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntInfixExpression(i.Operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(i.Operator, left, right)
	case i.Operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case i.Operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	}
	return newError(INFIXOP, i.Operator, left.Type(), right.Type())
}

func objectToNativeBoolean(o Object) bool {
	if r, ok := o.(*ReturnValue); ok {
		o = r.Value
	}
	switch obj := o.(type) {
	case *Boolean:
		return obj.Value
	case *Null:
		return false
	default:
		return true
	}
}

// Helpers for infix evaluation below
//TODO: [] == [] is false. change.
func evalIntInfixExpression(operator string, left Object, right Object) Object {
	l := left.(*Integer)
	r := right.(*Integer)

	switch operator {
	case "+":
		return &Integer{Value: l.Value + r.Value}
	case "-":
		return &Integer{Value: l.Value - r.Value}
	case "*":
		return &Integer{Value: l.Value * r.Value}
	case "/":
		return &Integer{Value: l.Value / r.Value}
	case ">":
		return nativeBoolToBooleanObject(l.Value > r.Value)
	case "<":
		return nativeBoolToBooleanObject(l.Value < r.Value)
	case "==":
		return nativeBoolToBooleanObject(l.Value == r.Value)
	case "!=":
		return nativeBoolToBooleanObject(l.Value != r.Value)
	case "%":
		return evalModuloExpression(l, r)
	}
	return NULL
}

func evalModuloExpression(left *Integer, right *Integer) Object {
	mod := left.Value % right.Value
	if mod < 0 {
		mod += right.Value
	}
	return &Integer{Value: mod}
}

func evalStringInfixExpression(operator string, left Object, right Object) Object {
	l := left.(*String)
	r := right.(*String)

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(l.Value == r.Value)
	case "!=":
		return nativeBoolToBooleanObject(l.Value != r.Value)
	case "+":
		return &String{Value: l.Value + r.Value}
	}
	return newError(INFIXOP, operator, l.Type(), r.Type())
}

// Back to evaluators called directly by Eval

// IF expressions, if (evaluates to boolean) True: { Block Statement } Optional Else: {Block Statement}
func evalIfExpression(ie *ast.IfExpression, s *Scope) Object {
	condition := Eval(ie.Condition, s)
	if isTrue(condition) {
		return evalBlockStatements(ie.Consequence.Statements, s)
	} else if ie.Alternative != nil {
		return evalBlockStatements(ie.Alternative.Statements, s)
	}
	return NULL
}

func evalDoLoopExpression(dl *ast.DoLoop, s *Scope) Object {
	for {
		e := Eval(dl.Block, NewScope(s))
		if _, ok := e.(*Break); ok {
			break
		}
		if v, ok := e.(*ReturnValue); ok {
			if v.Value != nil {
				return v.Value
			}
			break
		}
	}
	return NULL
}

// Helper function isTrue for IF evaluation - neccessity is dubious
func isTrue(obj Object) bool {
	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

// Block Statement Evaluation - The innards of both IF and Function calls
// very similar to parseProgram, but because we need to leave the return
// value wrapped in it's Object, it remains, for now.
func evalBlockStatements(block []ast.Statement, scope *Scope) (results Object) {
	for _, statement := range block {
		results = Eval(statement, scope)
		if results != nil && results.Type() == RETURN_VALUE_OBJ {
			return
		}
		if _, ok := results.(*Break); ok {
			return
		}
	}
	return
}

// Eval when a function is _called_, includes fn literal evaluation and calling builtins
func evalFunctionCall(call *ast.CallExpression, s *Scope) Object {
	fn, ok := s.Get(call.Function.String())
	if !ok {
		if f, ok := call.Function.(*ast.FunctionLiteral); ok {
			fn = &Function{Literal: f, Scope: s}
			s.Set(call.Function.String(), fn)
		} else if builtin, ok := builtins[call.Function.String()]; ok {
			return builtin.Fn(evalArgs(call.Arguments, s)...)
		} else {
			return newError(UNKNOWNIDENT, call.Function.String())
		}
	}
	f := fn.(*Function)
	f.Scope = NewScope(s)
	args := evalArgs(call.Arguments, f.Scope)
	// TODO: If not enough of arguments are passed a panic occur, if too few, no warning or error
	for i, v := range f.Literal.Parameters {
		f.Scope.Set(v.String(), args[i])
	}
	r := Eval(f.Literal.Body, f.Scope)
	if obj, ok := r.(*ReturnValue); ok {
		return obj.Value
	}
	return r
}

// Method calls for builtin Objects
func evalMethodCallExpression(call *ast.MethodCallExpression, scope *Scope) Object {
	obj := Eval(call.Object, scope)
	switch m := obj.(type) {
	case *IncludedObject:
		switch o := call.Call.(type) {
		case *ast.Identifier:
			if i, ok := m.Scope.Get(call.Call.String()); ok {
				return i
			}
		case *ast.CallExpression:
			if o.Function.String() == "Scope" {
				return obj.CallMethod("Scope")
			}
			return evalFunctionCall(o, m.Scope)
		}
	case *Struct:
		switch o := call.Call.(type) {
		case *ast.Identifier:
			if i, ok := m.Scope.Get(call.Call.String()); ok {
				return i
			}
		case *ast.CallExpression:
			args := evalArgs(o.Arguments, scope)
			return obj.CallMethod(o.Function.String(), args...)
		}
	default:
		if method, ok := call.Call.(*ast.CallExpression); ok {
			args := evalArgs(method.Arguments, scope)
			return obj.CallMethod(method.Function.String(), args...)
		}
	}
	return newError(NOMETHODERROR, call.String(), obj.Type())

}

func evalArgs(args []ast.Expression, scope *Scope) []Object {
	//TODO: Refactor this to accept the params and args, go ahead and
	// update scope while looping and return the Scope object.
	e := []Object{}
	for _, v := range args {
		e = append(e, Eval(v, scope))
	}
	return e
}

// Index Expressions, i.e. array[0], array[2:4] or hash["mykey"]

func evalIndexExpression(ie *ast.IndexExpression, s *Scope) Object {
	left := Eval(ie.Left, s)
	switch iterable := left.(type) {
	case *Array:
		return evalArrayIndex(iterable, ie, s)
	case *Hash:
		return evalHashKeyIndex(iterable, ie, s)
	case *String:
		return evalStringIndex(iterable, ie, s)
	}
	return newError(NOINDEXERROR, left.Type())
}

func evalStringIndex(str *String, ie *ast.IndexExpression, s *Scope) Object {
	var idx int64
	length := int64(len(str.Value))
	if exp, success := ie.Index.(*ast.SliceExpression); success {
		return evalStringSliceExpression(str, exp, s)
	}
	index := Eval(ie.Index, s)
	if index.Type() == ERROR_OBJ {
		return index
	}
	idx = index.(*Integer).Value
	if idx > length-1 {
		return newError(INDEXERROR, idx)
	}
	if idx < 0 {
		idx = length + idx
		if idx > length-1 || idx < 0 {
			return newError(INDEXERROR, idx)
		}
	}
	return &String{Value: string(str.Value[idx])}
}

func evalStringSliceExpression(str *String, se *ast.SliceExpression, s *Scope) Object {
	var idx int64
	var slice int64
	length := int64(len(str.Value))

	startIdx := Eval(se.StartIndex, s)
	if startIdx.Type() == ERROR_OBJ {
		return startIdx
	}
	idx = startIdx.(*Integer).Value
	if idx > length-1 {
		return newError(INDEXERROR, idx)
	}
	if idx < 0 {
		idx = length + idx
		if idx > length-1 || idx < 0 {
			return newError(INDEXERROR, idx)
		}
	}

	if se.EndIndex == nil {
		slice = length
	} else {
		slIndex := Eval(se.EndIndex, s)
		if slIndex.Type() == ERROR_OBJ {
			return slIndex
		}
		slice = slIndex.(*Integer).Value
		if slice > length-1 {
			return newError(SLICEERROR, idx, slice)
		}
		if slice < 0 {
			slice = length + slice
			if slice > length-1 || slice < idx {
				return newError(SLICEERROR, idx, slice)
			}
		}
	}
	if idx == 0 && slice == length {
		return str
	}
	if slice == length {
		return &String{Value: string(str.Value[idx:])}
	}
	return &String{Value: string(str.Value[idx:slice])}
}

func evalHashKeyIndex(hash *Hash, ie *ast.IndexExpression, s *Scope) Object {
	key := Eval(ie.Index, s)
	if key.Type() == ERROR_OBJ {
		return key
	}
	hashable, ok := key.(Hashable)
	if !ok {
		return newError(KEYERROR, key.Type())
	}
	hashPair, ok := hash.Pairs[hashable.HashKey()]
	// TODO: should we return an error here? If not, maybe arrays should return NULL as well?
	if !ok {
		return NULL
	}
	return hashPair.Value
}

func evalArraySliceExpression(array *Array, se *ast.SliceExpression, s *Scope) Object {
	var idx int64
	var slice int64
	length := int64(len(array.Members))

	startIdx := Eval(se.StartIndex, s)
	if startIdx.Type() == ERROR_OBJ {
		return startIdx
	}
	idx = startIdx.(*Integer).Value
	if idx > length-1 {
		return newError(INDEXERROR, idx)
	}
	if idx < 0 {
		idx = length + idx
		if idx > length-1 || idx < 0 {
			return newError(INDEXERROR, idx)
		}
	}

	if se.EndIndex == nil {
		slice = length
	} else {
		slIndex := Eval(se.EndIndex, s)
		if slIndex.Type() == ERROR_OBJ {
			return slIndex
		}
		slice = slIndex.(*Integer).Value
		if slice > length-1 {
			return newError(SLICEERROR, idx, slice)
		}
		if slice < 0 {
			slice = length + slice
			if slice > length-1 || slice < idx {
				return newError(SLICEERROR, idx, slice)
			}
		}
	}
	if idx == 0 && slice == length {
		return array
	}
	if slice == length {
		return &Array{Members: array.Members[idx:]}
	}
	return &Array{Members: array.Members[idx:slice]}
}

func evalArrayIndex(array *Array, ie *ast.IndexExpression, s *Scope) Object {
	var idx int64
	length := int64(len(array.Members))
	if exp, success := ie.Index.(*ast.SliceExpression); success {
		return evalArraySliceExpression(array, exp, s)
	}
	index := Eval(ie.Index, s)
	if index.Type() == ERROR_OBJ {
		return index
	}
	idx = index.(*Integer).Value
	if idx > length-1 {
		return newError(INDEXERROR, idx)
	}
	if idx < 0 {
		idx = length + idx
		if idx > length-1 || idx < 0 {
			return newError(INDEXERROR, idx)
		}
	}
	return array.Members[idx]
}
