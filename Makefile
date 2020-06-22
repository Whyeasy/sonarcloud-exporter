build:
	go build -o bin/sonarcloud-exporter ./cmd/sonarcloud-exporter

deps:
	go mod verify
	go mod tidy -v

tag:
	git fetch --tags
	git tag $(TAG)
	git push origin $(TAG)

untag:
	git fetch --tags
	git tag -d $(TAG)
	git push origin :refs/tags/$(TAG)
	curl --request DELETE --header "Authorization: token ${GITHUB_TOKEN}" "https://api.github.com/repos/whyeasy/sonarcloud-exporter/releases/:release_id/$(TAG)"

verify-goreleaser:
ifeq (, $(shell which goreleaser))
	$(error "No goreleaser in $(PATH), consider installing it from https://goreleaser.com/install")
endif
	goreleaser --version

verify-docker:
ifeq (, $(shell which docker))
	$(error "No docker in $(PATH), consider installing it from https://docs.docker.com/install")
endif
	docker --version

release: verify-goreleaser verify-docker
	goreleaser release --rm-dist

snapshot-release: verify-goreleaser
	goreleaser --snapshot --skip-publish --rm-dist