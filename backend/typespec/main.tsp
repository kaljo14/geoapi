import "@typespec/http";

using TypeSpec.Http;

@service({
  title: "Geopulse Backend",
})
@server("http://localhost:8081", "Local development server")
namespace GeopulseBackend;

model MessageResponse {
  message: string;
}

model ErrorResponse {
  error: string;
}

model LivezResponse {
  status: string;
}

model ReadyzResponse {
  status: string;
  db: string;
}

@route("/scrape")
interface Scrape {
  @post
  @summary("Trigger scraper")
  @doc("Starts a background Google Places scraper job.")
  create(): {
    @statusCode statusCode: 202;
    @body body: MessageResponse;
  };
}

@route("/enrich")
interface Enrich {
  @post
  @summary("Trigger enricher")
  @doc("Starts a background Google Places details enricher job.")
  create(): {
    @statusCode statusCode: 202;
    @body body: MessageResponse;
  };
}

@route("/livez")
interface Livez {
  @get
  @summary("Liveness probe")
  @doc("Always returns 200 OK.")
  get(): LivezResponse;
}

@route("/readyz")
interface Readyz {
  @get
  @summary("Readiness probe")
  @doc("Returns 200 if the service is ready. Returns 503 otherwise.")
  get(): ReadyzResponse | {
    @statusCode statusCode: 503;
    @body body: ReadyzResponse;
  };
}
