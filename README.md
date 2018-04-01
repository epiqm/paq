# Paq

Pack files/directories into a single binary.

## Getting Started

Paq is a command line utility that uses arguments for preforming tasks.

### Usage

Create a package using files in tests directory:

![tests directory files](https://epiqm.github.io/static/img/paq-files.png "tests directory files")

```
$ ./paq ./tests/* -o testarchive.pq
done.
```

Check newly created package named testarchive.pq:

![test package](https://epiqm.github.io/static/img/paq-archive.png "test package")

```
$ ls -i | grep "test"
4965616 testarchive.pq
```

Contents of testarchive.pq:

![package open in dhex](https://epiqm.github.io/static/img/paq-dhex.png "package open in dhex")

Unpack testarchive.pq package:

```
$ ./paq unpack testarchive.pq
testarchive.pq:
  note.txt (109 bytes, offset 46)
  rabbit.png (18170 bytes, offset 155)
done.
```

Check rabbit.png with image viewer or other application you use for browsing images.

After successful unpack you are able to see a pink rabbit picture.

Also check the text file named note.txt, content should be equal to:

![note text](https://epiqm.github.io/static/img/paq-notetxtbb.png "note text")

Delete files used for this example:

```
$ rm ./testarchive.pq ./note.txt ./rabbit.png
```

Unpack several packages into a single directory:

```
$ ./paq unpack package.pq package.pq2
package.pq:
  note.txt (109 bytes, offset 46)
  rabbit.png (18170 bytes, offset 155)
package.pq2:
  MIT.md (1070 bytes, offset 24)
done.
```

### Building

Install Go >=1.8, manage projects path environment variable, clone the repository to src directory.

```
$ cd ~/go/src
$ git clone git@github.com:epiqm/paq.git
$ cd ./paq
```

Build the utility.

```
$ go build
```

The go get command is not needed for building. It is written using native Golang imports.

After successful build run Paq with -v argument to display the version.

```
$ ./paq -v
0.1
```

## Versioning

Paq pack/unpack mechanism is compatible with packages created with previous releases of utility.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Author

Written by [Maxim R.](https://epiqm.github.io/)
