import Tasks from "pages/app/dashboard/Tasks"
import { render } from "utils/test-utils"

describe("<Tasks />", () => {

    it("Should keep the same UI", () => {
        const { container } = render(<Tasks />)

        //
        expect(container).toMatchSnapshot("pending-registrations")
    })
})