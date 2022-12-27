import { render, screen, fireEvent } from 'utils/test-utils';
import Print from '../Print'

describe('<Print />', () => {

    it('should call handlePrint', () => {
        const handlePrint = jest.fn()

        render(<Print onPrint={handlePrint} />)
        fireEvent.click(screen.getByTestId(/print-btn/i))

        expect(handlePrint).toHaveBeenCalledTimes(1)
        expect(screen.getByTestId(/print-btn/i).textContent).toBe('Print')
    })

    afterEach(() => {
        jest.resetAllMocks()
    })
})