import { Text } from "../../services/models/Source.ts"
import React, { ChangeEventHandler, useMemo, useState } from "react"
import { Form } from "react-router-dom"

const SourceCreate: React.FC = () => {
  const [text, setText] = useState<string>("")
  const handleText: ChangeEventHandler<HTMLTextAreaElement> = (e) => setText(e.target.value)
  const [translation, setTranslation] = useState<string>("")
  const handleTranslation: ChangeEventHandler<HTMLTextAreaElement> = (e) =>
    setTranslation(e.target.value)

  const [previewTexts, valid] = useMemo<[Text[], boolean]>((): [Text[], boolean] => {
    return textsFromTranslation(text, translation)
  }, [text, translation])

  return (
    <Form action="/sources" method="post">
      <div className="flex space-x-basic">
        <div className="flex flex-col items-center grow">
          <div>Source Language</div>
          <textarea name="parts[0].text" value={text} className="h-third" onChange={handleText} />
        </div>
        <div className="flex flex-col items-center grow">
          <div>Translation</div>
          <textarea
            name="parts[0].translation"
            value={translation}
            className="h-third"
            onChange={handleTranslation}
          />
        </div>
      </div>
      <div className="flex justify-end mt-half mb-std">
        <button type="submit" className="btn-primary" disabled={!valid}>
          Submit
        </button>
      </div>
      <PreviewTexts texts={previewTexts} />
    </Form>
  )
}

const PreviewTexts: React.FC<{ readonly texts: Text[] }> = ({ texts }) => {
  return (
    <div className="grid-std">
      {texts.map((text, index) => (
        // eslint-disable-next-line react/no-array-index-key
        <div key={`${text.text}-${text.translation}-${index}`}>
          {Boolean(text.previousBreak) && <br />}
          {text.text ? (
            <>
              {text.text}
              <br />
            </>
          ) : (
            <>
              <b>Missing Text</b>
              <br />
            </>
          )}
          {text.translation ? (
            <>
              {text.translation}
              <br />
            </>
          ) : (
            <>
              <b>Missing Translation</b>
              <br />
            </>
          )}
        </div>
      ))}
    </div>
  )
}

function textsFromTranslation(s: string, translation: string): [Text[], boolean] {
  const lines = split(s)
  const translations = splitClean(translation)
  const longestLength = lines.length > translations.length ? lines.length : translations.length

  const texts: Text[] = Array(longestLength)
  let i = 0
  let previousBreak = false
  let valid = true

  for (let a = 0; a < longestLength; a++) {
    const line = a < lines.length ? lines[a] : ""
    if (a < lines.length && line === "") {
      previousBreak = true
      continue
    }
    const translation = i < translations.length ? translations[i] : ""

    texts[i] = {
      text: line,
      translation,
      previousBreak,
    }
    i++
    previousBreak = false
    if (valid) valid = line !== "" && translation !== ""
  }
  return [texts.slice(0, i), valid]
}

function split(s: string): string[] {
  const lines = s.replace(/\r\n/gu, "\n").split("\n")
  const clean: string[] = Array(lines.length)

  let i = 0
  let previousBreak = false
  for (let line of lines) {
    line = line.trim()
    if (line === "") {
      if (previousBreak || i === 0) continue
    }
    clean[i] = line
    i++
    previousBreak = line === ""
  }

  if (previousBreak) i--
  return clean.slice(0, i)
}

function splitClean(s: string): string[] {
  const lines = s.replace(/\r\n/gu, "\n").split("\n")
  const clean: string[] = Array(lines.length)

  let i = 0
  for (let line of lines) {
    line = line.trim()
    if (line === "") continue
    clean[i] = line
    i++
  }
  return clean.slice(0, i)
}

export default SourceCreate
