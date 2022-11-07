import { render, screen } from "utils/test-utils"
import ReviewNoteForm from "pages/app/details/ReviewNotes/ReviewNoteForm"


const mockResponse = jest.fn(note => {
    return Promise.resolve({ note });
});

describe("ReviewNoteForm", () => {
    it("should have a disabled submit button", () => {
        render(<ReviewNoteForm handleReviewNoteSubmit={mockResponse} />);
        const submitEl = screen.getByText(/submit/i, {
            target: {
                value: 'Submit'
            }
        })
        expect(submitEl).toBeDisabled()
    })
})