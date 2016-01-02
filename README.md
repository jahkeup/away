# Away

Away is a simplistic linking utility designed to operate in a similar
fashion to GNU Stow. Unlike Stow, this program does not require Perl
as its statically compiled.

## Running

```sh
away .dotfiles/xmonad
```

## Compile

The process is the usual for golang based utilities.

```sh
go get github.com/jahkeup/away

go install github.com/jahkeup/away
```
