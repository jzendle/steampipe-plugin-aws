
install:
	#go build -o  ~/.steampipe/plugins/hub.steampipe.io/plugins/turbot/aws@latest/steampipe-plugin-aws.plugin  *.go
	go build -o  bin/steampipe-plugin-aws.plugin  main.go

clean:
	rm bin/*


