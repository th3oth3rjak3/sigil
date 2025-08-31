program           = { statement } ;

statement         = variable_def
                  | print_statement
                  | expression_statement ;

variable_def      = "let" IDENTIFIER [ ":" type ] "=" expression ";" ;
expression_statement = expression ";" ;

expression        = equality ;
equality          = comparison { ( "==" | "!=" ) comparison } ;
comparison        = term { ( ">" | ">=" | "<" | "<=" ) term } ;
term              = factor { ( "+" | "-" ) factor } ;
factor            = unary { ( "*" | "/" ) unary } ;
unary             = ( "!" | "-" ) unary | primary ;
primary           = NUMBER | IDENTIFIER | "(" expression ")" ;
