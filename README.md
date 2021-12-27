# clean_blog_server

A GraphQL / Postgres server managing a Blog CMS. Written in Golang with the Gin framework and schema-first code generation via GQLGen / SQLBoiler. For use with the clean_blog_client frontend based on the bootstrap template found at [startbootstrap clean blog].

### Core Dependencies

- [gin] - gin web framework
- [gqlgen] - schema-first GraphQL with codegen
- [sqlboiler] - type-safe ORM with DB introspection
- [go-redis] - redis for golang
- [postgres]- postgres for golang

[startbootstrap clean blog]: https://startbootstrap.com/theme/clean-blog
[gin]: https://gin-gonic.com
[gqlgen]: https://gqlgen.com
[sqlboiler]: https://github.com/volatiletech/sqlboiler#find
[go-redis]: https://github.com/go-redis/redis
[postgres]: https://github.com/lib/pq
