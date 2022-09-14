import { render, screen } from "utils/test-utils"
import { NameColumn } from ".."
import faker from 'faker'

describe('<NameColumn />', () => {

    it('should display data correctly', () => {
        const row = {
            original: {
                id: faker.datatype.uuid(),
                name: faker.name.findName(),
                common_name: faker.name.findName()
            }
        }
        render(<NameColumn row={row} />)

        expect(screen.getByTestId('commonName').textContent).toBe(row.original.common_name)
        expect(screen.getByTestId('name').textContent).toBe(row.original.name)

    })

    it('should display N/A', () => {
        const row = {
            original: {
                id: faker.datatype.uuid(),
                name: null,
                common_name: null
            }
        }
        render(<NameColumn row={row} />)

        expect(screen.getByTestId('commonName').textContent).toBe('N/A')
        expect(screen.getByTestId('name').textContent).toBe('N/A')

    })
})