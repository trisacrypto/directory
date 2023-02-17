import { render } from '@/utils/test-utils';
import { PendingAndRecentActivity } from '../dashboard';

describe('<Tasks />', () => {
    it('Should keep the same UI', () => {
        const { container } = render(<PendingAndRecentActivity />);

        //
        expect(container).toMatchSnapshot('pending-registrations');
    });
});
