module github.com/jo-hoe/google-sheets

go 1.19

require golang.org/x/oauth2 v0.19.0

require (
	cloud.google.com/go/compute v1.25.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
)

retract v1.0.0 // published erroneously
