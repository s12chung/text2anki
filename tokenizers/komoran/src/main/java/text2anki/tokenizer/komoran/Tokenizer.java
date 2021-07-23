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
 public static String getTokens(String strToAnalyze) {
   Komoran komoran = new Komoran(DEFAULT_MODEL.FULL);
   KomoranResult analyzeResultList = komoran.analyze(strToAnalyze);

   Map<String, Object> map = new HashMap<>();
   map.put("tokens", analyzeResultList.getTokenList());
   return new JSONObject(map).toString();
 }
}
