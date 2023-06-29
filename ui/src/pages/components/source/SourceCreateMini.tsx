import React, { ChangeEventHandler, useState } from "react"
import { Form, Link } from "react-router-dom"

const SourceCreateMini: React.FC = () => {
  const [text, setText] = useState<string>("")
  const handleText: ChangeEventHandler<HTMLTextAreaElement> = (e) => setText(e.target.value)
  return (
    <Form action="/sources" method="post">
      <textarea
        name="text"
        value={text}
        placeholder="You may also drag and drop here."
        className={text ? "h-third" : "h-20"}
        onChange={handleText}
      />
      <div className="flex-std mt-half mb-std">
        <div className="flex-grow">
          <Link to="sources/new" className="btn">
            Full
          </Link>
        </div>
        <div className="flex-shrink-0">
          <button type="submit" disabled={!text}>
            Submit
          </button>
        </div>
      </div>
    </Form>
  )
}

export default SourceCreateMini
