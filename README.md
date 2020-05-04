## What is this?

This is a gopher server written in golang.

## Why is this?

It felt appropriate.

## Is it secure?

Probably not, but I'm running it anyway.

## What does it support?

Right now it supports text documents and gophermaps.

## What is still on the todo list?

* Sanitizing content. The server currently assumes all of its input is valid.
* Images - I don't really need this for my use case, but it would be pretty
neat to have a go gopherhole full of gopher pictures.
* Compliance with rfc1436. I've tested it with a couple random clients,
but if it is actually compliant it was an accident.
* A basic multiplexer, probably similar to ServeMux in net/http.
