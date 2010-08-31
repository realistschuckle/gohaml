%package gohaml

%union {
	n inode
	s string
	i interface{}
	c icodenode
}

%type<n> Statement
%type<c> RightHandSide
%type<s> ident
%type<i> atom

%%

Statement : tfor ident ',' ident ':' '=' trange ident	{
															rn := new(rangenode)
															rn._first = $2
															rn._second = $4
															rn._third = $8
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
			  | ident									{
															dan := new(vdeclassnode)
															dan._rhs.value = $1
															dan._rhs.needsResolution = true
															$$ = dan			
														}
			  ;