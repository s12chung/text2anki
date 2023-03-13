package text2anki.tokenizer.komoran;

import text2anki.tokenizer.komoran.Tokenizer;
import org.json.JSONObject;

import com.sun.net.httpserver.Filter;
import com.sun.net.httpserver.Headers;
import com.sun.net.httpserver.HttpContext;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import com.sun.net.httpserver.spi.HttpServerProvider;

import java.io.IOException;
import java.io.InputStreamReader;
import java.io.Reader;
import java.net.InetSocketAddress;
import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.function.Function;
import java.util.logging.Logger;

public class Server {
  public static void main(String[] args) throws IOException {
    if (args.length == 2) {
      new Server(Integer.parseInt(args[0]), Integer.parseInt(args[1])).start();
      return;
    }
    new Server().start();
  }

  public static String tokenizeKey = "string";
  public static String pathHealthz = "/healthz";
  public static String pathTokenize = "/tokenize";

  public static int defaultPort = 9999;
  public static int defaultBacklog = 64;

  HttpServer server;

  public Server() throws IOException {
    this(defaultPort, defaultBacklog);
  }

  public Server(int port, int backlog) throws IOException {
    HttpServerProvider provider = HttpServerProvider.provider();
    server = provider.createHttpServer(new InetSocketAddress(port), backlog);

    HttpContext context = server.createContext(pathHealthz);
    context.getFilters().add(new TracingFilter());
    context.setHandler(respondWith(req -> HttpResponse.ok("ok\n\n" + Instant.now())));

    context = server.createContext(pathTokenize);
    context.getFilters().add(new TracingFilter());
    context.setHandler(respondWith((HttpRequest req) -> {
      if (!req.method().equals("POST")) {
        return HttpResponse.notFound();
      }

      String string;
      try {
        string = new JSONObject(req.bodyString()).getString(tokenizeKey);

      } catch (Exception e) {
        return HttpResponse.unprocessableEntity("Invalid JSON or key not found: " + tokenizeKey);
      }

      var tokens = Tokenizer.getTokens(string);
      return HttpResponse.okJSON(new JSONObject().put("tokens", tokens));
    }));
  }

  public void start() {
    server.start();
  }

  public void stop(int delay) {
    server.stop(delay);
  }

  static HttpHandler respondWith(HttpFunc hf) {
    return exchange -> {
      var req = HttpRequest.of(exchange);
      var res = hf.apply(req);

      var bytes = res.body().getBytes(StandardCharsets.UTF_8);
      exchange.getResponseHeaders().putAll(res.headers());

      try {
        exchange.sendResponseHeaders(res.status(), bytes.length);
        try (var os = exchange.getResponseBody()) {
          os.write(bytes);
        }
      } catch (IOException e) {
        e.printStackTrace();
      }
    };
  }

  interface HttpFunc extends Function<HttpRequest, HttpResponse> {
  }

  static class TracingFilter extends Filter {
    private final Logger LOG = Logger.getLogger(TracingFilter.class.getName());

    @Override
    public void doFilter(HttpExchange exchange, Chain chain) throws IOException {
      var req = HttpRequest.of(exchange);
      LOG.info(req.toString());
      chain.doFilter(exchange);
    }

    @Override
    public String description() {
      return "Trace";
    }
  }

  record HttpRequest(String method, URI requestUri, Headers headers, HttpExchange exchange) {
    static HttpRequest of(HttpExchange exchange) {
      return new HttpRequest(
          exchange.getRequestMethod(),
          exchange.getRequestURI(),
          exchange.getRequestHeaders(),
          exchange);
    }

    String bodyString() {
      int bufferSize = 1024;
      char[] buffer = new char[bufferSize];
      StringBuilder out = new StringBuilder();
      Reader in = new InputStreamReader(exchange().getRequestBody(), StandardCharsets.UTF_8);
      try {
        for (int numRead; (numRead = in.read(buffer, 0, buffer.length)) > 0;) {
          out.append(buffer, 0, numRead);
        }
      } catch (IOException e) {
        return "";
      }
      return out.toString();
    }
  }

  record HttpResponse(int status, Headers headers, String body) {
    static HttpResponse ok(String body) {
      return new HttpResponse(200, new Headers(), body).text();
    }

    static HttpResponse okJSON(JSONObject json) {
      return new HttpResponse(200, new Headers(), json.toString()).json();
    }

    static HttpResponse notFound() {
      return new HttpResponse(404, new Headers(), "404 Not Found").text();
    }

    static HttpResponse unprocessableEntity(String message) {
      return new HttpResponse(422, new Headers(), message).text();
    }

    HttpResponse header(String name, String value) {
      var res = new HttpResponse(status(), headers(), body());
      res.headers().add(name, value);
      return res;
    }

    private HttpResponse json() {
      return header("Content-type", "application/json");
    }

    private HttpResponse text() {
      return header("Content-type", "text/plain");
    }
  }
}
