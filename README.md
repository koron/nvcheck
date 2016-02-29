# Notation Variability Checker

Install with:

    $ go get -u github.com/koron/nvcheck

Next create a dict.yml file for your documents.

Then you can check variability like this:

    $ nvcheck your.txt

## Dictionary Examples

*   [vimdoc-ja-working](https://github.com/vim-jp/vimdoc-ja-working/blob/master/dict.yml)

## Replace Words

With `-r` option, nvcheck replace all words to correct and output to stdout.

    $ nvcheck -r your.txt

With `-i` option, nvcheck overwrite the file to correct.

    $ nvcheck -i your.txt

## LICENSE

MIT license.  See LICENSE.
