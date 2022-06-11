import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import Card from './Card';

describe('<Card />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });
  it('should render', () => {
    render(<Card>Test</Card>);

    expect(screen.getByText(/Test/i)).toBeInTheDocument();
    expect(screen.getByTestId('card')).toHaveStyle({
      border: '2px solid #E5EDF1',
      borderRadius: '10px',
      fontFamily: 'Open Sans'
    });
  });
});
