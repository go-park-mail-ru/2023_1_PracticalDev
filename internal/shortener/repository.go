package shortener

type ShortenerRepository interface {
	Get(hash string) (string, error)
	Create(url string) (string, error)
}
