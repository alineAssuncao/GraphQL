# GraphQL
https://gqlgen.com/

## Schema usado no GraphQL

mutation createCategory {
  createCategory (input: {name: "Tecnologia", description: "Cursos de Tecnologia"}) {
    id
    name
    description
  }
}

mutation createCourse {
  createCourse (input: {name: "Full Cycle", description: "The best", categoryId: "0c1737fc-19d3-4753-b304-eaab75b8b698"}) {
    id
    name
    description
  }
}

query queryCategories{
  categories {
    id
    name
    description
  }
}

query queryCategoriesWithCourses{
  categories {
    id
    name
    courses {
      id
      name
    }
  }
}

query queryCourses{
  courses {
    id
    name
    description
  }
}

query queryCoursesWithCategory{
  courses {
    id
    name
    description
    category {
      id
      name
      description
    }
  }
}


## rodando
go run server.go