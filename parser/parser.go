package parser

/*

GRAMMAR
=======

expression  -> query ;
query       -> ( transform )* ;
transform   -> IDENTITY | selector ;
selector    -> "[" ( STRING | INTEGER | range ) "]" ;
range       -> INTEGER ":" INTEGER ;

*/
