import { render, screen } from "utils/test-utils"
import LastUpdatedColumn from "../LastUpdatedColumn"
import faker from 'faker'
import dayjs from 'dayjs'

describe('<LastUpdatedColumn />', () => {

    it('should display data correctly', () => {
        const row = {
            original: {
                id: faker.datatype.uuid(),
                last_updated: faker.date.recent()
            }
        }
        render(<LastUpdatedColumn row={row} />)

        expect(screen.getByTestId('last_updated').textContent).toBe(
            dayjs(row?.original?.last_updated).fromNow()
        )

    })

    it('should display N/A', () => {
        const row = {
            original: {
                id: faker.datatype.uuid(),
                last_updated: null
            }
        }
        render(<LastUpdatedColumn row={row} />)

        expect(screen.getByTestId('last_updated').textContent).toBe('N/A')

    })
})