import { render, screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import CertificateDetails, { IdentityCertificateDropDown } from "pages/app/details/CertificateDetails"


describe('CertificateDetails', () => {

    beforeEach(() => {
        const now = new Date(1638305340000).toISOString()
        Date.now = jest.fn(() => now)

    })

    describe('Ceritifcate state badge', () => {

        it('Shoud be badge with valid label', () => {
            const data = { revoked: false, not_after: "2022-10-01T17:51:54Z" }
            render(<CertificateDetails data={data} />)

            const certificateState = screen.getByTestId('certificate-state')

            expect(certificateState).toHaveClass('bg-primary')
            expect(certificateState.textContent).toBe('valid')
        })

        it('Shoud be badge with revoked label', () => {
            const data = { revoked: true }
            render(<CertificateDetails data={data} />)

            const certificateState = screen.getByTestId('certificate-state')

            expect(certificateState).toHaveClass('bg-danger')
            expect(certificateState.textContent).toBe('revoked')
        })

        it('Shoud be badge with expiring soon label', () => {
            const data = { revoked: false, not_after: "2022-01-01T17:51:54Z" }
            render(<CertificateDetails data={data} />)

            const certificateState = screen.getByTestId('certificate-state')

            expect(certificateState).toHaveClass('bg-warning')
            expect(certificateState.textContent).toBe('expiring soon')

        })

        it('Shoud be badge with expired label', () => {
            const data = { revoked: false, not_after: "2019-01-01T17:51:54Z" }
            render(<CertificateDetails data={data} />)
            const certificateState = screen.getByTestId('certificate-state')

            expect(certificateState).toHaveClass('bg-danger')
            expect(certificateState.textContent).toBe('expired')

        })
    })

    describe('Expiriration Date', () => {

        it('Shoud be success text', () => {
            const data = { revoked: false, not_after: "2022-10-01T17:51:54Z" }
            render(<CertificateDetails data={data} />)

            const expires = screen.getByTestId('certificate-expiration-date')

            expect(expires).toHaveClass('text-success')
        })

        it('Should be warning text', () => {
            const data = { revoked: false, not_after: "2022-01-01T17:51:54Z" }
            render(<CertificateDetails data={data} />)

            const expires = screen.getByTestId('certificate-expiration-date')

            expect(expires).toHaveClass('text-warning')
        })

        it('Should be red text', () => {
            const data = { revoked: false, not_after: "2019-10-01T17:51:54Z" }
            render(<CertificateDetails data={data} />)

            const expires = screen.getByTestId('certificate-expiration-date')
            expect(expires).toHaveClass('text-danger')
        })
    })

    describe('IdentityCertificateDropDown', () => {
        Object.assign(navigator, {
            clipboard: {
                writeText: () => { },
            }
        })

        beforeEach(() => {
            jest.spyOn(navigator.clipboard, "writeText");
        })

        it('should copy signature into the clipboard', async () => {
            const data = {
                signature: 'qXW5p8Viu4MsY/KHQEdXU4XCscIx4DwtUpM9FmRGI6aEx1'
            }

            render(<CertificateDetails data={data} />)
            const threeDots = screen.getByTestId('certificate-details-3-dots')

            await waitFor(() => {
                userEvent.click(threeDots)
            })
            const copySignatureElement = screen.getByTestId('copy-signature')
            userEvent.click(copySignatureElement)
            expect(navigator.clipboard.writeText).toHaveBeenCalledWith(data.signature);

        })
        it('should copy serial number into the clipboard', async () => {
            const data = {
                serial_number: 'HCspYuEx68vw'
            }

            render(<CertificateDetails data={data} />)
            const threeDots = screen.getByTestId('certificate-details-3-dots')

            await waitFor(() => {
                userEvent.click(threeDots)
            })
            const serialNumberElement = screen.getByTestId('copy-serial-number')

            userEvent.click(serialNumberElement)
            expect(navigator.clipboard.writeText).toHaveBeenCalledWith(data.serial_number);
        })
    })
})