import { screen, waitFor } from "@testing-library/react"
import userEvent from '@testing-library/user-event'
import BasicDetails from "pages/app/details/BasicDetails"
import { render } from "utils/test-utils"
import { render as rtlRender } from '@testing-library/react'
import BasicDetailsDropDown from "pages/app/details/BasicDetails/components/BasicDetailsDropdown"
import countryCodeEmoji from "utils/country"
import VaspDetails from "pages/app/details"
import { createMemoryHistory } from 'history'
import { Router } from "react-router-dom"
import { Provider } from "react-redux";
import { configureStore } from "redux/store";
import { ModalProvider } from "contexts/modal"
import { Suspense } from "react"

describe("BasicDetails", () => {

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

    describe("Flag Emoji", () => {

        it("Should show country flag emoji", () => {
            const mockVaspData = {
                "name": "Guidehouse Inc.",
                "vasp": {
                    "entity": {
                        "country_of_registration": "EN",
                        "customer_number": "",
                    },
                },
            }

            render(<BasicDetails data={mockVaspData} />)

            const countryFlagEl = screen.getByTestId(/country-flag/i)
            expect(countryFlagEl.textContent).toBe(countryCodeEmoji(mockVaspData.vasp.entity.country_of_registration))
        })

        it("Should use IVMS101 country when country of registration is empty", () => {
            const mockVaspData = {
                "name": "Guidehouse Inc.",
                "vasp": {
                    "entity": {
                        "country_of_registration": "",
                        "customer_number": "",
                        "geographic_addresses": [
                            {
                                "address_line": [
                                    "150 North Riverside Plaza",
                                    "Suite 2100",
                                    "Chicago, IL 60606"
                                ],
                                "address_type": "ADDRESS_TYPE_CODE_BIZZ",
                                "building_name": "",
                                "building_number": "",
                                "country": "US",
                                "country_sub_division": "",
                                "department": "",
                                "district_name": "",
                                "floor": "",
                                "post_box": "",
                                "post_code": "",
                                "room": "",
                                "street_name": "",
                                "sub_department": "",
                                "town_location_name": "",
                                "town_name": ""
                            }
                        ]
                    },
                },
            }
            render(<BasicDetails data={mockVaspData} />)

            const countryFlagEl = screen.getByTestId(/country-flag/i)
            expect(countryFlagEl.textContent).toBe(countryCodeEmoji(mockVaspData.vasp.entity.geographic_addresses[0].country))
        })
    })

    describe("<VaspDetails />", () => {
        const history = createMemoryHistory()

        it("should render vasp detail correctly", () => {
            const initialValue = {
                VaspDetails: {
                    loading: false,
                    data: {}
                }
            }

            rtlRender(
                <Provider store={configureStore(initialValue)}>
                    <Suspense fallback="loading...">
                        <ModalProvider>
                            <Router history={history}>
                                <VaspDetails />
                            </Router>
                        </ModalProvider>
                    </Suspense>
                </Provider>
            )
        })
    })
})