module github.com/jo-hoe/google-sheets

go 1.19

require golang.org/x/oauth2 v0.30.0

require (
	cloud.google.com/go/compute/metadata v0.5.2 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

retract v1.0.0 // published erroneously
