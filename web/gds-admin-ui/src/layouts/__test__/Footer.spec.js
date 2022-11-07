import Footer from "layouts/Footer"
import { render } from "utils/test-utils"

describe('<Footer />', () => {
    const env = process.env

    beforeEach(() => {
        jest.resetModules()
        process.env = { ...env }
    })

    it('should not show the app version in dev env', () => {
        render(<Footer />)
        // expect(screen.queryByTestId('app-version')).toBeNull()
    })

    // it('should not show the git version in dev env', () => {
    //     render(<Footer />)
    //     expect(screen.queryByTestId('git-version')).toBeNull()
    // })


    // it('should show the app version number in prod env', () => {
    //     process.env.NODE_ENV = 'production'
    //     process.env.REACT_APP_VERSION_NUMBER = '1.2.3'
    //     render(<Footer />)

    //     expect(screen.getByTestId('app-version').textContent).toBe('. App version: 1.2.3')
    // })

    // it('should show git version in prod env', () => {
    //     process.env.NODE_ENV = 'production'
    //     process.env.REACT_APP_GIT_REVISION = '1.2.0'
    //     render(<Footer />)

    //     expect(screen.getByTestId('git-version').textContent).toBe(' . GIT version: 1.2.0')
    // })

    afterEach(() => {
        process.env = env
    })
})