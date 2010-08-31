%package gohaml

%import fmt

%union {
	n inode
	s string
	i interface{}
	c icodenode
}

%type<n> Statement
%type<c> RightHandSide
%type<s> ident ComplexIdent
%type<i> atom

%%

Statement : tfor ident ',' ident ':' '=' trange ident	{
															rn := new(rangenode)
															rn._lhs1 = $2
															rn._lhs2 = $4
															rn._rhs = res{$8, true}
															$$ = rn
														}
		  | ident ':' '=' RightHandSide					{
															$4.setLHS($1)
															$$ = $4
														}
          ;

RightHandSide : atom									{
															dan := new(declassnode)
															dan._rhs = $1
															$$ = dan
														}
			  | ident ComplexIdent						{
															dan := new(vdeclassnode)
															dan._rhs.value = $1 + $2
															dan._rhs.needsResolution = true
															$$ = dan			
														}
			  ;

ComplexIdent : '.' ident ComplexIdent					{ $$ = fmt.Sprintf(".%s%s", $2, $3)}
			 |											{ $$ = "" }
			 ;