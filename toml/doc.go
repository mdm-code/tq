/*
Package toml declares an adapter for external TOML packages to accommodate
different implementations of data encoders and decoders. This limits the
necessary code changes when swapping implementations. It can also leave room
for users to configure the implementation they want to use.
*/
package toml
