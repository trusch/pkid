package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/trusch/pkid/generator"
	"github.com/trusch/pkid/manager"
	"github.com/trusch/pkid/types"
)

type Server struct {
	mgr    manager.Manager
	ln     net.Listener
	server *http.Server
}

type entityType string

const (
	clientType entityType = "client"
	serverType entityType = "server"
	caType     entityType = "ca"
)

func New(addr string, mgr manager.Manager) *Server {
	srv := &http.Server{
		Addr:           addr,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server := &Server{mgr, nil, srv}
	server.constructRouter()
	return server
}

func (srv *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", srv.server.Addr)
	if err != nil {
		return err
	}
	srv.ln = ln
	return srv.server.Serve(ln)
}

func (srv *Server) Stop() error {
	return srv.ln.Close()
}

func (srv *Server) constructRouter() {
	router := mux.NewRouter()
	router.Path("/ca").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleCreateSelfSignedCA(w, r)
	})
	router.Path("/ca/{ca}/{typ}").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleCreateSigned(w, r)
	})
	router.Path("/ca/{ca}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleGetCA(w, r, "client")
	})
	router.Path("/ca/{ca}/client").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleList(w, r, "client")
	})
	router.Path("/ca/{ca}/server").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleList(w, r, "server")
	})
	router.Path("/ca/{ca}/ca").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleList(w, r, "ca")
	})
	router.Path("/ca/{ca}/{typ}/{id}/cert").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleGetCert(w, r)
	})
	router.Path("/ca/{ca}/{typ}/{id}/key").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleGetKey(w, r)
	})
	router.Path("/ca/{ca}/{typ}/{id}/revoke").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleRevoke(w, r)
	})
	router.Path("/ca/{ca}/cert").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleGetCACert(w, r)
	})
	router.Path("/ca/{ca}/key").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleGetCAKey(w, r)
	})
	router.Path("/ca/{ca}/crl").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.handleGetCRL(w, r)
	})
	srv.server.Handler = router
}

func (srv *Server) handleCreateSelfSignedCA(w http.ResponseWriter, r *http.Request) {
	options, err := srv.parseCreateOptionsFromRequest(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	id, err := srv.mgr.CreateCA("", options)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(id))
}

func (srv *Server) handleCreateSigned(w http.ResponseWriter, r *http.Request) {
	options, err := srv.parseCreateOptionsFromRequest(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	vars := mux.Vars(r)
	ca := vars["ca"]
	var id string
	switch entityType(vars["typ"]) {
	case caType:
		id, err = srv.mgr.CreateCA(ca, options)
	case clientType:
		id, err = srv.mgr.CreateClient(ca, options)
	case serverType:
		id, err = srv.mgr.CreateServer(ca, options)
	}
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(id))
}

func (srv *Server) handleGetCert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ca := vars["ca"]
	caEntity, err := srv.mgr.GetCA(ca)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var entity *types.Entity
	switch entityType(vars["typ"]) {
	case caType:
		{
			if _, ok := caEntity.CAs[vars["id"]]; !ok {
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
			ent, e := srv.mgr.GetCA(vars["id"])
			if ent != nil {
				entity = ent.Entity
			}
			if e != nil {
				err = e
			}
		}
	case clientType:
		{
			if _, ok := caEntity.Clients[vars["id"]]; !ok {
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
			entity, err = srv.mgr.GetClient(vars["id"])
		}
	case serverType:
		{
			if _, ok := caEntity.Servers[vars["id"]]; !ok {
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
			entity, err = srv.mgr.GetServer(vars["id"])
		}
	}
	if err != nil || entity == nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(entity.Cert))
}

func (srv *Server) handleGetKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ca := vars["ca"]
	caEntity, err := srv.mgr.GetCA(ca)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var entity *types.Entity
	switch entityType(vars["typ"]) {
	case caType:
		{
			if _, ok := caEntity.CAs[vars["id"]]; !ok {
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
			ent, e := srv.mgr.GetCA(vars["id"])
			if ent != nil {
				entity = ent.Entity
			}
			if e != nil {
				err = e
			}
		}
	case clientType:
		{
			if _, ok := caEntity.Clients[vars["id"]]; !ok {
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
			entity, err = srv.mgr.GetClient(vars["id"])
		}
	case serverType:
		{
			if _, ok := caEntity.Servers[vars["id"]]; !ok {
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}
			entity, err = srv.mgr.GetServer(vars["id"])
		}
	}
	if err != nil || entity == nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(entity.Key))
}

func (srv *Server) handleGetCAKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ca := vars["ca"]
	caEntity, err := srv.mgr.GetCA(ca)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(caEntity.Key))
}

func (srv *Server) handleGetCACert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ca := vars["ca"]
	caEntity, err := srv.mgr.GetCA(ca)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(caEntity.Cert))
}

func (srv *Server) handleGetCRL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ca := vars["ca"]
	crl, err := srv.mgr.GetCRL(ca)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(crl))
}

func (srv *Server) handleRevoke(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ca := vars["ca"]
	typ := vars["typ"]
	id := vars["id"]
	var err error
	switch entityType(typ) {
	case caType:
		err = srv.mgr.RevokeCA(ca, id)
	case clientType:
		err = srv.mgr.RevokeClient(ca, id)
	case serverType:
		err = srv.mgr.RevokeServer(ca, id)
	}
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("revoked"))
}

func (srv *Server) handleList(w http.ResponseWriter, r *http.Request, typ string) {
	vars := mux.Vars(r)
	ca := vars["ca"]
	caEntity, err := srv.mgr.GetCA(ca)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encoder := json.NewEncoder(w)

	switch entityType(typ) {
	case caType:
		encoder.Encode(caEntity.CAs)
	case clientType:
		encoder.Encode(caEntity.Clients)
	case serverType:
		encoder.Encode(caEntity.Servers)
	}
}

func (srv *Server) handleGetCA(w http.ResponseWriter, r *http.Request, typ string) {
	vars := mux.Vars(r)
	ca := vars["ca"]
	caEntity, err := srv.mgr.GetCA(ca)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	result := &types.CAEntity{
		Entity: &types.Entity{
			ID:        caEntity.Entity.ID,
			Name:      caEntity.Entity.Name,
			IsRevoked: caEntity.Entity.IsRevoked,
		},
		Revoked: caEntity.Revoked,
		Clients: caEntity.Clients,
		Servers: caEntity.Servers,
		CAs:     caEntity.CAs,
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

func (srv *Server) parseCreateOptionsFromRequest(r *http.Request) (*generator.Options, error) {
	options := &generator.Options{}
	if name := r.FormValue("name"); name != "" {
		options.Name = name
	} else {
		return nil, errors.New("Error in options parsing: no name given")
	}

	if rsaBitsStr := r.FormValue("rsaBits"); rsaBitsStr != "" {
		rsaBits, err := strconv.ParseInt(rsaBitsStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Error in options parsing: can not parse rsaBits (%v)", err)
		}
		options.RsaBits = int(rsaBits)
	}
	if curve := r.FormValue("curve"); curve != "" {
		options.Curve = curve
	}
	if notBeforeUnixStr := r.FormValue("notBefore"); notBeforeUnixStr != "" {
		notBeforeUnix, err := strconv.ParseInt(notBeforeUnixStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error in options parsing: can not parse notBefore (%v)", err)
		}
		options.NotBefore = time.Unix(notBeforeUnix, 0)
	}
	if validForStr := r.FormValue("validFor"); validForStr != "" {
		validFor, err := time.ParseDuration(validForStr)
		if err != nil {
			return nil, fmt.Errorf("Error in options parsing: can not parse validFor (%v)", err)
		}
		options.ValidFor = validFor
	}
	return options, nil
}
