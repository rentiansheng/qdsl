
%{

package main

import (
  )


%}

%union {
  obj []map[string]interface{}
  key string
  val interface{}
  list []interface{}
  op string
  errMsg string
}


%token <val> String Number Literal
%token <op> OpType JointType OpTypeBetween
%type <val> array between
%type <list> elements
%type <obj> object input
%type <val> value
%type <op> OpType  JointType



%%

input: object{
        yylex.(*lex).result = $$

}

object: String OpType value {
	$$ =  []map[string]interface{}{ {$1.(string): map[string]interface{}{ $2:   $3}}}

	yylex.(*lex).state = lexStateJoint
} | String '=' value {
	$$ =  []map[string]interface{}{{$1.(string): map[string]interface{}{ "=":   $3}}}
	yylex.(*lex).state = lexStateJoint

} |  String OpTypeBetween between {
	$$ =  []map[string]interface{}{{$1.(string): map[string]interface{}{ "between":   $3}}}
	yylex.(*lex).state = lexStateJoint
}| object JointType String OpType value {
   	obj := $$
	obj = append(obj, map[string]interface{}{
		$3.(string):map[string]interface{}{  $4:   $5},
		 "op": $2})
	$$ = obj
	yylex.(*lex).state = lexStateJoint

} | object JointType String '=' value {
    	obj := $$
  	obj = append(obj,map[string]interface{}{
  	$3.(string):map[string]interface{}{  "=":   $5},
  	"op": $2})
  	$$ = obj
  	yylex.(*lex).state = lexStateJoint

  }| object JointType String OpTypeBetween between{
       	obj := $$
     	obj = append(obj,map[string]interface{}{
     	$3.(string):map[string]interface{}{  "=":   $5},
     	"op": "between"})
     	$$ = obj
     	yylex.(*lex).state = lexStateJoint

     }




between : '(' elements ')' {
	yylex.(*lex).state = LexStateValueElement
	$$ = $2
  	yylex.(*lex).state = lexStateJoint
}

array: '[' elements ']' {
	yylex.(*lex).state = LexStateValueElement
	$$ = $2
	yylex.(*lex).state = lexStateJoint
}

elements: {
	$$ = []interface{}{}
}| value {
	$$ = []interface{}{$1}
}| elements ',' value {
	$$ = append($1, $3)
}

value: String
| Number
| Literal
| object {
        yylex.(*lex).pushState()
	$$ = $1
	yylex.(*lex).popState()
}
| array


%%
