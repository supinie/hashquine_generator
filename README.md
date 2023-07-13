# Go md5 Hashquine Generator

This is a script to generate a gif that will display its own md5sum.

It is a Go implementation of [Rogdham's python script](https://github.com/rogdham/gif-md5-hashquine). The script works by generating a collision for all possible characters at each possible position, for both visable and hidden. Once we have finished, we can then selectively choose to display characters based on the md5sum of the file, without changing it.

Here is an example generation:

![hashquine](./hashquine.gif)

which can be seen to be displaying its own md5 hash:

```
┌[supinie@ubuntuVM] [/dev/pts/2] [main ⚡] 
└[~/git/hashquine]> md5sum hashquine.gif 
18f6d7bae1ca24b0ea3224560f046cd9  hashquine.gif
```

## Prerequisites:

You must have fastcoll, this can be installed by the following:
```
$ git clone git@github.com:cr-marcstevens/hashclash.git
$ sudo apt-get install g++ autoconf automake libtool && sudo apt-get install zlib1g-dev libbz2-dev
$ cd hashclash && ./build.sh
$ cp bin/md5_fastcoll /usr/bin/fastcoll
$ cd .. && rm -rf hashclash     # optional
```

