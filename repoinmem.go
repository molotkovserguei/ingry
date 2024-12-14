package main

type repoInMem struct {
	books *booksInMem
}

func (r *repoInMem) Books() BooksDB {
	return r.books
}

func NewRepoInMem() (*repoInMem, error) {
	return &repoInMem{
		books: NewBooksInMem(),
	}, nil
}
