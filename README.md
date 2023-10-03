# catbox

`catbox` will pipe every message from an mbox file as an input to a shell
command. A message counter is available as a shell variable $NR.

If no file is specified, `catbox` will read from stdin.

Inspired by [caeml](https://github.com/ferdinandyb/caeml/).

### Installation

```sh
go install github.com/konimarti/catbox@latest
```

### Usage

Usage: `catbox [-h|-c <cmd>] <mbox>`

### Integration with aerc

Add this line to the [filters] section in your aerc.conf:

```
application/mbox=catbox -c caeml | colorize
```

### Examples

The following examples assume that you have a file `test.mbox` in a valid mbox
format in your local directory.

- Show the message number counter:

```
catbox -c 'echo $NR' test.mbox`
```

- Pipe every mbox message to caeml:

```
catbox -c caeml test.mbox
```

- Save every mbox message in a separate file

```
catbox -c 'cat > message_$NR' test.mbox
```

- Print only the first ten messages:

```
catbox -c "awk -v cbnr=\$NR '{if (cbnr>10) print}'" test.mbox
```

- Read from stdin and only display the full message headers:

```
cat test.mbox | catbox -c "sed -n '1,/^\\s*$/p'"
```
