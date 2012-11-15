jailer
======

Tiny golang application for parsing log files into json logs

Usage: jailer key[,key]=regexp [key[,key]=regexp]

Example:

    :$ echo 'INFO: [core_name] webapp=/solr path=/select' | jailer 'webapp,path=webapp=([^ ]+).*?path=([^ ]+)'
    { \"webapp\": \"/solr\", \"path\": \"/select\" }
