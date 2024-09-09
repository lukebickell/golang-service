package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	// RESTy routes for "articles" resource
	r.Route("/articles", func(r chi.Router) {
		// r.With(paginate).Get("/", listArticles)                           // GET /articles
		// r.With(paginate).Get("/{month}-{day}-{year}", listArticlesByDate) // GET /articles/01-16-2017

		// r.Post("/", createArticle)                                        // POST /articles
		// r.Get("/search", searchArticles)                                  // GET /articles/search

		// Regexp url parameters:
		// r.Get("/{articleSlug:[a-z-]+}", getArticleBySlug)                // GET /articles/home-is-toronto

		// Subrouters:
		r.Route("/{articleID}", func(r chi.Router) {
			r.Use(ArticleCtx)
			// r.Get("/", getArticle) // GET /articles/123
			// r.Put("/", updateArticle)                                       // PUT /articles/123
			// r.Delete("/", deleteArticle)                                    // DELETE /articles/123
		})
	})

	http.ListenAndServe(":3000", r)
}

func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID := chi.URLParam(r, "articleID")
		article, err := dbGetArticle(articleID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func getArticle(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	article, ok := ctx.Value("article").(*Article)
// 	if !ok {
// 		http.Error(w, http.StatusText(422), 422)
// 		return
// 	}
// 	w.Write([]byte(fmt.Sprintf("title:%s", article.Title)))
// }

func dbGetArticle(id string) (string, error) {
	return "article", nil
}
