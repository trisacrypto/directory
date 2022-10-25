import { render, screen } from 'utils/test-utils';
import Print from '../Print'
import userEvent from '@testing-library/user-event'

describe('<Print />', () => {

    it('should call handlePrint callback', () => {
        const handlePrint = jest.fn()

        render(<Print onPrint={handlePrint} />)
        userEvent.click(screen.getByTestId(/print-btn/i))

        expect(handlePrint).toHaveBeenCalledTimes(1)
    })

    it('should call handlePrint', () => {
        const handlePrint = jest.fn()

        render(<Print onPrint={handlePrint} />)
        userEvent.click(screen.getByTestId(/print-btn/i))

        expect(screen.getByTestId(/print-btn/i).textContent).toBe('Print')
    })
})