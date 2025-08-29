#!/bin/bash

(cd ../cmd/templ && go install .) && \
    (cd fork && templ generate) && \
    (cd orig && go run github.com/a-h/templ/cmd/templ@v0.3.943 generate)