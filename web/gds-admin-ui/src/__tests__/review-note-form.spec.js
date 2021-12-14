import { render, screen, waitFor } from "utils/test-utils"
import ReviewNoteForm from "pages/app/details/ReviewNotes/ReviewNoteForm"
import userEvent from "@testing-library/user-event";


const mockResponse = jest.fn(note => {
    return Promise.resolve({ note });
});

describe("ReviewNoteForm", () => {

    beforeEach(() => {
        render(<ReviewNoteForm handleReviewNoteSubmit={mockResponse} />);
    });

    it("should have a disabled submit button", () => {
        const submitEl = screen.getByText(/submit/i, {
            target: {
                value: 'Submit'
            }
        })
        expect(submitEl).toBeDisabled()
    })

    it("should have submit button enabled", () => {
        const textareaEl = screen.getByPlaceholderText(/Write a review note/i)
        const text = "A review note test"
        userEvent.type(textareaEl, text)

        const submitEl = screen.getByText(/submit/i, {
            target: {
                value: 'Submit'
            }
        })

        expect(textareaEl.value).toBe(text)
        expect(submitEl).not.toBeDisabled()
    })

    it("should submit the review note", async () => {
        const textareaEl = screen.getByPlaceholderText(/Write a review note/i)
        const text = "A review note test"
        userEvent.type(textareaEl, text)

        const submitEl = screen.getByText(/submit/i, {
            target: {
                value: 'Submit'
            }
        })

        await waitFor(() => {
            userEvent.click(submitEl)
        })

        expect(mockResponse).toBeCalledWith(text)
    })

})