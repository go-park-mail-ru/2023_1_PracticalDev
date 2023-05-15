package shortener

type ShortenerService interface {
	Get(hash string) (string, error)
	Create(url string) (string, error)
}
