package cache

import "k8s.io/apimachinery/pkg/util/sets"

type IndexFunc func(obj any) ([]string, error)

type Index map[string]sets.String

type Indices map[string]Index

type Indexers map[string]IndexFunc
