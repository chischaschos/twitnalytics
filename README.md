# twitnalytics

Analyze twitter user streams. By the way it is not finished yet lol.

# installation
Ideally installation is achieved by:

    go get https://github.com/chischaschos/twitnalytics
    export CONSUMER_KEY=ASDSADASD
    export CONSUMER_SECRET=ASDSADASD
    twitnalytics -u chischaschos,softr8,hecbuma

I haven't tested it yet, so just clone, and then run

    export CONSUMER_KEY=ASDSADASD
    export CONSUMER_SECRET=ASDSADASD
    go run main.go -u chischaschos,softr8,hecbuma

This won't produce the expected output either, but, you will now have a
terms table fully populated where you will be able to play with
twitter's term counts:

    cd $HOME
    sqlite3 .twitnalytics-db
    sqlite> SELECT a.term, a.count * b.count FROM tweet_terms a INNER
      JOIN tweet_terms b ON a.term = b.term AND a.username <> b.username
      ORDER by a.count*b.count desc;
