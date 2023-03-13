import java.util.function.Function;

import text2anki.tokenizer.komoran.Tokenizer;
import text2anki.tokenizer.komoran.Server;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

class TestTokenizer {
    public static void main(String args[]) throws IOException {
        testTokenizer();

        Server server = new Server();
        try {
            server.start();
            testServer();
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            server.stop(5);
        }
    }

    static String basicTestString = "대한민국은";

    static void testTokenizer() {
        String tokens = Tokenizer.getTokens(basicTestString).toString();
        Assert.equalsString("Tokenizer Output", tokens,
                "[Token [morph=대한민국, pos=NNP, beginIndex=0, endIndex=4], Token [morph=은, pos=JX, beginIndex=4, endIndex=5]]");
    }

    static String testServerURI = "http://localhost:" + Server.defaultPort;

    static void testServer() throws URISyntaxException, IOException, InterruptedException {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(new URI(testServerURI + Server.pathHealthz))
                .GET()
                .build();

        String response = stringResponse(request).split(System.lineSeparator(), 2)[0];
        String expected = "ok";
        Assert.equalsString("healthz Response", response, expected);

        request = HttpRequest.newBuilder()
                .uri(new URI(testServerURI + Server.pathTokenize))
                .POST(HttpRequest.BodyPublishers.ofString(String.format("{ \"string\": \"%s\" }", basicTestString)))
                .build();

        response = stringResponse(request);
        expected = "{\"tokens\":[{\"pos\":\"NNP\",\"endIndex\":4,\"beginIndex\":0,\"morph\":\"대한민국\"},{\"pos\":\"JX\",\"endIndex\":5,\"beginIndex\":4,\"morph\":\"은\"}]}";
        Assert.equalsString("tokenize Response", response, expected);

    }

    static String stringResponse(HttpRequest request) throws IOException, InterruptedException {
        HttpResponse<String> response = HttpClient.newBuilder()
                .build()
                .send(request, HttpResponse.BodyHandlers.ofString());
        return response.body();
    }

    static class Assert {
        public static void equalsString(String name, String str1, String str2) {
            System.out.print("[TEST] " + name + ": ");
            if (!str1.equals(str2)) {
                System.out.println();
                StringUtils.printDifference(str1, str2);
                System.exit(-1);
            }
            System.out.println("✓");
        }
    }

    // from:
    // https://stackoverflow.com/questions/12089967/find-difference-between-two-strings
    static class StringUtils {
        public static void printDifference(String str1, String str2) {
            System.out.println(str1);
            System.out.println("===");
            System.out.println(str2);
            System.out.println("---");
            System.out.println(difference(str1, str2));
        }

        public static final String EMPTY = "";
        public static final int INDEX_NOT_FOUND = -1;

        public static String difference(String str1, String str2) {
            if (str1 == null) {
                return str2;
            }
            if (str2 == null) {
                return str1;
            }
            int at = indexOfDifference(str1, str2);
            if (at == INDEX_NOT_FOUND) {
                return EMPTY;
            }
            return str2.substring(at);
        }

        public static int indexOfDifference(CharSequence cs1, CharSequence cs2) {
            if (cs1 == cs2) {
                return INDEX_NOT_FOUND;
            }
            if (cs1 == null || cs2 == null) {
                return 0;
            }
            int i;
            for (i = 0; i < cs1.length() && i < cs2.length(); ++i) {
                if (cs1.charAt(i) != cs2.charAt(i)) {
                    break;
                }
            }
            if (i < cs2.length() || i < cs1.length()) {
                return i;
            }
            return INDEX_NOT_FOUND;
        }
    }
}
