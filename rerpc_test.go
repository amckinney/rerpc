package rerpc_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/akshayjshah/rerpc"
	"github.com/akshayjshah/rerpc/internal/pingpb"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip" // registers gzip compressor
)

type Ping struct{}

func (p Ping) Ping(ctx context.Context, req *pingpb.PingRequest) (*pingpb.PingResponse, error) {
	return &pingpb.PingResponse{Number: req.Number}, nil
}

func (p Ping) Fail(_ context.Context, req *pingpb.FailRequest) (*pingpb.FailResponse, error) {
	return nil, nil
}

// FIXME: Generate this.
func NewPingHandler(s pingpb.PingServer) (string, http.Handler) {
	ping := rerpc.Handler{
		Implementation: func(ctx context.Context, req proto.Message) (proto.Message, error) {
			typed, ok := req.(*pingpb.PingRequest)
			if !ok {
				return nil, fmt.Errorf("can't call rerpc.ping.v0.Ping/Ping with a %T", req)
			}
			return s.Ping(ctx, typed)
		},
	}
	fail := rerpc.Handler{
		Implementation: func(ctx context.Context, req proto.Message) (proto.Message, error) {
			typed, ok := req.(*pingpb.FailRequest)
			if !ok {
				return nil, fmt.Errorf("can't call rerpc.ping.v0.Ping/Fail with a %T", req)
			}
			return s.Fail(ctx, typed)
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/rerpc.ping.v0.Ping/Ping", func(w http.ResponseWriter, r *http.Request) {
		ping.Serve(w, r, &pingpb.PingRequest{})
	})
	mux.HandleFunc("/rerpc.ping.v0.Ping/Fail", func(w http.ResponseWriter, r *http.Request) {
		fail.Serve(w, r, &pingpb.FailRequest{})
	})

	return "/rerpc.ping.v0.Ping/", mux
}

func NewPingClient(doer rerpc.Doer) pingpb.PingClient {
	return nil
}

type client struct {
	doer rerpc.Doer
}

func (c *client) Ping(ctx context.Context, req *pingpb.PingRequest) (*pingpb.PingResponse, error) {
	return nil, nil
}

func (c *client) Fail(ctx context.Context, req *pingpb.FailRequest) (*pingpb.FailResponse, error) {
	return nil, nil
}

func (c *client) Call(doer rerpc.Doer, ctx context.Context, req proto.Message, res proto.Message) error {
	return nil
}

func TestServerGRPC(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(NewPingHandler(Ping{}))

	num := rand.Int63()
	req := &pingpb.PingRequest{Number: num}

	assertPinged := func(t testing.TB, res *pingpb.PingResponse, err error) {
		if err != nil {
			t.Errorf("request failed: %v", err)
		} else if res.Number != num {
			t.Errorf("ping didn't echo input number")
		}
	}
	assertPingIdentity := func(t *testing.T, client pingpb.PingClient) {
		t.Run("identity", func(t *testing.T) {
			res, err := client.Ping(context.Background(), req)
			assertPinged(t, res, err)
		})
	}
	assertPingGzip := func(t *testing.T, client pingpb.PingClient) {
		t.Run("gzip", func(t *testing.T) {
			res, err := client.Ping(context.Background(), req, grpc.UseCompressor("gzip"))
			assertPinged(t, res, err)
		})
	}
	logBody := func(t testing.TB, response *http.Response) {
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			t.Fatalf("couldn't dump response: %v", err)
		}
		t.Log(string(dump))
	}
	_ = logBody

	t.Run("json", func(t *testing.T) {
		probe := `{"number":"42"}`
		server := httptest.NewServer(mux)
		defer server.Close()

		r, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/rerpc.ping.v0.Ping/Ping", server.URL),
			strings.NewReader(probe),
		)
		r.Header.Set("Content-Type", rerpc.TypeJSON)
		if err != nil {
			t.Fatalf("couldn't create request: %v", err)
		}
		response, err := server.Client().Do(r)
		if err != nil {
			t.Fatalf("error making request: %v", err)
		}
		if response.ProtoMajor != 1 {
			t.Fatalf("expected HTTP/1, got ProtoMajor %v", response.ProtoMajor)
		}
		if response.StatusCode != http.StatusOK {
			t.Fatalf("unexpected HTTP status code %v", response.StatusCode)
		}
		if me := response.Header.Get("Grpc-Encoding"); me != rerpc.EncodingIdentity {
			t.Errorf("expected grpc-encoding %q, got %q", rerpc.EncodingIdentity, me)
		}
		if response.Header.Get("Grpc-Accept-Encoding") == "" {
			t.Log("got empty grpc-accept-encoding header")
			t.Fail()
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("couldn't read response body: %v", err)
		}
		if string(body) != probe {
			t.Fatalf("expected %q, got %q", probe, string(body))
		}
	})

	t.Run("http1", func(t *testing.T) {
		// TODO replace with reRPC's own HTTP client?
		server := httptest.NewServer(mux)
		defer server.Close()

		body := &bytes.Buffer{}
		if err := rerpc.MarshalLPM(body, req, rerpc.EncodingIdentity); err != nil {
			t.Fatalf("couldn't write LPM: %v", err)
		}
		request, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/rerpc.ping.v0.Ping/Ping", server.URL),
			body,
		)
		if err != nil {
			t.Fatalf("couldn't create request: %v", err)
		}
		request.Header.Set("Content-Type", rerpc.TypeDefaultGRPC)
		request.Header.Set("Te", "trailers")

		response, err := server.Client().Do(request)
		if err != nil {
			t.Fatalf("error sending request: %v", err)
		}
		defer io.Copy(ioutil.Discard, response.Body) // TODO: necessary?
		defer response.Body.Close()

		if response.ProtoMajor != 1 {
			t.Fatalf("expected HTTP/1, got ProtoMajor %v", response.ProtoMajor)
		}
		if response.StatusCode != http.StatusOK {
			t.Fatalf("unexpected HTTP status code %v", response.StatusCode)
		}
		if me := response.Header.Get("Grpc-Encoding"); me != rerpc.EncodingIdentity {
			t.Errorf("expected grpc-encoding %q, got %q", rerpc.EncodingIdentity, me)
		}
		if response.Header.Get("Grpc-Accept-Encoding") == "" {
			t.Log("got empty grpc-accept-encoding header")
			t.Fail()
		}

		msg := &pingpb.PingResponse{}
		if err := rerpc.UnmarshalLPM(response.Body, msg, rerpc.EncodingIdentity); err != nil {
			t.Fatalf("can't unmarshal proto response: %v", err)
		}

		assertPinged(t, msg, nil)
	})

	t.Run("http2", func(t *testing.T) {
		server := httptest.NewTLSServer(h2c.NewHandler(
			mux,
			&http2.Server{},
		))
		defer server.Close()

		pool := x509.NewCertPool()
		pool.AddCert(server.Certificate())
		gconn, err := grpc.Dial(server.Listener.Addr().String(), grpc.WithTransportCredentials(
			credentials.NewClientTLSFromCert(pool, "" /* server name */),
		))
		if err != nil {
			t.Fatalf("gRPC client can't connect: %v", err)
		}
		client := pingpb.NewPingClient(gconn)

		assertPingIdentity(t, client)
		assertPingGzip(t, client)
	})

	t.Run("h2c", func(t *testing.T) {
		server := httptest.NewServer(h2c.NewHandler(
			mux,
			&http2.Server{},
		))
		defer server.Close()

		gconn, err := grpc.Dial(server.Listener.Addr().String(), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("gRPC client can't connect: %v", err)
		}
		client := pingpb.NewPingClient(gconn)

		assertPingIdentity(t, client)
		assertPingGzip(t, client)
	})

	t.Run("exploration", func(t *testing.T) {
		// Explore grpc-go's server behavior.
		gs := grpc.NewServer()
		pingpb.RegisterPingServer(gs, &Ping{})
		server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor != 2 {
				t.Fatalf("non-HTTP/2 request")
				return
			}
			gs.ServeHTTP(w, r)
		}))
		server.TLS = &tls.Config{
			CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
			NextProtos:   []string{http2.NextProtoTLS},
		}
		server.StartTLS()
		defer server.Close()

		pool := x509.NewCertPool()
		pool.AddCert(server.Certificate())
		client := &http.Client{Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool},
		}}

		body := &bytes.Buffer{}
		if err := rerpc.MarshalLPM(body, req, rerpc.EncodingIdentity); err != nil {
			t.Fatalf("couldn't write LPM: %v", err)
		}
		request, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/rerpc.ping.v0.Ping/Ping", server.URL),
			body,
		)
		if err != nil {
			t.Fatalf("couldn't create request: %v", err)
		}
		request.Header.Set("Content-Type", rerpc.TypeDefaultGRPC)
		request.Header.Set("Te", "trailers")

		response, err := client.Do(request)
		if err != nil {
			t.Fatalf("error sending request: %v", err)
		}
		defer io.Copy(ioutil.Discard, response.Body) // TODO: necessary?
		defer response.Body.Close()

		if response.ProtoMajor != 2 {
			t.Fatalf("expected HTTP/1, got ProtoMajor %v", response.ProtoMajor)
		}
		if response.StatusCode != http.StatusOK {
			logBody(t, response)
			t.Fatalf("unexpected HTTP status code %v", response.StatusCode)
		}

		// logBody(t, response)
		msg := &pingpb.PingResponse{}
		if err := rerpc.UnmarshalLPM(response.Body, msg, rerpc.EncodingIdentity); err != nil {
			t.Fatalf("can't unmarshal proto response: %v", err)
		}

		assertPinged(t, msg, nil)
	})
}
