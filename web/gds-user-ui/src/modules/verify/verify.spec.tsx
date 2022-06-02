import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render } from 'utils/test-utils';
import VerifyPage from '.';

describe('<VerifyPage />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  it('should match snapshot', () => {
    const { container } = render(<VerifyPage />);
    expect(container).toMatchSnapshot();
  });
});
