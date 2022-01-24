import { screen, waitFor } from "@testing-library/react"
import userEvent from '@testing-library/user-event'
import BasicDetails from "pages/app/details/BasicDetails"
import { render } from "utils/test-utils"
import BasicDetailsDropDown from "pages/app/details/BasicDetails/components/BasicDetailsDropdown"

describe("BasicDetailsDropDown", () => {

    describe("ReviewOption", () => {
        it("Should be enabled when status is PENDING_REVIEW", async () => {

            const mockVaspData = {
                "vasp": {
                    "verification_status": "PENDING_REVIEW",
                }
            }
            render(<BasicDetails data={mockVaspData} />)
            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

            await waitFor(() => userEvent.click(dropdownDots))
            const dropdown = screen.getByTestId(/reviewItem/i)

            expect(dropdown).not.toHaveClass('disabled')
        })

        it("Should be disabled when status is SUBMITTED", async () => {

            const mockVaspData = {
                "vasp": {
                    "verification_status": "SUBMITTED",
                }
            }
            render(<BasicDetails data={mockVaspData} />)
            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

            await waitFor(() => userEvent.click(dropdownDots))
            const dropdown = screen.getByTestId(/reviewItem/i)

            expect(dropdown).toHaveClass('disabled')
        })

        it("Should be disabled when status is EMAIL_VERIFIED", async () => {

            const mockVaspData = {
                "vasp": {
                    "verification_status": "EMAIL_VERIFIED",
                }
            }
            render(<BasicDetails data={mockVaspData} />)
            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

            await waitFor(() => userEvent.click(dropdownDots))
            const dropdown = screen.getByTestId(/reviewItem/i)

            expect(dropdown).toHaveClass('disabled')
        })

        it("Should be disabled when status is NO_VERIFICATION", async () => {

            const mockVaspData = {
                "vasp": {
                    "verification_status": "NO_VERIFICATION",
                }
            }
            render(<BasicDetails data={mockVaspData} />)
            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

            await waitFor(() => userEvent.click(dropdownDots))
            const dropdown = screen.getByTestId(/reviewItem/i)

            expect(dropdown).toHaveClass('disabled')
        })

        it("Should be disabled when status is VERIFIED", async () => {

            const mockVaspData = {
                "vasp": {
                    "verification_status": "VERIFIED",
                }
            }
            render(<BasicDetails data={mockVaspData} />)
            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

            await waitFor(() => userEvent.click(dropdownDots))
            const dropdown = screen.getByTestId(/reviewItem/i)

            expect(dropdown).toHaveClass('disabled')
        })

        it("Should be disabled when status is REJECTED", async () => {

            const mockVaspData = {
                "vasp": {
                    "verification_status": "REJECTED",
                }
            }
            render(<BasicDetails data={mockVaspData} />)
            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

            await waitFor(() => userEvent.click(dropdownDots))
            const dropdown = screen.getByTestId(/reviewItem/i)

            expect(dropdown).toHaveClass('disabled')
        })

        it("Should be disabled when status is ISSUING_CERTIFICATE", async () => {

            const mockVaspData = {
                "vasp": {
                    "verification_status": "ISSUING_CERTIFICATE",
                }
            }
            render(<BasicDetails data={mockVaspData} />)
            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

            await waitFor(() => userEvent.click(dropdownDots))
            const dropdown = screen.getByTestId(/reviewItem/i)

            expect(dropdown).toHaveClass('disabled')
        })

        it("Should be disabled when status is REVIEWED", async () => {

            const mockVaspData = {
                "vasp": {
                    "verification_status": "REVIEWED",
                }
            }
            render(<BasicDetails data={mockVaspData} />)
            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

            await waitFor(() => userEvent.click(dropdownDots))
            const dropdown = screen.getByTestId(/reviewItem/i)

            expect(dropdown).toHaveClass('disabled')
        })

        it("Should be disabled when status is APPEALED", async () => {

            const mockVaspData = {
                "vasp": {
                    "verification_status": "APPEALED",
                }
            }
            render(<BasicDetails data={mockVaspData} />)
            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

            await waitFor(() => userEvent.click(dropdownDots))
            const dropdown = screen.getByTestId(/reviewItem/i)

            expect(dropdown).toHaveClass('disabled')
        })
    })

    describe("DeleteButton", () => {

        it("should be enabled when status is SUBMITTED", async () => {
            const mockVaspData = {
                "vasp": {
                    "verification_status": "SUBMITTED",
                }
            }

            const isNotPendingReviewMock = jest.fn()
            render(<BasicDetailsDropDown vasp={mockVaspData} isNotPendingReview={isNotPendingReviewMock} />)

            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)
            await waitFor(() => userEvent.click(dropdownDots))

            expect(screen.getByRole('button', { name: /delete/i })).toBeEnabled()
        })

        it("should be enabled when status is PENDING_REVIEW", async () => {
            const mockVaspData = {
                "vasp": {
                    "verification_status": "PENDING_REVIEW",
                }
            }


            const isNotPendingReviewMock = jest.fn()
            render(<BasicDetailsDropDown vasp={mockVaspData} isNotPendingReview={isNotPendingReviewMock} />)

            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)
            await waitFor(() => userEvent.click(dropdownDots))

            expect(screen.getByRole('button', { name: /delete/i })).toBeEnabled()
            expect(screen.getByRole('button', { name: /delete/i })).not.toHaveClass("disabled")
        })

        it("should be enabled when status is EMAIL_VERIFIED", async () => {
            const mockVaspData = {
                "vasp": {
                    "verification_status": "EMAIL_VERIFIED",
                }
            }


            const isNotPendingReviewMock = jest.fn()
            render(<BasicDetailsDropDown vasp={mockVaspData} isNotPendingReview={isNotPendingReviewMock} />)

            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)
            await waitFor(() => userEvent.click(dropdownDots))

            expect(screen.getByRole('button', { name: /delete/i })).toBeEnabled()
        })

        it("should be enabled when status is NO_VERIFICATION", async () => {
            const mockVaspData = {
                "vasp": {
                    "verification_status": "NO_VERIFICATION",
                }
            }

            const isNotPendingReviewMock = jest.fn()
            render(<BasicDetailsDropDown vasp={mockVaspData} isNotPendingReview={isNotPendingReviewMock} />)

            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)
            await waitFor(() => userEvent.click(dropdownDots))
            expect(screen.getByRole('button', { name: /delete/i })).not.toHaveClass('disabled')
        })

        it("should be enabled when status is ERRORED", async () => {
            const mockVaspData = {
                "vasp": {
                    "verification_status": "ERRORED",
                }
            }

            const isNotPendingReviewMock = jest.fn()
            render(<BasicDetailsDropDown vasp={mockVaspData} isNotPendingReview={isNotPendingReviewMock} />)

            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)
            await waitFor(() => userEvent.click(dropdownDots))

            expect(screen.getByRole('button', { name: /delete/i })).not.toHaveClass('disabled')
        })

        it("should be disabled when status is VERIFIED", async () => {
            const mockVaspData = {
                "vasp": {
                    "verification_status": "VERIFIED",
                }
            }

            const isNotPendingReviewMock = jest.fn()
            render(<BasicDetailsDropDown vasp={mockVaspData} isNotPendingReview={isNotPendingReviewMock} />)

            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)
            await waitFor(() => userEvent.click(dropdownDots))
            expect(screen.getByRole('button', { name: /delete/i })).toHaveClass('disabled')
        })

        it("should be disabled when status is REJECTED", async () => {
            const mockVaspData = {
                "vasp": {
                    "verification_status": "REJECTED",
                }
            }

            const isNotPendingReviewMock = jest.fn()
            render(<BasicDetailsDropDown vasp={mockVaspData} isNotPendingReview={isNotPendingReviewMock} />)

            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)
            await waitFor(() => userEvent.click(dropdownDots))

            expect(screen.getByRole('button', { name: /delete/i })).toHaveClass('disabled')
        })

        it("should be disabled when status is ISSUING_CERTIFICATE", async () => {
            const mockVaspData = {
                "vasp": {
                    "verification_status": "ISSUING_CERTIFICATE",
                }
            }

            const isNotPendingReviewMock = jest.fn()
            render(<BasicDetailsDropDown vasp={mockVaspData} isNotPendingReview={isNotPendingReviewMock} />)

            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)
            await waitFor(() => userEvent.click(dropdownDots))

            expect(screen.getByRole('button', { name: /delete/i })).toHaveClass('disabled')
        })
        it("should be disabled when status is REVIEWED", async () => {
            const mockVaspData = {
                "vasp": {
                    "verification_status": "REVIEWED",
                }
            }

            const isNotPendingReviewMock = jest.fn()
            render(<BasicDetailsDropDown vasp={mockVaspData} isNotPendingReview={isNotPendingReviewMock} />)

            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)
            await waitFor(() => userEvent.click(dropdownDots))

            expect(screen.getByRole('button', { name: /delete/i })).toHaveClass('disabled')
        })
        it("should be disabled when status is APPEALED", async () => {
            const mockVaspData = {
                "vasp": {
                    "verification_status": "APPEALED",
                }
            }

            const isNotPendingReviewMock = jest.fn()
            render(<BasicDetailsDropDown vasp={mockVaspData} isNotPendingReview={isNotPendingReviewMock} />)

            const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)
            await waitFor(() => userEvent.click(dropdownDots))

            expect(screen.getByRole('button', { name: /delete/i })).toHaveClass('disabled')
        })

    })
})