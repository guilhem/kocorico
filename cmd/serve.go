// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/guilhem/kocorico/cluster"
	"github.com/guilhem/kocorico/etcd"
	"github.com/guilhem/kocorico/render"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: serve,
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntP("port", "p", 8080, "Listen port")
	if err := viper.BindPFlag("port", serveCmd.Flags().Lookup("port")); err != nil {
		log.Fatal(err)
	}
}

func serve(cmd *cobra.Command, args []string) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	r.Route("/cluster", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("cluster"))
		})
		r.Post("/", createCluster)

		r.Route("/:clusterUUID", func(r chi.Router) {
			r.Use(ClusterCtx)
			r.Get("/", getCluster) // GET /articles/123
			// r.Put("/", updateArticle)              // PUT /articles/123
			// r.Delete("/", deleteArticle)           // DELETE /articles/123
		})
	})

	if err := http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("port")), r); err != nil {
		log.Fatal(err)
	}
}

const clusterContext = "cluster"

func getCluster(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cluster, ok := ctx.Value(clusterContext).(cluster.Cluster)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	w.Write([]byte(fmt.Sprintf("title:%s", cluster.Name)))
}

// CreateArticle persists the posted Article and returns it
// back to the client as an acknowledgement.
func createCluster(w http.ResponseWriter, r *http.Request) {
	cluster := cluster.Cluster{}
	// ^ the above is a nifty trick for how to omit fields during json unmarshalling
	// through struct composition

	if err := render.YAMLBind(r.Body, &cluster); err != nil {
		render.YAML(w, r, err.Error())
		return
	}

	etcd.CreateCluster(cluster)

	render.YAML(w, r, cluster)
}

// ClusterCtx middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func ClusterCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clusterUUID := chi.URLParam(r, "clusterUUID")
		c, err := etcd.GetCluster(clusterUUID)
		if err != nil {
			//render.Status(r, http.StatusNotFound)
			render.YAML(w, r, http.StatusText(http.StatusNotFound))
			return
		}
		ctx := context.WithValue(r.Context(), clusterContext, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
