# ncomm

`ncomm` is like [`comm`](https://en.wikipedia.org/wiki/Comm), but for any number
of files, not just two files.

## Installation

```bash
go install github.com/ucarion/ncomm
```

## Usage

To compare the contents of two (sorted) files `a` and `b`, run:

```bash
ncomm a b
```

The output will be exactly the same as `comm a b`. Unlike `comm`, you can pass
additional files. To compare sorted files `a`, `b`, `c`, `d`, run:

```bash
ncomm a b c d
```

Whereas `ncomm a b` gave you three columns of output, `ncomm a b c d` will
produce 15 columns; each column represents a different combination of `a`, `b`,
`c`, or `d` that a line might appear in.

To make it more obvious what the columns mean, use `--show-header` for a legend:

```bash
ncomm --show-header a b c d
```

```text
x---	-x--	xx--	--x-	x-x-	-xx-	xxx-	---x	x--x	-x-x	xx-x	--xx	x-xx	-xxx	xxxx	
[...]
```

You can also run `ncomm --help` for detailed information:

```text
usage: ncomm [<options>] files...

like comm(1), but for any number of files

        --show-header    output a header to help explain what each column represents
    -h, --help           display this help and exit

```

## Example

Here's a more fleshed-out example of how `ncomm` can start to be useful in the
real world. In `examples` in this repo, there's a set of files:

* `eu.txt` lists members of the European Union
* `nato.txt` lists members of the North Atlantic Treaty Organization
* `oecd.txt` lists members of the Organisation for Economic Co-operation and Development

All three of these are clubs of countries. With `ncomm`, you can see how they overlap:

```bash
ncomm --show-header ./examples/eu.txt ./examples/nato.txt ./examples/oecd.txt
```

```text
x--	-x-	xx-	--x	x-x	-xx	xxx	
	Albania
			Australia
				Austria
						Belgium
		Bulgaria
					Canada
			Chile
			Colombia
			Costa Rica
		Croatia
						Czech Republic
						Denmark
						Estonia
				Finland
						France
						Germany
						Greece
						Hungary
					Iceland
				Ireland
			Israel
						Italy
			Japan
			Korea
						Latvia
						Lithuania
						Luxembourg
Malta
			Mexico
	Montenegro
						Netherlands
			New Zealand
	North Macedonia
					Norway
						Poland
						Portugal
Republic of Cyprus
		Romania
			Slovak Republic
		Slovakia
						Slovenia
						Spain
				Sweden
			Switzerland
					Turkey
					United Kingdom
					United States
```

The first row comes from `--show-header`, and makes it a bit easier to figure
out what each column represents. For instance, we can see that the 5th columns
has a header `x-x`, indicating lines that are in the 1st and 3rd file but not
the 2nd. In our example, that represents countries that are in the EU and OECD,
but not NATO. We can extract just those countries with a bit of `cut`, and then
remove blank lines with `grep`:

```bash
ncomm ./examples/eu.txt ./examples/nato.txt ./examples/oecd.txt | cut -s -f 5 | grep .
```

```text
Austria
Finland
Ireland
Sweden
```

Alternatively, if we wanted to count how many countries fell into each of the
seven possible "buckets", you could count how many tabs there are on each line,
and then count how many lines have that number of tabs:

```bash
ncomm ./examples/eu.txt ./examples/nato.txt ./examples/oecd.txt | awk -F$'\t' '{print NF}' | sort | uniq -c
```

```text
   2 1
   3 2
   4 3
  11 4
   4 5
   6 6
  17 7
```

As before, you could use `--show-header` to make sense of what each column
number means.
