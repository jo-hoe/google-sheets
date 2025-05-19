module github.com/jo-hoe/google-sheets

go 1.23.0

require golang.org/x/oauth2 v0.30.0

require (
	cloud.google.com/go/compute/metadata v0.7.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

retract v1.0.0 // published erroneously
