Releasing
=========

 1. Update the version in `cmd/http_to_nsq/main.go`.
 2. Update the version in `Dockerfile`.
 3. Update the `History.md` with the new release.
 4. `git commit -am "Release X.Y.Z."` (where X.Y.Z is the new version).
 5. `git tag -a X.Y.Z -m "Version X.Y.Z"` (where X.Y.Z is the new version).
 6. `git push && git push --tags`
 7. Build the binaries `gox -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./...`.
 8. Upload the binaries to [Github](https://github.com/segmentio/http_to_nsq/releases).
 9. Build the Docker container `docker build -t segment/http_to_nsq:X.Y.Z .` (where X.Y.Z is the new version).
10. Publish the Docker container `docker push segment/http_to_nsq:X.Y.Z` (where X.Y.Z is the new version).
