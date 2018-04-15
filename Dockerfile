# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM flat-migration

# Copy the local package files to the container's workspace.
ADD ./config.json .

# Build command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)

# to be removed after merge to master
RUN cd $GOPATH/src/flat-search; make pre-install; make install

# Run command by default when the container starts.
ENTRYPOINT while true; do flat-search -config=$GOPATH/config.json; sleep 600; done

