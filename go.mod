module github.com/jritsema/go-htmx-starter

go 1.20

require (
	github.com/jritsema/gotoolbox v0.8.0
	github.com/open-amt-cloud-toolkit/go-wsman-messages v0.0.0-00010101000000-000000000000
	go.etcd.io/bbolt v1.3.8
)

replace github.com/open-amt-cloud-toolkit/go-wsman-messages => ../go-wsman-messages

require golang.org/x/sys v0.4.0 // indirect
