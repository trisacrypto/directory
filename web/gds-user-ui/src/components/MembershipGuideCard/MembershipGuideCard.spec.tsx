import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import MembershipGuideCard from '.';

describe('<MembershipGuideCard />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render correctly props', () => {
    const { container } = render(
      <MembershipGuideCard
        stepNumber={2}
        header="Create Your Account"
        description="Create your TRISA account with your VASP email address"
        buttonText="Learn More"
        link="/auth/signup"
      />
    );

    expect(screen.getByTestId('step').textContent).toBe('Step 2');
    expect(screen.getByTestId('description').textContent).toBe(
      'Create your TRISA account with your VASP email address'
    );
    expect(screen.getByTestId('header').textContent).toBe('Create Your Account');
    expect(screen.getByRole('button', { name: 'Learn More' })).toBeInTheDocument();
  });

  it('should match snapshot', () => {
    const { container } = render(
      <MembershipGuideCard
        stepNumber={2}
        header="Create Your Account"
        description="Create your TRISA account with your VASP email address"
        buttonText="Learn More"
        link="/auth/signup"
      />
    );
    expect(container).toMatchSnapshot();
  });
});
