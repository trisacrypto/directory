import { render, screen } from "utils/test-utils"
import StatusColumn from "../StatusColumn"
import faker from 'faker'
import { StatusLabel, Status } from '../../../../constants/dashboard';

describe('<StatusColumn />', () => {

    it('should display data correctly', () => {
        const row = {
            original: {
                id: faker.datatype.uuid(),
                verification_status: faker.random.objectElement(Status),
            }
        }
        render(<StatusColumn row={row} />)

        expect(screen.getByTestId('verification_status').textContent).toBe(StatusLabel[row.original.verification_status])

    })

    it('should display N/A', () => {
        const row = {
            original: {
                id: faker.datatype.uuid(),
                verification_status: null
            }
        }
        render(<StatusColumn row={row} />)

        expect(screen.getByTestId('verification_status').textContent).toBe('N/A')

    })
})