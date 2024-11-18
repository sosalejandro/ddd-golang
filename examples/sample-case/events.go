package main

// addBookEvent is a domain event that represents adding a book to the library.
type addBookEvent struct {
	Title string
}

// removeBookEvent is a domain event that represents removing a book from the library.
type removeBookEvent struct {
	Title string
}
