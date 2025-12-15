module github.com/jo-hoe/google-sheets

go 1.24.2

require golang.org/x/oauth2 v0.34.0

require (
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

retract v1.0.0 // published erroneously
