# Techinical and Implementation Details

> This codebase is the central meta organization managing server for the Corpora collection project by IIIT Hyderabad's LTRC, Speech Processing Lab

We plan on implementing it in golang but golang might be a tough choice as it makes it hard to maintain it. Most people know python but not golang.

- Might use grpc for go/python interop
    - [Example](https://github.com/Jigar3/gRPC_101/blob/master/server/server.go)

Will add in details as I implement the features.

The server uses [cobra](https://github.com/spf13/cobra/) cli tool for command generation

Bash completions also work