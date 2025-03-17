module github.com/jo-hoe/google-sheets

go 1.19
toolchain go1.24.1

require golang.org/x/oauth2 v0.28.0

require (
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
)

retract v1.0.0 // published erroneously
