module github.com/jo-hoe/google-sheets

go 1.19

require golang.org/x/oauth2 v0.20.0

require (
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
)

retract v1.0.0 // published erroneously
