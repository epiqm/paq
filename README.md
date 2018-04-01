# Paq

Pack files/directories into a single binary.

## Getting Started

Paq is a command line utility that uses arguments for preforming tasks.

### Usage

Create a test package:

```
$ ./paq ./tests/* -o testarchive.pq
```

Check newly created package named testarchive.pq:

```
$ ls -i | grep "test"
4965616 testarchive.pq
```

Unpack testarchive.pq package:

```
$ ./paq unpack testarchive.pq
testarchive.pq:
  note.txt (109 bytes, offset 46)
  rabbit.png (18170 bytes, offset 155)
```

Check rabbit.png with image viewer or other application you use for browsing images.

After successful unpack you are able to see a pink rabbit picture.

Also check the text file named note.txt, content should be equal to:

```
A simple note to be packed.

For testing purposes.

Create a package and put this note and gif image inside.

```

Delete files used for this example:

```
$ rm ./testarchive.pq ./note.txt ./rabbit.png
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
