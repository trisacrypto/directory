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
      .emotion-0 {
        display: -webkit-box;
        display: -webkit-flex;
        display: -ms-flexbox;
        display: flex;
        -webkit-align-items: center;
        -webkit-box-align: center;
        -ms-flex-align: center;
        align-items: center;
        -webkit-flex-direction: column;
        -ms-flex-direction: column;
        flex-direction: column;
        width: 100%;
        text-align: center;
      }

      .emotion-0>*:not(style)~*:not(style) {
        margin-top: 0.5rem;
        -webkit-margin-end: 0px;
        margin-inline-end: 0px;
        margin-bottom: 0px;
        -webkit-margin-start: 0px;
        margin-inline-start: 0px;
      }

      .emotion-1 {
        width: 1em;
        height: 1em;
        display: inline-block;
        line-height: 1em;
        -webkit-flex-shrink: 0;
        -ms-flex-negative: 0;
        flex-shrink: 0;
        color: var(--chakra-colors-gray-300);
        vertical-align: middle;
        font-size: 5rem;
      }

      .emotion-2 {
        text-transform: capitalize;
      }

      <div>
        <div
          class="chakra-stack emotion-0"
        >
          <svg
            class="chakra-icon emotion-1"
            focusable="false"
            viewBox="0 0 24 24"
          >
            <path
              d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.42 0-8-3.58-8-8 0-1.85.63-3.55 1.69-4.9L16.9 18.31C15.55 19.37 13.85 20 12 20zm6.31-3.1L7.1 5.69C8.45 4.63 10.15 4 12 4c4.42 0 8 3.58 8 8 0 1.85-.63 3.55-1.69 4.9z"
              fill="currentColor"
            />
          </svg>
          <p
            class="chakra-text emotion-2"
            data-testid="label"
          >
            No Data
          </p>
        </div>
      </div>
    `);
  });
});
