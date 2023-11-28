package lexer

/*
The purpose of token tests is the verify how the Lexeme() value receiver
function behaves given different values for the buffer, start and end.

Scenarios:
1. Happy path: The buffer has a larger query buffer, the token is of type
   string, both start and end point to the right elements in the buffer,
   start and end are non-negative, start > end. This is how the it is used
   in code.
2. Given the buffer, when the buffer is nil, then Lexeme returns an empty
   string.
3. Given the buffer when the buffer is empty then Lexeme returns an empty
   string.
4. Given start and end, when start > end, then Lexeme returns an empty string
   without an error or panic.
5. Given start and end, start < end, when end > len(buffer), then the function
   goes only to the end of the buffer even when that would seemingly return
   an incomplete lexeme.
6. A default bare token with a nil buffer and start, end equal to 0 should
   return an empty string as a lexeme without an error.
*/
