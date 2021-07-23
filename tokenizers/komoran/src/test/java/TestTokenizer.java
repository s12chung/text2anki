import text2anki.tokenizer.komoran.Tokenizer;

class TestTokenizer{  
    public static void main(String args[]){  
        String tokens = Tokenizer.getTokens("대한민국은");
        String expectedTokens = "{\"tokens\":[{\"pos\":\"NNP\",\"endIndex\":4,\"beginIndex\":0,\"morph\":\"대한민국\"},{\"pos\":\"JX\",\"endIndex\":5,\"beginIndex\":4,\"morph\":\"은\"}]}";
        if (!tokens.equals(expectedTokens)) {
            System.exit(-1);
        }
    }  
}
