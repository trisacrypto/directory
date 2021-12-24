import faker from "faker"
import ReviewNote from "pages/app/details/ReviewNotes/ReviewNote"
import { render, screen } from "utils/test-utils"


describe("ReviewNote", () => {
    const vaspId = faker.datatype.uuid()
    const note = {
        id: faker.datatype.uuid(),
        created: faker.date.recent(),
        modified: "",
        author: faker.internet.email(),
        editor: "",
        text: faker.lorem.text()
    }

    beforeAll(() => {
        global.confirm = () => true
    })

    it("should render review note", () => {

        render(<ReviewNote note={note} vaspId={vaspId} />)

        const noteEl = screen.getByTestId("note")
        const authorEl = screen.getByTestId("author")

        expect(note.text).toBe(noteEl.textContent)
        expect(note.author).toBe(authorEl.textContent)
    })
})