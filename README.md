# vectorugo [![Build Status](https://travis-ci.com/Comonut/vectorugo.svg?branch=master)](https://travis-ci.com/Comonut/vectorugo)

Description
===============

A vector store built in Go that can be used for storing and retrieving vectors by ids and for querying for K Nearest Neighbours. Saving and retrieving vectors works on a key-value principle where keys are strings and values are float arrays and it supports both in memory and on disk storage. The KNN search uses an Approximate Nearest Neighbour index that is built incrementally, so adding new objects doesn't mean it needs to be rebuilt. For now the only way to interact is throught the provided REST API.


Usage
===========

## Running
You will need go to build it from source. Once built you can start it in memory using `./main` and for disk you will need to provide the storage names (which will be then used as prefix for the storage files) and the dimension of the vectors that will be used, ex `./main -dim 300 -name words -persistance`

## REST API

To insert vectors, you need to use POST on `/vectors` with a json containing vector ids as fields with their corresponding arrays ex.
```json
{
  "v1" : [0.0,1.0],
  "v2" : [1.0,0.0]
}
```

To retrieve a vector's value by it's id you can use a GET on `/vectors` with the vector id as `id` query paramater.
Example for `/vectors?id=v1` would be:

```json
{ 
   "ID":"v1",
   "Array":[ 
      0.0,
      1.0
   ]
}
```
To use KNN search you can use a GET on `/search` and you can specify the target with the `id` query parameter and the number of results with the `k` query paramater. 
Here is the result of `/search?k=5&id=news` on a store containing 30k word embeddings
```json
[ 
   { 
      "ID":"news",
      "Distance":0
   },
   { 
      "ID":"press",
      "Distance":0.9802878802581415
   },
   { 
      "ID":"media",
      "Distance":1.002220129130821
   },
   { 
      "ID":"reports",
      "Distance":1.0628062539955245
   },
   { 
      "ID":"reporters",
      "Distance":1.0690725863443509
   }
]
```

### Disclaimer
While this projects can be useful, it's far from being production grade - it is still a WIP and there are a couple of bugs. 
