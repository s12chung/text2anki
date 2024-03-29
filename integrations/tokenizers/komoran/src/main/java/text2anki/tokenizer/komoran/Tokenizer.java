package text2anki.tokenizer.komoran;

import kr.co.shineware.nlp.komoran.constant.DEFAULT_MODEL;
import kr.co.shineware.nlp.komoran.core.Komoran;
import kr.co.shineware.nlp.komoran.model.KomoranResult;
import kr.co.shineware.nlp.komoran.model.Token;

import java.util.List;

public class Tokenizer {
  public static List getTokens(String string) {
    Komoran komoran = new Komoran(DEFAULT_MODEL.FULL);
    KomoranResult analyzeResultList = komoran.analyze(string);
    return analyzeResultList.getTokenList();
  }
}
