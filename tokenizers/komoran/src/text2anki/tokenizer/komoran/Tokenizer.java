package text2anki.tokenizer.komoran;

import kr.co.shineware.nlp.komoran.constant.DEFAULT_MODEL;
import kr.co.shineware.nlp.komoran.core.Komoran;
import kr.co.shineware.nlp.komoran.model.KomoranResult;
import kr.co.shineware.nlp.komoran.model.Token;

import org.json.JSONObject;

import java.util.List;
import java.util.Map;
import java.util.HashMap;

public class Tokenizer {
 public static String tokenList(String strToAnalyze) {
   Komoran komoran = new Komoran(DEFAULT_MODEL.FULL);
   KomoranResult analyzeResultList = komoran.analyze(strToAnalyze);

   System.out.println(analyzeResultList.getPlainText());
   
   List<Token> tokenList = analyzeResultList.getTokenList();
   for (Token token : tokenList) {
     System.out.format("(%2d, %2d) %s/%s\n", token.getBeginIndex(), token.getEndIndex(), token.getMorph(), token.getPos());
   }
   Map<String, Object> map = new HashMap<>();
   map.put("token_list", tokenList);
   return new JSONObject(map).toString();
 }

 public static String testy() {
   return tokenList("대한민국은 민주공화국이다.");
 }
}
