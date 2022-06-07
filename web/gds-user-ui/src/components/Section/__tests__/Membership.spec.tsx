import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import MembershipGuide from '../MembershipGuide';

describe('<MembershipGuide />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should match snapshot', () => {
    const { container } = render(<MembershipGuide />);

    expect(container).toMatchSnapshot();
  });
});
