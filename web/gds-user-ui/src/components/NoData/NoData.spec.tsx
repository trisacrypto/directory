import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import NoData from './NoData';

describe('<NoData />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });
  it('should render label', () => {
    const label = 'No Data Available';
    render(<NoData label={label} />);

    expect(screen.getByTestId(/label/i).textContent).toBe(label);
  });

  it('should render the default label', () => {
    render(<NoData />);
    expect(screen.getByTestId(/label/i).textContent).toBe('No Data');
  });

  it('should match inline snapshot', () => {
    const { container } = render(<NoData />);
    expect(container).toMatchInlineSnapshot(`
      <div>
        <div
          class="chakra-stack css-1bppy43"
        >
          <svg
            class="chakra-icon css-1p10qzu"
            focusable="false"
            viewBox="0 0 24 24"
          >
            <path
              d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.42 0-8-3.58-8-8 0-1.85.63-3.55 1.69-4.9L16.9 18.31C15.55 19.37 13.85 20 12 20zm6.31-3.1L7.1 5.69C8.45 4.63 10.15 4 12 4c4.42 0 8 3.58 8 8 0 1.85-.63 3.55-1.69 4.9z"
              fill="currentColor"
            />
          </svg>
          <p
            class="chakra-text css-10iahqc"
            data-testid="label"
          >
            No Data
          </p>
        </div>
      </div>
    `);
  });
});
