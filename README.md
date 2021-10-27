# text2anki

Create Anki word based Flashcards from text. Card format is:

- Front: Sentence, Text-to-Speech mp3 of sentence
- Back: Front, Word(s) in sentence, word rarity, definition of word/setnence

Procress:

1. Take a string as text
1. Tokenize the text into parts of speech tokens
1. Select token to find definition via. Dictionary API
1. Auto-fill in flashcard fields with Dictionary selection
1. Repeat from above, after a CSV for the flashcards and a folder of Text-to-Speech mp3s are generated

In the future, `text2anki` will have a UI to select photos and match subtitles with existing audio/video.

## Support

Langauages:

- Korean

Systems Requirements: macOS with Java
